// Command buildtags shows build constraints: which of debug_on.go and
// debug_off.go is compiled depends on the -tags flag.
//
//	go run ./cmd/buildtags              # debug=false
//	go run -tags debug ./cmd/buildtags  # debug=true
package main

import (
	"fmt"
	"io"
	"os"
)

func main() { run(os.Stdout) }

func run(w io.Writer) { fmt.Fprintf(w, "debug=%t mode=%s\n", debug, mode()) }
