// Package temperature converts between temperature scales.
//
// It exists to show what good doc comments look like and how `go doc`
// renders them.
//
// # Scales
//
// Two scales are supported:
//   - [Celsius], the SI scale used in most of the world
//   - [Fahrenheit], used mainly in the United States
//
// Convert a value with its method, for example:
//
//	f := temperature.Celsius(100).Fahrenheit() // 212
//
// Values are plain float64s underneath, so they print with [fmt.Println].
package temperature

// AbsoluteZero is the lowest possible temperature, in degrees Celsius.
const AbsoluteZero Celsius = -273.15

// Celsius is a temperature in degrees Celsius.
type Celsius float64

// Fahrenheit is a temperature in degrees Fahrenheit.
type Fahrenheit float64

// Fahrenheit converts c to degrees Fahrenheit.
func (c Celsius) Fahrenheit() Fahrenheit {
	return Fahrenheit(c*9/5 + 32)
}

// Celsius converts f to degrees Celsius.
func (f Fahrenheit) Celsius() Celsius {
	return Celsius((f - 32) * 5 / 9)
}

// CToF converts c to degrees Fahrenheit.
//
// Deprecated: Use [Celsius.Fahrenheit] instead.
func CToF(c float64) float64 {
	return float64(Celsius(c).Fahrenheit())
}
