package examples

import "os"

// TypeSizes assumes a 64-bit platform (int, uint, uintptr: 8 bytes).

func ExampleTypeSizes() {
	TypeSizes(os.Stdout)
	// Output:
	// bool        1 bytes
	// int8        1 bytes
	// int16       2 bytes
	// int32       4 bytes
	// int64       8 bytes
	// int         8 bytes
	// uint        8 bytes
	// uintptr     8 bytes
	// float32     4 bytes
	// float64     8 bytes
	// complex64   8 bytes
	// complex128 16 bytes
	// byte        1 bytes
	// rune        4 bytes
	// string     16 bytes
	// int:   -9223372036854775808 to 9223372036854775807
	// int8:  -128 to 127 | uint8: 0 to 255
}

func ExampleIntegerOverflow() {
	IntegerOverflow(os.Stdout)
	// Output:
	// int8 127+1 = -128
	// uint8 0-1 = 255
	// MaxInt64+1 = -9223372036854775808
	// -(-128) as int8 = -128
	// bits: 127=01111111 -128=10000000 -1=11111111
}

func ExampleDivisionAndRemainder() {
	DivisionAndRemainder(os.Stdout)
	// Output:
	// 7/2 = 3  -7/2 = -3  -7%2 = -1  7%-2 = 1
	// -7>>1 = -4 (a shift rounds toward minus infinity)
	// -1%5 = -1  wrapped: 4
}

func ExampleFloatPrecision() {
	FloatPrecision(os.Stdout)
	// Output:
	// 0.1+0.2 = 0.30000000000000004
	// 0.1 is really 0.1000000000000000055511151
	// float32 2^24+1 == 2^24: true
	// float64 2^53+1 == 2^53: true
	// ten 0.1s: 0.9999999999999999 == 1: false
	// ten 10-cent ints == 100: true
}

func ExampleFloatSpecialValues() {
	FloatSpecialValues(os.Stdout)
	// Output:
	// 1/0 = +Inf  -1/0 = -Inf  0/0 = NaN
	// NaN == NaN: false  NaN != NaN: true  IsNaN: true
	// -0: -0  -0 == 0: true  1/-0 = -Inf
	// MaxFloat64*2 = +Inf  Inf-Inf = NaN
}

func ExampleConversions() {
	Conversions(os.Stdout)
	// Output:
	// int(-2.7) = -2  Round: -3  Floor: -3
	// int8(300) = 44  uint8(300) = 44
	// uint8(-1) = 255  uint64(-1) = 18446744073709551615
	// float64(2^53+1) == float64(2^53): true
	// int16(0x12345678) = 0x5678
}

func ExampleAliasVsDefined() {
	AliasVsDefined(os.Stdout)
	// Output:
	// float64 float64 examples.Celsius
	// c + 1.5 = 23
	// uint8 int32 65 233
}

func ExampleOperatorPrecedence() {
	OperatorPrecedence(os.Stdout)
	// Output:
	// 1+2*3 = 7
	// 1<<2+1 = 5
	// 6&3 == 2: true
	// 2^3 = 1
	// t || f && g: true
	// 10-4-3 = 3  2*3/4*4 = 4
	// 7/2*2.0 = 6
}

func ExampleBitwiseOps() {
	BitwiseOps(os.Stdout)
	// Output:
	// a&b=1000 a|b=1110 a^b=0110 a&^b=0100
	// ^a=11110011 a<<2=00110000 a>>2=0011
	// uint8 1<<7 = 128  1<<8 = 0
	// int8 -8>>1 = -4  ^-8 = 7
	// 6 is even: true  7 is even: false
}

func ExampleShortCircuit() {
	ShortCircuit(os.Stdout)
	// Output:
	// eval(a) => false && ... = false
	// eval(a) => true || ... = true
	// empty line handled without a panic
}

func ExampleAssignmentOps() {
	AssignmentOps(os.Stdout)
	// Output:
	// x = 2
	// flags = 0110
	// i = 1
	// s = gopher
}

func ExampleComparisonRules() {
	ComparisonRules(os.Stdout)
	// Output:
	// "apple" < "banana": true  "Zebra" < "apple": true  "é" > "z": true
	// int64(a) == b: true
	// complex ==: true
	// NaN < 1: false  NaN > 1: false  NaN == NaN: false
	// min(3, 1, 2) = 1  max(2.5, 1) = 2.5  min("b", "a") = a
}

func ExampleStringBasics() {
	StringBasics(os.Stdout)
	// Output:
	// len: 6 runes: 5
	// s[1]=195 (uint8) s[1:3]="é"
	//  0:'h'(int32) 1:'é'(int32) 3:'!'(int32)
	// Héllo héllo
	// gopher true true
}

func ExampleStringConversions() {
	StringConversions(os.Stdout)
	// Output:
	// string(rune(65)) = A  strconv.Itoa(65) = 65
	// Atoi("42")+1 = 43 err: <nil>
	// Atoi("4 2") err: strconv.Atoi: parsing "4 2": invalid syntax
	// 255 in base 2 and 16: 11111111 ff
	// ParseInt("300", 10, 8) err: strconv.ParseInt: parsing "300": value out of range
	// []byte("hi") = [104 105]  []rune("héllo") has 5 runes
}

func ExampleComplexNumbers() {
	ComplexNumbers(os.Stdout)
	// Output:
	// (3+4i) complex128 real=3 imag=4 abs=5
	// cmplx.Sqrt(-1) = (0+1i) math.Sqrt(-1) = NaN
	// c*c = (-7+24i)
}

func ExampleOverflowDetection() {
	OverflowDetection(os.Stdout)
	// Output:
	// MaxUint64+1: sum = 0 carry = 1
	// 2^40*2^40: hi = 65536 lo = 0
	// ones: 3 len: 8 trailing zeros: 4
}

func ExampleFloatComparison() {
	FloatComparison(os.Stdout)
	// Output:
	// a+b == 0.3: false
	// almostEqual(a+b, 0.3, 1e-9): true
	// almostEqual(1e-20, 2e-20, 1e-9): false
}

func ExampleUnsignedUnderflow() {
	UnsignedUnderflow(os.Stdout)
	// Output:
	// uint len-1 = 18446744073709551615
	// int  len-1 = -1
}
func ExampleNumberLiterals() {
	NumberLiterals(os.Stdout)
	// Output:
	// 1000000 10 15 31 0.25 1000 0.0025
	// 97 10 65 233 128512
	// float64 float64 int32
	// 0644 == 0o644 == 420: true true
}
