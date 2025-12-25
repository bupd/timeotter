package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bupd/timeotter/pkg/config"
	"github.com/bupd/timeotter/pkg/oauth"
	"golang.org/x/oauth2"
)

// E2E tests for timeotter main application
// These tests verify the application initialization and config loading

func TestE2E_ConfigInitialization(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	// Setup config
	configDir := filepath.Join(tmpDir, ".config", "timeotter")
	if err := os.MkdirAll(configDir, 0750); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}

	configContent := `
CalendarID = "e2e@calendar.google.com"
CmdToExec = "echo 'Meeting starting'"
TokenFile = "` + tmpDir + `/token.json"
CredentialsFile = "` + tmpDir + `/credentials.json"
MaxRes = 10
TriggerBeforeMinutes = 5
`
	configPath := filepath.Join(configDir, "config.toml")
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	// Test config loading
	v, err := config.ReadConfig()
	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}

	var conf config.Config
	if err := v.Unmarshal(&conf); err != nil {
		t.Fatalf("failed to unmarshal config: %v", err)
	}

	if err := config.ValidateConfig(&conf); err != nil {
		t.Fatalf("failed to validate config: %v", err)
	}

	// Verify config values are loaded correctly
	if conf.CalendarID != "e2e@calendar.google.com" {
		t.Errorf("CalendarID: got %s, want e2e@calendar.google.com", conf.CalendarID)
	}
	if conf.CmdToExec != "echo 'Meeting starting'" {
		t.Errorf("CmdToExec: got %s", conf.CmdToExec)
	}
	if conf.MaxRes != 10 {
		t.Errorf("MaxRes: got %d, want 10", conf.MaxRes)
	}
}

func TestE2E_TokenFileHandling(t *testing.T) {
	tmpDir := t.TempDir()
	tokenPath := filepath.Join(tmpDir, "test_token.json")

	// Create a token
	token := &oauth2.Token{
		AccessToken:  "e2e_access_token",
		TokenType:    "Bearer",
		RefreshToken: "e2e_refresh_token",
	}

	// Save using oauth package
	oauth.SaveToken(tokenPath, token)

	// Verify file exists
	if _, err := os.Stat(tokenPath); os.IsNotExist(err) {
		t.Fatal("token file was not created")
	}

	// Load token
	loaded, err := oauth.TokenFromFile(tokenPath)
	if err != nil {
		t.Fatalf("failed to load token: %v", err)
	}

	if loaded.AccessToken != "e2e_access_token" {
		t.Errorf("AccessToken: got %s, want e2e_access_token", loaded.AccessToken)
	}
}

func TestE2E_MissingCredentialsFile(t *testing.T) {
	tmpDir := t.TempDir()
	credPath := filepath.Join(tmpDir, "nonexistent_credentials.json")

	// Try to read non-existent credentials file
	_, err := os.ReadFile(credPath)
	if err == nil {
		t.Error("expected error for missing credentials file")
	}
}

func TestE2E_ConfigWithAllFields(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	configDir := filepath.Join(tmpDir, ".config", "timeotter")
	if err := os.MkdirAll(configDir, 0750); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}

	// Full config with all fields
	configContent := `
CalendarID = "full@calendar.google.com"
CmdToExec = "notify-send 'Event'"
TokenFile = "~/token.json"
CredentialsFile = "~/credentials.json"
BackupFile = "~/backup.txt"
MaxRes = 25
TriggerBeforeMinutes = 15
CronMarker = "# e2e test crons"
ShowDeleted = true
`
	configPath := filepath.Join(configDir, "config.toml")
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	v, err := config.ReadConfig()
	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}

	var conf config.Config
	if err := v.Unmarshal(&conf); err != nil {
		t.Fatalf("failed to unmarshal config: %v", err)
	}

	if err := config.ValidateConfig(&conf); err != nil {
		t.Fatalf("failed to validate config: %v", err)
	}

	// Verify all fields
	tests := []struct {
		name string
		got  interface{}
		want interface{}
	}{
		{"CalendarID", conf.CalendarID, "full@calendar.google.com"},
		{"CmdToExec", conf.CmdToExec, "notify-send 'Event'"},
		{"TokenFile", conf.TokenFile, tmpDir + "/token.json"},
		{"CredentialsFile", conf.CredentialsFile, tmpDir + "/credentials.json"},
		{"BackupFile", conf.BackupFile, tmpDir + "/backup.txt"},
		{"MaxRes", conf.MaxRes, int64(25)},
		{"TriggerBeforeMinutes", conf.TriggerBeforeMinutes, 15},
		{"CronMarker", conf.CronMarker, "# e2e test crons"},
		{"ShowDeleted", conf.ShowDeleted, true},
	}

	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s: got %v, want %v", tt.name, tt.got, tt.want)
		}
	}
}

func TestE2E_ConfigMissingRequiredFields(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	configDir := filepath.Join(tmpDir, ".config", "timeotter")
	if err := os.MkdirAll(configDir, 0750); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}

	// Config missing CalendarID
	configContent := `
CmdToExec = "echo hello"
TokenFile = "/tmp/token.json"
`
	configPath := filepath.Join(configDir, "config.toml")
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	v, err := config.ReadConfig()
	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}

	var conf config.Config
	if err := v.Unmarshal(&conf); err != nil {
		t.Fatalf("failed to unmarshal config: %v", err)
	}

	// ValidateConfig should fail
	err = config.ValidateConfig(&conf)
	if err == nil {
		t.Error("expected validation error for missing CalendarID")
	}
}

func TestE2E_GlobalVariablesAssignment(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	configDir := filepath.Join(tmpDir, ".config", "timeotter")
	if err := os.MkdirAll(configDir, 0750); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}

	configContent := `
CalendarID = "globals@calendar.google.com"
CmdToExec = "echo test"
TokenFile = "/tmp/token.json"
CredentialsFile = "/tmp/creds.json"
BackupFile = "/tmp/backup.txt"
MaxRes = 5
TriggerBeforeMinutes = 5
CronMarker = "# marker"
ShowDeleted = false
`
	configPath := filepath.Join(configDir, "config.toml")
	if err := os.WriteFile(configPath, []byte(configContent), 0600); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	v, err := config.ReadConfig()
	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}

	var conf config.Config
	if err := v.Unmarshal(&conf); err != nil {
		t.Fatalf("failed to unmarshal config: %v", err)
	}

	if err := config.ValidateConfig(&conf); err != nil {
		t.Fatalf("failed to validate config: %v", err)
	}

	// Simulate global variable assignment (as done in main.go)
	calendarID := conf.CalendarID
	cmdToExec := conf.CmdToExec
	maxRes := conf.MaxRes
	tokenFile := conf.TokenFile
	credentialsFile := conf.CredentialsFile
	backupFile := conf.BackupFile
	triggerBeforeMinutes := conf.TriggerBeforeMinutes
	cronMarker := conf.CronMarker
	showDeleted := conf.ShowDeleted

	// Verify assignments
	if calendarID != "globals@calendar.google.com" {
		t.Errorf("calendarID: got %s", calendarID)
	}
	if cmdToExec != "echo test" {
		t.Errorf("cmdToExec: got %s", cmdToExec)
	}
	if maxRes != 5 {
		t.Errorf("maxRes: got %d", maxRes)
	}
	if tokenFile != "/tmp/token.json" {
		t.Errorf("tokenFile: got %s", tokenFile)
	}
	if credentialsFile != "/tmp/creds.json" {
		t.Errorf("credentialsFile: got %s", credentialsFile)
	}
	if backupFile != "/tmp/backup.txt" {
		t.Errorf("backupFile: got %s", backupFile)
	}
	if triggerBeforeMinutes != 5 {
		t.Errorf("triggerBeforeMinutes: got %d", triggerBeforeMinutes)
	}
	if cronMarker != "# marker" {
		t.Errorf("cronMarker: got %s", cronMarker)
	}
	if showDeleted {
		t.Error("showDeleted should be false")
	}
}
