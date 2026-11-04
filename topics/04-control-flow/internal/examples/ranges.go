package examples

import (
	"fmt"
	"io"
	"iter"
	"maps"
	"slices"
)

// RangeOverInt shows range over an integer (Go 1.22): 0..n-1, typed like n,
// zero iterations for n <= 0, and "for range n" when the index isn't needed.
func RangeOverInt(w io.Writer) {
	for i := range 3 {
		fmt.Fprint(w, i, " ")
	}
	fmt.Fprintln(w, "| range 3")

	var n uint8 = 2
	for i := range n {
		fmt.Fprintf(w, "%d:%T ", i, i)
	}
	fmt.Fprintln(w, "| i has n's type")

	count := 0
	for range -5 {
		count++
	}
	fmt.Fprintln(w, count, "| range -5 runs zero times")

	for range 2 {
		fmt.Fprint(w, "hi ")
	}
	fmt.Fprintln(w, "| no variable at all")
}

// RangeOverSlice shows the index/value pairs of a slice, and the forms that
// drop one of them.
func RangeOverSlice(w io.Writer) {
	langs := []string{"go", "python", "java"}
	for i, s := range langs {
		fmt.Fprintf(w, "%d=%s ", i, s)
	}
	fmt.Fprintln(w, "| index and value")
	for i := range langs {
		fmt.Fprint(w, i, " ")
	}
	fmt.Fprintln(w, "| index only")
	for _, s := range langs {
		fmt.Fprint(w, len(s), " ")
	}
	fmt.Fprintln(w, "| value only")
}

// RangeCopiesValues shows that the iteration value is a copy of the element:
// changing it doesn't change the slice. Write through the index instead.
func RangeCopiesValues(w io.Writer) {
	type item struct {
		name  string
		price int
	}
	items := []item{{"pen", 100}, {"book", 900}}

	for _, it := range items {
		it.price *= 2 // changes the copy only
	}
	fmt.Fprintln(w, items, "| after doubling the copy")

	for i := range items {
		items[i].price *= 2 // changes the element
	}
	fmt.Fprintln(w, items, "| after doubling through the index")
}

// RangeEvaluatedOnce shows that the range expression is evaluated once,
// before the loop starts: a slice's length is fixed then, and an array is
// copied.
func RangeEvaluatedOnce(w io.Writer) {
	s := []int{1, 2, 3}
	for _, v := range s {
		s = append(s, v*10) // the loop still runs 3 times
	}
	fmt.Fprintln(w, s, "| append during range")

	arr := [3]int{1, 2, 3}
	for i, v := range arr { // ranges over a copy of arr
		if i == 0 {
			arr[2] = 100
		}
		if i == 2 {
			fmt.Fprintln(w, v, "| array: the copy still has 3")
		}
	}

	sl := []int{1, 2, 3}
	for i, v := range sl { // copies the slice header, shares the elements
		if i == 0 {
			sl[2] = 100
		}
		if i == 2 {
			fmt.Fprintln(w, v, "| slice: the shared element is 100")
		}
	}
}

// RangeOverString shows that ranging over a string decodes UTF-8: the index
// is the byte offset where each rune starts, and invalid bytes come out as
// U+FFFD (the replacement character) one byte at a time.
func RangeOverString(w io.Writer) {
	s := "héllo"
	fmt.Fprintln(w, "len:", len(s))
	for i, r := range s {
		fmt.Fprintf(w, "%d:%c ", i, r)
	}
	fmt.Fprintln(w, "| byte offset : rune")

	for i := 0; i < len(s); i++ {
		fmt.Fprintf(w, "%02x ", s[i])
	}
	fmt.Fprintln(w, "| indexing gives bytes")

	bad := "a\xffb"
	for i, r := range bad {
		fmt.Fprintf(w, "%d:%U ", i, r)
	}
	fmt.Fprintln(w, "| invalid UTF-8")
}

// RangeOverMap shows that map iteration order is unspecified, and the
// idiomatic way to get a stable order: sort the keys first.
func RangeOverMap(w io.Writer) {
	stock := map[string]int{"pear": 3, "apple": 5, "fig": 0}
	for _, k := range slices.Sorted(maps.Keys(stock)) {
		fmt.Fprintf(w, "%s=%d ", k, stock[k])
	}
	fmt.Fprintln(w, "| sorted keys")
}

// MapDeleteDuringRange shows that deleting entries while ranging over a map
// is safe: a deleted entry that hasn't been reached is never produced.
func MapDeleteDuringRange(w io.Writer) {
	m := map[int]bool{}
	for i := range 10 {
		m[i] = true
	}
	for k := range m {
		if k%2 == 1 {
			delete(m, k)
		}
	}
	fmt.Fprintln(w, len(m), slices.Sorted(maps.Keys(m)))
}

// RangeOverChannel shows a range loop receiving values until the channel is
// closed. Channels are topic 23; here the channel is buffered and filled
// before the loop, so there's no goroutine.
func RangeOverChannel(w io.Writer) {
	ch := make(chan string, 3)
	ch <- "a"
	ch <- "b"
	ch <- "c"
	close(ch) // without close, the loop would wait forever after "c"
	for v := range ch {
		fmt.Fprint(w, v, " ")
	}
	fmt.Fprintln(w, "| received until closed")
}

// RangeOverFunc shows range over an iterator function (Go 1.23). A
// break in the loop body makes yield return false. Iterators are topic 15.
func RangeOverFunc(w io.Writer) {
	for n := range countdown(5) {
		if n == 2 {
			break
		}
		fmt.Fprint(w, n, " ")
	}
	fmt.Fprintln(w, "| break stops the iterator")
}

func countdown(from int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := from; i > 0; i-- {
			if !yield(i) {
				return
			}
		}
	}
}
