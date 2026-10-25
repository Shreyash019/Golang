// Exercise E4 — refactor to idiomatic Go.
//
// This program works: it parses key=value arguments into a server config and
// prints a startup line. But it is written like Java, with an interface that
// has one implementation, an "Impl" type, getters, a factory, a utils file
// and panics used as exceptions. Rewrite it as idiomatic Go (notes §8):
// a plain struct, a function that returns (Config, error), errors instead of
// panics, and no utils. It trains the AGENTS.md §1.1 habits.
//
// behaviour_test.go must keep passing. idiom_test.go fails until the Java
// habits are gone.
//
//	go run ./exercises/04-refactor-config port=9090 env=prod debug=true
//	make ex E=04-refactor-config
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

// run is the program. Like a Java "try/catch", it turns the panics thrown
// by the config manager back into an error.
func run(args []string, w io.Writer) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if ce, ok := r.(*ConfigException); ok {
				err = fmt.Errorf("%s", ce.Message)
				return
			}
			panic(r)
		}
	}()

	var manager IConfigManager = NewConfigManager()
	manager.LoadFromArgs(args)
	fmt.Fprintf(w, "starting server on :%d (env=%s, debug=%t)\n", manager.GetPort(), manager.GetEnv(), manager.GetDebug())
	return nil
}
