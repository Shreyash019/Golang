package shipping

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
	"testing"
)

// TestIdioms fails while the habits are still in the code. Each subtest
// names one habit and the idiomatic alternative (notes §8).
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
	inspect := func(fn func(n ast.Node) bool) {
		for _, f := range files {
			ast.Inspect(f, fn)
		}
	}

	t.Run("no else after return", func(t *testing.T) {
		inspect(func(n ast.Node) bool {
			is, ok := n.(*ast.IfStmt)
			if !ok {
				return true
			}
			if _, isBlock := is.Else.(*ast.BlockStmt); isBlock && endsInReturn(is.Body) {
				t.Errorf("%s: the if block ends in return, so drop the else and outdent its body (Code Review Comments: Indent Error Flow)", fset.Position(is.Else.Pos()))
			}
			return true
		})
	})

	t.Run("ifs nested at most two deep", func(t *testing.T) {
		var walk func(n ast.Node, depth int)
		walk = func(n ast.Node, depth int) {
			ast.Inspect(n, func(c ast.Node) bool {
				if c == n {
					return true
				}
				switch c := c.(type) {
				case *ast.IfStmt:
					if depth+1 > 2 {
						t.Errorf("%s: if nested %d deep: handle the failure cases first with guard clauses and return early", fset.Position(c.Pos()), depth+1)
					}
					walk(c.Body, depth+1)
					if c.Else != nil {
						walk(c.Else, depth+1)
					}
					return false
				case *ast.FuncLit:
					walk(c.Body, 0)
					return false
				}
				return true
			})
		}
		for _, f := range files {
			for _, d := range f.Decls {
				if fd, ok := d.(*ast.FuncDecl); ok && fd.Body != nil {
					walk(fd.Body, 0)
				}
			}
		}
	})

	t.Run("switch instead of an if-else-if ladder", func(t *testing.T) {
		counted := map[*ast.IfStmt]bool{} // the else-if links of a ladder already reported
		inspect(func(n ast.Node) bool {
			is, ok := n.(*ast.IfStmt)
			if !ok || counted[is] {
				return true
			}
			branches := 1
			for e := is.Else; e != nil; {
				next, ok := e.(*ast.IfStmt)
				if !ok {
					break
				}
				counted[next] = true
				branches++
				e = next.Else
			}
			if branches >= 3 {
				t.Errorf("%s: an if/else-if ladder with %d conditions: use a switch (Effective Go: Switch)", fset.Position(is.Pos()), branches)
			}
			return true
		})
	})

	t.Run("range instead of C-style index loops", func(t *testing.T) {
		inspect(func(n ast.Node) bool {
			fs, ok := n.(*ast.ForStmt)
			if !ok || fs.Init == nil || fs.Cond == nil || fs.Post == nil {
				return true
			}
			cond, ok := fs.Cond.(*ast.BinaryExpr)
			if !ok || cond.Op != token.LSS {
				return true
			}
			if call, ok := cond.Y.(*ast.CallExpr); ok {
				if id, ok := call.Fun.(*ast.Ident); ok && id.Name == "len" {
					t.Errorf("%s: for i := 0; i < len(x); i++: use for i, v := range x (or for i := range x)", fset.Position(fs.Pos()))
				}
			}
			return true
		})
	})

	t.Run("no for true", func(t *testing.T) {
		inspect(func(n ast.Node) bool {
			if fs, ok := n.(*ast.ForStmt); ok {
				if id, ok := fs.Cond.(*ast.Ident); ok && id.Name == "true" {
					t.Errorf("%s: for true: write for { ... }, or better, a loop whose condition says when it ends", fset.Position(fs.Pos()))
				}
			}
			return true
		})
	})

	t.Run("no comparisons with true or false", func(t *testing.T) {
		inspect(func(n ast.Node) bool {
			be, ok := n.(*ast.BinaryExpr)
			if !ok || (be.Op != token.EQL && be.Op != token.NEQ) {
				return true
			}
			for _, side := range []ast.Expr{be.X, be.Y} {
				if id, ok := side.(*ast.Ident); ok && (id.Name == "true" || id.Name == "false") {
					t.Errorf("%s: comparing with %s: a bool is already a condition, write `if ok` or `if !ok`", fset.Position(be.Pos()), id.Name)
				}
			}
			return true
		})
	})
}

func endsInReturn(b *ast.BlockStmt) bool {
	if len(b.List) == 0 {
		return false
	}
	_, ok := b.List[len(b.List)-1].(*ast.ReturnStmt)
	return ok
}
