// Package progress reports how far along a batch job is.
//
// Exercise E2 — fix the bug. Percent compiles, passes go vet, and is right
// for 0% and 100%, yet a progress bar built on it jumps straight from 0 to
// 100. Find the root cause, fix it with the smallest change, and explain the
// gotcha in one sentence (notes §5.3, §7.6).
//
//	make ex E=02-fix-percentage
package progress

// Percent returns done as a whole-number percentage of total, rounded
// down: Percent(1, 3) is 33. A total of zero (nothing to do) counts as
// 100%. Callers guarantee 0 <= done <= total.
func Percent(done, total int) int {
	if total == 0 {
		return 100
	}
	return done / total * 100
}
