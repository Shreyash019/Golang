package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
	"testing"
	"unicode"
)

// TestIdioms fails while the Java/Python habits are still in the code. Each
// subtest names one habit and the idiomatic alternative (notes §8).
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

	t.Run("MixedCaps names", func(t *testing.T) {
		eachIdent(files, func(id *ast.Ident) {
			if isScreamingSnake(id.Name) {
				t.Errorf("%s: %s: Go names are MixedCaps, constants included (maxRetries, not MAX_RETRIES)", fset.Position(id.Pos()), id.Name)
			}
		})
	})

	t.Run("no package-level variables", func(t *testing.T) {
		for _, f := range files {
			for _, d := range f.Decls {
				gd, ok := d.(*ast.GenDecl)
				if !ok || gd.Tok != token.VAR {
					continue
				}
				for _, spec := range gd.Specs {
					for _, name := range spec.(*ast.ValueSpec).Names {
						t.Errorf("%s: package-level var %s: settings that never change are constants; running totals are local variables", fset.Position(name.Pos()), name.Name)
					}
				}
			}
		}
	})

	t.Run("no explicit zero-value initializers", func(t *testing.T) {
		for _, f := range files {
			ast.Inspect(f, func(n ast.Node) bool {
				vs, ok := n.(*ast.ValueSpec)
				if !ok || vs.Type == nil {
					return true
				}
				for i, v := range vs.Values {
					if isZeroLiteral(v) {
						t.Errorf("%s: var %s %s = ...: a var declaration already zeroes the variable, write `var %s %s` (or use := with a real value)",
							fset.Position(vs.Pos()), vs.Names[i].Name, exprString(vs.Type), vs.Names[i].Name, exprString(vs.Type))
					}
				}
				return true
			})
		}
	})

	t.Run("no getters or constants types", func(t *testing.T) {
		eachDecl(files, func(name string, pos token.Pos) {
			if rest, ok := strings.CutPrefix(name, "Get"); ok && rest != "" && unicode.IsUpper([]rune(rest)[0]) {
				t.Errorf("%s: %s: no Get prefix, and constants need no getter at all: use the constant", fset.Position(pos), name)
			}
			if strings.HasSuffix(name, "Constants") {
				t.Errorf("%s: type %s: a type that only groups constants is a Java habit; a const block does that", fset.Position(pos), name)
			}
		})
	})

	t.Run("durations as time.Duration", func(t *testing.T) {
		eachIdent(files, func(id *ast.Ident) {
			// Case-sensitive, so delayMs and DELAY_MS match but items doesn't.
			if strings.HasSuffix(id.Name, "Ms") || strings.HasSuffix(id.Name, "MS") || strings.HasSuffix(id.Name, "Millis") {
				t.Errorf("%s: %s: keep durations as time.Duration (100 * time.Millisecond), not ints with a unit in the name", fset.Position(id.Pos()), id.Name)
			}
		})
	})
}

// eachIdent calls fn for every identifier declared in files: constants,
// variables, types, functions, parameters and results.
func eachIdent(files []*ast.File, fn func(*ast.Ident)) {
	for _, f := range files {
		ast.Inspect(f, func(n ast.Node) bool {
			switch n := n.(type) {
			case *ast.ValueSpec:
				for _, id := range n.Names {
					fn(id)
				}
			case *ast.TypeSpec:
				fn(n.Name)
			case *ast.FuncDecl:
				fn(n.Name)
			case *ast.Field:
				for _, id := range n.Names {
					fn(id)
				}
			case *ast.AssignStmt:
				if n.Tok == token.DEFINE {
					for _, lhs := range n.Lhs {
						if id, ok := lhs.(*ast.Ident); ok {
							fn(id)
						}
					}
				}
			}
			return true
		})
	}
}

// eachDecl calls fn for every function, method and type name in files.
func eachDecl(files []*ast.File, fn func(name string, pos token.Pos)) {
	for _, f := range files {
		ast.Inspect(f, func(n ast.Node) bool {
			switch n := n.(type) {
			case *ast.FuncDecl:
				fn(n.Name.Name, n.Name.Pos())
			case *ast.TypeSpec:
				fn(n.Name.Name, n.Name.Pos())
			}
			return true
		})
	}
}

// isScreamingSnake reports whether name looks like MAX_RETRIES: upper case
// with at least one underscore.
func isScreamingSnake(name string) bool {
	return strings.Contains(name, "_") && name != "_" && strings.ToUpper(name) == name
}

// isZeroLiteral reports whether e is 0, "", false, nil or an empty
// composite literal such as T{}.
func isZeroLiteral(e ast.Expr) bool {
	switch e := e.(type) {
	case *ast.BasicLit:
		return e.Value == "0" || e.Value == `""` || e.Value == "0.0"
	case *ast.Ident:
		return e.Name == "false" || e.Name == "nil"
	case *ast.CompositeLit:
		return len(e.Elts) == 0
	}
	return false
}

func exprString(e ast.Expr) string {
	switch e := e.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.SelectorExpr:
		return exprString(e.X) + "." + e.Sel.Name
	}
	return "T"
}
