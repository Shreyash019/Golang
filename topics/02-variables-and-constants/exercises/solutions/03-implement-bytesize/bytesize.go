// Package bytesize formats and parses byte counts with binary units
// (KiB, MiB, GiB, TiB), the way a CLI flag or a config file might write them.
//
// Solution E3.
package bytesize

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// The binary units, in bytes. They are untyped constants, so the same KiB
// works as an int64 in Parse and as a float64 divisor in Format without a
// conversion.
const (
	_   = iota // skip 0: 1 << (10 * 0) would be one byte
	KiB = 1 << (10 * iota)
	MiB
	GiB
	TiB
)

// Format returns n as a human-readable size:
//   - below 1 KiB (in absolute value), the exact count and "B": "512 B", "0 B";
//   - otherwise the largest unit that n reaches, up to TiB, with at most one
//     decimal place and no trailing ".0": "1 KiB", "1.5 KiB", "10 MiB",
//     "2048 TiB";
//   - negative sizes keep their sign: "-2 KiB".
func Format(n int64) string {
	// The units from largest to smallest, so the first one n reaches wins.
	// A local array, not a package variable: nothing outside Format can
	// change it (structs: topic 10).
	units := [...]struct {
		name string
		size int64
	}{
		{"TiB", TiB},
		{"GiB", GiB},
		{"MiB", MiB},
		{"KiB", KiB},
	}
	abs := n
	if abs < 0 {
		abs = -abs
	}
	for _, u := range units {
		if abs >= u.size {
			v := strconv.FormatFloat(float64(n)/float64(u.size), 'f', 1, 64)
			return strings.TrimSuffix(v, ".0") + " " + u.name
		}
	}
	return strconv.FormatInt(n, 10) + " B"
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
	s = strings.TrimSpace(s)
	end := strings.IndexFunc(s, func(r rune) bool { return (r < '0' || r > '9') && r != '.' })
	if end == -1 {
		end = len(s)
	}
	num, unit := s[:end], strings.TrimSpace(s[end:])
	if num == "" {
		return 0, fmt.Errorf("parse size %q: missing number", s)
	}
	v, err := strconv.ParseFloat(num, 64)
	if err != nil {
		return 0, fmt.Errorf("parse size %q: %w", s, err)
	}

	var size int64
	switch unit {
	case "", "B":
		size = 1
	case "KiB":
		size = KiB
	case "MiB":
		size = MiB
	case "GiB":
		size = GiB
	case "TiB":
		size = TiB
	default:
		return 0, fmt.Errorf("parse size %q: unknown unit %q (want B, KiB, MiB, GiB or TiB)", s, unit)
	}

	total := v * float64(size)
	// float64(math.MaxInt64) rounds up to 2^63, so >= rejects every value
	// that int64 can't hold.
	if total >= math.MaxInt64 {
		return 0, fmt.Errorf("parse size %q: too large", s)
	}
	return int64(total), nil
}
