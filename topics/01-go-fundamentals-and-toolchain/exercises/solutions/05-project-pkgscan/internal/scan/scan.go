// Package scan finds the Go packages in a file tree, applying the same
// directory rules as the go command.
package scan

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"path"
	"slices"
	"strings"
)

// Package describes one Go package (one directory) found by Scan.
type Package struct {
	Dir       string   // slash-separated path relative to the scan root; "." for the root
	Name      string   // from the package clause of the non-test files
	Files     int      // non-test .go files
	TestFiles int      // _test.go files
	Imports   []string // unique import paths of the non-test files, sorted
	HasDoc    bool     // some non-test file has a package doc comment
}

// IsCommand reports whether p builds an executable (package main).
func (p Package) IsCommand() bool { return p.Name == "main" }

// Scan walks fsys and returns its Go packages sorted by Dir.
//
// Like the go command, it skips directories named testdata or vendor and
// any file or directory whose name starts with "." or "_". A directory is a
// package if it has at least one non-test .go file. Build constraints are
// not evaluated: every .go file counts.
func Scan(fsys fs.FS) ([]Package, error) {
	byDir := make(map[string]*Package)
	fset := token.NewFileSet()

	err := fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		name := d.Name()
		if d.IsDir() {
			if p != "." && ignored(name) {
				return fs.SkipDir
			}
			return nil
		}
		if ignored(name) || !strings.HasSuffix(name, ".go") {
			return nil
		}

		dir := path.Dir(p)
		pkg := byDir[dir]
		if pkg == nil {
			pkg = &Package{Dir: dir}
			byDir[dir] = pkg
		}
		if strings.HasSuffix(name, "_test.go") {
			pkg.TestFiles++
			return nil
		}
		return addFile(fsys, fset, p, pkg)
	})
	if err != nil {
		return nil, err
	}

	var pkgs []Package
	for _, pkg := range byDir {
		if pkg.Files == 0 { // only test files: not a package the go command builds
			continue
		}
		slices.Sort(pkg.Imports)
		pkg.Imports = slices.Compact(pkg.Imports)
		pkgs = append(pkgs, *pkg)
	}
	slices.SortFunc(pkgs, func(a, b Package) int { return strings.Compare(a.Dir, b.Dir) })
	return pkgs, nil
}

// addFile parses the package clause, doc comment and imports of one
// non-test file and records them on pkg.
func addFile(fsys fs.FS, fset *token.FileSet, p string, pkg *Package) error {
	src, err := fs.ReadFile(fsys, p)
	if err != nil {
		return err
	}
	f, err := parser.ParseFile(fset, p, src, parser.ImportsOnly|parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parse %s: %w", p, err)
	}
	if pkg.Name != "" && pkg.Name != f.Name.Name {
		return fmt.Errorf("%s: multiple packages %s and %s", pkg.Dir, pkg.Name, f.Name.Name)
	}
	pkg.Name = f.Name.Name
	pkg.Files++
	if f.Doc != nil {
		pkg.HasDoc = true
	}
	for _, imp := range f.Imports {
		pkg.Imports = append(pkg.Imports, strings.Trim(imp.Path.Value, `"`))
	}
	return nil
}

// ignored reports whether the go command ignores a file or directory name.
func ignored(name string) bool {
	return strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_") || name == "testdata" || name == "vendor"
}
