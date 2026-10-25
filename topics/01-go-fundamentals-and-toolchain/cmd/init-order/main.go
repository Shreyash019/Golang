// Command init-order prints the order in which Go initializes a program:
// imported packages first, then this package's variables (in dependency
// order), then its init functions (in source order), then main.
package main

import (
	"fmt"

	"github.com/Shreyash019/Golang/topics/01-go-fundamentals-and-toolchain/internal/greet"
)

//nolint:unused // these variables exist only so their initializers run and print
var (
	a = trace("a", b+1) // depends on b, so b is initialized first
	b = trace("b", 1)
	c = trace("c", 3)
)

//nolint:unused // called only from the initializers above
func trace(name string, v int) int {
	fmt.Printf("main: var %s = %d\n", name, v)
	return v
}

func init() { fmt.Println("main: init #1") }

func init() { fmt.Println("main: init #2") }

func main() {
	fmt.Println("main: main() starts, greet.Message =", greet.Message)
}
