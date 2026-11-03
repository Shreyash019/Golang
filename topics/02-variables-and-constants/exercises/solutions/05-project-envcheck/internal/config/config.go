// Package config loads a service's settings from environment variables,
// with a default for anything that isn't set.
//
// Solution E5.
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

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

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

// The defaults and limits are constants: they never change at runtime, and
// naming them puts every default in one place for a reader (and a reviewer).
const (
	defaultPort         = 8080
	defaultTimeout      = 5 * time.Second // typed: time.Duration
	defaultLevel        = LevelInfo
	defaultMaxBodyBytes = 1 << 20 // 1 MiB, untyped: fits the int64 field
	maxBodyBytesLimit   = 1 << 30 // 1 GiB
	minPort, maxPort    = 1, 65535
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
	switch strings.ToLower(s) {
	case "debug":
		return LevelDebug, nil
	case "info":
		return LevelInfo, nil
	case "warn":
		return LevelWarn, nil
	case "error":
		return LevelError, nil
	}
	return LevelUnset, fmt.Errorf("unknown log level %q (want debug, info, warn or error)", s)
}

// Load builds a Config from the variables that lookup reports, using the
// defaults for the ones that aren't set. See the package comment for the
// rules.
//
// Each variable is parsed into a variable declared with := inside its own
// if block, and only copied into cfg once it is valid. Nothing in an inner
// block is named like a variable outside it, so nothing can be shadowed.
func Load(lookup LookupFunc) (Config, error) {
	cfg := Config{
		Port:         defaultPort,
		Timeout:      defaultTimeout,
		Level:        defaultLevel,
		MaxBodyBytes: defaultMaxBodyBytes,
	}

	if s, ok := value(lookup, "APP_PORT"); ok {
		port, err := strconv.Atoi(s)
		if err != nil {
			return Config{}, fmt.Errorf("APP_PORT: %q is not an integer", s)
		}
		if port < minPort || port > maxPort {
			return Config{}, fmt.Errorf("APP_PORT: %d is outside %d-%d", port, minPort, maxPort)
		}
		cfg.Port = port
	}

	if s, ok := value(lookup, "APP_TIMEOUT"); ok {
		timeout, err := time.ParseDuration(s)
		if err != nil {
			return Config{}, fmt.Errorf("APP_TIMEOUT: %q is not a duration such as 5s or 300ms", s)
		}
		if timeout <= 0 {
			return Config{}, fmt.Errorf("APP_TIMEOUT: %v must be greater than zero", timeout)
		}
		cfg.Timeout = timeout
	}

	if s, ok := value(lookup, "APP_LOG_LEVEL"); ok {
		level, err := ParseLevel(s)
		if err != nil {
			return Config{}, fmt.Errorf("APP_LOG_LEVEL: %w", err)
		}
		cfg.Level = level
	}

	if s, ok := value(lookup, "APP_MAX_BODY_BYTES"); ok {
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return Config{}, fmt.Errorf("APP_MAX_BODY_BYTES: %q is not an integer", s)
		}
		if n < 1 || n > maxBodyBytesLimit {
			return Config{}, fmt.Errorf("APP_MAX_BODY_BYTES: %d is outside 1-%d", n, maxBodyBytesLimit)
		}
		cfg.MaxBodyBytes = n
	}

	return cfg, nil
}

// value returns the variable's value and true when it is set to a non-empty
// string: an empty value means "use the default", like an unset one.
func value(lookup LookupFunc, key string) (string, bool) {
	s, ok := lookup(key)
	return s, ok && s != ""
}
