// Package report renders scan results as an aligned text table.
package report

import (
	"io"

	"github.com/Shreyash019/Golang/topics/01-go-fundamentals-and-toolchain/exercises/05-project-pkgscan/internal/scan"
)

// Write prints one row per package, then a summary line:
//
//	DIR             PACKAGE  KIND     FILES  TESTS  IMPORTS  DOC
//	cmd/api         main     command  1      0      2        no
//	internal/order  order    library  2      1      2        yes
//
//	2 packages (1 command, 1 library)
//
// Columns are aligned with text/tabwriter (minwidth 0, tabwidth 0,
// padding 2, padchar ' ', no flags). Return any write error.
func Write(w io.Writer, pkgs []scan.Package) error {
	return nil
}
