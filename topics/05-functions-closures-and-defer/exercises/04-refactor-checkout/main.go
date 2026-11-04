// Exercise E4 — refactor to idiomatic Go.
//
// Goal: this command works (behaviour_test.go passes), but it's Java
// translated into Go: a strategy interface with one "Impl" type per coupon,
// a factory object, exceptions emulated with panic/recover, and an
// out-parameter instead of a second result. (The methods on the Impl types
// are topic 10; you won't need any.) Rewrite it as idiomatic Go without
// changing its behaviour: functions as values, a constructor function that
// returns (value, error), and errors returned instead of panics (notes §8).
// idiom_test.go fails until the mechanically detectable habits are gone.
//
//	make ex E=04-refactor-checkout
//	go run ./exercises/04-refactor-checkout -coupon SAVE10 1999 500
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
)

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// DiscountStrategy computes the discount, in cents, for a list of prices.
type DiscountStrategy interface {
	CalculateDiscount(prices []int64) int64
}

// NoDiscountStrategyImpl gives no discount.
type NoDiscountStrategyImpl struct{}

// CalculateDiscount implements DiscountStrategy.
func (s *NoDiscountStrategyImpl) CalculateDiscount(prices []int64) int64 {
	return 0
}

// PercentDiscountStrategyImpl takes a percentage off the subtotal.
type PercentDiscountStrategyImpl struct {
	percent int64
}

// CalculateDiscount implements DiscountStrategy.
func (s *PercentDiscountStrategyImpl) CalculateDiscount(prices []int64) int64 {
	return calculateSubtotal(prices) * s.percent / 100
}

// FlatDiscountStrategyImpl takes a fixed amount off, never more than the
// subtotal.
type FlatDiscountStrategyImpl struct {
	amount int64
}

// CalculateDiscount implements DiscountStrategy.
func (s *FlatDiscountStrategyImpl) CalculateDiscount(prices []int64) int64 {
	subtotal := calculateSubtotal(prices)
	if s.amount > subtotal {
		return subtotal
	}
	return s.amount
}

// BuyOneGetOneStrategyImpl makes the cheapest item free when there are at
// least two items.
type BuyOneGetOneStrategyImpl struct{}

// CalculateDiscount implements DiscountStrategy.
func (s *BuyOneGetOneStrategyImpl) CalculateDiscount(prices []int64) int64 {
	if len(prices) < 2 {
		return 0
	}
	cheapest := prices[0]
	for _, p := range prices {
		if p < cheapest {
			cheapest = p
		}
	}
	return cheapest
}

// DiscountStrategyFactory creates discount strategies from coupon codes.
type DiscountStrategyFactory struct{}

// NewDiscountStrategyFactory returns a new DiscountStrategyFactory.
func NewDiscountStrategyFactory() *DiscountStrategyFactory {
	return &DiscountStrategyFactory{}
}

// CreateStrategy returns the strategy for a coupon code. It panics for an
// unknown code.
func (f *DiscountStrategyFactory) CreateStrategy(code string) DiscountStrategy {
	switch code {
	case "":
		return &NoDiscountStrategyImpl{}
	case "SAVE10":
		return &PercentDiscountStrategyImpl{percent: 10}
	case "FLAT500":
		return &FlatDiscountStrategyImpl{amount: 500}
	case "BOGO":
		return &BuyOneGetOneStrategyImpl{}
	}
	panic(fmt.Sprintf("unknown coupon %q", code))
}

func calculateSubtotal(prices []int64) int64 {
	var subtotal int64
	for _, p := range prices {
		subtotal += p
	}
	return subtotal
}

// parsePrice parses s into *result and reports whether it succeeded.
func parsePrice(s string, result *int64) bool {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil || v < 0 {
		return false
	}
	*result = v
	return true
}

// run is the whole program: it parses the arguments and prints the receipt.
func run(args []string, w io.Writer) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("checkout: %v", r)
		}
	}()

	fs := flag.NewFlagSet("checkout", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	coupon := fs.String("coupon", "", "coupon code: SAVE10, FLAT500 or BOGO")
	if parseErr := fs.Parse(args); parseErr != nil {
		panic(parseErr)
	}
	if fs.NArg() == 0 {
		panic("no items")
	}

	prices := make([]int64, 0, fs.NArg())
	for _, arg := range fs.Args() {
		var price int64
		if !parsePrice(arg, &price) {
			panic(fmt.Sprintf("invalid price %q: want a whole, non-negative number of cents", arg))
		}
		prices = append(prices, price)
	}

	strategy := NewDiscountStrategyFactory().CreateStrategy(*coupon)
	subtotal := calculateSubtotal(prices)
	discount := strategy.CalculateDiscount(prices)

	fmt.Fprintf(w, "items:    %d\n", len(prices))
	fmt.Fprintf(w, "subtotal: %d\n", subtotal)
	fmt.Fprintf(w, "discount: %d\n", discount)
	fmt.Fprintf(w, "total:    %d\n", subtotal-discount)
	return nil
}
