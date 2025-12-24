package config

import (
	"os"
	"path/filepath"
	"testing"
)

// Integration tests for config package
// These tests verify the full config loading cycle

func TestIntegration_FullConfigCycle(t *testing.T) {
	// Save and restore HOME
	origHome := os.Getenv("HOME")
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	// Create config directory and file
	configDir := filepath.Join(tmpDir, ".config", "timeotter")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}

	configContent := `
CalendarID = "integration@calendar.google.com"
CmdToExec = "notify-send 'Meeting'"
TokenFile = "~/token.json"
CredentialsFile = "~/credentials.json"
BackupFile = "~/backup.txt"
MaxRes = 20
TriggerBeforeMinutes = 10
CronMarker = "# integration test crons"
ShowDeleted = true
`
	configPath := filepath.Join(configDir, "config.toml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	// Read config
	v, err := ReadConfig()
	if err != nil {
		t.Fatalf("ReadConfig failed: %v", err)
	}

	// Unmarshal
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Validate
	if err := ValidateConfig(&config); err != nil {
		t.Fatalf("ValidateConfig failed: %v", err)
	}

	// Verify values
	if config.CalendarID != "integration@calendar.google.com" {
		t.Errorf("CalendarID mismatch: %s", config.CalendarID)
	}
	if config.CmdToExec != "notify-send 'Meeting'" {
		t.Errorf("CmdToExec mismatch: %s", config.CmdToExec)
	}
	if config.MaxRes != 20 {
		t.Errorf("MaxRes mismatch: %d", config.MaxRes)
	}
	if config.TriggerBeforeMinutes != 10 {
		t.Errorf("TriggerBeforeMinutes mismatch: %d", config.TriggerBeforeMinutes)
	}
	if config.CronMarker != "# integration test crons" {
		t.Errorf("CronMarker mismatch: %s", config.CronMarker)
	}
	if !config.ShowDeleted {
		t.Error("ShowDeleted should be true")
	}

	// Verify path expansion
	expectedTokenFile := tmpDir + "/token.json"
	if config.TokenFile != expectedTokenFile {
		t.Errorf("TokenFile not expanded: got %s, want %s", config.TokenFile, expectedTokenFile)
	}
}

func TestIntegration_ConfigWithDefaults(t *testing.T) {
	// Save and restore HOME
	origHome := os.Getenv("HOME")
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	// Create config with minimal required fields
	configDir := filepath.Join(tmpDir, ".config", "timeotter")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}

	configContent := `
CalendarID = "minimal@calendar.google.com"
CmdToExec = "echo hello"
TokenFile = "/tmp/token.json"
`
	configPath := filepath.Join(configDir, "config.toml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	v, err := ReadConfig()
	if err != nil {
		t.Fatalf("ReadConfig failed: %v", err)
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Verify defaults are applied
	if config.MaxRes != 5 {
		t.Errorf("default MaxRes should be 5, got %d", config.MaxRes)
	}
	if config.TriggerBeforeMinutes != 5 {
		t.Errorf("default TriggerBeforeMinutes should be 5, got %d", config.TriggerBeforeMinutes)
	}
	if config.ShowDeleted {
		t.Error("default ShowDeleted should be false")
	}
	if config.CronMarker != "# custom crons below this can be deleted." {
		t.Errorf("default CronMarker mismatch: %s", config.CronMarker)
	}
}

func TestIntegration_ConfigValidationConstraints(t *testing.T) {
	// Save and restore HOME
	origHome := os.Getenv("HOME")
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	configDir := filepath.Join(tmpDir, ".config", "timeotter")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}

	// Config with out-of-bounds values
	configContent := `
CalendarID = "bounds@calendar.google.com"
CmdToExec = "echo test"
TokenFile = "/tmp/token.json"
MaxRes = 500
TriggerBeforeMinutes = -10
`
	configPath := filepath.Join(configDir, "config.toml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	v, err := ReadConfig()
	if err != nil {
		t.Fatalf("ReadConfig failed: %v", err)
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if err := ValidateConfig(&config); err != nil {
		t.Fatalf("ValidateConfig failed: %v", err)
	}

	// Verify constraints were applied
	if config.MaxRes != 100 {
		t.Errorf("MaxRes should be clamped to 100, got %d", config.MaxRes)
	}
	if config.TriggerBeforeMinutes != 0 {
		t.Errorf("TriggerBeforeMinutes should be clamped to 0, got %d", config.TriggerBeforeMinutes)
	}
}

func TestIntegration_ConfigFilePermissions(t *testing.T) {
	// Save and restore HOME
	origHome := os.Getenv("HOME")
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	configDir := filepath.Join(tmpDir, ".config", "timeotter")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}

	configContent := `
CalendarID = "test@calendar.google.com"
CmdToExec = "echo hello"
TokenFile = "/tmp/token.json"
`
	configPath := filepath.Join(configDir, "config.toml")

	// Create config file with restricted permissions
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	// Should still be able to read
	_, err := ReadConfig()
	if err != nil {
		t.Errorf("failed to read config with 0600 permissions: %v", err)
	}
}
