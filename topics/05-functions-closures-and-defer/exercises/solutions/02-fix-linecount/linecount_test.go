package linecount

import (
	"errors"
	"io"
	"io/fs"
	"strings"
	"testing"
)

// fakeFS serves files from memory and records how many are open.
type fakeFS struct {
	files   map[string]string
	open    int // currently open
	maxOpen int // most open at the same time
}

func (f *fakeFS) Open(name string) (io.ReadCloser, error) {
	content, ok := f.files[name]
	if !ok {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}
	f.open++
	f.maxOpen = max(f.maxOpen, f.open)
	return &fakeFile{Reader: strings.NewReader(content), fs: f}, nil
}

type fakeFile struct {
	io.Reader
	fs     *fakeFS
	closed bool
}

func (ff *fakeFile) Close() error {
	if !ff.closed {
		ff.closed = true
		ff.fs.open--
	}
	return nil
}

func newFS() *fakeFS {
	return &fakeFS{files: map[string]string{
		"a.txt":     "one\ntwo\nthree\n",
		"b.txt":     "four\n",
		"c.txt":     "five\nsix",
		"empty.txt": "",
	}}
}

func TestCountLines(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		files []string
		want  int
	}{
		{"none", nil, 0},
		{"one file", []string{"a.txt"}, 3},
		{"several files", []string{"a.txt", "b.txt", "c.txt", "empty.txt"}, 6},
		{"same file twice", []string{"b.txt", "b.txt"}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fsys := newFS()
			got, err := CountLines(tt.files, fsys.Open)
			if err != nil {
				t.Fatalf("CountLines(%q) error = %v", tt.files, err)
			}
			if got != tt.want {
				t.Errorf("CountLines(%q) = %d, want %d", tt.files, got, tt.want)
			}
			if fsys.open != 0 {
				t.Errorf("CountLines(%q) left %d files open, want 0", tt.files, fsys.open)
			}
			if len(tt.files) > 0 && fsys.maxOpen != 1 {
				t.Errorf("CountLines(%q) had %d files open at the same time, want 1: each file must be closed before the next is opened", tt.files, fsys.maxOpen)
			}
		})
	}
}

func TestCountLinesMissingFile(t *testing.T) {
	t.Parallel()
	fsys := newFS()
	files := []string{"a.txt", "b.txt", "missing.txt"}
	_, err := CountLines(files, fsys.Open)
	if !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("CountLines(%q) error = %v, want an error wrapping fs.ErrNotExist", files, err)
	}
	if fsys.open != 0 {
		t.Errorf("CountLines(%q) left %d files open after the error, want 0", files, fsys.open)
	}
	if fsys.maxOpen != 1 {
		t.Errorf("CountLines(%q) had %d files open at the same time, want 1", files, fsys.maxOpen)
	}
}
