//go:build debug

package main

const debug = true

func mode() string { return "verbose logging, pprof enabled" }
