// Package report renders scan results as an aligned text table.
package report

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/Shreyash019/Golang/topics/01-go-fundamentals-and-toolchain/exercises/solutions/05-project-pkgscan/internal/scan"
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
// padding 2, padchar ' ', no flags).
func Write(w io.Writer, pkgs []scan.Package) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "DIR\tPACKAGE\tKIND\tFILES\tTESTS\tIMPORTS\tDOC")
	commands := 0
	for _, p := range pkgs {
		kind := "library"
		if p.IsCommand() {
			kind = "command"
			commands++
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%d\t%d\t%d\t%s\n", p.Dir, p.Name, kind, p.Files, p.TestFiles, len(p.Imports), yesNo(p.HasDoc))
	}
	if err := tw.Flush(); err != nil {
		return err
	}
	_, err := fmt.Fprintf(w, "\n%s (%s, %s)\n",
		plural(len(pkgs), "package", "packages"),
		plural(commands, "command", "commands"),
		plural(len(pkgs)-commands, "library", "libraries"))
	return err
}

func yesNo(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

func plural(n int, one, many string) string {
	if n == 1 {
		return fmt.Sprintf("%d %s", n, one)
	}
	return fmt.Sprintf("%d %s", n, many)
}
