package main

import (
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name string
		vars map[string]string
		want string
	}{
		{
			name: "defaults",
			want: "port=8080 timeout=5s log_level=info max_body_bytes=1048576\n",
		},
		{
			name: "overrides",
			vars: map[string]string{"APP_PORT": "9090", "APP_LOG_LEVEL": "debug", "APP_TIMEOUT": "1m30s"},
			want: "port=9090 timeout=1m30s log_level=debug max_body_bytes=1048576\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sb strings.Builder
			if err := run(&sb, lookup(tt.vars)); err != nil {
				t.Fatalf("run returned error %v", err)
			}
			if got := sb.String(); got != tt.want {
				t.Errorf("run printed %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRunInvalid(t *testing.T) {
	var sb strings.Builder
	err := run(&sb, lookup(map[string]string{"APP_PORT": "abc"}))
	if err == nil {
		t.Fatalf("run with APP_PORT=abc = nil, printed %q; want an error", sb.String())
	}
	if sb.Len() != 0 {
		t.Errorf("run printed %q on error; want nothing on stdout", sb.String())
	}
}

func lookup(vars map[string]string) func(string) (string, bool) {
	return func(key string) (string, bool) {
		v, ok := vars[key]
		return v, ok
	}
}
