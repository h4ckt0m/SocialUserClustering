package config

import (
	"encoding/json"
	"errors"
	"os"
)

// Configuration struct
type Configuration struct {
	Usernames []string `json:"usernames"`
	DBServer string `json:"dbserver"`
	APIKEY string `json:"apikey"`
}

// ParseConf loads config file, parses options and stores it in configuration struct
func ParseConf(cFile string) (Configuration, error) {
	var conf Configuration

	// If config file does not exist creates one
	if _, err := os.Stat(cFile); err != nil {
		return conf, err
	}

	// Open config and parses json into
	file, _ := os.Open(cFile)
	defer file.Close()
	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&conf); err != nil {
		return conf, errors.New("Failed to decode config file")
	}

	return conf, nil
}
