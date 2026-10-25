package gomod

import (
	"errors"
	"slices"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Module
		wantErr string // substring of the expected error; "" means no error
	}{
		{
			name:  "minimal",
			input: "module example.com/hello\n\ngo 1.22\n",
			want:  Module{Path: "example.com/hello", GoVersion: "1.22"},
		},
		{
			name: "toolchain and single-line require",
			input: `module github.com/acme/api

go 1.24.0

toolchain go1.27.1

require github.com/google/uuid v1.6.0
`,
			want: Module{
				Path:      "github.com/acme/api",
				GoVersion: "1.24.0",
				Toolchain: "go1.27.1",
				Require:   []Require{{Path: "github.com/google/uuid", Version: "v1.6.0"}},
			},
		},
		{
			name: "require block with comments and indirect",
			input: `// Service module.
module github.com/acme/orders // the module path

go 1.23

require (
	github.com/jackc/pgx/v5 v5.7.1
	// a comment line inside the block
	golang.org/x/sync v0.8.0 // indirect

	golang.org/x/text v0.18.0
)
`,
			want: Module{
				Path:      "github.com/acme/orders",
				GoVersion: "1.23",
				Require: []Require{
					{Path: "github.com/jackc/pgx/v5", Version: "v5.7.1"},
					{Path: "golang.org/x/sync", Version: "v0.8.0", Indirect: true},
					{Path: "golang.org/x/text", Version: "v0.18.0"},
				},
			},
		},
		{
			name: "directives it does not model are skipped, blocks too",
			input: `module example.com/app

go 1.22

replace (
	example.com/lib => ../lib
	example.com/other v1.0.0 => example.com/fork v1.0.1
)

exclude example.com/bad v1.2.3

tool golang.org/x/tools/cmd/stringer

require example.com/lib v1.0.0
`,
			want: Module{
				Path:      "example.com/app",
				GoVersion: "1.22",
				Require:   []Require{{Path: "example.com/lib", Version: "v1.0.0"}},
			},
		},
		{
			name:    "missing module directive",
			input:   "go 1.22\n",
			wantErr: "missing module directive",
		},
		{
			name:    "require without a version names the line",
			input:   "module example.com/app\n\ngo 1.22\n\nrequire (\n\texample.com/lib\n)\n",
			wantErr: "line 6",
		},
		{
			name:    "module with two arguments",
			input:   "module example.com/a example.com/b\n",
			wantErr: "line 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse([]byte(tt.input))
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("Parse() error = %v, want an error containing %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Parse() unexpected error: %v", err)
			}
			assertModule(t, got, tt.want)
		})
	}
}

func TestParseMissingModuleIsErrNoModule(t *testing.T) {
	_, err := Parse([]byte("go 1.22\n"))
	if !errors.Is(err, ErrNoModule) {
		t.Errorf("Parse() error = %v, want errors.Is(err, ErrNoModule)", err)
	}
}

func assertModule(t *testing.T, got, want Module) {
	t.Helper()
	if got.Path != want.Path || got.GoVersion != want.GoVersion || got.Toolchain != want.Toolchain {
		t.Errorf("Parse() = {Path:%q GoVersion:%q Toolchain:%q}, want {Path:%q GoVersion:%q Toolchain:%q}",
			got.Path, got.GoVersion, got.Toolchain, want.Path, want.GoVersion, want.Toolchain)
	}
	if !slices.Equal(got.Require, want.Require) {
		t.Errorf("Parse().Require = %+v, want %+v", got.Require, want.Require)
	}
}
