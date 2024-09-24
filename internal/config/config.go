package config

import (
	"encoding/json"
	"os"
	"path"
)

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func Read() (Config, error) {
	configPath, err := getConfigFilePath()
	data, err := os.ReadFile(configPath)

	config := Config{}
	err = json.Unmarshal(data, &config)

	return config, err
}

func (cfg *Config) SetUser(user string) error {
	cfg.CurrentUserName = user
	return write(*cfg)
}

func getConfigFilePath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	filePath := path.Join(dir, configFileName)
	return filePath, nil
}

func write(cfg Config) error {
	configPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}
