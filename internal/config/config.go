package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	Current_user_name string `json:"current_user_name"`
	Db_url            string `json:"db_url"`
}

func (cfg *Config) SetUser(username string) (err error) {
	cfg.Current_user_name = username

	err = write(*cfg)
	if err != nil {
		return fmt.Errorf("could not write file, error: %w", err)
	}

	return nil
}

func ReadConfig() (config Config, err error) {
	address, err := getConfigFilePath()
	if err != nil {
		return config, err
	}

	data, err := os.ReadFile(address)
	if err != nil {
		return config, fmt.Errorf("could not read file, error: %w", err)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("could not read file, error: %w", err)
	}

	return config, nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not find home directory, error: %w", err)
	}

	fullPath := filepath.Join(homeDir, configFileName)
	return fullPath, nil
}

func write(cfg Config) error {
	fullPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(cfg)
	if err != nil {
		return err
	}

	return nil
}
