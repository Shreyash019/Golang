package temperature

import "testing"

// A table-driven test: one table of cases, one loop, one subtest per case.
// Adding a case is one line, and a failure names the case that broke.
func TestCelsiusFahrenheit(t *testing.T) {
	tests := []struct {
		name string
		in   Celsius
		want Fahrenheit
	}{
		{name: "freezing", in: 0, want: 32},
		{name: "boiling", in: 100, want: 212},
		{name: "same on both scales", in: -40, want: -40},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.in.Fahrenheit()
			if got != tt.want {
				t.Errorf("Celsius(%v).Fahrenheit() = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}
