// Package scan finds the Go packages in a file tree, applying the same
// directory rules as the go command.
package scan

import "io/fs"

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
func (p Package) IsCommand() bool {
	return false
}

// Scan walks fsys and returns its Go packages sorted by Dir.
//
// Like the go command, it skips directories named testdata or vendor and
// any file or directory whose name starts with "." or "_". A directory is a
// package if it has at least one non-test .go file. Build constraints are
// not evaluated: every .go file counts.
//
// Errors: a file that doesn't parse returns an error starting with
// "parse <path>: ", and two different package names in one directory
// return "<dir>: multiple packages <a> and <b>".
//
// Tools from the standard library: fs.WalkDir (and fs.SkipDir), fs.ReadFile,
// go/parser.ParseFile with parser.ImportsOnly|parser.ParseComments (the
// *ast.File has Name, Doc and Imports), path.Dir, slices.Sort/Compact.
func Scan(fsys fs.FS) ([]Package, error) {
	return nil, nil
}
