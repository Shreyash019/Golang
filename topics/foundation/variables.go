package foundation

import "fmt"

// const variable
const version float32 = 1.0

func variables() {
	// 1. Variable declaration
	var firstName string
	var lastName string
	var isAdult bool
	var adult string
	var verification float32

	// Zero Value declaration
	var boolen bool
	var integer int
	var float float32
	var strings string
	var pointer *int
	var slices []int
	var maps map[int]string
	var structs struct{}
	var arryas [5]int
	fmt.Printf("ZERO values: Boolean: %t, Integer: %d, Float: %f, String: %s, Pointer: %v, Slices: %d, Maps: %v\n", boolen, integer, float, strings, pointer, slices, maps)
	fmt.Printf("Structs: %v, Arrays: %v\n", structs, arryas)
	fmt.Printf("Nil Slices: %t, Nil Maps: %t\n", slices == nil, maps == nil)

	// 2. Initialization
	age := 25 // Declaration + Initialization {:= short hand declaration}
	// Short hand declaration: type getting infered based on value assigned
	// Type inference means the compiler determines a variable's type from its initializer when the type is not explicitly provided.
	// Type inference is not dynamic typing

	// 3. Value Assignment
	firstName, lastName, age, verification = "John", "Doe", 25, 90.0
	isAdult = age >= 18
	if isAdult {
		adult = "Major"
	} else {
		adult = "Minor"
	}
	fmt.Printf("Hello MR %s %s!, Your verification is %f completed as you are %s with age %d.\n", firstName, lastName, verification, adult, age)
}
