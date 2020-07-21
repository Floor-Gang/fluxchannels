package internal

import (
	utilConfig "github.com/Floor-Gang/utilpkg/config"
	"log"
)

// Config structure.
type Config struct {
	Auth              string                  `yaml:"auth_server"`
	Prefix            string                  `yaml:"prefix"`
	Categories        map[string]FluxCategory `yaml:"categories"`
	ChannelNamePrefix string                  `yaml:"channel_name_prefix"`
}

const configPath = "./config.yml"

// GetConfig retrieves a configuration.
func GetConfig() Config {
	config := Config{
		Prefix:     ".flux",
		Categories: make(map[string]FluxCategory),
	}
	err := utilConfig.GetConfig(configPath, &config)

	if err != nil {
		log.Fatalln(err)
	}

	return config
}

// Save saves configuration
func (config *Config) Save() error {
	return utilConfig.Save(configPath, config)
}
