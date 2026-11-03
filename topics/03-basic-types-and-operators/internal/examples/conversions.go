package examples

import (
	"fmt"
	"io"
	"math"
	"strconv"
)

// Conversions shows what T(x) does between numeric types: it never fails
// at runtime, and it silently truncates, wraps or loses precision.
func Conversions(w io.Writer) {
	f := -2.7
	fmt.Fprintln(w, "int(-2.7) =", int(f), " Round:", int(math.Round(f)), " Floor:", int(math.Floor(f)))

	x := int64(300)
	fmt.Fprintln(w, "int8(300) =", int8(x), " uint8(300) =", uint8(x))

	y := int64(-1)
	fmt.Fprintln(w, "uint8(-1) =", uint8(y), " uint64(-1) =", uint64(y))

	big := int64(1<<53 + 1)
	fmt.Fprintln(w, "float64(2^53+1) == float64(2^53):", float64(big) == float64(big-1))

	v := int64(0x1234_5678)
	fmt.Fprintf(w, "int16(%#x) = %#x\n", v, int16(v))
}

// StringConversions contrasts string(rune) with strconv, the package that
// converts numbers to and from their text form (strconv in depth: topic 31).
func StringConversions(w io.Writer) {
	n := 65
	fmt.Fprintln(w, "string(rune(65)) =", string(rune(n)), " strconv.Itoa(65) =", strconv.Itoa(n))

	v, err := strconv.Atoi("42")
	fmt.Fprintln(w, "Atoi(\"42\")+1 =", v+1, "err:", err)
	_, err = strconv.Atoi("4 2")
	fmt.Fprintln(w, "Atoi(\"4 2\") err:", err)

	fmt.Fprintln(w, "255 in base 2 and 16:", strconv.FormatInt(255, 2), strconv.FormatInt(255, 16))
	_, err = strconv.ParseInt("300", 10, 8) // parse into an 8-bit range
	fmt.Fprintln(w, "ParseInt(\"300\", 10, 8) err:", err)

	r := []rune("héllo")
	fmt.Fprintln(w, "[]byte(\"hi\") =", []byte("hi"), " []rune(\"héllo\") has", len(r), "runes") //nolint:staticcheck // QF1010: printing the bytes is the point
}
