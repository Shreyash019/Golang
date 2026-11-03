// Package ipv4 does IPv4 subnet arithmetic on plain uint32 addresses: the
// kind of math behind firewall rules, VPC subnets and "is this client in
// our office network?" checks.
//
// Solution E5. Production code would use net/netip (netip.ParsePrefix,
// Prefix.Contains, Prefix.Masked), which also handles IPv6; this package
// shows the arithmetic underneath.
//
// An address is a uint32 whose most significant byte is the first number
// of the dotted form: 192.168.1.10 is 0xC0A8010A.
package ipv4

import (
	"fmt"
	"strconv"
	"strings"
)

// Prefix is an address with a prefix length, as in 192.168.1.10/24.
// Addr keeps the host bits as written; Bits is from 0 to 32.
type Prefix struct {
	Addr uint32
	Bits int
}

// ParseAddr parses a dotted-decimal address such as "192.168.1.10": exactly
// four decimal numbers from 0 to 255, separated by dots. Leading zeros
// ("01") and anything else are errors.
func ParseAddr(s string) (uint32, error) {
	parts := strings.Split(s, ".")
	if len(parts) != 4 {
		return 0, fmt.Errorf("address %q: want four numbers separated by dots", s)
	}
	var addr uint32
	for _, part := range parts {
		// ParseUint with bitSize 8 rejects values over 255 for us, and
		// rejects signs and spaces. Leading zeros are a separate check,
		// because some tools read "010" as octal.
		if len(part) > 1 && part[0] == '0' {
			return 0, fmt.Errorf("address %q: leading zero in %q", s, part)
		}
		n, err := strconv.ParseUint(part, 10, 8)
		if err != nil {
			return 0, fmt.Errorf("address %q: %q is not a number from 0 to 255", s, part)
		}
		addr = addr<<8 | uint32(n)
	}
	return addr, nil
}

// FormatAddr returns the dotted-decimal form of a: FormatAddr(0xC0A8010A)
// is "192.168.1.10".
func FormatAddr(a uint32) string {
	// uint8(x) keeps the low 8 bits, which is exactly one octet.
	return fmt.Sprintf("%d.%d.%d.%d", uint8(a>>24), uint8(a>>16), uint8(a>>8), uint8(a))
}

// ParsePrefix parses "ADDR/BITS", such as "10.0.0.0/8", with BITS from 0 to
// 32. A missing "/BITS" is an error.
func ParsePrefix(s string) (Prefix, error) {
	addrText, bitsText, ok := strings.Cut(s, "/")
	if !ok {
		return Prefix{}, fmt.Errorf("prefix %q: want ADDR/BITS", s)
	}
	addr, err := ParseAddr(addrText)
	if err != nil {
		return Prefix{}, fmt.Errorf("prefix %q: %w", s, err)
	}
	bits, err := strconv.Atoi(bitsText)
	if err != nil || bits < 0 || bits > 32 {
		return Prefix{}, fmt.Errorf("prefix %q: length %q is not a number from 0 to 32", s, bitsText)
	}
	return Prefix{Addr: addr, Bits: bits}, nil
}

// Mask returns the netmask for a prefix length: Mask(24) is 0xFFFFFF00
// (255.255.255.0), Mask(0) is 0 and Mask(32) is 0xFFFFFFFF.
func Mask(bits int) uint32 {
	// For bits == 0 this shifts a uint32 by 32. In C that is undefined
	// behaviour; in Go, shifting by the width or more is defined: every
	// bit is shifted out and the result is 0, exactly the /0 mask.
	return ^uint32(0) << (32 - bits)
}

// Network returns p's network address: its address with the host bits
// cleared.
func Network(p Prefix) uint32 {
	return p.Addr & Mask(p.Bits)
}

// Broadcast returns p's broadcast address: its address with the host bits
// set.
func Broadcast(p Prefix) uint32 {
	return p.Addr | ^Mask(p.Bits)
}

// Contains reports whether addr is inside p's network.
func Contains(p Prefix, addr uint32) bool {
	m := Mask(p.Bits)
	return addr&m == p.Addr&m // & binds tighter than ==: no parentheses needed
}

// Hosts returns the number of usable host addresses in p and the first and
// last of them. Normally the network and broadcast addresses are not
// usable, so a /24 has 254 hosts, from .1 to .254. Two lengths are special
// (RFC 3021 and single hosts): a /31 has 2 usable hosts (both addresses)
// and a /32 has 1 (the address itself). A /0 has 4294967294 hosts, which
// is why the count is a uint64.
func Hosts(p Prefix) (count uint64, first, last uint32) {
	network, broadcast := Network(p), Broadcast(p)
	switch p.Bits {
	case 32:
		return 1, network, network
	case 31:
		return 2, network, broadcast
	}
	// 1 << (32 - Bits) addresses, computed in uint64 because a /0 has 2^32,
	// one more than a uint32 can hold.
	total := uint64(1) << (32 - p.Bits)
	return total - 2, network + 1, broadcast - 1
}
