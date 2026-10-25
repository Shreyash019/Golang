//go:build !debug

package main

const debug = false

func mode() string { return "production" }
