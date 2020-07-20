package internal

import (
	utilConfig "github.com/Floor-Gang/utilpkg/config"
	"log"
)

// Config structure.
type Config struct {
	Auth   string `yaml:"auth_server"`
	Prefix string `yaml:"prefix"`
}

const configPath = "./config.yml"

// GetConfig retrieves a configuration.
func GetConfig() Config {
	config := Config{
		Prefix: ".<command name>",
	}
	err := utilConfig.GetConfig(configPath, &config)

	if err != nil {
		log.Fatalln(err)
	}

	return config
}

// Save saves configuration
func (config *Config) Save() {
	if err := utilConfig.Save(configPath, config); err != nil {
		log.Println("Failed to save config", err)
	}
}
