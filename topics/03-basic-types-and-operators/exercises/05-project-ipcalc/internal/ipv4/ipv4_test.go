package ipv4

import "testing"

func TestParseAddr(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in   string
		want uint32
	}{
		{"0.0.0.0", 0},
		{"192.168.1.10", 0xC0A8010A},
		{"10.0.0.1", 0x0A000001},
		{"255.255.255.255", 0xFFFFFFFF},
		{"127.0.0.1", 0x7F000001},
	}
	for _, tt := range tests {
		got, err := ParseAddr(tt.in)
		if err != nil || got != tt.want {
			t.Errorf("ParseAddr(%q) = %#08x, %v; want %#08x, nil", tt.in, got, err, tt.want)
		}
	}
	for _, in := range []string{"", "1.2.3", "1.2.3.4.5", "256.0.0.1", "1.2.3.-4", "a.b.c.d", "1..2.3", "01.2.3.4", "1.2.3.4 ", "+1.2.3.4"} {
		if got, err := ParseAddr(in); err == nil {
			t.Errorf("ParseAddr(%q) = %#08x, nil; want an error", in, got)
		}
	}
}

func TestFormatAddr(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in   uint32
		want string
	}{
		{0, "0.0.0.0"},
		{0xC0A8010A, "192.168.1.10"},
		{0xFFFFFFFF, "255.255.255.255"},
		{0x0A000001, "10.0.0.1"},
	}
	for _, tt := range tests {
		if got := FormatAddr(tt.in); got != tt.want {
			t.Errorf("FormatAddr(%#08x) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestParsePrefix(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in   string
		want Prefix
	}{
		{"192.168.1.10/24", Prefix{0xC0A8010A, 24}},
		{"0.0.0.0/0", Prefix{0, 0}},
		{"10.1.2.3/32", Prefix{0x0A010203, 32}},
	}
	for _, tt := range tests {
		got, err := ParsePrefix(tt.in)
		if err != nil || got != tt.want {
			t.Errorf("ParsePrefix(%q) = %+v, %v; want %+v, nil", tt.in, got, err, tt.want)
		}
	}
	for _, in := range []string{"", "10.0.0.0", "10.0.0.0/", "10.0.0.0/33", "10.0.0.0/-1", "10.0.0.0/x", "10.0.0/8", "/8"} {
		if got, err := ParsePrefix(in); err == nil {
			t.Errorf("ParsePrefix(%q) = %+v, nil; want an error", in, got)
		}
	}
}

func TestMask(t *testing.T) {
	t.Parallel()
	tests := []struct {
		bits int
		want uint32
	}{
		{0, 0},
		{1, 0x80000000},
		{8, 0xFF000000},
		{24, 0xFFFFFF00},
		{31, 0xFFFFFFFE},
		{32, 0xFFFFFFFF},
	}
	for _, tt := range tests {
		if got := Mask(tt.bits); got != tt.want {
			t.Errorf("Mask(%d) = %#08x, want %#08x", tt.bits, got, tt.want)
		}
	}
}

func TestNetworkBroadcastHosts(t *testing.T) {
	t.Parallel()
	tests := []struct {
		p                  Prefix
		network, broadcast uint32
		count              uint64
		first, last        uint32
	}{
		{Prefix{0xC0A8010A, 24}, 0xC0A80100, 0xC0A801FF, 254, 0xC0A80101, 0xC0A801FE},
		{Prefix{0x0A123456, 8}, 0x0A000000, 0x0AFFFFFF, 16777214, 0x0A000001, 0x0AFFFFFE},
		{Prefix{0xC0A80107, 30}, 0xC0A80104, 0xC0A80107, 2, 0xC0A80105, 0xC0A80106},
		{Prefix{0xC0A80107, 31}, 0xC0A80106, 0xC0A80107, 2, 0xC0A80106, 0xC0A80107},
		{Prefix{0xC0A80107, 32}, 0xC0A80107, 0xC0A80107, 1, 0xC0A80107, 0xC0A80107},
		{Prefix{0x01020304, 0}, 0, 0xFFFFFFFF, 4294967294, 1, 0xFFFFFFFE},
	}
	for _, tt := range tests {
		if got := Network(tt.p); got != tt.network {
			t.Errorf("Network(%+v) = %#08x, want %#08x", tt.p, got, tt.network)
		}
		if got := Broadcast(tt.p); got != tt.broadcast {
			t.Errorf("Broadcast(%+v) = %#08x, want %#08x", tt.p, got, tt.broadcast)
		}
		count, first, last := Hosts(tt.p)
		if count != tt.count || first != tt.first || last != tt.last {
			t.Errorf("Hosts(%+v) = %d, %#08x, %#08x; want %d, %#08x, %#08x", tt.p, count, first, last, tt.count, tt.first, tt.last)
		}
	}
}

func TestContains(t *testing.T) {
	t.Parallel()
	tests := []struct {
		p    Prefix
		addr uint32
		want bool
	}{
		{Prefix{0x0A000000, 8}, 0x0A010203, true},
		{Prefix{0x0A000000, 8}, 0x0B000000, false},
		{Prefix{0xC0A80100, 24}, 0xC0A801FF, true},
		{Prefix{0xC0A80100, 24}, 0xC0A80200, false},
		{Prefix{0, 0}, 0xDEADBEEF, true},
		{Prefix{0xC0A80107, 32}, 0xC0A80107, true},
		{Prefix{0xC0A80107, 32}, 0xC0A80106, false},
	}
	for _, tt := range tests {
		if got := Contains(tt.p, tt.addr); got != tt.want {
			t.Errorf("Contains(%+v, %#08x) = %t, want %t", tt.p, tt.addr, got, tt.want)
		}
	}
}
