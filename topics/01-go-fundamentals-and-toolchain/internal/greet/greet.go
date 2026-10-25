// Package greet exists to show that an imported package is fully
// initialized before the package that imports it.
package greet

import "fmt"

// Message is set during this package's initialization.
var Message = build()

func build() string {
	fmt.Println("greet: var Message")
	return "hello"
}

func init() { fmt.Println("greet: init") }
