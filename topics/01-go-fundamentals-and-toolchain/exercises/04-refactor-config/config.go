package main

import "strconv"

// IConfigManager is the configuration service contract.
type IConfigManager interface {
	LoadFromArgs(args []string)
	GetPort() int
	SetPort(port int)
	GetEnv() string
	SetEnv(env string)
	GetDebug() bool
	SetDebug(debug bool)
}

// ConfigManagerImpl is the default IConfigManager implementation.
type ConfigManagerImpl struct {
	port  int
	env   string
	debug bool
}

// ConfigException is thrown when the configuration is invalid.
type ConfigException struct {
	Message string
}

// NewConfigManager is the factory for IConfigManager.
func NewConfigManager() IConfigManager {
	manager := &ConfigManagerImpl{}
	manager.SetPort(8080)
	manager.SetEnv("dev")
	manager.SetDebug(false)
	return manager
}

// LoadFromArgs loads key=value arguments, throwing ConfigException on bad input.
func (c *ConfigManagerImpl) LoadFromArgs(args []string) {
	utils := StringUtils{}
	for _, arg := range args {
		key, value := utils.SplitKeyValue(arg)
		if utils.IsNullOrEmpty(key) || utils.IsNullOrEmpty(value) {
			panic(&ConfigException{Message: "invalid argument " + strconv.Quote(arg) + ": want key=value"})
		}
		switch key {
		case "port":
			port, err := strconv.Atoi(value)
			if err != nil || port < 1 || port > 65535 {
				panic(&ConfigException{Message: "invalid port " + strconv.Quote(value)})
			}
			c.SetPort(port)
		case "env":
			c.SetEnv(value)
		case "debug":
			debug, err := strconv.ParseBool(value)
			if err != nil {
				panic(&ConfigException{Message: "invalid debug " + strconv.Quote(value)})
			}
			c.SetDebug(debug)
		default:
			panic(&ConfigException{Message: "unknown key " + strconv.Quote(key)})
		}
	}
}

// GetPort returns the port.
func (c *ConfigManagerImpl) GetPort() int { return c.port }

// SetPort sets the port.
func (c *ConfigManagerImpl) SetPort(port int) { c.port = port }

// GetEnv returns the environment.
func (c *ConfigManagerImpl) GetEnv() string { return c.env }

// SetEnv sets the environment.
func (c *ConfigManagerImpl) SetEnv(env string) { c.env = env }

// GetDebug returns the debug flag.
func (c *ConfigManagerImpl) GetDebug() bool { return c.debug }

// SetDebug sets the debug flag.
func (c *ConfigManagerImpl) SetDebug(debug bool) { c.debug = debug }
