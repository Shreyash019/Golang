// Command checkout prints a receipt for a list of item prices, in cents,
// with an optional coupon.
//
//	checkout [-coupon SAVE10|FLAT500|BOGO] PRICE_CENTS...
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
)

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "checkout:", err)
		os.Exit(1)
	}
}

// discount returns the discount, in cents, for a list of prices.
//
// A one-method interface with stateless implementations is a function in
// disguise. A func type replaces the interface, the four Impl types and
// their methods; closures (percentOff, flatOff) carry the little state the
// structs held.
type discount func(prices []int64) int64

// discountFor returns the discount for a coupon code. It replaces the
// factory object: a plain function, returning an error for unknown codes
// instead of panicking, because a mistyped coupon is user input, not a bug.
func discountFor(code string) (discount, error) {
	switch code {
	case "":
		return noDiscount, nil
	case "SAVE10":
		return percentOff(10), nil
	case "FLAT500":
		return flatOff(500), nil
	case "BOGO":
		return cheapestFree, nil
	default:
		return nil, fmt.Errorf("unknown coupon %q", code)
	}
}

func noDiscount([]int64) int64 { return 0 }

func percentOff(percent int64) discount {
	return func(prices []int64) int64 {
		return subtotal(prices) * percent / 100
	}
}

func flatOff(amount int64) discount {
	return func(prices []int64) int64 {
		return min(amount, subtotal(prices))
	}
}

func cheapestFree(prices []int64) int64 {
	if len(prices) < 2 {
		return 0
	}
	return slices.Min(prices)
}

func subtotal(prices []int64) int64 {
	var total int64
	for _, p := range prices {
		total += p
	}
	return total
}

// parsePrice returns the price in s. Returning (value, error) replaces the
// out-parameter and the bool: the caller gets the value and a reason it
// can show, and can't forget which argument gets written.
func parsePrice(s string) (int64, error) {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil || v < 0 {
		return 0, fmt.Errorf("invalid price %q: want a whole, non-negative number of cents", s)
	}
	return v, nil
}

// run is the whole program: it parses the arguments and prints the receipt.
// Each failure returns an error where it's detected; there's nothing to
// recover, so the deferred recover and its named result are gone.
func run(args []string, w io.Writer) error {
	fs := flag.NewFlagSet("checkout", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	coupon := fs.String("coupon", "", "coupon code: SAVE10, FLAT500 or BOGO")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() == 0 {
		return errors.New("no items")
	}

	prices := make([]int64, 0, fs.NArg())
	for _, arg := range fs.Args() {
		price, err := parsePrice(arg)
		if err != nil {
			return err
		}
		prices = append(prices, price)
	}

	off, err := discountFor(*coupon)
	if err != nil {
		return err
	}
	sub := subtotal(prices)
	d := off(prices)

	fmt.Fprintf(w, "items:    %d\n", len(prices))
	fmt.Fprintf(w, "subtotal: %d\n", sub)
	fmt.Fprintf(w, "discount: %d\n", d)
	fmt.Fprintf(w, "total:    %d\n", sub-d)
	return nil
}
