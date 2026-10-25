// Command version prints a version string stamped in at build time:
//
//	go run -ldflags "-X main.version=1.4.2" ./cmd/version
package main

import (
	"fmt"
	"io"
	"os"
)

// version is overwritten by the linker flag -X main.version=...
// It must be a package-level string variable (not a constant).
var version = "dev"

func main() { run(os.Stdout) }

func run(w io.Writer) { fmt.Fprintln(w, "version", version) }
