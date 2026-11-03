package checked

import (
	"math"
	"testing"
)

const (
	maxI = math.MaxInt64
	minI = math.MinInt64
)

type binCase struct {
	a, b   int64
	want   int64
	wantOK bool
}

func testBinary(t *testing.T, name string, fn func(a, b int64) (int64, bool), tests []binCase) {
	t.Helper()
	for _, tt := range tests {
		got, ok := fn(tt.a, tt.b)
		if got != tt.want || ok != tt.wantOK {
			t.Errorf("%s(%d, %d) = %d, %t; want %d, %t", name, tt.a, tt.b, got, ok, tt.want, tt.wantOK)
		}
	}
}

func TestAdd(t *testing.T) {
	t.Parallel()
	testBinary(t, "Add", Add, []binCase{
		{1, 2, 3, true},
		{-5, 3, -2, true},
		{0, 0, 0, true},
		{maxI, 0, maxI, true},
		{maxI - 1, 1, maxI, true},
		{maxI, 1, 0, false},
		{minI, -1, 0, false},
		{minI, maxI, -1, true},
		{maxI, maxI, 0, false},
		{minI, minI, 0, false},
	})
}

func TestSub(t *testing.T) {
	t.Parallel()
	testBinary(t, "Sub", Sub, []binCase{
		{5, 3, 2, true},
		{3, 5, -2, true},
		{0, 0, 0, true},
		{minI, 1, 0, false},
		{maxI, -1, 0, false},
		{0, minI, 0, false},
		{-1, minI, maxI, true},
		{minI, minI, 0, true},
	})
}

func TestMul(t *testing.T) {
	t.Parallel()
	testBinary(t, "Mul", Mul, []binCase{
		{6, 7, 42, true},
		{-6, 7, -42, true},
		{0, maxI, 0, true},
		{minI, 0, 0, true},
		{1 << 31, 1 << 31, 1 << 62, true},
		{1 << 32, 1 << 31, 0, false},
		{maxI, 2, 0, false},
		{minI, 1, minI, true},
		{minI, -1, 0, false},
		{-1, minI, 0, false},
		{maxI, -1, -maxI, true},
		{3_037_000_499, 3_037_000_499, 9_223_372_030_926_249_001, true},
		{3_037_000_500, 3_037_000_500, 0, false},
	})
}

func TestAbs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in     int64
		want   int64
		wantOK bool
	}{
		{5, 5, true},
		{-5, 5, true},
		{0, 0, true},
		{maxI, maxI, true},
		{-maxI, maxI, true},
		{minI, 0, false},
	}
	for _, tt := range tests {
		if got, ok := Abs(tt.in); got != tt.want || ok != tt.wantOK {
			t.Errorf("Abs(%d) = %d, %t; want %d, %t", tt.in, got, ok, tt.want, tt.wantOK)
		}
	}
}

func TestToInt32(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in     int64
		want   int32
		wantOK bool
	}{
		{0, 0, true},
		{-42, -42, true},
		{math.MaxInt32, math.MaxInt32, true},
		{math.MinInt32, math.MinInt32, true},
		{math.MaxInt32 + 1, 0, false},
		{math.MinInt32 - 1, 0, false},
		{1 << 40, 0, false},
	}
	for _, tt := range tests {
		if got, ok := ToInt32(tt.in); got != tt.want || ok != tt.wantOK {
			t.Errorf("ToInt32(%d) = %d, %t; want %d, %t", tt.in, got, ok, tt.want, tt.wantOK)
		}
	}
}

func TestRoundToInt64(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in     float64
		want   int64
		wantOK bool
	}{
		{0, 0, true},
		{2.4, 2, true},
		{2.5, 3, true},
		{-2.5, -3, true},
		{-2.7, -3, true},
		{1e15 + 0.5, 1_000_000_000_000_001, true},
		{-9.223372036854775808e18, minI, true},
		{9.223372036854775808e18, 0, false},
		{1e19, 0, false},
		{-1e19, 0, false},
		{math.Inf(1), 0, false},
		{math.Inf(-1), 0, false},
		{math.NaN(), 0, false},
	}
	for _, tt := range tests {
		if got, ok := RoundToInt64(tt.in); got != tt.want || ok != tt.wantOK {
			t.Errorf("RoundToInt64(%v) = %d, %t; want %d, %t", tt.in, got, ok, tt.want, tt.wantOK)
		}
	}
}
