// Package shipping quotes parcel shipping costs.
package shipping

import (
	"fmt"
	"slices"
	"time"
)

// MaxWeightGrams is the heaviest parcel that can be shipped.
const MaxWeightGrams = 30000

// Parcel is a parcel to ship.
type Parcel struct {
	WeightGrams int
	Country     string // ISO 3166-1 alpha-2 code, e.g. "IN"
	Express     bool
}

// Quote returns the cost of shipping p, in cents: a base price for the
// destination plus 200 cents per started kilogram, doubled for express.
//
// Guard clauses handle each failure and return, so the happy path runs down
// the left margin instead of three levels deep, and each error sits next to
// the condition that causes it. The country ladder becomes a switch with
// several values per case, and its default case is the unsupported country.
func Quote(p Parcel) (int64, error) {
	if p.WeightGrams <= 0 {
		return 0, fmt.Errorf("weight must be positive, got %d g", p.WeightGrams)
	}
	if p.WeightGrams > MaxWeightGrams {
		return 0, fmt.Errorf("parcel too heavy: %d g, max %d g", p.WeightGrams, MaxWeightGrams)
	}

	var base int64
	switch p.Country {
	case "IN":
		base = 500
	case "US", "CA":
		base = 1500
	case "DE", "FR", "NL":
		base = 1200
	default:
		return 0, fmt.Errorf("unsupported country %q", p.Country)
	}

	kg := int64((p.WeightGrams + 999) / 1000) // started kilograms: ceiling division
	cost := base + kg*200
	if p.Express {
		cost *= 2
	}
	return cost, nil
}

// CheapestIndex returns the index of the lowest quote (the first one if
// there are ties), or -1 if quotes is empty.
//
// The empty case is a guard clause, so result never needs a placeholder
// value, and range removes the index bookkeeping.
func CheapestIndex(quotes []int64) int {
	if len(quotes) == 0 {
		return -1
	}
	best := 0
	for i, q := range quotes {
		if q < quotes[best] {
			best = i
		}
	}
	return best
}

// HasRestricted reports whether any item of any shipment is in restricted.
//
// Returning as soon as the answer is known replaces the found flag and the
// three breaks that propagated it outward. slices.Contains states the
// innermost loop's intent in one call.
func HasRestricted(shipments [][]string, restricted []string) bool {
	for _, items := range shipments {
		for _, item := range items {
			if slices.Contains(restricted, item) {
				return true
			}
		}
	}
	return false
}

// Backoff returns the delays to wait before each of n retries: base, then
// doubling each time, never above limit. It returns nil when n < 1.
//
// The number of iterations is known, so range over n says so, and the
// do-while emulation (for true + break) disappears. min (Go 1.21) caps the
// delay without an if.
func Backoff(base, limit time.Duration, n int) []time.Duration {
	if n < 1 {
		return nil
	}
	delays := make([]time.Duration, 0, n)
	d := base
	for range n {
		delays = append(delays, d)
		d = min(d*2, limit)
	}
	return delays
}
