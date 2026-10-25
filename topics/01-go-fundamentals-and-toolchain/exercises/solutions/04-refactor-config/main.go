// Exercise E4 — reference solution.
//
// The Java version needed an interface, an Impl, a factory, six accessors,
// an exception type, a recover block and a utils type. In Go it's a struct
// with exported fields and one function that returns (Config, error).
package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "config:", err)
		os.Exit(2)
	}
}

func run(args []string, w io.Writer) error {
	cfg, err := parseArgs(args)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "starting server on :%d (env=%s, debug=%t)\n", cfg.Port, cfg.Env, cfg.Debug)
	return nil
}
