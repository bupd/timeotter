package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

// Config structure to match the TOML structure
type Config struct {
	CalendarID string `mapstructure:"CalendarID"`
	CmdToExec  string `mapstructure:"CmdToExec"`
	MaxRes     int64  `mapstructure:"MaxRes"`
	TokenFile  string `mapstructure:"TokenFile"`
}

// Helper function to read the config file using Viper
func ReadConfig() (*viper.Viper, error) {
	// Resolve the tilde (~) in the path
	configPath := os.ExpandEnv("~/.config/timeotter/config.toml")
	// Initialize a new Viper instance
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("toml")

	// Read the configuration file
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w. Please ensure the config file exists.", err)
	}

	return v, nil
}

// Reads config and gives the config
func GetConfig() Config {
	// Load the config file using Viper
	v, err := ReadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Map the values from Viper into the Config struct
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		log.Fatalf("Error unmarshalling config: %v", err)
	}

	// Now use the loaded config values
	fmt.Printf("Calendar ID: %s\n", config.CalendarID)
	fmt.Printf("Command to Execute: %s\n", config.CmdToExec)
	fmt.Printf("Max Resolution: %d\n", config.MaxRes)
	fmt.Printf("Token File: %s\n", config.TokenFile)

	return config
}
