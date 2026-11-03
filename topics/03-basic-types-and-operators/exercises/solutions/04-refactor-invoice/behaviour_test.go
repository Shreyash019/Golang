package main

import (
	"strings"
	"testing"
)

// These tests pin the program's behaviour. They pass before the refactor
// and must still pass after it.
func TestRun(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "discounted order",
			args: []string{"Keyboard=49.99x2", "Mouse=19.95", "Cable=4.50x3"},
			want: `items:    3 lines, 6 units
subtotal: $133.43
discount: -$13.34
tax (8%): $9.61
total:    $129.70
boxes:    2 (4 units per box)
`,
		},
		{
			name: "small order",
			args: []string{"Pen=1.25x3"},
			want: `items:    1 lines, 3 units
subtotal: $3.75
discount: -$0.00
tax (8%): $0.30
total:    $4.05
boxes:    1 (4 units per box)
`,
		},
		{
			name: "over the threshold",
			args: []string{"Desk=100", "Lamp=0.5x8"},
			want: `items:    2 lines, 9 units
subtotal: $104.00
discount: -$10.40
tax (8%): $7.49
total:    $101.09
boxes:    3 (4 units per box)
`,
		},
		{
			name: "exactly 100",
			args: []string{"Desk=100.00"},
			want: `items:    1 lines, 1 units
subtotal: $100.00
discount: -$0.00
tax (8%): $8.00
total:    $108.00
boxes:    1 (4 units per box)
`,
		},
		{
			name: "boxes divide exactly",
			args: []string{"Mug=7.10x8"},
			want: `items:    1 lines, 8 units
subtotal: $56.80
discount: -$0.00
tax (8%): $4.54
total:    $61.34
boxes:    2 (4 units per box)
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sb strings.Builder
			if err := run(&sb, tt.args); err != nil {
				t.Fatalf("run(%q) returned error %v", tt.args, err)
			}
			if got := sb.String(); got != tt.want {
				t.Errorf("run(%q) printed:\n%s\nwant:\n%s", tt.args, got, tt.want)
			}
		})
	}
}

func TestRunErrors(t *testing.T) {
	for _, args := range [][]string{
		nil,
		{"Pen"},
		{"=1.00"},
		{"Pen=abc"},
		{"Pen=-1"},
		{"Pen=1.00x0"},
		{"Pen=1.00x1.5"},
		{"Pen=1.00xabc"},
	} {
		var sb strings.Builder
		if err := run(&sb, args); err == nil {
			t.Errorf("run(%q) = nil; want an error", args)
		}
	}
}
