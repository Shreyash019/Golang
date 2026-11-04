package funcs

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"
)

func TestDivide(t *testing.T) {
	t.Parallel()
	tests := []struct {
		a, b    int
		q, r    int
		wantErr error
	}{
		{17, 5, 3, 2, nil},
		{-7, 2, -3, -1, nil},
		{7, -2, -3, 1, nil},
		{0, 3, 0, 0, nil},
		{5, 0, 0, 0, ErrDivideByZero},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d/%d", tt.a, tt.b), func(t *testing.T) {
			t.Parallel()
			q, r, err := Divide(tt.a, tt.b)
			if q != tt.q || r != tt.r || !errors.Is(err, tt.wantErr) {
				t.Errorf("Divide(%d, %d) = %d, %d, %v, want %d, %d, %v", tt.a, tt.b, q, r, err, tt.q, tt.r, tt.wantErr)
			}
		})
	}
}

func TestSum(t *testing.T) {
	t.Parallel()
	tests := []struct {
		nums []int
		want int
	}{
		{nil, 0},
		{[]int{5}, 5},
		{[]int{1, 2, 3, 4}, 10},
		{[]int{-3, 3}, 0},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprint(tt.nums), func(t *testing.T) {
			t.Parallel()
			if got := Sum(tt.nums...); got != tt.want {
				t.Errorf("Sum(%v...) = %d, want %d", tt.nums, got, tt.want)
			}
		})
	}
	if got := Sum(1, 2); got != 3 {
		t.Errorf("Sum(1, 2) = %d, want 3", got)
	}
}

func TestPipeline(t *testing.T) {
	t.Parallel()
	inc := func(n int) int { return n + 1 }
	double := func(n int) int { return n * 2 }
	square := func(n int) int { return n * n }
	tests := []struct {
		name string
		fs   []func(int) int
		in   int
		want int
	}{
		{"none is identity", nil, 7, 7},
		{"one", []func(int) int{inc}, 7, 8},
		{"left to right", []func(int) int{inc, double}, 3, 8},
		{"order matters", []func(int) int{double, inc}, 3, 7},
		{"three", []func(int) int{inc, double, square}, 1, 16},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			p := Pipeline(tt.fs...)
			if p == nil {
				t.Fatal("Pipeline returned a nil function")
			}
			if got := p(tt.in); got != tt.want {
				t.Errorf("Pipeline(%s)(%d) = %d, want %d", tt.name, tt.in, got, tt.want)
			}
		})
	}
}

func TestCounter(t *testing.T) {
	t.Parallel()
	next, reset := Counter(10, 5)
	if next == nil || reset == nil {
		t.Fatal("Counter returned a nil function")
	}
	other, _ := Counter(0, 1)

	var got []int
	for range 3 {
		got = append(got, next())
	}
	got = append(got, other())
	reset()
	got = append(got, next(), next(), other())

	want := []int{10, 15, 20, 0, 10, 15, 1}
	if !slices.Equal(got, want) {
		t.Errorf("next, next, next, other, reset, next, next, other = %v, want %v", got, want)
	}
}

func TestRetry(t *testing.T) {
	t.Parallel()
	errTemporary := errors.New("temporary failure")
	tests := []struct {
		name      string
		attempts  int
		succeedOn int // the attempt that returns nil; 0 = never
		wantCalls []int
		wantErr   string // "" = nil
	}{
		{"first try", 3, 1, []int{1}, ""},
		{"third try", 5, 3, []int{1, 2, 3}, ""},
		{"never", 3, 0, []int{1, 2, 3}, "after 3 attempts: temporary failure"},
		{"zero attempts means one", 0, 0, []int{1}, "after 1 attempts: temporary failure"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var calls []int
			err := Retry(tt.attempts, func(attempt int) error {
				calls = append(calls, attempt)
				if attempt == tt.succeedOn {
					return nil
				}
				return errTemporary
			})
			if !slices.Equal(calls, tt.wantCalls) {
				t.Errorf("Retry(%d, fn) called fn with attempts %v, want %v", tt.attempts, calls, tt.wantCalls)
			}
			switch {
			case tt.wantErr == "" && err != nil:
				t.Errorf("Retry(%d, fn) = %v, want nil", tt.attempts, err)
			case tt.wantErr != "" && (err == nil || err.Error() != tt.wantErr):
				t.Errorf("Retry(%d, fn) = %v, want %q", tt.attempts, err, tt.wantErr)
			case tt.wantErr != "" && !errors.Is(err, errTemporary):
				t.Errorf("Retry(%d, fn) = %v, which doesn't wrap the last error", tt.attempts, err)
			}
		})
	}
}

func TestWalk(t *testing.T) {
	t.Parallel()
	tree := &Node{Name: "root", Children: []*Node{
		{Name: "a", Children: []*Node{{Name: "x"}, {Name: "y"}}},
		{Name: "skip", Children: []*Node{{Name: "hidden"}}},
		{Name: "b"},
	}}
	tests := []struct {
		name string
		root *Node
		want []string
	}{
		{"nil tree", nil, nil},
		{"single node", &Node{Name: "solo"}, []string{"0 solo"}},
		{"whole tree, skipping a subtree", tree, []string{
			"0 root", "1 root/a", "2 root/a/x", "2 root/a/y", "1 root/skip", "1 root/b",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var got []string
			Walk(tt.root, func(path string, depth int) bool {
				got = append(got, fmt.Sprint(depth, " ", path))
				return !strings.HasSuffix(path, "skip")
			})
			if !slices.Equal(got, tt.want) {
				t.Errorf("Walk(%s) visited\n%q\nwant\n%q", tt.name, got, tt.want)
			}
		})
	}
}
