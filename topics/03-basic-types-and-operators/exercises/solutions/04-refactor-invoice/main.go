// Command invoice prices an order: subtotal, a 10% discount over $100,
// 8% tax on the discounted amount, the total, and how many boxes it ships
// in. Items are NAME=PRICE or NAME=PRICExQUANTITY:
//
//	go run ./exercises/solutions/04-refactor-invoice Keyboard=49.99x2 Mouse=19.95 Cable=4.50x3
//
// Solution E4. What changed, and why:
//   - Money is an int64 number of cents, parsed exactly from the text. A
//     float64 can't represent 0.10 or 19.95, so every step needed a round2
//     patch; integers are exact, and rounding happens once, on purpose,
//     where a percentage is applied.
//   - Rates are integer percentages, and percentOf rounds half up with
//     integer arithmetic: (x*p + 50) / 100.
//   - Counts are ints, and the box count is a ceiling division in integers:
//     (units + perBox - 1) / perBox. No float64, no math.Ceil.
//   - Booleans are used directly (if hasQty), and the getter is a plain
//     function named for what it returns.
package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Prices and rates. Money is in cents; rates are percentages.
const (
	taxPercent        = 8
	discountPercent   = 10
	discountThreshold = 100_00 // cents: $100.00
	unitsPerBox       = 4
)

// cents is an amount of money in cents. A defined type keeps it from being
// mixed up with a plain count (defined types: topic 13).
type cents int64

// String formats c as dollars and cents: 12970 is "$129.70".
func (c cents) String() string {
	sign := ""
	if c < 0 {
		sign, c = "-", -c
	}
	return fmt.Sprintf("%s$%d.%02d", sign, c/100, c%100)
}

// LineItem is one line of the order.
type LineItem struct {
	Name      string
	UnitPrice cents
	Quantity  int
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
	items := make([]LineItem, 0, len(args))
	for _, arg := range args {
		item, err := parseItem(arg)
		if err != nil {
			return err
		}
		items = append(items, item)
	}

	sub := subtotal(items)
	units := 0
	for _, it := range items {
		units += it.Quantity
	}
	var discount cents
	if sub > discountThreshold {
		discount = percentOf(sub, discountPercent)
	}
	taxable := sub - discount
	tax := percentOf(taxable, taxPercent)
	total := taxable + tax
	boxes := (units + unitsPerBox - 1) / unitsPerBox

	fmt.Fprintf(w, "items:    %d lines, %d units\n", len(items), units)
	fmt.Fprintf(w, "subtotal: %v\n", sub)
	fmt.Fprintf(w, "discount: -%v\n", discount)
	fmt.Fprintf(w, "tax (%d%%): %v\n", taxPercent, tax)
	fmt.Fprintf(w, "total:    %v\n", total)
	fmt.Fprintf(w, "boxes:    %d (%d units per box)\n", boxes, unitsPerBox)
	return nil
}

// parseItem parses NAME=PRICE or NAME=PRICExQUANTITY.
func parseItem(s string) (LineItem, error) {
	name, rest, ok := strings.Cut(s, "=")
	if !ok || name == "" {
		return LineItem{}, fmt.Errorf("item %q: want NAME=PRICE or NAME=PRICExQUANTITY", s)
	}
	priceText, qtyText, hasQty := strings.Cut(rest, "x")
	price, err := parseCents(priceText)
	if err != nil {
		return LineItem{}, fmt.Errorf("item %q: %w", s, err)
	}
	qty := 1
	if hasQty {
		qty, err = strconv.Atoi(qtyText)
		if err != nil || qty < 1 {
			return LineItem{}, fmt.Errorf("item %q: invalid quantity %q", s, qtyText)
		}
	}
	return LineItem{Name: name, UnitPrice: price, Quantity: qty}, nil
}

// parseCents parses a non-negative dollar amount with up to two decimal
// places ("19", "4.5", "19.95") into cents, exactly: no float involved.
func parseCents(s string) (cents, error) {
	whole, frac, hasFrac := strings.Cut(s, ".")
	if whole == "" || len(frac) > 2 || (hasFrac && frac == "") {
		return 0, fmt.Errorf("invalid price %q", s)
	}
	dollars, err := strconv.ParseUint(whole, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid price %q", s)
	}
	var c uint64
	if frac != "" {
		if c, err = strconv.ParseUint(frac, 10, 8); err != nil {
			return 0, fmt.Errorf("invalid price %q", s)
		}
		if len(frac) == 1 { // "4.5" means 50 cents, not 5
			c *= 10
		}
	}
	return cents(dollars*100 + c), nil
}

// subtotal returns the sum of price × quantity over all items.
func subtotal(items []LineItem) cents {
	var sum cents
	for _, it := range items {
		sum += it.UnitPrice * cents(it.Quantity)
	}
	return sum
}

// percentOf returns p percent of c, rounded half up to the nearest cent.
// Amounts are never negative here, so half up and half away from zero
// agree.
func percentOf(c cents, p int64) cents {
	return (c*cents(p) + 50) / 100
}
