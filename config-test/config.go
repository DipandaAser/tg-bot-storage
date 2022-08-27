//
package config_test

import (
	"github.com/DipandaAser/tg-bot-storage/internal/config"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	config.Config `yaml:",inline"`
	//ChatID is the ID of the chat to store message, this is use in some test
	ChatID int64 `yaml:"chat_id"`
	//DraftChatID is the ID of the chat used as a draft chat when downloading file
	DraftChatID int64 `yaml:"draft_chat_id"`
}

var defaultConfig *Config

//GetDefaultConfig returns a config with default values from the yaml configFilePath
func GetConfig(configFilePath string) Config {
	if defaultConfig == nil {
		defaultConfig = &Config{Config: config.Config{Tokens: make([]string, 0)}}
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
