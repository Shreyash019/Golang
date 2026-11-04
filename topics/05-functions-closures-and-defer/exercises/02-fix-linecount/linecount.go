// Exercise E2 — fix the bug.
//
// Goal: CountLines returns the right totals, but it keeps every file open
// until it returns. Given 10 000 files, a real program fails with "too many
// open files". The test counts how many files are open at once. Find the
// root cause, fix it with a small change, and explain the gotcha in one
// sentence. It trains when deferred calls run (notes §5.7, §7.4).
//
//	make ex E=02-fix-linecount

// Package linecount counts the lines in a set of files.
package linecount

import (
	"bufio"
	"fmt"
	"io"
)

// Opener opens the named file for reading. os.Open fits after wrapping:
//
//	func(name string) (io.ReadCloser, error) { return os.Open(name) }
type Opener func(name string) (io.ReadCloser, error)

// CountLines returns the total number of lines in the named files. It opens
// one file at a time and closes each file before opening the next.
func CountLines(names []string, open Opener) (int, error) {
	total := 0
	for _, name := range names {
		rc, err := open(name)
		if err != nil {
			return 0, fmt.Errorf("count lines: %w", err)
		}
		defer rc.Close()

		sc := bufio.NewScanner(rc)
		for sc.Scan() {
			total++
		}
		if err := sc.Err(); err != nil {
			return 0, fmt.Errorf("count lines in %s: %w", name, err)
		}
	}
	return total, nil
}
