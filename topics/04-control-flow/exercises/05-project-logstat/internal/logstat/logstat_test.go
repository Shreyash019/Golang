package logstat

import (
	"bufio"
	"errors"
	"strings"
	"testing"
)

func TestParseLevel(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in     string
		want   Level
		wantOK bool
	}{
		{"debug", LevelDebug, true},
		{"INFO", LevelInfo, true},
		{"Warn", LevelWarn, true},
		{"warning", LevelWarn, true},
		{"error", LevelError, true},
		{"FATAL", LevelFatal, true},
		{"trace", LevelUnknown, false},
		{"", LevelUnknown, false},
		{"err", LevelUnknown, false},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			t.Parallel()
			got, ok := ParseLevel(tt.in)
			if got != tt.want || ok != tt.wantOK {
				t.Errorf("ParseLevel(%q) = %d, %t, want %d, %t", tt.in, got, ok, tt.want, tt.wantOK)
			}
		})
	}
}

const sample = `# service log
2026-11-04T10:00:00Z INFO server started
2026-11-04T10:00:01Z DEBUG cache warmed

2026-11-04T10:00:02Z WARN slow query took 1.2s
garbage
  2026-11-04T10:00:03Z ERROR db: connection refused
2026-11-04T10:00:04Z info request served
2026-11-04T10:00:05Z FATAL out of memory
2026-11-04T10:00:06Z TRACE very chatty
2026-11-04T10:00:07Z ERROR
`

func TestFilter(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    string
		minLevel Level
		limit    int
		wantOut  string
		want     Summary
	}{
		{
			name:     "info and above",
			input:    sample,
			minLevel: LevelInfo,
			wantOut: `2026-11-04T10:00:00Z INFO server started
2026-11-04T10:00:02Z WARN slow query took 1.2s
  2026-11-04T10:00:03Z ERROR db: connection refused
2026-11-04T10:00:04Z info request served
2026-11-04T10:00:05Z FATAL out of memory
`,
			want: Summary{Lines: 11, Matched: 5, Skipped: 2, Malformed: 3, Debug: 1, Info: 2, Warn: 1, Error: 1, Fatal: 1},
		},
		{
			name:     "error and above",
			input:    sample,
			minLevel: LevelError,
			wantOut: `  2026-11-04T10:00:03Z ERROR db: connection refused
2026-11-04T10:00:05Z FATAL out of memory
`,
			want: Summary{Lines: 11, Matched: 2, Skipped: 2, Malformed: 3, Debug: 1, Info: 2, Warn: 1, Error: 1, Fatal: 1},
		},
		{
			name:     "limit stops reading after the 2nd match",
			input:    sample,
			minLevel: LevelInfo,
			limit:    2,
			wantOut: `2026-11-04T10:00:00Z INFO server started
2026-11-04T10:00:02Z WARN slow query took 1.2s
`,
			want: Summary{Lines: 5, Matched: 2, Skipped: 2, Debug: 1, Info: 1, Warn: 1},
		},
		{
			name:     "limit larger than the matches",
			input:    sample,
			minLevel: LevelFatal,
			limit:    5,
			wantOut:  "2026-11-04T10:00:05Z FATAL out of memory\n",
			want:     Summary{Lines: 11, Matched: 1, Skipped: 2, Malformed: 3, Debug: 1, Info: 2, Warn: 1, Error: 1, Fatal: 1},
		},
		{
			name:     "LevelUnknown prints every valid line",
			input:    "t DEBUG a\nt INFO b\n",
			minLevel: LevelUnknown,
			wantOut:  "t DEBUG a\nt INFO b\n",
			want:     Summary{Lines: 2, Matched: 2, Debug: 1, Info: 1},
		},
		{
			name:  "empty input",
			input: "",
			want:  Summary{},
		},
		{
			name:     "no trailing newline",
			input:    "t WARN disk at 91%",
			minLevel: LevelWarn,
			wantOut:  "t WARN disk at 91%\n",
			want:     Summary{Lines: 1, Matched: 1, Warn: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var out strings.Builder
			got, err := Filter(strings.NewReader(tt.input), &out, tt.minLevel, tt.limit)
			if err != nil {
				t.Fatalf("Filter() error = %v, want nil", err)
			}
			if got != tt.want {
				t.Errorf("Filter() summary = %+v, want %+v", got, tt.want)
			}
			if out.String() != tt.wantOut {
				t.Errorf("Filter() wrote:\n%s\nwant:\n%s", out.String(), tt.wantOut)
			}
		})
	}
}

func TestFilterLineTooLong(t *testing.T) {
	t.Parallel()
	input := "t INFO ok\nt INFO " + strings.Repeat("x", bufio.MaxScanTokenSize) + "\n"
	var out strings.Builder
	s, err := Filter(strings.NewReader(input), &out, LevelInfo, 0)
	if !errors.Is(err, bufio.ErrTooLong) {
		t.Fatalf("Filter() error = %v, want an error wrapping bufio.ErrTooLong", err)
	}
	if s.Matched != 1 || out.String() != "t INFO ok\n" {
		t.Errorf("Filter() before the long line: matched %d, wrote %q; want 1 and %q", s.Matched, out.String(), "t INFO ok\n")
	}
}

type failingWriter struct{}

var errDiskFull = errors.New("disk full")

func (failingWriter) Write([]byte) (int, error) { return 0, errDiskFull }

func TestFilterWriteError(t *testing.T) {
	t.Parallel()
	_, err := Filter(strings.NewReader("t ERROR boom\n"), failingWriter{}, LevelInfo, 0)
	if !errors.Is(err, errDiskFull) {
		t.Errorf("Filter() error = %v, want an error wrapping %v", err, errDiskFull)
	}
}
