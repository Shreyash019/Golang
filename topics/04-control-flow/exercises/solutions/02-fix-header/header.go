// Package header parses the header block of a small text format:
//
//	# comment
//	title = Weekly report
//	owner = ops
//	---
//	the body starts here and may contain anything, even a = b
package header

import "strings"

// Field is one "key = value" line of a header.
type Field struct {
	Key, Value string
}

// Parse returns the header fields in the order they appear. It reads lines
// up to the "---" separator; blank lines, lines starting with # and lines
// without "=" are skipped. Everything after the separator is the body and is
// ignored.
//
// The bug was an unlabeled break inside the switch: in Go, break leaves the
// innermost for, switch or select, so it ended the switch and the loop went
// on to parse the body (continue has no such problem: it always applies to a
// loop). Returning is the smallest fix and says exactly what we mean: the
// header is complete. A labeled break (break lines) would also work, and is
// the right tool when code after the loop still has to run.
func Parse(lines []string) []Field {
	var fields []Field
	for _, line := range lines {
		line = strings.TrimSpace(line)
		switch {
		case line == "" || strings.HasPrefix(line, "#"):
			continue
		case line == "---":
			return fields
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		fields = append(fields, Field{Key: strings.TrimSpace(key), Value: strings.TrimSpace(value)})
	}
	return fields
}
