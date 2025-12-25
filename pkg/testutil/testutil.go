package testutil

import (
	"os"
	"path/filepath"
	"testing"

	"google.golang.org/api/calendar/v3"
)

// CreateTempConfig creates a temporary config directory with a valid TOML config file.
// Returns the temp directory path and a cleanup function.
func CreateTempConfig(t *testing.T, content string) string {
	t.Helper()

	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".config", "timeotter")
	if err := os.MkdirAll(configDir, 0750); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}

	configPath := filepath.Join(configDir, "config.toml")
	if err := os.WriteFile(configPath, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	return tmpDir
}

// CreateTempTokenFile creates a temporary OAuth token file.
func CreateTempTokenFile(t *testing.T, content string) string {
	t.Helper()

	tmpDir := t.TempDir()
	tokenPath := filepath.Join(tmpDir, "token.json")
	if err := os.WriteFile(tokenPath, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write token file: %v", err)
	}

	return tokenPath
}

// ValidConfigTOML returns a minimal valid config TOML string.
func ValidConfigTOML() string {
	return `
CalendarID = "test@calendar.google.com"
CmdToExec = "echo hello"
TokenFile = "/tmp/token.json"
MaxRes = 5
`
}

// FullConfigTOML returns a config TOML with all fields set.
func FullConfigTOML() string {
	return `
CalendarID = "test@calendar.google.com"
CmdToExec = "echo hello"
TokenFile = "/tmp/token.json"
CredentialsFile = "/tmp/credentials.json"
BackupFile = "/tmp/backup.txt"
MaxRes = 10
TriggerBeforeMinutes = 10
CronMarker = "# timeotter crons"
ShowDeleted = true
`
}

// ValidTokenJSON returns a valid OAuth token JSON string.
func ValidTokenJSON() string {
	return `{
	"access_token": "test_access_token",
	"token_type": "Bearer",
	"refresh_token": "test_refresh_token",
	"expiry": "2025-12-31T23:59:59Z"
}`
}

// MockCalendarEvent creates a mock calendar event with the given summary and start time.
func MockCalendarEvent(summary, dateTime string) *calendar.Event {
	return &calendar.Event{
		Summary: summary,
		Start: &calendar.EventDateTime{
			DateTime: dateTime,
		},
	}
}

// MockCalendarEventAllDay creates a mock all-day calendar event.
func MockCalendarEventAllDay(summary, date string) *calendar.Event {
	return &calendar.Event{
		Summary: summary,
		Start: &calendar.EventDateTime{
			Date: date,
		},
	}
}

// MockCalendarEvents creates a mock calendar.Events with the given events.
func MockCalendarEvents(events ...*calendar.Event) *calendar.Events {
	return &calendar.Events{
		Items: events,
	}
}

// SetHomeDir temporarily sets HOME environment variable for the test.
// Uses t.Setenv which automatically restores the value after the test.
func SetHomeDir(t *testing.T, dir string) {
	t.Helper()
	t.Setenv("HOME", dir)
}
