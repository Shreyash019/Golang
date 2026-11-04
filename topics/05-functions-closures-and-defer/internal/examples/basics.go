package examples

import (
	"fmt"
	"io"
	"strconv"
)

// PassByValue shows that every argument is copied into the parameter. What
// the copy shares with the caller depends on the type: an int or an array
// shares nothing, a slice header shares its backing array, a map value is a
// handle to the same map, and a pointer points at the caller's variable.
func PassByValue(w io.Writer) {
	n := 1
	arr := [2]int{1, 2}
	sl := []int{1, 2}
	m := map[string]int{"a": 1}
	x := 1

	mutate(n, arr, sl, m, &x)

	fmt.Fprintln(w, "n:  ", n, "    (the int was copied)")
	fmt.Fprintln(w, "x:  ", x, "   (written through the pointer)")
	fmt.Fprintln(w, "arr:", arr, "(the whole array was copied)")
	fmt.Fprintln(w, "sl: ", sl, "(element write shared; the append was not)")
	fmt.Fprintln(w, "m:  ", m, "(same map)")
}

//nolint:staticcheck // SA4009: overwriting the copies is what the example shows
func mutate(n int, arr [2]int, sl []int, m map[string]int, p *int) {
	n = 100            //nolint:staticcheck // changes the local copy
	arr[0] = 100       // changes the local copy of the array
	sl[0] = 100        // writes through the copied header into the shared array
	sl = append(sl, 3) //nolint:ineffassign // grows the local header only
	m["b"] = 2         // the map is shared
	*p = 42            // writes the caller's variable
	_, _ = n, arr
}

// MultipleReturns shows functions returning several values, the comma-ok
// shape, and discarding results with the blank identifier.
func MultipleReturns(w io.Writer) {
	q, r := divmod(17, 5)
	fmt.Fprintln(w, "divmod(17, 5) =", q, r)

	_, r = divmod(-17, 5)
	fmt.Fprintln(w, "remainder only:", r)

	if n, err := strconv.Atoi("12a"); err != nil {
		fmt.Fprintln(w, "Atoi failed:", err, "| n is the zero value:", n)
	}

	// A call returning several values can feed a call that takes exactly
	// those parameters.
	fmt.Fprintln(w, "sum of divmod:", add(divmod(17, 5)))
}

func divmod(a, b int) (int, int) { return a / b, a % b }

func add(a, b int) int { return a + b }

// NamedResults shows named result parameters: they start at their zero
// values, a bare "return" returns their current values, and they document
// what each result means.
func NamedResults(w io.Writer) {
	lo, hi := minMax([]int{4, -2, 9, 0})
	fmt.Fprintln(w, "minMax:", lo, hi)

	lo, hi = minMax(nil)
	fmt.Fprintln(w, "minMax(nil):", lo, hi, "(zero values)")
}

func minMax(xs []int) (lo, hi int) {
	if len(xs) == 0 {
		return // lo and hi are still 0
	}
	lo, hi = xs[0], xs[0]
	for _, x := range xs[1:] {
		lo = min(lo, x)
		hi = max(hi, x)
	}
	return lo, hi // explicit, even though "return" alone would work
}

// Variadic shows a ...T parameter: inside the function it's a []T. Callers
// pass zero or more values, or an existing slice with s...
func Variadic(w io.Writer) {
	fmt.Fprintln(w, sum(), sum(1), sum(1, 2, 3))
	nums := []int{10, 20, 30}
	fmt.Fprintln(w, sum(nums...))
	fmt.Fprintln(w, describe("ids"), describe("ids", 7, 8))
}

func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

func describe(label string, nums ...int) string {
	return fmt.Sprintf("%s: len=%d nil=%t", label, len(nums), nums == nil)
}

// VariadicAliasing shows that f(s...) passes the caller's slice itself, not a
// copy: the function can modify the caller's elements.
func VariadicAliasing(w io.Writer) {
	scores := []int{1, 2, 3}
	double(scores...)
	fmt.Fprintln(w, "after double(scores...):", scores)

	double(4, 5, 6) // a new slice is built for these; nothing to alias
	fmt.Fprintln(w, "after double(4, 5, 6): nothing visible to change")
}

func double(nums ...int) {
	for i := range nums {
		nums[i] *= 2
	}
}
