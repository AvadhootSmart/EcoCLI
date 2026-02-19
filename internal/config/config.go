package config

import (
	"encoding/json"
	"fmt"

	// "fmt"
	"os"
	"path/filepath"
)

const (
	ConfigDir  = ".config/eco"
	ConfigFile = "config.json"
)

// Config holds all configuration for the eco daemon
type Config struct {
	DeviceID     string
	SharedSecret string
	Port         int
}

// ConfigPath returns the full path to the config file
func ConfigPath() (configPath string, err error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(userHomeDir, ConfigDir, ConfigFile), nil
}

// Load reads the config from disk
func Load() (*Config, error) {
	cfgPath, err := ConfigPath()
	if err != nil {
		return nil, err
	}

	cfgData, err := os.ReadFile(cfgPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, err
	}

	var config Config
	err = json.Unmarshal(cfgData, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// Save writes the config to disk
func (c *Config) Save() error {
	cfgPath, err := ConfigPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(cfgPath), 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(cfgPath, data, 0600)
}

// IsInitialized checks if the config has been set up
func (c *Config) IsInitialized() bool {
	if c.DeviceID != "" && c.SharedSecret != "" {
		return true
	}
	return false
}

// GetDeviceCredentials returns the device ID and secret for auth
func (c *Config) GetDeviceCredentials() (deviceID, secret string) {
	return c.DeviceID, c.SharedSecret
}

func (c *Config) DeleteConfig() error {
	if !c.IsInitialized() {
		return fmt.Errorf("Config does not exist to delete")
	}

	cfgPath, err := ConfigPath()
	if err != nil {
		return err
	}

	err = os.Remove(cfgPath)
	if err != nil {
		return err
	}

	return nil
}
