package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode"
)

// TestIdioms fails while the Java habits are still in the code. Each subtest
// names one habit and the idiomatic alternative.
func TestIdioms(t *testing.T) {
	fset := token.NewFileSet()
	paths, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	var files []*ast.File
	for _, p := range paths {
		if strings.HasSuffix(p, "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, p, nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		files = append(files, f)
	}

	t.Run("no utils file", func(t *testing.T) {
		if _, err := os.Stat("utils.go"); err == nil {
			t.Error("utils.go exists: put each helper next to the code that uses it, or use the standard library (strings.Cut)")
		}
	})

	t.Run("no Get prefix", func(t *testing.T) {
		eachFunc(files, func(name string, pos token.Pos) {
			if hasWordPrefix(name, "Get") {
				t.Errorf("%s: %s has a Get prefix: Go getters are named after the field (Port, not GetPort), and plain exported fields need no getter at all", fset.Position(pos), name)
			}
		})
	})

	t.Run("no I-prefixed interfaces or Impl types", func(t *testing.T) {
		eachType(files, func(name string, pos token.Pos) {
			if hasWordPrefix(name, "I") || strings.HasSuffix(name, "Impl") {
				t.Errorf("%s: type %s: name types for what they are (Config), and declare interfaces where they are used, only when you need one", fset.Position(pos), name)
			}
			if strings.HasSuffix(name, "Utils") || strings.HasSuffix(name, "Manager") || strings.HasSuffix(name, "Exception") {
				t.Errorf("%s: type %s: Utils/Manager/Exception are Java names; Go uses plain functions, data types and error values", fset.Position(pos), name)
			}
		})
	})

	t.Run("no panic or recover for errors", func(t *testing.T) {
		for _, f := range files {
			ast.Inspect(f, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}
				if id, ok := call.Fun.(*ast.Ident); ok && (id.Name == "panic" || id.Name == "recover") {
					t.Errorf("%s: %s used for error handling: return an error value instead", fset.Position(call.Pos()), id.Name)
				}
				return true
			})
		}
	})
}

func eachFunc(files []*ast.File, fn func(name string, pos token.Pos)) {
	for _, f := range files {
		for _, d := range f.Decls {
			if fd, ok := d.(*ast.FuncDecl); ok {
				fn(fd.Name.Name, fd.Pos())
			}
		}
		ast.Inspect(f, func(n ast.Node) bool {
			if it, ok := n.(*ast.InterfaceType); ok {
				for _, m := range it.Methods.List {
					for _, name := range m.Names {
						fn(name.Name, name.Pos())
					}
				}
			}
			return true
		})
	}
}

func eachType(files []*ast.File, fn func(name string, pos token.Pos)) {
	for _, f := range files {
		ast.Inspect(f, func(n ast.Node) bool {
			if ts, ok := n.(*ast.TypeSpec); ok {
				fn(ts.Name.Name, ts.Pos())
			}
			return true
		})
	}
}

// hasWordPrefix reports whether name is prefix followed by an upper-case
// letter, so "GetPort" and "IConfig" match but "Getter" and "Index" don't.
func hasWordPrefix(name, prefix string) bool {
	rest, ok := strings.CutPrefix(name, prefix)
	if !ok || rest == "" {
		return false
	}
	return unicode.IsUpper([]rune(rest)[0])
}
