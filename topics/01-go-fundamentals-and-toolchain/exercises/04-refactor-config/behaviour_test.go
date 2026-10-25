package main

import (
	"strings"
	"testing"
)

// These tests pin the program's behaviour. They pass before your refactor
// and must still pass after it.
func TestRun(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    string
		wantErr string
	}{
		{name: "defaults", args: nil, want: "starting server on :8080 (env=dev, debug=false)\n"},
		{name: "all keys", args: []string{"port=9090", "env=prod", "debug=true"}, want: "starting server on :9090 (env=prod, debug=true)\n"},
		{name: "later value wins", args: []string{"port=1", "port=2"}, want: "starting server on :2 (env=dev, debug=false)\n"},
		{name: "port not a number", args: []string{"port=abc"}, wantErr: `invalid port "abc"`},
		{name: "port out of range", args: []string{"port=70000"}, wantErr: `invalid port "70000"`},
		{name: "bad debug", args: []string{"debug=yes"}, wantErr: `invalid debug "yes"`},
		{name: "unknown key", args: []string{"host=localhost"}, wantErr: `unknown key "host"`},
		{name: "not key=value", args: []string{"verbose"}, wantErr: `invalid argument "verbose"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sb strings.Builder
			err := run(tt.args, &sb)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("run(%q) error = %v, want an error containing %q", tt.args, err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("run(%q) unexpected error: %v", tt.args, err)
			}
			if got := sb.String(); got != tt.want {
				t.Errorf("run(%q) printed %q, want %q", tt.args, got, tt.want)
			}
		})
	}
}
