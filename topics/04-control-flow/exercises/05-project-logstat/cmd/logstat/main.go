// Command logstat prints the log lines from standard input at or above a
// level, then a summary of what it read on standard error.
//
// Usage:
//
//	logstat [-level info] [-n 0] < app.log
//
//	-level  the minimum level to print: debug, info, warn, error or fatal
//	-n      stop after this many matching lines (0: no limit)
//
// Exercise E5 — mini-project. main is complete: it only wires things
// together. Implement internal/logstat until all tests pass:
//
//	make ex E=05-project-logstat
//	go run ./exercises/05-project-logstat/cmd/logstat -level warn < exercises/05-project-logstat/testdata/app.log
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/Shreyash019/Golang/topics/04-control-flow/exercises/05-project-logstat/internal/logstat"
)

func main() {
	if err := run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "logstat:", err)
		os.Exit(2)
	}
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("logstat", flag.ContinueOnError)
	fs.SetOutput(stderr)
	levelName := fs.String("level", "info", "minimum `level` to print: debug, info, warn, error or fatal")
	limit := fs.Int("n", 0, "stop after `N` matching lines (0: no limit)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() > 0 {
		return errors.New("logstat reads standard input and takes no arguments")
	}
	if *limit < 0 {
		return fmt.Errorf("-n must be 0 or more, got %d", *limit)
	}
	minLevel, ok := logstat.ParseLevel(*levelName)
	if !ok {
		return fmt.Errorf("unknown level %q", *levelName)
	}

	s, err := logstat.Filter(stdin, stdout, minLevel, *limit)
	if err != nil {
		return err
	}
	fmt.Fprintf(stderr, "lines=%d matched=%d skipped=%d malformed=%d debug=%d info=%d warn=%d error=%d fatal=%d\n",
		s.Lines, s.Matched, s.Skipped, s.Malformed, s.Debug, s.Info, s.Warn, s.Error, s.Fatal)
	return nil
}
