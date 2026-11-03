// Package bytesize formats and parses byte counts with binary units
// (KiB, MiB, GiB, TiB), the way a CLI flag or a config file might write them.
//
// Exercise E3 — implement the unit constants, Format and Parse so that
// bytesize_test.go passes. It trains iota and constant expressions (notes
// §5.5, §6.4), untyped constants mixing with int64 and float64 (§5.3), and
// declaring variables in the narrowest scope without shadowing (§5.7, §7.1).
//
//	make ex E=03-implement-bytesize
package bytesize

// The binary units, in bytes. Replace the zeros with a const block that
// uses iota, so that each unit is 1024 times the previous one. Keep them
// untyped so they work in int64 and float64 arithmetic alike.
const (
	KiB = 0
	MiB = 0
	GiB = 0
	TiB = 0
)

// Format returns n as a human-readable size:
//   - below 1 KiB (in absolute value), the exact count and "B": "512 B", "0 B";
//   - otherwise the largest unit that n reaches, up to TiB, with at most one
//     decimal place and no trailing ".0": "1 KiB", "1.5 KiB", "10 MiB",
//     "2048 TiB";
//   - negative sizes keep their sign: "-2 KiB".
func Format(n int64) string {
	return ""
}

// Parse is the inverse of Format. It accepts a non-negative number,
// optionally with a decimal part, then optional spaces, then an optional
// unit: B, KiB, MiB, GiB or TiB (case-sensitive). No unit means bytes.
// Leading and trailing spaces are ignored. Fractions of a byte are
// truncated: "1.5 B" is 1.
//
// Parse returns an error for an empty string, a missing or malformed
// number, a negative number, an unknown unit ("KB", "kib") and a result
// that doesn't fit in an int64.
func Parse(s string) (int64, error) {
	return 0, nil
}
