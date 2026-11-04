// Command examples runs the worked examples from this topic's notes.
//
// Usage:
//
//	go run ./cmd/examples          # list the examples
//	go run ./cmd/examples <name>   # run one example
//	go run ./cmd/examples all      # run every example that doesn't crash on purpose
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/Shreyash019/Golang/topics/05-functions-closures-and-defer/internal/examples"
)

func main() {
	if err := run(os.Args[1:], os.Stdout, examples.All()); err != nil {
		fmt.Fprintln(os.Stderr, "examples:", err)
		os.Exit(1)
	}
}

func run(args []string, w io.Writer, all []examples.Example) error {
	if len(args) == 0 {
		if len(all) == 0 {
			fmt.Fprintln(w, "no examples yet")
			return nil
		}
		for _, ex := range all {
			note := ""
			if ex.Crashes {
				note = " (crashes on purpose)"
			}
			fmt.Fprintf(w, "%-28s §%-5s %s%s\n", ex.Name, ex.Section, ex.Summary, note)
		}
		return nil
	}

	name := args[0]
	if name == "all" {
		for _, ex := range all {
			if ex.Crashes {
				continue
			}
			fmt.Fprintf(w, "=== %s (§%s)\n", ex.Name, ex.Section)
			ex.Run(w)
			fmt.Fprintln(w)
		}
		return nil
	}

	for _, ex := range all {
		if ex.Name == name {
			ex.Run(w)
			return nil
		}
	}
	return fmt.Errorf("unknown example %q (run with no arguments to list them)", name)
}
