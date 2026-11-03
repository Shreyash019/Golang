// Package config loads a service's settings from environment variables,
// with a default for anything that isn't set.
//
// Exercise E5 — mini-project. Implement ParseLevel and Load until
// config_test.go and ../../cmd/envcheck/main_test.go pass. It trains
// constants for defaults (a typed time.Duration constant, untyped size
// constants), an iota enum whose zero value means "not set", and parsing
// with := and = without shadowing (notes §5.4, §5.5, §7.1, §7.9).
//
//	make ex E=05-project-envcheck
//	APP_PORT=9090 go run ./exercises/05-project-envcheck/cmd/envcheck
//
// The variables, their defaults and their rules:
//
//	APP_PORT            8080     an integer from 1 to 65535
//	APP_TIMEOUT         5s       a time.ParseDuration string, greater than zero
//	APP_LOG_LEVEL       info     debug, info, warn or error, in any letter case
//	APP_MAX_BODY_BYTES  1048576  an integer from 1 to 1073741824 (1 MiB default, 1 GiB limit)
//
// A variable that isn't set, or is set to the empty string, takes its
// default. Load checks the variables in the order above and stops at the
// first invalid one, returning an error whose message starts with the
// variable's name, e.g. `APP_PORT: "abc" is not an integer`.
package config

import "time"

// Level is a log level. Its zero value, LevelUnset, means "not set", so a
// Level nobody assigned can't be mistaken for a real level.
type Level int

// The log levels, in increasing severity.
const (
	LevelUnset Level = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
)

// String returns the level's name as written in APP_LOG_LEVEL, or "unset".
// (A method on a defined type: methods are taught in topic 10.)
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	default:
		return "unset"
	}
}

// Config is the effective configuration: every field is set after a
// successful Load.
type Config struct {
	Port         int
	Timeout      time.Duration
	Level        Level
	MaxBodyBytes int64
}

// LookupFunc returns the value of an environment variable and whether it is
// set. os.LookupEnv is a LookupFunc; tests pass a fake one.
type LookupFunc func(key string) (string, bool)

// ParseLevel converts a level name (debug, info, warn, error; any letter
// case) into a Level. Any other string, including "", is an error.
func ParseLevel(s string) (Level, error) {
	return LevelUnset, nil
}

// Load builds a Config from the variables that lookup reports, using the
// defaults for the ones that aren't set. See the package comment for the
// rules.
func Load(lookup LookupFunc) (Config, error) {
	return Config{}, nil
}
