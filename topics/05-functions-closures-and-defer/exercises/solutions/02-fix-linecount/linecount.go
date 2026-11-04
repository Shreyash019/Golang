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
//
// The bug was defer rc.Close() directly in the loop: a deferred call runs
// when the surrounding *function* returns, not at the end of the iteration,
// so every file stayed open until CountLines returned. Moving the body of
// the loop into countFile gives each file its own function, so its defer
// runs as soon as that file is counted.
func CountLines(names []string, open Opener) (int, error) {
	total := 0
	for _, name := range names {
		n, err := countFile(name, open)
		if err != nil {
			return 0, fmt.Errorf("count lines: %w", err)
		}
		total += n
	}
	return total, nil
}

// countFile counts the lines in one file. The error from Close is ignored
// on purpose: the file was only read, so a failed close can't lose data
// (closing a file you *wrote* is different; see the notes, §7.6).
func countFile(name string, open Opener) (int, error) {
	rc, err := open(name)
	if err != nil {
		return 0, err
	}
	defer rc.Close()

	n := 0
	sc := bufio.NewScanner(rc)
	for sc.Scan() {
		n++
	}
	if err := sc.Err(); err != nil {
		return 0, fmt.Errorf("read %s: %w", name, err)
	}
	return n, nil
}
