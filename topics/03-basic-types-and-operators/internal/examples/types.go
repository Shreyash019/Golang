package examples

import (
	"fmt"
	"io"
	"math"
	"math/cmplx"
	"unicode/utf8"
	"unsafe"
)

// TypeSizes prints the size in bytes of every basic type and the range of
// the platform-sized int. unsafe.Sizeof is a compile-time constant and is
// safe to call; the unsafe package itself is topic 30.
func TypeSizes(w io.Writer) {
	rows := []struct {
		name string
		size uintptr
	}{
		{"bool", unsafe.Sizeof(false)},
		{"int8", unsafe.Sizeof(int8(0))},
		{"int16", unsafe.Sizeof(int16(0))},
		{"int32", unsafe.Sizeof(int32(0))},
		{"int64", unsafe.Sizeof(int64(0))},
		{"int", unsafe.Sizeof(int(0))},
		{"uint", unsafe.Sizeof(uint(0))},
		{"uintptr", unsafe.Sizeof(uintptr(0))},
		{"float32", unsafe.Sizeof(float32(0))},
		{"float64", unsafe.Sizeof(float64(0))},
		{"complex64", unsafe.Sizeof(complex64(0))},
		{"complex128", unsafe.Sizeof(complex128(0))},
		{"byte", unsafe.Sizeof(byte(0))},
		{"rune", unsafe.Sizeof(rune(0))},
		{"string", unsafe.Sizeof("")},
	}
	for _, r := range rows {
		fmt.Fprintf(w, "%-10s %2d bytes\n", r.name, r.size)
	}
	fmt.Fprintln(w, "int:  ", math.MinInt, "to", math.MaxInt)
	fmt.Fprintln(w, "int8: ", math.MinInt8, "to", math.MaxInt8, "| uint8: 0 to", math.MaxUint8)
}

// Celsius is a defined type: a new type whose underlying type is float64.
type Celsius float64

// Temperature is an alias: another name for float64, not a new type.
type Temperature = float64

// AliasVsDefined shows that an alias is identical to its type, while a
// defined type is distinct and needs an explicit conversion.
func AliasVsDefined(w io.Writer) {
	f := 21.5
	var t Temperature = f //nolint:staticcheck // QF1011: the explicit type is the point: Temperature IS float64
	c := Celsius(f)       // a conversion is required: Celsius is a different type
	fmt.Fprintf(w, "%T %T %T\n", f, t, c)
	fmt.Fprintln(w, "c + 1.5 =", c+1.5) // untyped constants convert to Celsius

	var b byte = 'A' // byte is an alias for uint8
	var u uint8 = b  //nolint:staticcheck // QF1011: shows byte and uint8 are the same type
	var r rune = 'é' //nolint:staticcheck // QF1011: rune is an alias for int32
	var i int32 = r  //nolint:staticcheck // QF1011: shows rune and int32 are the same type
	fmt.Fprintf(w, "%T %T %v %v\n", b, r, u, i)
}

// StringBasics is a first look at strings: immutable byte sequences,
// usually UTF-8. Strings in depth: topic 08.
func StringBasics(w io.Writer) {
	s := "héllo"
	fmt.Fprintln(w, "len:", len(s), "runes:", utf8.RuneCountInString(s))
	fmt.Fprintf(w, "s[1]=%d (%T) s[1:3]=%q\n", s[1], s[1], s[1:3])
	for i, r := range "hé!" {
		fmt.Fprintf(w, " %d:%q(%T)", i, r, r)
	}
	fmt.Fprintln(w)

	b := []byte(s) // a mutable copy
	b[0] = 'H'
	fmt.Fprintln(w, string(b), s)
	fmt.Fprintln(w, "go"+"pher", "a" < "b", "Z" < "a")
}

// ComplexNumbers shows Go's built-in complex types.
func ComplexNumbers(w io.Writer) {
	c := complex(3, 4)
	fmt.Fprintf(w, "%v %T real=%v imag=%v abs=%v\n", c, c, real(c), imag(c), cmplx.Abs(c))
	fmt.Fprintln(w, "cmplx.Sqrt(-1) =", cmplx.Sqrt(-1), "math.Sqrt(-1) =", math.Sqrt(-1))
	fmt.Fprintln(w, "c*c =", c*c)
}

// NumberLiterals shows Go's literal forms: digit separators, binary, octal
// and hex integers, hex floats, exponents and rune escapes.
func NumberLiterals(w io.Writer) {
	fmt.Fprintln(w, 1_000_000, 0b1010, 0o17, 0x1F, 0x1p-2, 1e3, 2.5e-3)
	fmt.Fprintln(w, 'a', '\n', '\x41', 'é', '\U0001F600')
	fmt.Fprintf(w, "%T %T %T\n", 1e3, 0x1p-2, 'é')
	fmt.Fprintln(w, "0644 == 0o644 == 420:", 0644 == 0o644, 0o644 == 420) // a leading 0 alone means octal
}
