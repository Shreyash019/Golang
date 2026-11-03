package examples

import (
	"fmt"
	"io"
	"math"
)

// OperatorPrecedence shows Go's five binary precedence levels, which differ
// from C, Java and Python for shifts and bitwise operators.
func OperatorPrecedence(w io.Writer) {
	fmt.Fprintln(w, "1+2*3 =", 1+2*3)
	fmt.Fprintln(w, "1<<2+1 =", 1<<2+1) // << binds like *, so (1<<2)+1
	fmt.Fprintln(w, "6&3 == 2:", 6&3 == 2)
	fmt.Fprintln(w, "2^3 =", 2^3) // ^ is XOR, not a power
	t, f, g := true, false, false
	fmt.Fprintln(w, "t || f && g:", t || f && g)               // && binds tighter: t || (f && g)
	fmt.Fprintln(w, "10-4-3 =", 10-4-3, " 2*3/4*4 =", 2*3/4*4) // left to right
	fmt.Fprintln(w, "7/2*2.0 =", 7/2*2.0)                      // (7/2) is integer division first
}

// BitwiseOps shows &, |, ^, &^, <<, >> and unary ^ on unsigned and signed
// integers.
func BitwiseOps(w io.Writer) {
	a, b := uint8(0b1100), uint8(0b1010)
	fmt.Fprintf(w, "a&b=%04b a|b=%04b a^b=%04b a&^b=%04b\n", a&b, a|b, a^b, a&^b)
	fmt.Fprintf(w, "^a=%08b a<<2=%08b a>>2=%04b\n", ^a, a<<2, a>>2)

	shift := 8
	var one uint8 = 1
	fmt.Fprintln(w, "uint8 1<<7 =", one<<(shift-1), " 1<<8 =", one<<shift) // shifted-out bits are gone

	var s int8 = -8
	fmt.Fprintln(w, "int8 -8>>1 =", s>>1, " ^-8 =", ^s) // arithmetic shift; ^x == -x-1
	fmt.Fprintln(w, "6 is even:", 6&1 == 0, " 7 is even:", 7&1 == 0)
}

// trace reports that it was evaluated, so short-circuiting is visible.
func trace(w io.Writer, name string, v bool) bool {
	fmt.Fprintf(w, "eval(%s) ", name)
	return v
}

// ShortCircuit shows that && and || evaluate their right operand only when
// needed.
func ShortCircuit(w io.Writer) {
	r := trace(w, "a", false) && trace(w, "b", true)
	fmt.Fprintln(w, "=> false && ... =", r)
	r = trace(w, "a", true) || trace(w, "b", true)
	fmt.Fprintln(w, "=> true || ... =", r)

	var line string
	if len(line) > 0 && line[0] == '#' { // safe: line[0] is never evaluated for ""
		fmt.Fprintln(w, "comment")
	}
	fmt.Fprintln(w, "empty line handled without a panic")
}

// AssignmentOps shows the op= assignment forms and the ++/-- statements.
func AssignmentOps(w io.Writer) {
	x := 10
	x += 5
	x -= 3
	x *= 2
	x /= 4
	x %= 4
	fmt.Fprintln(w, "x =", x)

	flags := uint8(0b0110)
	flags |= 0b0001
	flags &^= 0b0100
	flags <<= 1
	fmt.Fprintf(w, "flags = %04b\n", flags)

	i := 0
	i++
	i++
	i--
	fmt.Fprintln(w, "i =", i)

	s := "go"
	s += "pher"
	fmt.Fprintln(w, "s =", s)
}

// ComparisonRules shows what can be compared and ordered, and the min/max
// builtins (Go 1.21+).
func ComparisonRules(w io.Writer) {
	fmt.Fprintln(w, `"apple" < "banana":`, "apple" < "banana", ` "Zebra" < "apple":`, "Zebra" < "apple", ` "é" > "z":`, "é" > "z")

	var a int32 = 5
	var b int64 = 5
	fmt.Fprintln(w, "int64(a) == b:", int64(a) == b)

	c := complex(1, 2)
	fmt.Fprintln(w, "complex ==:", c == complex(1, 2))

	nan := math.NaN()
	fmt.Fprintln(w, "NaN < 1:", nan < 1, " NaN > 1:", nan > 1, " NaN == NaN:", nan == nan) //nolint:staticcheck // SA4012 on purpose
	fmt.Fprintln(w, "min(3, 1, 2) =", min(3, 1, 2), " max(2.5, 1) =", max(2.5, 1), ` min("b", "a") =`, min("b", "a"))
}
