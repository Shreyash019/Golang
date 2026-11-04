// Command exit-skips-defer shows that os.Exit ends the program immediately:
// deferred calls don't run. log.Fatal calls os.Exit too.
//
//	go run ./cmd/exit-skips-defer; echo "exit status: $?"
package main

import (
	"fmt"
	"os"
)

func main() {
	defer fmt.Println("deferred: flush and close") // never printed
	fmt.Println("work done, exiting with status 1")
	os.Exit(1)
}
