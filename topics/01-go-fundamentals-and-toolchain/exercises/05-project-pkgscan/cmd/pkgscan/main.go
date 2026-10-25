// Command pkgscan lists the Go packages under a directory: their names,
// whether they are commands (package main), how many files and test files
// they have, how many packages they import, and whether they have a package
// doc comment.
//
// Usage:
//
//	pkgscan [-version] [dir]
//
// Exercise E5 — mini-project. main is complete: it only wires things together.
// Implement internal/scan and internal/report until all tests pass:
//
//	make ex E=05-project-pkgscan
//	go run ./exercises/05-project-pkgscan/cmd/pkgscan .
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/Shreyash019/Golang/topics/01-go-fundamentals-and-toolchain/exercises/05-project-pkgscan/internal/report"
	"github.com/Shreyash019/Golang/topics/01-go-fundamentals-and-toolchain/exercises/05-project-pkgscan/internal/scan"
)

// version is stamped at build time: go build -ldflags "-X main.version=v1.0.0"
var version = "dev"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

// run is main without the process-level side effects, so tests can call it.
// It returns the exit code: 0 success, 1 runtime error, 2 usage error.
func run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("pkgscan", flag.ContinueOnError)
	flags.SetOutput(stderr)
	flags.Usage = func() {
		fmt.Fprintln(stderr, "usage: pkgscan [-version] [dir]")
		flags.PrintDefaults()
	}
	showVersion := flags.Bool("version", false, "print the version and exit")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if *showVersion {
		fmt.Fprintln(stdout, "pkgscan", version)
		return 0
	}
	if flags.NArg() > 1 {
		flags.Usage()
		return 2
	}
	root := "."
	if flags.NArg() == 1 {
		root = flags.Arg(0)
	}

	pkgs, err := scan.Scan(os.DirFS(root))
	if err != nil {
		fmt.Fprintln(stderr, "pkgscan:", err)
		return 1
	}
	if err := report.Write(stdout, pkgs); err != nil {
		fmt.Fprintln(stderr, "pkgscan:", err)
		return 1
	}
	return 0
}
