package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "go.mod", "module example.com/demo\n\ngo 1.22\n")
	writeFile(t, root, "cmd/demo/main.go", "package main\n\nimport \"fmt\"\n\nfunc main() { fmt.Println(\"hi\") }\n")
	writeFile(t, root, "internal/store/store.go", "// Package store keeps things.\npackage store\n")

	tests := []struct {
		name     string
		args     []string
		wantCode int
		wantOut  []string // substrings of stdout
		wantErr  string   // substring of stderr
	}{
		{name: "scan a tree", args: []string{root}, wantCode: 0, wantOut: []string{"cmd/demo", "command", "internal/store", "2 packages (1 command, 1 library)"}},
		{name: "version", args: []string{"-version"}, wantCode: 0, wantOut: []string{"pkgscan dev"}},
		{name: "missing dir", args: []string{filepath.Join(root, "nope")}, wantCode: 1, wantErr: "pkgscan:"},
		{name: "too many args", args: []string{root, root}, wantCode: 2, wantErr: "usage: pkgscan"},
		{name: "unknown flag", args: []string{"-x"}, wantCode: 2, wantErr: "flag provided but not defined"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr strings.Builder
			code := run(tt.args, &stdout, &stderr)
			if code != tt.wantCode {
				t.Errorf("run(%q) = %d, want %d (stderr: %q)", tt.args, code, tt.wantCode, stderr.String())
			}
			for _, want := range tt.wantOut {
				if !strings.Contains(stdout.String(), want) {
					t.Errorf("run(%q) stdout = %q, want it to contain %q", tt.args, stdout.String(), want)
				}
			}
			if !strings.Contains(stderr.String(), tt.wantErr) {
				t.Errorf("run(%q) stderr = %q, want it to contain %q", tt.args, stderr.String(), tt.wantErr)
			}
		})
	}
}

func writeFile(t *testing.T, root, name, content string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
