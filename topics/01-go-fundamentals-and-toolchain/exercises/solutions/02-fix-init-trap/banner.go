// Package banner builds the startup line a service logs when it boots.
//
// Exercise E2 — fix the bug. It trains the initialization-order gotcha
// (notes §7.6): package-level variables are initialized before any init
// function runs, so a variable computed from values that init sets later
// sees their zero values.
package banner

import "fmt"

// The defaults are plain initializers, so they are set during variable
// initialization, before banner is computed from them. No init function is
// needed: Effective Go reserves init for work that can't be expressed as a
// declaration.
var (
	defaultName = "orders-api"
	defaultPort = 8080
)

var banner = fmt.Sprintf("%s listening on :%d", defaultName, defaultPort)

// Banner returns the startup line, e.g. "orders-api listening on :8080".
func Banner() string { return banner }
