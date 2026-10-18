package foundation

import (
	"fmt"
)

func dataTypes() {
	fmt.Println("======= Data Types...")
	// 1. Boolean
	var boolean bool
	fmt.Println("Boolean", boolean)

	// 2. Integer
	var intt8 int8
	var intt16 int16
	var intt32 int32
	var intt64 int64
	intt16 = -900
	fmt.Println("Integer (INT): ", intt8, intt16, intt32, intt64)
	var uintt8 uint8
	var uintt16 uint16
	var uintt32 uint32
	var uintt64 uint64
	uintt16 = 9000
	fmt.Println("Integer (UINT): ", uintt8, uintt16, uintt32, uintt64)
	var intt int
	var uintt uint
	var uinttptr uintptr
	fmt.Println("Integer without base", intt, uintt, uinttptr)

	// 3. Float
	var temperature float32 = 36.5
	var price float64 = 199.99
	fmt.Println(temperature)
	fmt.Println(price)

	// 4. Complex
	var a complex64 = 1 + 2i
	var b complex128 = 3 + 4i
	fmt.Println(a)
	fmt.Println(b)
	z := complex(3.0, 4.0)
	fmt.Println(z)
	fmt.Println(real(z)) // 3
	fmt.Println(imag(z)) // 4

	// 5. String
	name := "Go"
	message := "Hello, World!"
	fmt.Println(name, message, len(message))

	// 6. byte, rune
	s := "é"
	fmt.Println(s[0], s[1])
	fmt.Println(len(s))         // 2 bytes
	fmt.Println(len([]rune(s))) // 1 Unicode code point
}
