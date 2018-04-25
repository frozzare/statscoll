package config

import (
	"errors"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

var (
	// ErrUnmarshal is returned when config can't be unmarshaled.
	ErrUnmarshal = errors.New("can't unmarshal config value")
)

// Config represents a config file.
type Config struct {
	Port int    `yaml:"port"`
	DSN  string `yaml:"dsn"`
}

// Default sets the default config.
func (c *Config) Default() error {
	if c.Port == 0 {
		c.Port = 9300
	}

	if len(c.DSN) == 0 {
		c.DSN = "root@/statscoll?charset=utf8&parseTime=true"
	}

	return nil
}

// ReadFile creates a new config struct from a yaml file.
func ReadFile(file string) (*Config, error) {
	var config *Config

	dat, err := ioutil.ReadFile(file)
	if err == nil {
		if err := yaml.Unmarshal(dat, &config); err != nil {
			return &Config{}, err
		}
	} else {
		config = &Config{}
	}

	if err := config.Default(); err != nil {
		return nil, err
	}

	return config, nil
}
