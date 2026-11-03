package examples

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"
)

// ZeroValues declares one variable of every kind without an initializer and
// prints what each one holds: memory is always zeroed, so there is no
// "uninitialized" state in Go.
func ZeroValues(w io.Writer) {
	var (
		i   int
		f   float64
		b   bool
		s   string
		r   rune
		p   *int
		sl  []int
		m   map[string]int
		ch  chan int
		fn  func()
		err error
		arr [3]int
		st  struct {
			Name string
			Age  int
		}
	)
	fmt.Fprintf(w, "int=%d float64=%g bool=%t string=%q rune=%d\n", i, f, b, s, r)
	fmt.Fprintf(w, "pointer=%v slice=%v (nil? %t) map=%v (nil? %t)\n", p, sl, sl == nil, m, m == nil)
	fmt.Fprintf(w, "chan nil? %t func nil? %t error=%v\n", ch == nil, fn == nil, err)
	fmt.Fprintf(w, "array=%v struct=%+v\n", arr, st)
	fmt.Fprintf(w, "len(slice)=%d len(map)=%d len(string)=%d\n", len(sl), len(m), len(s))
	fmt.Fprintln(w, "read from nil map:", m["missing"])
}

// UsefulZeroValues uses standard-library types straight from their zero
// value, with no constructor: the design goal "make the zero value useful".
func UsefulZeroValues(w io.Writer) {
	var sb strings.Builder // ready to use
	sb.WriteString("built ")
	sb.WriteString("without a constructor")
	fmt.Fprintln(w, sb.String())

	var buf bytes.Buffer // an empty buffer, ready to write to
	fmt.Fprintf(&buf, "%d bytes buffered", 16)
	fmt.Fprintln(w, buf.String())

	var mu sync.Mutex // an unlocked mutex
	mu.Lock()
	fmt.Fprintln(w, "locked a zero-value mutex")
	mu.Unlock()

	var nums []int // nil slice: append allocates on first use
	nums = append(nums, 1, 2, 3)
	fmt.Fprintln(w, nums, len(nums))
}
