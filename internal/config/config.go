package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	URL  string `json:"db_url"`
	User string `json:"current_user_name"`
}

func Read() (Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	conf, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("unable to read config file: %v", err)
	}

	config := Config{}
	err = json.Unmarshal(conf, &config)
	if err != nil {
		return Config{}, fmt.Errorf("unable to parse config file: %v", err)
	}

	return config, nil
}

func (c *Config) SetUser(username string) error {
	c.User = username
	return write(*c)
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to resolve home directory: %v", err)
	}
	fulldir := filepath.Join(home, configFileName)
	return fulldir, nil
}

func write(cfg Config) error {
	file, err := getConfigFilePath()
	if err != nil {
		return err
	}
	raw, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(file, raw, 0644)
}
