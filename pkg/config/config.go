// Package config handles configuration loading and validation for timeotter.
package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Config structure to match the TOML structure
type Config struct {
	CalendarID           string `mapstructure:"CalendarID"`
	CmdToExec            string `mapstructure:"CmdToExec"`
	MaxRes               int64  `mapstructure:"MaxRes"`
	TokenFile            string `mapstructure:"TokenFile"`
	CredentialsFile      string `mapstructure:"CredentialsFile"`
	BackupFile           string `mapstructure:"BackupFile"`
	TriggerBeforeMinutes int    `mapstructure:"TriggerBeforeMinutes"`
	CronMarker           string `mapstructure:"CronMarker"`
	ShowDeleted          bool   `mapstructure:"ShowDeleted"`
}

// ReadConfig reads the configuration file using Viper and returns the config instance.
func ReadConfig() (*viper.Viper, error) {
	dirname := GetHomeDir()
	configPath := fmt.Sprintf("%s/.config/timeotter/config.toml", dirname)
	// Initialize a new Viper instance
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("toml")

	// Set defaults
	v.SetDefault("MaxRes", 5)
	v.SetDefault("CredentialsFile", fmt.Sprintf("%s/.cal-credentials.json", dirname))
	v.SetDefault("BackupFile", fmt.Sprintf("%s/.crontab_backup.txt", dirname))
	v.SetDefault("TriggerBeforeMinutes", 5)
	v.SetDefault("CronMarker", "# custom crons below this can be deleted.")
	v.SetDefault("ShowDeleted", false)

	// Read the configuration file
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	return v, nil
}

// ExpandPath expands ~ to the user's home directory
func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		return strings.Replace(path, "~", GetHomeDir(), 1)
	}
	if path == "~" {
		return GetHomeDir()
	}
	return path
}

// ValidateConfig validates config values and applies constraints
func ValidateConfig(config *Config) error {
	// Validate required fields
	if config.CalendarID == "" {
		return fmt.Errorf("CalendarID is required")
	}
	if config.CmdToExec == "" {
		return fmt.Errorf("CmdToExec is required")
	}
	if config.TokenFile == "" {
		return fmt.Errorf("TokenFile is required")
	}

	// Validate MaxRes: min 1, max 100
	if config.MaxRes < 1 {
		config.MaxRes = 1
	}
	if config.MaxRes > 100 {
		config.MaxRes = 100
	}

	// Validate TriggerBeforeMinutes: must be non-negative
	if config.TriggerBeforeMinutes < 0 {
		config.TriggerBeforeMinutes = 0
	}

	// Expand ~ in file paths
	config.CredentialsFile = ExpandPath(config.CredentialsFile)
	config.BackupFile = ExpandPath(config.BackupFile)
	config.TokenFile = ExpandPath(config.TokenFile)

	return nil
}

// GetConfig loads, validates and returns the application configuration.
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

	// Validate and apply constraints
	if err := ValidateConfig(&config); err != nil {
		log.Fatalf("Error validating config: %v", err)
	}

	// Now use the loaded config values
	fmt.Printf("Token File: %s\n", config.TokenFile)
	fmt.Printf("Calendar ID: %s\n", config.CalendarID)
	fmt.Printf("Max Events: %d\n", config.MaxRes)
	fmt.Printf("Command to Execute: %s\n", config.CmdToExec)

	return config
}

// GetHomeDir returns the current user's home directory path.
func GetHomeDir() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("unable to get user home directory: %v", err)
	}
	// fmt.Println(dirname)

	return dirname
}
