package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetHomeDir(t *testing.T) {
	dir := GetHomeDir()
	if dir == "" {
		t.Error("expected non-empty home directory")
	}
	// Should be an absolute path
	if !filepath.IsAbs(dir) {
		t.Errorf("expected absolute path, got %s", dir)
	}
}

func TestExpandPath(t *testing.T) {
	homeDir := GetHomeDir()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "expand tilde with path",
			input:    "~/some/path",
			expected: homeDir + "/some/path",
		},
		{
			name:     "expand tilde only",
			input:    "~",
			expected: homeDir,
		},
		{
			name:     "absolute path unchanged",
			input:    "/absolute/path",
			expected: "/absolute/path",
		},
		{
			name:     "relative path unchanged",
			input:    "relative/path",
			expected: "relative/path",
		},
		{
			name:     "empty string unchanged",
			input:    "",
			expected: "",
		},
		{
			name:     "tilde in middle unchanged",
			input:    "/path/~/something",
			expected: "/path/~/something",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExpandPath(tt.input)
			if result != tt.expected {
				t.Errorf("ExpandPath(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestValidateConfig_RequiredFields(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid config",
			config: Config{
				CalendarID: "test@calendar.google.com",
				CmdToExec:  "echo hello",
				TokenFile:  "/path/to/token.json",
			},
			expectError: false,
		},
		{
			name: "missing CalendarID",
			config: Config{
				CmdToExec: "echo hello",
				TokenFile: "/path/to/token.json",
			},
			expectError: true,
			errorMsg:    "CalendarID is required",
		},
		{
			name: "missing CmdToExec",
			config: Config{
				CalendarID: "test@calendar.google.com",
				TokenFile:  "/path/to/token.json",
			},
			expectError: true,
			errorMsg:    "CmdToExec is required",
		},
		{
			name: "missing TokenFile",
			config: Config{
				CalendarID: "test@calendar.google.com",
				CmdToExec:  "echo hello",
			},
			expectError: true,
			errorMsg:    "TokenFile is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfig(&tt.config)
			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}

func TestValidateConfig_MaxResClamping(t *testing.T) {
	tests := []struct {
		name        string
		inputMaxRes int64
		wantMaxRes  int64
	}{
		{
			name:        "valid MaxRes",
			inputMaxRes: 50,
			wantMaxRes:  50,
		},
		{
			name:        "MaxRes too low clamped to 1",
			inputMaxRes: 0,
			wantMaxRes:  1,
		},
		{
			name:        "negative MaxRes clamped to 1",
			inputMaxRes: -5,
			wantMaxRes:  1,
		},
		{
			name:        "MaxRes too high clamped to 100",
			inputMaxRes: 150,
			wantMaxRes:  100,
		},
		{
			name:        "MaxRes at lower bound",
			inputMaxRes: 1,
			wantMaxRes:  1,
		},
		{
			name:        "MaxRes at upper bound",
			inputMaxRes: 100,
			wantMaxRes:  100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := Config{
				CalendarID: "test@calendar.google.com",
				CmdToExec:  "echo hello",
				TokenFile:  "/path/to/token.json",
				MaxRes:     tt.inputMaxRes,
			}
			err := ValidateConfig(&config)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if config.MaxRes != tt.wantMaxRes {
				t.Errorf("MaxRes = %d, want %d", config.MaxRes, tt.wantMaxRes)
			}
		})
	}
}

func TestValidateConfig_TriggerBeforeMinutes(t *testing.T) {
	tests := []struct {
		name                       string
		inputTriggerBeforeMinutes  int
		wantTriggerBeforeMinutes   int
	}{
		{
			name:                       "positive value unchanged",
			inputTriggerBeforeMinutes:  10,
			wantTriggerBeforeMinutes:   10,
		},
		{
			name:                       "zero unchanged",
			inputTriggerBeforeMinutes:  0,
			wantTriggerBeforeMinutes:   0,
		},
		{
			name:                       "negative clamped to 0",
			inputTriggerBeforeMinutes:  -5,
			wantTriggerBeforeMinutes:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := Config{
				CalendarID:           "test@calendar.google.com",
				CmdToExec:            "echo hello",
				TokenFile:            "/path/to/token.json",
				TriggerBeforeMinutes: tt.inputTriggerBeforeMinutes,
			}
			err := ValidateConfig(&config)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if config.TriggerBeforeMinutes != tt.wantTriggerBeforeMinutes {
				t.Errorf("TriggerBeforeMinutes = %d, want %d", config.TriggerBeforeMinutes, tt.wantTriggerBeforeMinutes)
			}
		})
	}
}

func TestValidateConfig_PathExpansion(t *testing.T) {
	homeDir := GetHomeDir()

	config := Config{
		CalendarID:      "test@calendar.google.com",
		CmdToExec:       "echo hello",
		TokenFile:       "~/token.json",
		CredentialsFile: "~/credentials.json",
		BackupFile:      "~/backup.txt",
	}

	err := ValidateConfig(&config)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if config.TokenFile != homeDir+"/token.json" {
		t.Errorf("TokenFile not expanded, got %s", config.TokenFile)
	}
	if config.CredentialsFile != homeDir+"/credentials.json" {
		t.Errorf("CredentialsFile not expanded, got %s", config.CredentialsFile)
	}
	if config.BackupFile != homeDir+"/backup.txt" {
		t.Errorf("BackupFile not expanded, got %s", config.BackupFile)
	}
}

func TestReadConfig_MissingFile(t *testing.T) {
	// Save and restore HOME
	origHome := os.Getenv("HOME")
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	// Try to read config without creating the file
	_, err := ReadConfig()
	if err == nil {
		t.Error("expected error for missing config file, got nil")
	}
}

func TestReadConfig_ValidConfig(t *testing.T) {
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
CalendarID = "test@calendar.google.com"
CmdToExec = "echo hello"
TokenFile = "/path/to/token.json"
MaxRes = 10
`
	configPath := filepath.Join(configDir, "config.toml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	v, err := ReadConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if v.GetString("CalendarID") != "test@calendar.google.com" {
		t.Errorf("CalendarID mismatch, got %s", v.GetString("CalendarID"))
	}
	if v.GetString("CmdToExec") != "echo hello" {
		t.Errorf("CmdToExec mismatch, got %s", v.GetString("CmdToExec"))
	}
	if v.GetInt64("MaxRes") != 10 {
		t.Errorf("MaxRes mismatch, got %d", v.GetInt64("MaxRes"))
	}
}

func TestReadConfig_DefaultValues(t *testing.T) {
	// Save and restore HOME
	origHome := os.Getenv("HOME")
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	// Create config directory and file with minimal config
	configDir := filepath.Join(tmpDir, ".config", "timeotter")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}

	configContent := `
CalendarID = "test@calendar.google.com"
CmdToExec = "echo hello"
TokenFile = "/path/to/token.json"
`
	configPath := filepath.Join(configDir, "config.toml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	v, err := ReadConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check default values
	if v.GetInt64("MaxRes") != 5 {
		t.Errorf("default MaxRes should be 5, got %d", v.GetInt64("MaxRes"))
	}
	if v.GetInt("TriggerBeforeMinutes") != 5 {
		t.Errorf("default TriggerBeforeMinutes should be 5, got %d", v.GetInt("TriggerBeforeMinutes"))
	}
	if !v.GetBool("ShowDeleted") == true {
		// ShowDeleted default is false, so !false == true
	}
	if v.GetString("CronMarker") != "# custom crons below this can be deleted." {
		t.Errorf("default CronMarker mismatch, got %s", v.GetString("CronMarker"))
	}
}
