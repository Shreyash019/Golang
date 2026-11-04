// Exercise E4 — refactor to idiomatic Go.
//
// Goal: this package works (behaviour_test.go passes) but its control flow
// reads like Java: arrow-shaped nesting with else after return, an if/else-if
// ladder instead of a switch, C-style index loops, a flag variable threaded
// through three loops, a do-while emulated with for true, and comparisons
// with true. Rewrite it as idiomatic Go (notes §8) without changing
// behaviour. idiom_test.go fails until the mechanically detectable habits
// are gone; the review covers the rest.
//
//	make ex E=04-refactor-shipping

// Package shipping quotes parcel shipping costs.
package shipping

import (
	"errors"
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
func Quote(p Parcel) (int64, error) {
	if p.WeightGrams > 0 {
		if p.WeightGrams <= MaxWeightGrams {
			var base int64
			if p.Country == "IN" {
				base = 500
			} else if p.Country == "US" || p.Country == "CA" {
				base = 1500
			} else if p.Country == "DE" || p.Country == "FR" || p.Country == "NL" {
				base = 1200
			} else {
				return 0, errors.New("unsupported country " + p.Country)
			}
			kg := int64((p.WeightGrams + 999) / 1000)
			cost := base + kg*200
			if p.Express == true {
				cost = cost * 2
			}
			return cost, nil
		} else {
			return 0, errors.New("parcel too heavy")
		}
	} else {
		return 0, errors.New("weight must be positive")
	}
}

// CheapestIndex returns the index of the lowest quote (the first one if
// there are ties), or -1 if quotes is empty.
func CheapestIndex(quotes []int64) int {
	result := -1
	if len(quotes) != 0 {
		result = 0
		for i := 0; i < len(quotes); i++ {
			if quotes[i] < quotes[result] {
				result = i
			}
		}
	}
	return result
}

// HasRestricted reports whether any item of any shipment is in restricted.
func HasRestricted(shipments [][]string, restricted []string) bool {
	found := false
	for i := 0; i < len(shipments); i++ {
		for j := 0; j < len(shipments[i]); j++ {
			for k := 0; k < len(restricted); k++ {
				if shipments[i][j] == restricted[k] {
					found = true
					break
				}
			}
			if found == true {
				break
			}
		}
		if found {
			break
		}
	}
	return found
}

// Backoff returns the delays to wait before each of n retries: base, then
// doubling each time, never above limit. It returns nil when n < 1.
func Backoff(base, limit time.Duration, n int) []time.Duration {
	if n < 1 {
		return nil
	} else {
		var delays []time.Duration
		d := base
		for true {
			delays = append(delays, d)
			if len(delays) == n {
				break
			}
			d = d * 2
			if d > limit {
				d = limit
			}
		}
		return delays
	}
}
