// Command calc evaluates an integer arithmetic expression.
//
// Usage:
//
//	calc [-trace] EXPRESSION...
//	calc '2 * (3 + 4)'          # 14
//	calc -trace 1 + 2           # the parse, rule by rule, then 3
//	calc -- '-(2 + 3) * 2'      # -- ends the flags, so a leading - is not a flag
//
// The arguments are joined with spaces, so quoting is optional unless the
// shell would interpret * or the parentheses.
//
// Exercise E5 — mini-project. main is complete: it only wires things
// together. Implement internal/expr until all tests pass:
//
//	make ex E=05-project-calc
//	go run ./exercises/05-project-calc/cmd/calc -trace '2 * (3 + 4)'
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Shreyash019/Golang/topics/05-functions-closures-and-defer/exercises/05-project-calc/internal/expr"
)

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "calc:", err)
		os.Exit(1)
	}
}

func run(args []string, w io.Writer) error {
	fs := flag.NewFlagSet("calc", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	trace := fs.Bool("trace", false, "print each grammar rule as it's parsed")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() == 0 {
		return errors.New("usage: calc [-trace] EXPRESSION")
	}
	src := strings.Join(fs.Args(), " ")

	var (
		v   int64
		err error
	)
	if *trace {
		v, err = expr.Trace(src, w)
	} else {
		v, err = expr.Eval(src)
	}
	if err != nil {
		return err
	}
	fmt.Fprintln(w, v)
	return nil
}
