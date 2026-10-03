package config

import (
	"encoding/json"
	// "fmt"
	"log"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbUrl string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (C *Config) SetUser(user string) {
	C.CurrentUserName= user

	// Write to file
	err := write(*C)
	if err != nil {
		log.Fatal(err)
	}
}

func Read() (Config, error) {
	filePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	filePath := home + "/" + configFileName

	// fmt.Printf("Home directory: %s\n", home)
	// fmt.Printf("File path: %s\n", filePath)

	return filePath, nil
}

func write(cfg Config) error {
	// data, err := json.MarshalIndent(cfg, "", "  ")
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	
	err = os.WriteFile(filePath, data, 0666)
	if err != nil {
		return err
	}

	return nil
}
