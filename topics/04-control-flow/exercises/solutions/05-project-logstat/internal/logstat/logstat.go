// Package logstat filters plain-text log lines by level and counts what it
// saw. A log line is "<timestamp> <LEVEL> <message>", for example
//
//	2026-11-04T10:00:03Z ERROR db: connection refused
package logstat

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

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
//
// One switch on the normalised name, with several values in a case for the
// alias, replaces a chain of if/else-if comparisons.
func ParseLevel(s string) (l Level, ok bool) {
	switch strings.ToUpper(s) {
	case "DEBUG":
		return LevelDebug, true
	case "INFO":
		return LevelInfo, true
	case "WARN", "WARNING":
		return LevelWarn, true
	case "ERROR":
		return LevelError, true
	case "FATAL":
		return LevelFatal, true
	default:
		return LevelUnknown, false
	}
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
// The loop is the condition-only form of for: sc.Scan() is the condition.
// continue handles the lines that need no more work, a switch counts the
// level, and break leaves the loop when the limit is reached (the break is
// in an if, not a switch, so it applies to the for).
func Filter(r io.Reader, w io.Writer, minLevel Level, limit int) (Summary, error) {
	var s Summary
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := sc.Text()
		s.Lines++

		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			s.Skipped++
			continue
		}
		fields := strings.Fields(trimmed)
		if len(fields) < 3 {
			s.Malformed++
			continue
		}
		lvl, ok := ParseLevel(fields[1])
		if !ok {
			s.Malformed++
			continue
		}

		switch lvl {
		case LevelDebug:
			s.Debug++
		case LevelInfo:
			s.Info++
		case LevelWarn:
			s.Warn++
		case LevelError:
			s.Error++
		case LevelFatal:
			s.Fatal++
		}

		if lvl < minLevel {
			continue
		}
		if _, err := fmt.Fprintln(w, line); err != nil {
			return s, fmt.Errorf("write line %d: %w", s.Lines, err)
		}
		s.Matched++
		if limit > 0 && s.Matched == limit {
			break
		}
	}
	if err := sc.Err(); err != nil {
		return s, fmt.Errorf("read line %d: %w", s.Lines+1, err)
	}
	return s, nil
}
