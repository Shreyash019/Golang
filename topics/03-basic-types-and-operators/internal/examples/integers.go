package examples

import (
	"fmt"
	"io"
	"math"
	"math/bits"
)

// IntegerOverflow shows that integer arithmetic on variables wraps around
// silently (two's complement), at every size, signed or unsigned.
func IntegerOverflow(w io.Writer) {
	var i8 int8 = 127
	i8++
	fmt.Fprintln(w, "int8 127+1 =", i8)

	var u8 uint8
	u8--
	fmt.Fprintln(w, "uint8 0-1 =", u8)

	big := math.MaxInt64
	big++
	fmt.Fprintln(w, "MaxInt64+1 =", big)

	low := int8(-128)
	fmt.Fprintln(w, "-(-128) as int8 =", -low)
	fmt.Fprintf(w, "bits: 127=%08b -128=%08b -1=%08b\n", uint8(127), uint8(low), uint8(low+127))
}

// DivisionAndRemainder shows that integer division truncates toward zero
// and the remainder takes the sign of the dividend.
func DivisionAndRemainder(w io.Writer) {
	a, b := -7, 2
	fmt.Fprintln(w, "7/2 =", 7/b, " -7/2 =", a/b, " -7%2 =", a%b, " 7%-2 =", 7%-b)
	fmt.Fprintln(w, "-7>>1 =", a>>1, "(a shift rounds toward minus infinity)")

	n, idx := 5, -1
	fmt.Fprintln(w, "-1%5 =", idx%n, " wrapped:", (idx%n+n)%n)
}

// UnsignedUnderflow shows why lengths and indexes are ints in Go: an
// unsigned "length minus one" wraps to a huge number instead of -1.
func UnsignedUnderflow(w io.Writer) {
	var items []string
	n := uint(len(items))
	fmt.Fprintln(w, "uint len-1 =", n-1)
	fmt.Fprintln(w, "int  len-1 =", len(items)-1)
}

// OverflowDetection uses math/bits, which reports carries and full
// 128-bit products for unsigned arithmetic, and counts bits.
func OverflowDetection(w io.Writer) {
	sum, carry := bits.Add64(math.MaxUint64, 1, 0)
	fmt.Fprintln(w, "MaxUint64+1: sum =", sum, "carry =", carry)

	hi, lo := bits.Mul64(1<<40, 1<<40) // 2^80 doesn't fit in 64 bits
	fmt.Fprintln(w, "2^40*2^40: hi =", hi, "lo =", lo)

	var x uint64 = 0b1011_0000
	fmt.Fprintln(w, "ones:", bits.OnesCount64(x), "len:", bits.Len64(x), "trailing zeros:", bits.TrailingZeros64(x))
}

// IntDivisionByZero crashes: integer division by a zero variable panics.
func IntDivisionByZero(w io.Writer) {
	total, count := 10, 0
	fmt.Fprintln(w, "average:", total/count)
}

// NegativeShift crashes: a negative shift count panics at runtime.
func NegativeShift(w io.Writer) {
	n := -1
	fmt.Fprintln(w, 1<<n)
}
