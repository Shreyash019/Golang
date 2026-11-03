// Package progress reports how far along a batch job is.
//
// Solution E2. The bug: done / total * 100 divides first, and integer
// division truncates, so any done < total gives 0 (and 0 * 100 is 0).
// Multiplying first keeps the precision: done*100 / total. The product
// can't overflow for any realistic job size in a 64-bit int (it would need
// more than 9*10^16 items).
package progress

// Percent returns done as a whole-number percentage of total, rounded
// down: Percent(1, 3) is 33. A total of zero (nothing to do) counts as
// 100%. Callers guarantee 0 <= done <= total.
func Percent(done, total int) int {
	if total == 0 {
		return 100
	}
	return done * 100 / total
}
