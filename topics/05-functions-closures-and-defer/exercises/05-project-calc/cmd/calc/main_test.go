package main

import (
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	t.Parallel()
	tests := []struct {
		args []string
		want string
	}{
		{[]string{"2 * (3 + 4)"}, "14\n"},
		{[]string{"1", "+", "2", "*", "3"}, "7\n"},
		{[]string{"--", "-(2 + 3) * 2"}, "-10\n"},
		{[]string{"-trace", "4"}, "> expr\n  > term\n    > factor\n    < factor = 4\n  < term = 4\n< expr = 4\n4\n"},
	}
	for _, tt := range tests {
		t.Run(strings.Join(tt.args, " "), func(t *testing.T) {
			t.Parallel()
			var sb strings.Builder
			if err := run(tt.args, &sb); err != nil {
				t.Fatalf("run(%q) error = %v", tt.args, err)
			}
			if got := sb.String(); got != tt.want {
				t.Errorf("run(%q) printed %q, want %q", tt.args, got, tt.want)
			}
		})
	}
}

func TestRunErrors(t *testing.T) {
	t.Parallel()
	for _, args := range [][]string{
		nil,
		{"-trace"},
		{"1 +"},
		{"8 / 0"},
		{"-unknown", "1"},
	} {
		var sb strings.Builder
		if err := run(args, &sb); err == nil {
			t.Errorf("run(%q) = nil, printed %q; want an error", args, sb.String())
		}
	}
}
