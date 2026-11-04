// Exercise E2 — fix the bug.
//
// Goal: Parse must stop at the "---" separator, but body lines that look like
// "key = value" end up in the header. Find the root cause and fix it with the
// smallest change, then explain the gotcha in one sentence. It trains break
// and continue inside a switch (notes §7.1) and labels (§5.8).
//
//	make ex E=02-fix-header

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
func Parse(lines []string) []Field {
	var fields []Field
	for _, line := range lines {
		line = strings.TrimSpace(line)
		switch {
		case line == "" || strings.HasPrefix(line, "#"):
			continue
		case line == "---":
			break // end of the header
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		fields = append(fields, Field{Key: strings.TrimSpace(key), Value: strings.TrimSpace(value)})
	}
	return fields
}
