package examples

import (
	"fmt"
	"io"
)

// Weekday is a small enumeration built with iota.
type Weekday int

// The days of the week, numbered from 0 by iota.
const (
	Sunday Weekday = iota
	Monday
	Tuesday
)

// IotaBasics shows that iota is the index of the line (ConstSpec) inside its
// const block, and that an omitted expression repeats the previous one.
func IotaBasics(w io.Writer) {
	fmt.Fprintf(w, "Sunday=%d Monday=%d Tuesday=%d (%T)\n", Sunday, Monday, Tuesday, Tuesday)

	const (
		a = iota * 10 // 0
		b             // the expression repeats: iota is now 1, so 10
		_             // 20, discarded
		c             // 30
	)
	fmt.Fprintln(w, "a b c:", a, b, c)

	const (
		x, y = iota, iota + 10 // both on line 0: iota is 0 for both
		z, v                   // line 1: 1, 11
	)
	fmt.Fprintln(w, "x y z v:", x, y, z, v)

	const restart = iota // a new const block starts again at 0
	fmt.Fprintln(w, "restart:", restart)
}

// Perm is a set of permission bits.
type Perm uint8

// Permission bits, one bit each.
const (
	Read  Perm = 1 << iota // 1
	Write                  // 2
	Exec                   // 4
)

// IotaBitFlags builds a bit set with 1 << iota and tests, sets and clears
// bits with &, | and &^.
func IotaBitFlags(w io.Writer) {
	p := Read | Write
	fmt.Fprintf(w, "%03b write? %t exec? %t\n", p, p&Write != 0, p&Exec != 0)
	p &^= Write // clear the Write bit
	p |= Exec   // set the Exec bit
	fmt.Fprintf(w, "%03b write? %t exec? %t\n", p, p&Write != 0, p&Exec != 0)
}

// Binary size units, the classic iota pattern from Effective Go.
const (
	_   = iota // skip 0: 1 << (10*0) would be 1 byte
	KiB = 1 << (10 * iota)
	MiB
	GiB
	TiB
)

// IotaSizes prints the size units and uses them, untyped, in both integer
// and floating-point arithmetic.
func IotaSizes(w io.Writer) {
	fmt.Fprintln(w, KiB, MiB, GiB, TiB)
	var fileSize int64 = 3 * MiB // an untyped constant fits any numeric type
	fmt.Fprintf(w, "%d bytes = %.1f MiB = %.4f GiB\n", fileSize, float64(fileSize)/MiB, float64(fileSize)/GiB)
}

// IotaZeroValueTrap shows why the zero value of an iota enum should mean
// "unknown": otherwise a variable nobody set looks like a real state.
func IotaZeroValueTrap(w io.Writer) {
	type Status int
	const (
		Active Status = iota // 0: also the zero value!
		Suspended
	)
	var s Status // never set: missing JSON field, forgotten struct field…
	fmt.Fprintln(w, "unset status == Active:", s == Active)

	type SafeStatus int
	const (
		StatusUnknown SafeStatus = iota // 0 means "not set"
		StatusActive
		StatusSuspended
	)
	var safe SafeStatus
	fmt.Fprintln(w, "unset safe status == StatusUnknown:", safe == StatusUnknown)
	_, _ = Suspended, StatusSuspended
}

// IotaReorder shows how inserting a constant in the middle of an iota block
// renumbers everything after it, which breaks values already stored outside
// the program (a database row, a message, a file).
func IotaReorder(w io.Writer) {
	const ( // version 1 of the code
		v1Pending = iota
		v1Paid
		v1Shipped
	)
	const ( // version 2: someone inserts Cancelled in the middle
		v2Pending = iota
		v2Cancelled
		v2Paid
		v2Shipped
	)
	stored := v1Paid // version 1 wrote 1 to the database
	fmt.Fprintln(w, "stored value:", stored)
	fmt.Fprintln(w, "v2 reads it as Paid?", stored == v2Paid, "as Cancelled?", stored == v2Cancelled)
	_, _, _, _ = v1Pending, v1Shipped, v2Pending, v2Shipped
}
