package main

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/Shreyash019/Golang/topics/02-variables-and-constants/internal/examples"
)

func TestRun(t *testing.T) {
	fake := []examples.Example{
		{Name: "ok", Section: "1", Summary: "prints ok", Run: func(w io.Writer) { fmt.Fprint(w, "ok") }},
		{Name: "boom", Section: "2", Summary: "panics", Crashes: true, Run: func(io.Writer) { panic("boom") }},
	}

	tests := []struct {
		name    string
		args    []string
		want    string
		wantErr bool
	}{
		{name: "list", args: nil, want: "crashes on purpose"},
		{name: "one", args: []string{"ok"}, want: "ok"},
		{name: "all skips crashing examples", args: []string{"all"}, want: "=== ok (§1)\nok\n"},
		{name: "unknown", args: []string{"nope"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sb strings.Builder
			err := run(tt.args, &sb, fake)
			if (err != nil) != tt.wantErr {
				t.Fatalf("run(%v) error = %v, wantErr %v", tt.args, err, tt.wantErr)
			}
			if !strings.Contains(sb.String(), tt.want) {
				t.Errorf("run(%v) output = %q, want it to contain %q", tt.args, sb.String(), tt.want)
			}
		})
	}
}
