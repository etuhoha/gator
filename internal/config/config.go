package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	DbUrl string `json:"db_url"`
	User  string `json:"current_user_name"`
}

func Read() (Config, error) {
	var config Config

	data, err := os.ReadFile(configFile())
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func (c *Config) SetUser(user string) error {
	c.User = user

	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = os.WriteFile(configFile(), data, 0x666)
	if err != nil {
		return err
	}

	return nil
}

func configFile() string {
	fname := ".gatorconfig.json"

	home, err := os.UserHomeDir()
	if err != nil {
		return fname
	}

	return filepath.Join(home, fname)
}
