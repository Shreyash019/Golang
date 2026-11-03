package examples

import (
	"fmt"
	"io"
	"strconv"
)

// scopeLevel is declared in the package block, visible in every file of the
// package. Package-level mutable variables are a code smell; this one is
// never written, it exists only to be shadowed.
var scopeLevel = "package"

// ScopeLevels declares the same name in nested blocks. Each use resolves to
// the innermost declaration that encloses it.
func ScopeLevels(w io.Writer) {
	fmt.Fprintln(w, "1:", scopeLevel) // the package variable
	scopeLevel := "function"          // shadows it for the rest of the function
	fmt.Fprintln(w, "2:", scopeLevel)
	{
		scopeLevel := "explicit block"
		fmt.Fprintln(w, "3:", scopeLevel)
	}
	if scopeLevel := "if statement"; scopeLevel != "" {
		fmt.Fprintln(w, "4:", scopeLevel) // the if's init scope covers if and else
	}
	for i := 0; i < 1; i++ {
		scopeLevel := "loop body"
		fmt.Fprintln(w, "5:", scopeLevel, "i =", i)
	}
	fmt.Fprintln(w, "6:", scopeLevel) // back to the function's variable
}

// ScopeStartsAfterDeclaration shows that a local variable's scope begins
// after its declaration ends, so the right-hand side of x := x + 10 still
// refers to the outer x.
func ScopeStartsAfterDeclaration(w io.Writer) {
	x := 1
	{
		x := x + 10 // the right side reads the OUTER x
		fmt.Fprintln(w, "inner:", x)
	}
	fmt.Fprintln(w, "outer:", x)
}

// ShortRedeclaration shows that := reuses an existing variable when it was
// declared in the SAME scope (here err), but declares a new one in an inner
// scope. Comparing addresses proves which is which (pointers: topic 09).
func ShortRedeclaration(w io.Writer) {
	a, err := strconv.Atoi("1") //nolint:staticcheck // SA4006: err is overwritten on purpose
	first := &err
	b, err := strconv.Atoi("x") // b is new, err is only assigned
	fmt.Fprintln(w, "a b err:", a, b, err)
	fmt.Fprintln(w, "same err variable:", first == &err)
	{
		c, err := strconv.Atoi("2") // inner block: a NEW err
		fmt.Fprintln(w, "inner c err:", c, err, "same err variable:", first == &err)
	}
	fmt.Fprintln(w, "outer err after the block:", err)
}

// portFrom has the classic shadowing bug: inside the if, := declares new
// port and err variables, so the parsed value never reaches the outer port.
func portFrom(s string) (int, error) {
	port := 8080
	if s != "" {
		port, err := strconv.Atoi(s) // NEW port and err, scoped to the if block
		if err != nil {
			return 0, err
		}
		if port < 1 || port > 65535 {
			return 0, fmt.Errorf("port %d out of range", port)
		}
	}
	return port, nil // the OUTER port: still 8080
}

// portFromFixed assigns with = to the outer port.
func portFromFixed(s string) (int, error) {
	port := 8080
	if s != "" {
		var err error
		port, err = strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		if port < 1 || port > 65535 {
			return 0, fmt.Errorf("port %d out of range", port)
		}
	}
	return port, nil
}

// ShadowedErr runs the buggy and the fixed parser on the same input.
func ShadowedErr(w io.Writer) {
	port, err := portFrom("9090")
	fmt.Fprintf(w, "buggy: port=%d err=%v\n", port, err)
	port, err = portFromFixed("9090")
	fmt.Fprintf(w, "fixed: port=%d err=%v\n", port, err)
	port, err = portFrom("abc")
	fmt.Fprintf(w, "buggy, bad input: port=%d err=%v\n", port, err)
}

// currentEnv is meant to be set by setupEnv.
var currentEnv string

func detectEnv() (string, error) { return "production", nil }

// setupEnv means to set the package variable, but := declares a local one.
func setupEnv(w io.Writer) error {
	currentEnv, err := detectEnv()
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "inside setupEnv: %q\n", currentEnv)
	return nil
}

// ShadowedPackageVar shows the package variable is still empty afterwards.
func ShadowedPackageVar(w io.Writer) {
	if err := setupEnv(w); err != nil {
		fmt.Fprintln(w, "error:", err)
		return
	}
	fmt.Fprintf(w, "package currentEnv: %q\n", currentEnv)
}
