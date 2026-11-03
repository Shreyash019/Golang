package main

import (
	"strings"
	"testing"
)

// These tests pin the program's behaviour. They pass before the refactor
// and must still pass after it.
func TestRun(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "default attempts",
			want: `attempt 1: wait 100ms
attempt 2: wait 200ms
attempt 3: wait 400ms
attempt 4: wait 800ms
attempt 5: wait 1.6s
total wait: 3.1s
`,
		},
		{
			name: "fewer attempts",
			args: []string{"3"},
			want: `attempt 1: wait 100ms
attempt 2: wait 200ms
attempt 3: wait 400ms
total wait: 700ms
`,
		},
		{
			name: "capped attempts",
			args: []string{"9"},
			want: `requested 9 attempts, capped at 5
attempt 1: wait 100ms
attempt 2: wait 200ms
attempt 3: wait 400ms
attempt 4: wait 800ms
attempt 5: wait 1.6s
total wait: 3.1s
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sb strings.Builder
			if err := run(&sb, tt.args); err != nil {
				t.Fatalf("run(%q) returned error %v", tt.args, err)
			}
			if got := sb.String(); got != tt.want {
				t.Errorf("run(%q) printed:\n%s\nwant:\n%s", tt.args, got, tt.want)
			}
		})
	}
}

// Running twice must give the same result: no state may leak between runs.
func TestRunIsRepeatable(t *testing.T) {
	var first, second strings.Builder
	if err := run(&first, []string{"2"}); err != nil {
		t.Fatal(err)
	}
	if err := run(&second, []string{"2"}); err != nil {
		t.Fatal(err)
	}
	if first.String() != second.String() {
		t.Errorf("second run printed:\n%s\nfirst run printed:\n%s", second.String(), first.String())
	}
}

func TestRunErrors(t *testing.T) {
	for _, arg := range []string{"x", "0", "-2"} {
		var sb strings.Builder
		if err := run(&sb, []string{arg}); err == nil {
			t.Errorf("run([%q]) = nil; want an error", arg)
		}
	}
}
