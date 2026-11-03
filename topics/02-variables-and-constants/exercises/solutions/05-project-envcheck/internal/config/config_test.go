package config

import (
	"strings"
	"testing"
	"time"
)

// env returns a LookupFunc backed by a map, standing in for os.LookupEnv.
func env(vars map[string]string) LookupFunc {
	return func(key string) (string, bool) {
		v, ok := vars[key]
		return v, ok
	}
}

func TestParseLevel(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in   string
		want Level
	}{
		{"debug", LevelDebug},
		{"info", LevelInfo},
		{"warn", LevelWarn},
		{"error", LevelError},
		{"WARN", LevelWarn},
		{"Error", LevelError},
	}
	for _, tt := range tests {
		got, err := ParseLevel(tt.in)
		if err != nil || got != tt.want {
			t.Errorf("ParseLevel(%q) = %v, %v; want %v, nil", tt.in, got, err, tt.want)
		}
	}
	for _, in := range []string{"", "verbose", " info", "warning"} {
		if got, err := ParseLevel(in); err == nil {
			t.Errorf("ParseLevel(%q) = %v, nil; want an error", in, got)
		}
	}
}

func TestLoad(t *testing.T) {
	t.Parallel()
	defaults := Config{Port: 8080, Timeout: 5 * time.Second, Level: LevelInfo, MaxBodyBytes: 1 << 20}
	tests := []struct {
		name string
		vars map[string]string
		want Config
	}{
		{name: "nothing set", vars: nil, want: defaults},
		{
			name: "empty values take defaults",
			vars: map[string]string{"APP_PORT": "", "APP_TIMEOUT": "", "APP_LOG_LEVEL": "", "APP_MAX_BODY_BYTES": ""},
			want: defaults,
		},
		{
			name: "everything set",
			vars: map[string]string{"APP_PORT": "9090", "APP_TIMEOUT": "250ms", "APP_LOG_LEVEL": "WARN", "APP_MAX_BODY_BYTES": "2048"},
			want: Config{Port: 9090, Timeout: 250 * time.Millisecond, Level: LevelWarn, MaxBodyBytes: 2048},
		},
		{
			name: "one set, the rest default",
			vars: map[string]string{"APP_TIMEOUT": "1m"},
			want: Config{Port: 8080, Timeout: time.Minute, Level: LevelInfo, MaxBodyBytes: 1 << 20},
		},
		{
			name: "limits are inclusive",
			vars: map[string]string{"APP_PORT": "65535", "APP_MAX_BODY_BYTES": "1073741824"},
			want: Config{Port: 65535, Timeout: 5 * time.Second, Level: LevelInfo, MaxBodyBytes: 1 << 30},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := Load(env(tt.vars))
			if err != nil {
				t.Fatalf("Load(%v) returned error %v", tt.vars, err)
			}
			if got != tt.want {
				t.Errorf("Load(%v) = %+v, want %+v", tt.vars, got, tt.want)
			}
		})
	}
}

func TestLoadErrors(t *testing.T) {
	t.Parallel()
	tests := []struct {
		vars       map[string]string
		wantPrefix string
	}{
		{map[string]string{"APP_PORT": "abc"}, "APP_PORT: "},
		{map[string]string{"APP_PORT": "0"}, "APP_PORT: "},
		{map[string]string{"APP_PORT": "70000"}, "APP_PORT: "},
		{map[string]string{"APP_TIMEOUT": "soon"}, "APP_TIMEOUT: "},
		{map[string]string{"APP_TIMEOUT": "0s"}, "APP_TIMEOUT: "},
		{map[string]string{"APP_TIMEOUT": "-1s"}, "APP_TIMEOUT: "},
		{map[string]string{"APP_LOG_LEVEL": "verbose"}, "APP_LOG_LEVEL: "},
		{map[string]string{"APP_MAX_BODY_BYTES": "big"}, "APP_MAX_BODY_BYTES: "},
		{map[string]string{"APP_MAX_BODY_BYTES": "0"}, "APP_MAX_BODY_BYTES: "},
		{map[string]string{"APP_MAX_BODY_BYTES": "1073741825"}, "APP_MAX_BODY_BYTES: "},
		{map[string]string{"APP_PORT": "abc", "APP_TIMEOUT": "soon"}, "APP_PORT: "},
	}
	for _, tt := range tests {
		got, err := Load(env(tt.vars))
		if err == nil {
			t.Errorf("Load(%v) = %+v, nil; want an error starting with %q", tt.vars, got, tt.wantPrefix)
			continue
		}
		if !strings.HasPrefix(err.Error(), tt.wantPrefix) {
			t.Errorf("Load(%v) error = %q, want it to start with %q", tt.vars, err, tt.wantPrefix)
		}
	}
}
