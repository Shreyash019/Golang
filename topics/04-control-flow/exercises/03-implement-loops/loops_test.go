package loops

import (
	"fmt"
	"slices"
	"testing"
)

func TestStatusClass(t *testing.T) {
	t.Parallel()
	tests := []struct {
		code int
		want string
	}{
		{100, "informational"},
		{199, "informational"},
		{200, "success"},
		{204, "success"},
		{301, "redirect"},
		{399, "redirect"},
		{404, "client error"},
		{499, "client error"},
		{500, "server error"},
		{599, "server error"},
		{99, "invalid"},
		{600, "invalid"},
		{0, "invalid"},
		{-200, "invalid"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprint(tt.code), func(t *testing.T) {
			t.Parallel()
			if got := StatusClass(tt.code); got != tt.want {
				t.Errorf("StatusClass(%d) = %q, want %q", tt.code, got, tt.want)
			}
		})
	}
}

func TestCollatzSteps(t *testing.T) {
	t.Parallel()
	tests := []struct {
		n         int
		wantSteps int
		wantOK    bool
	}{
		{1, 0, true},
		{2, 1, true},
		{6, 8, true},
		{7, 16, true},
		{27, 111, true},
		{0, 0, false},
		{-5, 0, false},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprint(tt.n), func(t *testing.T) {
			t.Parallel()
			steps, ok := CollatzSteps(tt.n)
			if steps != tt.wantSteps || ok != tt.wantOK {
				t.Errorf("CollatzSteps(%d) = %d, %t, want %d, %t", tt.n, steps, ok, tt.wantSteps, tt.wantOK)
			}
		})
	}
}

func TestPrimes(t *testing.T) {
	t.Parallel()
	tests := []struct {
		limit int
		want  []int
	}{
		{-3, nil},
		{1, nil},
		{2, []int{2}},
		{10, []int{2, 3, 5, 7}},
		{30, []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29}},
		{49, []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47}},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprint(tt.limit), func(t *testing.T) {
			t.Parallel()
			if got := Primes(tt.limit); !slices.Equal(got, tt.want) {
				t.Errorf("Primes(%d) = %v, want %v", tt.limit, got, tt.want)
			}
		})
	}
}

func TestFind(t *testing.T) {
	t.Parallel()
	grid := []string{
		"héllo",
		"wörld",
		"go!",
	}
	tests := []struct {
		name    string
		target  rune
		wantRow int
		wantCol int
		wantOK  bool
	}{
		{"first rune", 'h', 0, 0, true},
		{"after a 2-byte rune", 'l', 0, 2, true},
		{"multi-byte target", 'ö', 1, 1, true},
		{"first occurrence wins", 'o', 0, 4, true},
		{"last row", '!', 2, 2, true},
		{"missing", 'z', 0, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			row, col, ok := Find(grid, tt.target)
			if row != tt.wantRow || col != tt.wantCol || ok != tt.wantOK {
				t.Errorf("Find(grid, %q) = %d, %d, %t, want %d, %d, %t", tt.target, row, col, ok, tt.wantRow, tt.wantCol, tt.wantOK)
			}
		})
	}
	t.Run("empty grid", func(t *testing.T) {
		t.Parallel()
		if row, col, ok := Find(nil, 'a'); row != 0 || col != 0 || ok {
			t.Errorf("Find(nil, 'a') = %d, %d, %t, want 0, 0, false", row, col, ok)
		}
	})
}

func TestCompress(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in, want string
	}{
		{"", ""},
		{"a", "a1"},
		{"aaabcc", "a3b1c2"},
		{"abc", "a1b1c1"},
		{"aaaaaaaaaaaa", "a12"},
		{"ééé!", "é3!1"},
		{"aabaa", "a2b1a2"},
		{"日日本", "日2本1"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			t.Parallel()
			if got := Compress(tt.in); got != tt.want {
				t.Errorf("Compress(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
