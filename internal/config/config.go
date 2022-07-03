package config

import (
	"github.com/joho/godotenv"
	"os"
	"strings"
)

var defaultConfig *Config

//Config describes the server configuration
type Config struct {
	Tokens []string
	ApiKey string
}

//GetDefaultConfig returns a config with default values  and env variables
func GetDefaultConfig() Config {
	if defaultConfig == nil {
		//load postgres env variables
		_ = godotenv.Load()

		defaultConfig = &Config{
			Tokens: strings.Split(os.Getenv("TOKENS"), ","),
			ApiKey: os.Getenv("API_KEY"),
		}
	}
	return *defaultConfig
}
