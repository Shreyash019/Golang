// Command envcheck prints the configuration a service would run with, from
// its APP_* environment variables and their defaults, so a bad deployment
// setting fails here instead of at startup in production.
//
// Usage:
//
//	APP_PORT=9090 APP_LOG_LEVEL=debug envcheck
//
// Exercise E5 — mini-project. main is complete: it only wires things
// together. Implement internal/config until all tests pass:
//
//	make ex E=05-project-envcheck
//	APP_PORT=9090 go run ./exercises/05-project-envcheck/cmd/envcheck
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/Shreyash019/Golang/topics/02-variables-and-constants/exercises/05-project-envcheck/internal/config"
)

func main() {
	if err := run(os.Stdout, os.LookupEnv); err != nil {
		fmt.Fprintln(os.Stderr, "envcheck:", err)
		os.Exit(1)
	}
}

func run(w io.Writer, lookup config.LookupFunc) error {
	cfg, err := config.Load(lookup)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "port=%d timeout=%s log_level=%s max_body_bytes=%d\n",
		cfg.Port, cfg.Timeout, cfg.Level, cfg.MaxBodyBytes)
	return nil
}
