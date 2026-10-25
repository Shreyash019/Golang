package scan

import (
	"slices"
	"strings"
	"testing"
	"testing/fstest"
)

func file(src string) *fstest.MapFile { return &fstest.MapFile{Data: []byte(src)} }

func TestScan(t *testing.T) {
	tests := []struct {
		name string
		fsys fstest.MapFS
		want []Package
	}{
		{
			name: "empty tree",
			fsys: fstest.MapFS{"README.md": file("# hi\n")},
			want: nil,
		},
		{
			name: "command and libraries",
			fsys: fstest.MapFS{
				"go.mod":                       file("module example.com/app\n"),
				"cmd/api/main.go":              file("package main\n\nimport (\n\t\"fmt\"\n\t\"example.com/app/internal/order\"\n)\n\nfunc main() { fmt.Println(order.X) }\n"),
				"internal/order/order.go":      file("// Package order handles orders.\npackage order\n\nimport \"fmt\"\n\nvar X = fmt.Sprint(1)\n"),
				"internal/order/store.go":      file("package order\n\nimport (\n\t\"fmt\"\n\t\"sync\"\n)\n\nvar _ sync.Mutex\nvar _ = fmt.Sprint\n"),
				"internal/order/order_test.go": file("package order\n"),
			},
			want: []Package{
				{Dir: "cmd/api", Name: "main", Files: 1, Imports: []string{"example.com/app/internal/order", "fmt"}},
				{Dir: "internal/order", Name: "order", Files: 2, TestFiles: 1, Imports: []string{"fmt", "sync"}, HasDoc: true},
			},
		},
		{
			name: "root package and ignored names",
			fsys: fstest.MapFS{
				"doc.go":                file("// Package lib is a library.\npackage lib\n"),
				"testdata/fixture.go":   file("package fixture\n"),
				"vendor/x/x.go":         file("package x\n"),
				".git/hooks/h.go":       file("package hooks\n"),
				"_old/old.go":           file("package old\n"),
				"_scratch.go":           file("package scratch\n"),
				"only_tests/a_test.go":  file("package onlytests\n"),
				"internal/util/util.go": file("package util\n"),
			},
			want: []Package{
				{Dir: ".", Name: "lib", Files: 1, HasDoc: true},
				{Dir: "internal/util", Name: "util", Files: 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Scan(tt.fsys)
			if err != nil {
				t.Fatalf("Scan() unexpected error: %v", err)
			}
			if !slices.EqualFunc(got, tt.want, equalPackage) {
				t.Errorf("Scan() =\n  %+v\nwant\n  %+v", got, tt.want)
			}
		})
	}
}

func TestScanErrors(t *testing.T) {
	tests := []struct {
		name    string
		fsys    fstest.MapFS
		wantErr string
	}{
		{
			name:    "two package names in one directory",
			fsys:    fstest.MapFS{"a/one.go": file("package one\n"), "a/two.go": file("package two\n")},
			wantErr: "a: multiple packages one and two",
		},
		{
			name:    "syntax error names the file",
			fsys:    fstest.MapFS{"bad/bad.go": file("pakage bad\n")},
			wantErr: "parse bad/bad.go",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Scan(tt.fsys)
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("Scan() error = %v, want an error containing %q", err, tt.wantErr)
			}
		})
	}
}

func TestIsCommand(t *testing.T) {
	if !(Package{Name: "main"}).IsCommand() {
		t.Error(`Package{Name: "main"}.IsCommand() = false, want true`)
	}
	if (Package{Name: "order"}).IsCommand() {
		t.Error(`Package{Name: "order"}.IsCommand() = true, want false`)
	}
}

func equalPackage(a, b Package) bool {
	return a.Dir == b.Dir && a.Name == b.Name && a.Files == b.Files && a.TestFiles == b.TestFiles &&
		a.HasDoc == b.HasDoc && slices.Equal(a.Imports, b.Imports)
}
