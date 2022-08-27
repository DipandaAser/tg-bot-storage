package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

const (
	configFilePath = "./config/config.yaml"
)

var defaultConfig *Config

//Config describes the server configuration
type Config struct {
	Tokens []string `yaml:"tokens"`
	ApiKey string   `yaml:"api_key"`
}

//GetDefaultConfig returns a config with default values from the yaml configFilePath
func GetDefaultConfig() Config {
	if defaultConfig == nil {
		defaultConfig = &Config{Tokens: make([]string, 0)}
		// read the config file
		file, err := os.Open(configFilePath)
		if err != nil {
			return Config{}
		}
		defer file.Close()

		err = yaml.NewDecoder(file).Decode(defaultConfig)
		if err != nil {
			return Config{}
		}
	}
	return *defaultConfig
}
