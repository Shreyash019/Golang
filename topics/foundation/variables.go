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

	// 4. Variable Scope
	/*
		Scope determines where an identifier can be referenced in source code.
		Go has several important scopes:
		- Universe scope: predeclared identifiers such as int, string, and true.
		- Package-level scope: declarations outside functions, available within their package according to the language's scope rules.
		- Function-local scope: variables declared inside a function.
		- Block scope: variables introduced inside {} blocks, including conditional and loop bodies.
		- File-specific import scope: imported package names are available in the file that imports them, not automatically throughout every file in the package.
	*/
	port := 4001
	var appName string = "VARIABLE CHECK"
	if port > 0 {
		message := "Server Running"
		fmt.Println(message, appName, port)
	}
	// fmt.Println(message, appName, port) // Compile error: message is out of scope
	/*
		Scope is not the same as lifetime
		This distinction is especially important when you eventually study pointers, closures, goroutines, and garbage collection.
		- Scope answers: where can the identifier be used in the source code?
		- Lifetime answers: how long does the underlying value need to remain available during execution?
		A local variable can escape its function through a returned pointer or a closure. The compiler decides how storage is implemented, including whether it needs heap allocation. Scope alone does not determine stack-versus-heap allocation.
	*/

	// 5. Variable Shadowing
	count := 10
	if count > 5 {
		count := 30
		count += 30
		fmt.Println("Inner Count: ", count)
	}
	fmt.Println("Outer Count: ", count)
}
