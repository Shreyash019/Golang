// Package banner builds the startup line a service logs when it boots.
//
// Exercise E2 — fix the bug. Banner should return
// "orders-api listening on :8080", but the test shows it doesn't.
// Find the root cause, fix it with the smallest idiomatic change, and
// explain the gotcha in one sentence (notes §7).
//
//	make ex E=02-fix-init-trap
package banner

import "fmt"

var (
	defaultName string
	defaultPort int
)

var banner = fmt.Sprintf("%s listening on :%d", defaultName, defaultPort)

func init() {
	defaultName = "orders-api"
	defaultPort = 8080
}

// Banner returns the startup line, e.g. "orders-api listening on :8080".
func Banner() string { return banner }
