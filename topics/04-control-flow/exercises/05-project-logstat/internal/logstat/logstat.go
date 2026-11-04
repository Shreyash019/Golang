// Package logstat filters plain-text log lines by level and counts what it
// saw. A log line is "<timestamp> <LEVEL> <message>", for example
//
//	2026-11-04T10:00:03Z ERROR db: connection refused
package logstat

import "io"

// Level is the severity of a log line. The zero value, LevelUnknown, is not a
// real level, so a Level that was never set can't pass for LevelDebug.
type Level int

// The levels, from least to most severe.
const (
	LevelUnknown Level = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

// ParseLevel parses a level name, ignoring case: "debug", "info", "warn" (or
// "warning"), "error" and "fatal". ok is false for any other name.
func ParseLevel(s string) (l Level, ok bool) {
	return LevelUnknown, false
}

// Summary counts the lines Filter read.
type Summary struct {
	Lines     int // every line read, including skipped and malformed ones
	Matched   int // lines at or above the minimum level, written to the output
	Skipped   int // blank lines and comments (starting with #)
	Malformed int // other lines with fewer than 3 fields or an unknown level

	// Lines per level, matched or not.
	Debug, Info, Warn, Error, Fatal int
}

// Filter reads log lines from r and writes those at or above minLevel to w,
// unchanged, one per line. Blank lines and lines starting with # (after
// leading spaces) are skipped; lines with fewer than three fields or an
// unknown level are counted as malformed. When limit > 0, Filter stops
// reading right after the limit-th matching line, like grep -m. Lines
// longer than bufio.MaxScanTokenSize (64 KiB) are an error.
//
// Tools you'll want: bufio.NewScanner (its Scan method is a loop condition),
// strings.TrimSpace, strings.HasPrefix, strings.Fields, fmt.Fprintln, and
// the scanner's Err method after the loop. Wrap errors with fmt.Errorf and %w.
func Filter(r io.Reader, w io.Writer, minLevel Level, limit int) (Summary, error) {
	return Summary{}, nil
}
