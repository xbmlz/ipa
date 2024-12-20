package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	DefaultConfigPath         = "./configs"
	DefaultConfigName         = ".env"
	DefaultConfigOverrideName = ".local.env"
)

type Config interface {
	GetString(key, defaultValue string) string
	GetInt(key string, defaultValue int) int
	GetFloat64(key string, defaultValue float64) float64
	GetBool(key string, defaultValue bool) bool
}

type logger interface {
	Warnf(format string, a ...interface{})
	Infof(format string, a ...interface{})
	Debugf(format string, a ...interface{})
	Fatalf(format string, a ...interface{})
}

type config struct {
	logger logger
}

func New(configPath string, logger logger) Config {
	c := &config{logger: logger}
	c.load(configPath)
	return c
}

func (c *config) load(configPath string) {
	var (
		defaultFile         = configPath + DefaultConfigName
		defaultOverrideFile = configPath + DefaultConfigOverrideName
		env                 = c.GetString("APP_ENV", "")
	)

	err := godotenv.Load(defaultFile)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			c.logger.Fatalf("Failed to load config from file: %v, Err: %v", defaultFile, err)
		}
		c.logger.Warnf("Failed to load config from file: %v, Err: %v", defaultFile, err)
	} else {
		c.logger.Infof("Loaded config from file: %v", defaultFile)
	}

	if env != "" {
		// If 'APP_ENV' is set to x, then GoFr will read '.env' from configs directory, and then it will be overwritten
		// by configs present in file '.x.env'
		defaultOverrideFile = fmt.Sprintf("%s/.%s.env", configPath, env)
	}

	err = godotenv.Overload(defaultOverrideFile)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			c.logger.Fatalf("Failed to load config from file: %v, Err: %v", defaultOverrideFile, err)
		}
	} else {
		c.logger.Infof("Loaded config from file: %v", defaultOverrideFile)
	}
}

// GetString returns the env variable for the given key
// and falls back to the given defaultValue if not set
func (c *config) GetString(key, defaultValue string) string {
	v, ok := os.LookupEnv(key)
	if ok {
		return v
	}
	return defaultValue
}

// GetInt returns the env variable (parsed as integer) for
// the given key and falls back to the given defaultValue if not set
func (c *config) GetInt(key string, defaultValue int) int {
	v, ok := os.LookupEnv(key)
	if ok {
		value, err := strconv.Atoi(v)
		if err != nil {
			return defaultValue
		}
		return value
	}
	return defaultValue
}

// GetFloat64 returns the env variable (parsed as float64) for
// the given key and falls back to the given defaultValue if not set
func (c *config) GetFloat64(key string, defaultValue float64) float64 {
	v, ok := os.LookupEnv(key)
	if ok {
		value, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return defaultValue
		}
		return value
	}
	return defaultValue
}

// GetBool returns the env variable (parsed as bool) for
// the given key and falls back to the given defaultValue if not set
func (c *config) GetBool(key string, defaultValue bool) bool {
	v, ok := os.LookupEnv(key)
	if ok {
		value, err := strconv.ParseBool(v)
		if err != nil {
			return defaultValue
		}
		return value
	}
	return defaultValue
}
