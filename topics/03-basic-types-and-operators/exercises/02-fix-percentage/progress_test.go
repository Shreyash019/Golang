package progress

import "testing"

func TestPercent(t *testing.T) {
	t.Parallel()
	tests := []struct {
		done, total, want int
	}{
		{0, 10, 0},
		{10, 10, 100},
		{0, 0, 100},
		{1, 4, 25},
		{1, 3, 33},
		{2, 3, 66},
		{5, 200, 2},
		{999, 1000, 99},
		{3_000_000_000, 4_000_000_000, 75},
	}
	for _, tt := range tests {
		if got := Percent(tt.done, tt.total); got != tt.want {
			t.Errorf("Percent(%d, %d) = %d, want %d", tt.done, tt.total, got, tt.want)
		}
	}
}
