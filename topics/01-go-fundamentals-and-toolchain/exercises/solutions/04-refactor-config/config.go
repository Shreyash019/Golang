package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Config is the server configuration. The fields have no invariants beyond
// what parseArgs checks, so they are plain exported fields: no getters.
type Config struct {
	Port  int
	Env   string
	Debug bool
}

// defaultConfig returns the configuration used when no arguments are given.
func defaultConfig() Config {
	return Config{Port: 8080, Env: "dev"}
}

// parseArgs applies key=value arguments over the defaults. Later values win.
func parseArgs(args []string) (Config, error) {
	cfg := defaultConfig()
	for _, arg := range args {
		key, value, ok := strings.Cut(arg, "=")
		if !ok || key == "" || value == "" {
			return Config{}, fmt.Errorf("invalid argument %q: want key=value", arg)
		}
		switch key {
		case "port":
			port, err := strconv.Atoi(value)
			if err != nil || port < 1 || port > 65535 {
				return Config{}, fmt.Errorf("invalid port %q: want 1-65535", value)
			}
			cfg.Port = port
		case "env":
			cfg.Env = value
		case "debug":
			debug, err := strconv.ParseBool(value)
			if err != nil {
				return Config{}, fmt.Errorf("invalid debug %q: want true or false", value)
			}
			cfg.Debug = debug
		default:
			return Config{}, fmt.Errorf("unknown key %q", key)
		}
	}
	return cfg, nil
}
