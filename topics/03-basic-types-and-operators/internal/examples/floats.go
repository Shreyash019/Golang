package examples

import (
	"fmt"
	"io"
	"math"
)

// FloatSpecialValues shows IEEE 754's infinities, NaN and negative zero:
// floating-point division by zero doesn't panic.
func FloatSpecialValues(w io.Writer) {
	var zero float64
	fmt.Fprintln(w, "1/0 =", 1/zero, " -1/0 =", -1/zero, " 0/0 =", zero/zero)

	nan := math.NaN()
	fmt.Fprintln(w, "NaN == NaN:", nan == nan, " NaN != NaN:", nan != nan, " IsNaN:", math.IsNaN(nan)) //nolint:staticcheck // SA4012 on purpose

	negZero := math.Copysign(0, -1)
	fmt.Fprintln(w, "-0:", negZero, " -0 == 0:", negZero == 0, " 1/-0 =", 1/negZero)

	huge := math.MaxFloat64
	fmt.Fprintln(w, "MaxFloat64*2 =", huge*2, " Inf-Inf =", math.Inf(1)-math.Inf(1))
}

// FloatPrecision shows what binary floating point can't represent exactly.
func FloatPrecision(w io.Writer) {
	a, b := 0.1, 0.2
	fmt.Fprintln(w, "0.1+0.2 =", a+b)
	fmt.Fprintf(w, "0.1 is really %.25f\n", a)

	var f32 float32 = 16_777_216 // 2^24: float32 has a 24-bit significand
	fmt.Fprintln(w, "float32 2^24+1 == 2^24:", f32+1 == f32)
	f64 := float64(1 << 53) // float64 has a 53-bit significand
	fmt.Fprintln(w, "float64 2^53+1 == 2^53:", f64+1 == f64)

	sum := 0.0
	for range 10 {
		sum += 0.1
	}
	fmt.Fprintln(w, "ten 0.1s:", sum, "== 1:", sum == 1)

	cents := 0
	for range 10 {
		cents += 10
	}
	fmt.Fprintln(w, "ten 10-cent ints == 100:", cents == 100)
}

// almostEqual reports whether a and b are equal within a relative
// tolerance. There is no single right tolerance: it depends on the
// magnitudes and the computation.
func almostEqual(a, b, relTol float64) bool {
	return math.Abs(a-b) <= relTol*math.Max(math.Abs(a), math.Abs(b))
}

// FloatComparison compares floats with == and with a tolerance.
func FloatComparison(w io.Writer) {
	a, b := 0.1, 0.2
	fmt.Fprintln(w, "a+b == 0.3:", a+b == 0.3)
	fmt.Fprintln(w, "almostEqual(a+b, 0.3, 1e-9):", almostEqual(a+b, 0.3, 1e-9))
	fmt.Fprintln(w, "almostEqual(1e-20, 2e-20, 1e-9):", almostEqual(1e-20, 2e-20, 1e-9))
}
