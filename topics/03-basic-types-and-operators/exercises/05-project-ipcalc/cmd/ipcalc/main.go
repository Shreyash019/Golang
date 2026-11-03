// Command ipcalc prints the details of an IPv4 subnet and, optionally,
// whether an address belongs to it.
//
// Usage:
//
//	ipcalc ADDR/BITS [ADDR]
//	ipcalc 192.168.1.10/24
//	ipcalc 10.0.0.0/8 10.20.30.40
//
// Exercise E5 — mini-project. main is complete: it only wires things
// together. Implement internal/ipv4 until all tests pass:
//
//	make ex E=05-project-ipcalc
//	go run ./exercises/05-project-ipcalc/cmd/ipcalc 192.168.1.10/24
package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/Shreyash019/Golang/topics/03-basic-types-and-operators/exercises/05-project-ipcalc/internal/ipv4"
)

func main() {
	if err := run(os.Stdout, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "ipcalc:", err)
		os.Exit(2)
	}
}

func run(w io.Writer, args []string) error {
	if len(args) < 1 || len(args) > 2 {
		return errors.New("usage: ipcalc ADDR/BITS [ADDR]")
	}
	p, err := ipv4.ParsePrefix(args[0])
	if err != nil {
		return err
	}
	var check uint32
	if len(args) == 2 {
		if check, err = ipv4.ParseAddr(args[1]); err != nil {
			return err
		}
	}

	count, first, last := ipv4.Hosts(p)
	fmt.Fprintf(w, "address:   %s\n", ipv4.FormatAddr(p.Addr))
	fmt.Fprintf(w, "network:   %s/%d\n", ipv4.FormatAddr(ipv4.Network(p)), p.Bits)
	fmt.Fprintf(w, "netmask:   %s\n", ipv4.FormatAddr(ipv4.Mask(p.Bits)))
	fmt.Fprintf(w, "broadcast: %s\n", ipv4.FormatAddr(ipv4.Broadcast(p)))
	fmt.Fprintf(w, "hosts:     %d (%s - %s)\n", count, ipv4.FormatAddr(first), ipv4.FormatAddr(last))
	if len(args) == 2 {
		verdict := "is not in"
		if ipv4.Contains(p, check) {
			verdict = "is in"
		}
		fmt.Fprintf(w, "%s %s %s/%d\n", args[1], verdict, ipv4.FormatAddr(ipv4.Network(p)), p.Bits)
	}
	return nil
}
