package temperature_test

import (
	"fmt"

	"github.com/Shreyash019/Golang/topics/01-go-fundamentals-and-toolchain/internal/temperature"
)

// A testable example: `go test` checks the output, and `go doc` and
// pkg.go.dev show it as usage documentation for Celsius.Fahrenheit.
func ExampleCelsius_Fahrenheit() {
	fmt.Println(temperature.Celsius(100).Fahrenheit())
	fmt.Println(temperature.AbsoluteZero.Fahrenheit())
	// Output:
	// 212
	// -459.66999999999996
}
