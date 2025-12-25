package cron

import (
	"runtime"
	"strings"
	"testing"
)

func TestGetCronLocation(t *testing.T) {
	location := GetCronLocation()

	// Should return a non-empty string
	if location == "" {
		t.Error("expected non-empty cron location")
	}

	// Should contain the username
	if runtime.GOOS == "linux" {
		if !strings.HasPrefix(location, "/var/spool/cron/") {
			t.Errorf("expected Linux cron path, got %s", location)
		}
	}

	// Should contain a username (path should have more than just the directory)
	parts := strings.Split(location, "/")
	lastPart := parts[len(parts)-1]
	if lastPart == "" || lastPart == "cron" {
		t.Errorf("expected username in path, got %s", location)
	}
}

func TestEscapeAwkRegex(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "plain text unchanged",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "escape period",
			input:    "file.txt",
			expected: "file\\.txt",
		},
		{
			name:     "escape asterisk",
			input:    "*.go",
			expected: "\\*\\.go",
		},
		{
			name:     "escape plus",
			input:    "a+b",
			expected: "a\\+b",
		},
		{
			name:     "escape question mark",
			input:    "what?",
			expected: "what\\?",
		},
		{
			name:     "escape square brackets",
			input:    "[abc]",
			expected: "\\[abc\\]",
		},
		{
			name:     "escape caret",
			input:    "^start",
			expected: "\\^start",
		},
		{
			name:     "escape dollar",
			input:    "end$",
			expected: "end\\$",
		},
		{
			name:     "escape parentheses",
			input:    "(group)",
			expected: "\\(group\\)",
		},
		{
			name:     "escape curly braces",
			input:    "{1,2}",
			expected: "\\{1,2\\}",
		},
		{
			name:     "escape pipe",
			input:    "a|b",
			expected: "a\\|b",
		},
		{
			name:     "escape backslash",
			input:    "path\\to",
			expected: "path\\\\to",
		},
		{
			name:     "cron marker string",
			input:    "# custom crons below this can be deleted.",
			expected: "# custom crons below this can be deleted\\.",
		},
		{
			name:     "multiple special chars",
			input:    "[*+?]",
			expected: "\\[\\*\\+\\?\\]",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := escapeAwkRegex(tt.input)
			if result != tt.expected {
				t.Errorf("escapeAwkRegex(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetCronLocation_ConsistentReturns(t *testing.T) {
	// Call multiple times and verify consistency
	loc1 := GetCronLocation()
	loc2 := GetCronLocation()

	if loc1 != loc2 {
		t.Errorf("GetCronLocation returned inconsistent values: %s vs %s", loc1, loc2)
	}
}

func TestEscapeAwkRegex_Idempotent(t *testing.T) {
	// Escaping already escaped strings should escape the backslashes
	input := "a.b"
	escaped := escapeAwkRegex(input)
	doubleEscaped := escapeAwkRegex(escaped)

	// Double escaping should escape the backslash from first escape
	if doubleEscaped != "a\\\\\\.b" {
		t.Errorf("double escape failed: got %q", doubleEscaped)
	}
}
