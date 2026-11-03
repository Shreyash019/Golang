// Exercise E4 — refactor to idiomatic Go.
//
// Command invoice prices an order: subtotal, a 10% discount over $100,
// 8% tax on the discounted amount, the total, and how many boxes it ships
// in. Items are NAME=PRICE or NAME=PRICExQUANTITY:
//
//	go run ./exercises/04-refactor-invoice Keyboard=49.99x2 Mouse=19.95 Cable=4.50x3
//
// It works, and behaviour_test.go passes. But it's written the way a lot of
// Python and JavaScript code handles numbers: money in float64 (parsed with
// ParseFloat and patched with a round2 helper), counts in float64,
// math.Ceil on floats for a ceiling division, comparisons with == true,
// and a Java-style getter.
//
// Rewrite it as idiomatic Go (notes §8, §10): money as an integer number of
// cents, parsed exactly from the text; integer arithmetic with explicit
// rounding; no float64 and no math package. behaviour_test.go must keep
// passing; idiom_test.go fails until the habits are gone:
//
//	make ex E=04-refactor-invoice
package main

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
	taxRate           = 0.08
	discountRate      = 0.10
	discountThreshold = 100.0
	unitsPerBox       = 4.0
)

// LineItem is one line of the order.
type LineItem struct {
	Name      string
	UnitPrice float64 // dollars
	Quantity  float64
}

func main() {
	if err := run(os.Stdout, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "invoice:", err)
		os.Exit(1)
	}
}

func run(w io.Writer, args []string) error {
	if len(args) == 0 {
		return errors.New("no items: pass NAME=PRICE or NAME=PRICExQUANTITY")
	}
	var items []LineItem
	for _, arg := range args {
		item, err := parseItem(arg)
		if err != nil {
			return err
		}
		items = append(items, item)
	}

	subtotal := GetSubtotal(items)
	units := 0.0
	for _, it := range items {
		units = units + it.Quantity
	}
	discount := 0.0
	isLarge := subtotal > discountThreshold
	if isLarge == true {
		discount = round2(subtotal * discountRate)
	}
	taxable := round2(subtotal - discount)
	tax := round2(taxable * taxRate)
	total := round2(taxable + tax)
	boxes := int(math.Ceil(units / unitsPerBox))

	fmt.Fprintf(w, "items:    %d lines, %d units\n", len(items), int(units))
	fmt.Fprintf(w, "subtotal: $%.2f\n", subtotal)
	fmt.Fprintf(w, "discount: -$%.2f\n", discount)
	fmt.Fprintf(w, "tax (%d%%): $%.2f\n", int(math.Round(taxRate*100)), tax)
	fmt.Fprintf(w, "total:    $%.2f\n", total)
	fmt.Fprintf(w, "boxes:    %d (%d units per box)\n", boxes, int(unitsPerBox))
	return nil
}

// parseItem parses NAME=PRICE or NAME=PRICExQUANTITY.
func parseItem(s string) (LineItem, error) {
	name, rest, ok := strings.Cut(s, "=")
	if !ok || name == "" {
		return LineItem{}, fmt.Errorf("item %q: want NAME=PRICE or NAME=PRICExQUANTITY", s)
	}
	priceText, qtyText, hasQty := strings.Cut(rest, "x")
	price, err := strconv.ParseFloat(priceText, 64)
	if err != nil || price < 0 {
		return LineItem{}, fmt.Errorf("item %q: invalid price %q", s, priceText)
	}
	qty := 1.0
	if hasQty == true {
		qty, err = strconv.ParseFloat(qtyText, 64)
		if err != nil || qty < 1 || qty != math.Trunc(qty) {
			return LineItem{}, fmt.Errorf("item %q: invalid quantity %q", s, qtyText)
		}
	}
	return LineItem{Name: name, UnitPrice: price, Quantity: qty}, nil
}

// GetSubtotal returns the sum of price × quantity over all items.
func GetSubtotal(items []LineItem) float64 {
	sum := 0.0
	for _, it := range items {
		sum = sum + it.UnitPrice*it.Quantity
	}
	return round2(sum)
}

// round2 rounds to cents.
func round2(x float64) float64 { return math.Round(x*100) / 100 }
