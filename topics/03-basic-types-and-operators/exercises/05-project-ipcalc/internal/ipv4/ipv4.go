// Package ipv4 does IPv4 subnet arithmetic on plain uint32 addresses: the
// kind of math behind firewall rules, VPC subnets and "is this client in
// our office network?" checks.
//
// Exercise E5 — mini-project. Implement every function until ipv4_test.go
// and ../../cmd/ipcalc/main_test.go pass. It trains unsigned integers and
// their ranges (notes §5.1, §5.2), shifts and bit masks (§6.1), conversions
// between integer sizes (§5.5) and Go's defined behaviour for shifts of 32
// or more bits (§5.7). Do the math yourself with operators; don't use
// net or net/netip (in real code you would: see the notes, §10).
//
//	make ex E=05-project-ipcalc
//	go run ./exercises/05-project-ipcalc/cmd/ipcalc 192.168.1.10/24
//
// An address is a uint32 whose most significant byte is the first number
// of the dotted form: 192.168.1.10 is 0xC0A8010A.
package ipv4

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
	return 0, nil
}

// FormatAddr returns the dotted-decimal form of a: FormatAddr(0xC0A8010A)
// is "192.168.1.10".
func FormatAddr(a uint32) string {
	return ""
}

// ParsePrefix parses "ADDR/BITS", such as "10.0.0.0/8", with BITS from 0 to
// 32. A missing "/BITS" is an error.
func ParsePrefix(s string) (Prefix, error) {
	return Prefix{}, nil
}

// Mask returns the netmask for a prefix length: Mask(24) is 0xFFFFFF00
// (255.255.255.0), Mask(0) is 0 and Mask(32) is 0xFFFFFFFF.
func Mask(bits int) uint32 {
	return 0
}

// Network returns p's network address: its address with the host bits
// cleared.
func Network(p Prefix) uint32 {
	return 0
}

// Broadcast returns p's broadcast address: its address with the host bits
// set.
func Broadcast(p Prefix) uint32 {
	return 0
}

// Contains reports whether addr is inside p's network.
func Contains(p Prefix, addr uint32) bool {
	return false
}

// Hosts returns the number of usable host addresses in p and the first and
// last of them. Normally the network and broadcast addresses are not
// usable, so a /24 has 254 hosts, from .1 to .254. Two lengths are special
// (RFC 3021 and single hosts): a /31 has 2 usable hosts (both addresses)
// and a /32 has 1 (the address itself). A /0 has 4294967294 hosts, which
// is why the count is a uint64.
func Hosts(p Prefix) (count uint64, first, last uint32) {
	return 0, 0, 0
}
