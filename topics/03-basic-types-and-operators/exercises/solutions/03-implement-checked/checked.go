// Package checked does integer arithmetic and conversions that report
// overflow instead of silently wrapping around.
//
// Solution E3. The technique throughout: do the operation (it wraps, which
// is well defined in Go), then detect the wrap from the signs, or check the
// range before converting. No operation here can panic.
package checked

import "math"

// Add returns a+b and true, or 0 and false if the sum overflows int64.
func Add(a, b int64) (int64, bool) {
	s := a + b
	// Adding a positive number must make the result bigger, and adding a
	// negative one smaller. A wrapped result moves the other way.
	if (b > 0 && s < a) || (b < 0 && s > a) {
		return 0, false
	}
	return s, true
}

// Sub returns a-b and true, or 0 and false if the difference overflows
// int64.
func Sub(a, b int64) (int64, bool) {
	d := a - b
	// The mirror of Add. Note that Add(a, -b) would be wrong for
	// b == MinInt64, because -MinInt64 wraps back to MinInt64.
	if (b > 0 && d > a) || (b < 0 && d < a) {
		return 0, false
	}
	return d, true
}

// Mul returns a*b and true, or 0 and false if the product overflows int64.
func Mul(a, b int64) (int64, bool) {
	if a == 0 || b == 0 {
		return 0, true
	}
	// MinInt64 * -1 is the one product where the division check below
	// would itself overflow (MinInt64 / -1 wraps to MinInt64), so it is
	// handled first.
	if (a == -1 && b == math.MinInt64) || (b == -1 && a == math.MinInt64) {
		return 0, false
	}
	p := a * b
	if p/b != a {
		return 0, false
	}
	return p, true
}

// Abs returns the absolute value of v and true, or 0 and false for
// math.MinInt64, whose absolute value doesn't fit in an int64.
func Abs(v int64) (int64, bool) {
	if v == math.MinInt64 {
		return 0, false
	}
	if v < 0 {
		return -v, true
	}
	return v, true
}

// ToInt32 converts v to int32, or returns 0 and false if v is outside the
// int32 range.
func ToInt32(v int64) (int32, bool) {
	if v < math.MinInt32 || v > math.MaxInt32 {
		return 0, false
	}
	return int32(v), true
}

// RoundToInt64 rounds f to the nearest integer, with halves rounded away
// from zero (2.5 → 3, -2.5 → -3), and converts it to int64. It returns 0
// and false if f is NaN, infinite, or rounds to a value outside the int64
// range.
func RoundToInt64(f float64) (int64, bool) {
	r := math.Round(f)
	// The int64 range is [-2^63, 2^63). Both bounds are exact float64
	// values: math.MinInt64 is -2^63, and math.MaxInt64 (2^63-1) rounds up
	// to 2^63 when converted to float64, which is why the upper check is
	// >= and not >. Every comparison with NaN is false, so NaN needs its
	// own check; the infinities fail the range checks.
	if math.IsNaN(r) || r < math.MinInt64 || r >= math.MaxInt64 {
		return 0, false
	}
	return int64(r), true
}
