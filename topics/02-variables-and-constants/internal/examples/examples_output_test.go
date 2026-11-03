package examples

import "os"

func ExampleZeroValues() {
	ZeroValues(os.Stdout)
	// Output:
	// int=0 float64=0 bool=false string="" rune=0
	// pointer=<nil> slice=[] (nil? true) map=map[] (nil? true)
	// chan nil? true func nil? true error=<nil>
	// array=[0 0 0] struct={Name: Age:0}
	// len(slice)=0 len(map)=0 len(string)=0
	// read from nil map: 0
}

func ExampleDeclarationForms() {
	DeclarationForms(os.Stdout)
	// Output:
	// a=0 b=42 c=x d=3.5 e=true
	// int int string float64 bool
	// localhost:8080
}

func ExampleDefaultTypes() {
	DefaultTypes(os.Stdout)
	// Output:
	// int float64 int32 complex128 string bool
	// mixed=3 float64
	// any(1<<10) is int
}

func ExampleUntypedConstants() {
	UntypedConstants(os.Stdout)
	// Output:
	// small: 4
	// huge as float64: 1.2676506002282294e+30
	// duration: 3ns byte: 65
}

func ExampleConstantFloatExactness() {
	ConstantFloatExactness(os.Stdout)
	// Output:
	// constants: true
	// variables: false
	// x+y = 0.30000000000000004
}

func ExampleTypedVsUntyped() {
	TypedVsUntyped(os.Stdout)
	// Output:
	// timeout: 2s (time.Duration)
	// retries*timeout: 6s
	// time.Duration(n)*timeout: 6s
}

func ExampleIotaBasics() {
	IotaBasics(os.Stdout)
	// Output:
	// Sunday=0 Monday=1 Tuesday=2 (examples.Weekday)
	// a b c: 0 10 30
	// x y z v: 0 10 1 11
	// restart: 0
}

func ExampleScopeLevels() {
	ScopeLevels(os.Stdout)
	// Output:
	// 1: package
	// 2: function
	// 3: explicit block
	// 4: if statement
	// 5: loop body i = 0
	// 6: function
}

func ExampleScopeStartsAfterDeclaration() {
	ScopeStartsAfterDeclaration(os.Stdout)
	// Output:
	// inner: 11
	// outer: 1
}

func ExampleShortRedeclaration() {
	ShortRedeclaration(os.Stdout)
	// Output:
	// a b err: 1 0 strconv.Atoi: parsing "x": invalid syntax
	// same err variable: true
	// inner c err: 2 <nil> same err variable: false
	// outer err after the block: strconv.Atoi: parsing "x": invalid syntax
}

func ExampleUsefulZeroValues() {
	UsefulZeroValues(os.Stdout)
	// Output:
	// built without a constructor
	// 16 bytes buffered
	// locked a zero-value mutex
	// [1 2 3] 3
}

func ExampleMultipleAssignment() {
	MultipleAssignment(os.Stdout)
	// Output:
	// swapped: 2 1
	// i: 2 nums: [99 20 30]
	// fraction: 0.75
}

func ExampleIotaBitFlags() {
	IotaBitFlags(os.Stdout)
	// Output:
	// 011 write? true exec? false
	// 101 write? false exec? true
}

func ExampleIotaSizes() {
	IotaSizes(os.Stdout)
	// Output:
	// 1024 1048576 1073741824 1099511627776
	// 3145728 bytes = 3.0 MiB = 0.0029 GiB
}

func ExampleShadowedErr() {
	ShadowedErr(os.Stdout)
	// Output:
	// buggy: port=8080 err=<nil>
	// fixed: port=9090 err=<nil>
	// buggy, bad input: port=0 err=strconv.Atoi: parsing "abc": invalid syntax
}

func ExampleShadowedPackageVar() {
	ShadowedPackageVar(os.Stdout)
	// Output:
	// inside setupEnv: "production"
	// package currentEnv: ""
}

func ExampleConstantIntDivision() {
	ConstantIntDivision(os.Stdout)
	// Output:
	// 0 0.5 0
}

func ExampleIotaZeroValueTrap() {
	IotaZeroValueTrap(os.Stdout)
	// Output:
	// unset status == Active: true
	// unset safe status == StatusUnknown: true
}

func ExampleIotaReorder() {
	IotaReorder(os.Stdout)
	// Output:
	// stored value: 1
	// v2 reads it as Paid? false as Cancelled? true
}
