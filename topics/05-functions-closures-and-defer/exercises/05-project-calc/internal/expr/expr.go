// Package expr evaluates integer arithmetic expressions such as
// "2 * (3 + 4) - -1".
//
// The grammar, lowest precedence first (operators are left-associative):
//
//	expr   = term { ("+" | "-") term }
//	term   = factor { ("*" | "/" | "%") factor }
//	factor = "-" factor | number | "(" expr ")"
//
// Numbers are decimal int64 literals. Spaces are ignored. Arithmetic is
// int64 arithmetic: it wraps on overflow and / truncates toward zero.
package expr

import (
	"errors"
	"io"
)

// ErrDivisionByZero is wrapped by the error Eval returns when the right
// operand of / or % is zero.
var ErrDivisionByZero = errors.New("division by zero")

// Eval returns the value of the expression src.
//
// Errors (offsets count bytes from 0):
//
//	unexpected character '$' at offset 3
//	number "99999999999999999999" at offset 0 is out of range
//	unexpected ')' at offset 4      (a token that can't appear there;
//	                                 numbers appear as digits: unexpected 2 at offset 2)
//	unexpected end of input
//	missing ')' at offset 4         (offset of the token found instead)
//	offset 2: division by zero      (offset of the operator; wraps ErrDivisionByZero)
func Eval(src string) (int64, error) {
	return 0, nil
}

// Trace is Eval, but also writes a line to w when each grammar rule starts
// ("> term") and when it returns ("< term = 6", or "< term error"),
// indented two spaces per level of nesting.
//
// Suggested shape (notes §5.5, §5.6, §5.7): tokenize src first, then write
// the three rules as mutually recursive closures declared with
//
//	var expr, term, factor func() (int64, error)
//
// sharing the current position by capturing it. Give each rule named
// results, and trace it with a call on entry plus a deferred closure that
// reads those named results on the way out. Eval can then share the code
// by tracing to io.Discard.
func Trace(src string, w io.Writer) (int64, error) {
	return 0, nil
}
