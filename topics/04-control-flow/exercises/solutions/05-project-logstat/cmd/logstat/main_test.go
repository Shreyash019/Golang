package main

import (
	"os"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	t.Parallel()
	log, err := os.ReadFile("../../testdata/app.log")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		args       []string
		wantStdout string
		wantStderr string
	}{
		{
			args: nil,
			wantStdout: `2026-11-04T10:00:00Z INFO server started
2026-11-04T10:00:02Z WARN slow query took 1.2s
2026-11-04T10:00:03Z ERROR db: connection refused
2026-11-04T10:00:04Z info request served
2026-11-04T10:00:05Z FATAL out of memory
`,
			wantStderr: "lines=11 matched=5 skipped=2 malformed=3 debug=1 info=2 warn=1 error=1 fatal=1\n",
		},
		{
			args: []string{"-level", "warn", "-n", "2"},
			wantStdout: `2026-11-04T10:00:02Z WARN slow query took 1.2s
2026-11-04T10:00:03Z ERROR db: connection refused
`,
			wantStderr: "lines=7 matched=2 skipped=2 malformed=1 debug=1 info=1 warn=1 error=1 fatal=0\n",
		},
	}
	for _, tt := range tests {
		t.Run(strings.Join(tt.args, " "), func(t *testing.T) {
			t.Parallel()
			var stdout, stderr strings.Builder
			if err := run(tt.args, strings.NewReader(string(log)), &stdout, &stderr); err != nil {
				t.Fatalf("run(%q) error = %v", tt.args, err)
			}
			if stdout.String() != tt.wantStdout {
				t.Errorf("run(%q) stdout:\n%s\nwant:\n%s", tt.args, stdout.String(), tt.wantStdout)
			}
			if stderr.String() != tt.wantStderr {
				t.Errorf("run(%q) stderr = %q, want %q", tt.args, stderr.String(), tt.wantStderr)
			}
		})
	}
}

func TestRunErrors(t *testing.T) {
	t.Parallel()
	for _, args := range [][]string{
		{"-level", "loud"},
		{"-n", "-1"},
		{"app.log"},
		{"-unknown-flag"},
	} {
		var stdout, stderr strings.Builder
		if err := run(args, strings.NewReader(""), &stdout, &stderr); err == nil {
			t.Errorf("run(%q) = nil; want an error", args)
		}
	}
}
