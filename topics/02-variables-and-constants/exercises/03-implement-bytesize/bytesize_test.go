package bytesize

import "testing"

func TestUnits(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		got, want int64
	}{
		{"KiB", KiB, 1 << 10},
		{"MiB", MiB, 1 << 20},
		{"GiB", GiB, 1 << 30},
		{"TiB", TiB, 1 << 40},
	}
	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s = %d, want %d", tt.name, tt.got, tt.want)
		}
	}
}

func TestFormat(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in   int64
		want string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1023, "1023 B"},
		{1024, "1 KiB"},
		{1536, "1.5 KiB"},
		{10 << 20, "10 MiB"},
		{3<<30 + 300<<20, "3.3 GiB"},
		{5 << 40, "5 TiB"},
		{2048 << 40, "2048 TiB"},
		{-2048, "-2 KiB"},
		{-100, "-100 B"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			t.Parallel()
			if got := Format(tt.in); got != tt.want {
				t.Errorf("Format(%d) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestParse(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in   string
		want int64
	}{
		{"0", 0},
		{"512", 512},
		{"512B", 512},
		{"512 B", 512},
		{"1KiB", 1024},
		{"1.5 KiB", 1536},
		{"  10 MiB  ", 10 << 20},
		{"2GiB", 2 << 30},
		{"1 TiB", 1 << 40},
		{"0.5 KiB", 512},
		{"1.5 B", 1},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			t.Parallel()
			got, err := Parse(tt.in)
			if err != nil {
				t.Fatalf("Parse(%q) returned error %v; want %d", tt.in, err, tt.want)
			}
			if got != tt.want {
				t.Errorf("Parse(%q) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}

func TestParseErrors(t *testing.T) {
	t.Parallel()
	for _, in := range []string{
		"",
		"   ",
		"KiB",
		"abc",
		"1.2.3 KiB",
		"-1 KiB",
		"10 KB",
		"10 kib",
		"10 KiB extra",
		"9000000 TiB",
	} {
		t.Run(in, func(t *testing.T) {
			t.Parallel()
			if got, err := Parse(in); err == nil {
				t.Errorf("Parse(%q) = %d, nil; want an error", in, got)
			}
		})
	}
}

// Format and Parse are inverses for every size Format prints exactly.
func TestRoundTrip(t *testing.T) {
	t.Parallel()
	for _, n := range []int64{0, 1, 1023, 1024, 1536, 7 << 20, 1 << 40} {
		s := Format(n)
		got, err := Parse(s)
		if err != nil || got != n {
			t.Errorf("Parse(Format(%d)) = Parse(%q) = %d, %v; want %d, nil", n, s, got, err, n)
		}
	}
}
