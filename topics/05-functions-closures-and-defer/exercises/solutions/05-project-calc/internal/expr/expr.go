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
	"fmt"
	"io"
	"strconv"
	"strings"
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
	return evaluate(src, io.Discard)
}

// Trace is Eval, but also writes a line to w when each grammar rule starts
// ("> term") and when it returns ("< term = 6", or "< term error"),
// indented two spaces per level of nesting.
func Trace(src string, w io.Writer) (int64, error) {
	return evaluate(src, w)
}

type token struct {
	kind byte // '0' number, one of + - * / % ( ), or 0 at the end
	val  int64
	pos  int
}

// describe returns how a token appears in error messages.
func describe(t token) string {
	if t.kind == 0 {
		return "end of input"
	}
	if t.kind == '0' {
		return strconv.FormatInt(t.val, 10)
	}
	return fmt.Sprintf("'%c'", t.kind)
}

// tokenize splits src into tokens, ending with a kind-0 token.
func tokenize(src string) ([]token, error) {
	var toks []token
	for i := 0; i < len(src); {
		c := src[i]
		switch {
		case c == ' ' || c == '\t':
			i++
		case strings.IndexByte("+-*/%()", c) >= 0:
			toks = append(toks, token{kind: c, pos: i})
			i++
		case c >= '0' && c <= '9':
			start := i
			for i < len(src) && src[i] >= '0' && src[i] <= '9' {
				i++
			}
			v, err := strconv.ParseInt(src[start:i], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("number %q at offset %d is out of range", src[start:i], start)
			}
			toks = append(toks, token{kind: '0', val: v, pos: start})
		default:
			return nil, fmt.Errorf("unexpected character %q at offset %d", c, i)
		}
	}
	return append(toks, token{pos: len(src)}), nil
}

// evaluate is a recursive-descent parser built from closures: expr, term
// and factor are mutually recursive, so the variables are declared first
// and the function literals assigned afterwards (a literal can't refer to
// a name declared in the same := statement). They share pos and depth by
// capturing them. Each rule traces itself with enter, and a deferred
// closure that reads the rule's named results after return has set them.
func evaluate(src string, w io.Writer) (int64, error) {
	toks, err := tokenize(src)
	if err != nil {
		return 0, err
	}
	pos, depth := 0, 0
	peek := func() token { return toks[pos] }
	next := func() token {
		t := toks[pos]
		pos++
		return t
	}

	enter := func(rule string) {
		fmt.Fprintf(w, "%s> %s\n", strings.Repeat("  ", depth), rule)
		depth++
	}
	leave := func(rule string, v int64, err error) {
		depth--
		if err != nil {
			fmt.Fprintf(w, "%s< %s error\n", strings.Repeat("  ", depth), rule)
			return
		}
		fmt.Fprintf(w, "%s< %s = %d\n", strings.Repeat("  ", depth), rule, v)
	}

	var expr, term, factor func() (int64, error)

	expr = func() (v int64, err error) {
		enter("expr")
		defer func() { leave("expr", v, err) }()
		if v, err = term(); err != nil {
			return 0, err
		}
		for k := peek().kind; k == '+' || k == '-'; k = peek().kind {
			next()
			rhs, err := term()
			if err != nil {
				return 0, err
			}
			if k == '+' {
				v += rhs
			} else {
				v -= rhs
			}
		}
		return v, nil
	}

	term = func() (v int64, err error) {
		enter("term")
		defer func() { leave("term", v, err) }()
		if v, err = factor(); err != nil {
			return 0, err
		}
		for k := peek().kind; k == '*' || k == '/' || k == '%'; k = peek().kind {
			op := next()
			rhs, err := factor()
			if err != nil {
				return 0, err
			}
			switch {
			case k == '*':
				v *= rhs
			case rhs == 0:
				return 0, fmt.Errorf("offset %d: %w", op.pos, ErrDivisionByZero)
			case k == '/':
				v /= rhs
			default:
				v %= rhs
			}
		}
		return v, nil
	}

	factor = func() (v int64, err error) {
		enter("factor")
		defer func() { leave("factor", v, err) }()
		t := next()
		switch t.kind {
		case '-':
			if v, err = factor(); err != nil {
				return 0, err
			}
			return -v, nil
		case '0':
			return t.val, nil
		case '(':
			if v, err = expr(); err != nil {
				return 0, err
			}
			if closing := next(); closing.kind != ')' {
				return 0, fmt.Errorf("missing ')' at offset %d", closing.pos)
			}
			return v, nil
		case 0:
			return 0, errors.New("unexpected end of input")
		default:
			return 0, fmt.Errorf("unexpected %s at offset %d", describe(t), t.pos)
		}
	}

	v, err := expr()
	if err != nil {
		return 0, err
	}
	if t := peek(); t.kind != 0 {
		return 0, fmt.Errorf("unexpected %s at offset %d", describe(t), t.pos)
	}
	return v, nil
}
