// Package checked does integer arithmetic and conversions that report
// overflow instead of silently wrapping around. Code that sums file sizes,
// quotas or money in an int64 needs this: Go's + and * wrap without a
// trace (notes §5.2, §7.1).
//
// Exercise E3 — implement every function so that checked_test.go passes.
// It trains two's complement ranges (notes §5.2), conversions (§5.5),
// floating-point special values (§5.4) and the asymmetry of signed
// integers (§7.11). Use only Go operators and the math package's
// constants and functions: no math/big.
//
//	make ex E=03-implement-checked
package checked

// Add returns a+b and true, or 0 and false if the sum overflows int64.
func Add(a, b int64) (int64, bool) {
	return 0, false
}

// Sub returns a-b and true, or 0 and false if the difference overflows
// int64.
func Sub(a, b int64) (int64, bool) {
	return 0, false
}

// Mul returns a*b and true, or 0 and false if the product overflows int64.
func Mul(a, b int64) (int64, bool) {
	return 0, false
}

// Abs returns the absolute value of v and true, or 0 and false for
// math.MinInt64, whose absolute value doesn't fit in an int64.
func Abs(v int64) (int64, bool) {
	return 0, false
}

// ToInt32 converts v to int32, or returns 0 and false if v is outside the
// int32 range.
func ToInt32(v int64) (int32, bool) {
	return 0, false
}

// RoundToInt64 rounds f to the nearest integer, with halves rounded away
// from zero (2.5 → 3, -2.5 → -3), and converts it to int64. It returns 0
// and false if f is NaN, infinite, or rounds to a value outside the int64
// range.
func RoundToInt64(f float64) (int64, bool) {
	return 0, false
}
