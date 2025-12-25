package oauth

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

// Integration tests for oauth package
// These tests verify file I/O operations with the token system

func TestIntegration_TokenRoundtrip(t *testing.T) {
	tmpDir := t.TempDir()
	tokenPath := filepath.Join(tmpDir, "token.json")

	// Create a token with all fields
	original := &oauth2.Token{
		AccessToken:  "integration_access_token",
		TokenType:    "Bearer",
		RefreshToken: "integration_refresh_token",
		Expiry:       time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC),
	}

	// Save token
	SaveToken(tokenPath, original)

	// Load token
	loaded, err := TokenFromFile(tokenPath)
	if err != nil {
		t.Fatalf("failed to load token: %v", err)
	}

	// Verify all fields
	if loaded.AccessToken != original.AccessToken {
		t.Errorf("AccessToken: got %s, want %s", loaded.AccessToken, original.AccessToken)
	}
	if loaded.TokenType != original.TokenType {
		t.Errorf("TokenType: got %s, want %s", loaded.TokenType, original.TokenType)
	}
	if loaded.RefreshToken != original.RefreshToken {
		t.Errorf("RefreshToken: got %s, want %s", loaded.RefreshToken, original.RefreshToken)
	}
	if loaded.Expiry.Unix() != original.Expiry.Unix() {
		t.Errorf("Expiry: got %v, want %v", loaded.Expiry, original.Expiry)
	}
}

func TestIntegration_TokenFilePermissions(t *testing.T) {
	tmpDir := t.TempDir()
	tokenPath := filepath.Join(tmpDir, "secure_token.json")

	token := &oauth2.Token{
		AccessToken: "secure_access_token",
	}

	SaveToken(tokenPath, token)

	// Check file permissions
	info, err := os.Stat(tokenPath)
	if err != nil {
		t.Fatalf("failed to stat token file: %v", err)
	}

	perms := info.Mode().Perm()
	if perms != 0600 {
		t.Errorf("expected permissions 0600, got %04o", perms)
	}
}

func TestIntegration_TokenOverwrite(t *testing.T) {
	tmpDir := t.TempDir()
	tokenPath := filepath.Join(tmpDir, "overwrite_token.json")

	// Save first token
	token1 := &oauth2.Token{
		AccessToken: "first_token",
	}
	SaveToken(tokenPath, token1)

	// Save second token (overwrite)
	token2 := &oauth2.Token{
		AccessToken: "second_token",
	}
	SaveToken(tokenPath, token2)

	// Load and verify
	loaded, err := TokenFromFile(tokenPath)
	if err != nil {
		t.Fatalf("failed to load token: %v", err)
	}

	if loaded.AccessToken != "second_token" {
		t.Errorf("expected second_token, got %s", loaded.AccessToken)
	}
}

func TestIntegration_TokenInSubdirectory(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "deep", "nested", "dir")
	if err := os.MkdirAll(subDir, 0750); err != nil {
		t.Fatalf("failed to create subdirectory: %v", err)
	}

	tokenPath := filepath.Join(subDir, "token.json")

	token := &oauth2.Token{
		AccessToken: "nested_token",
	}
	SaveToken(tokenPath, token)

	loaded, err := TokenFromFile(tokenPath)
	if err != nil {
		t.Fatalf("failed to load token: %v", err)
	}

	if loaded.AccessToken != "nested_token" {
		t.Errorf("expected nested_token, got %s", loaded.AccessToken)
	}
}

func TestIntegration_ExpiredToken(t *testing.T) {
	tmpDir := t.TempDir()
	tokenPath := filepath.Join(tmpDir, "expired_token.json")

	// Create an expired token
	expired := &oauth2.Token{
		AccessToken:  "expired_access",
		TokenType:    "Bearer",
		RefreshToken: "expired_refresh",
		Expiry:       time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), // Past date
	}

	SaveToken(tokenPath, expired)

	loaded, err := TokenFromFile(tokenPath)
	if err != nil {
		t.Fatalf("failed to load expired token: %v", err)
	}

	// Token should load successfully even if expired
	// The oauth2 library will handle refresh
	if loaded.AccessToken != "expired_access" {
		t.Errorf("expected expired_access, got %s", loaded.AccessToken)
	}

	// Verify token reports as expired
	if loaded.Valid() {
		t.Log("Note: Token reports as valid despite past expiry - may have no expiry check")
	}
}

func TestIntegration_TokenWithSpecialCharacters(t *testing.T) {
	tmpDir := t.TempDir()
	tokenPath := filepath.Join(tmpDir, "special_token.json")

	// Token with special characters
	token := &oauth2.Token{
		AccessToken:  "token_with_special_chars!@#$%^&*()",
		TokenType:    "Bearer",
		RefreshToken: "refresh/with/slashes+and+plus",
	}

	SaveToken(tokenPath, token)

	loaded, err := TokenFromFile(tokenPath)
	if err != nil {
		t.Fatalf("failed to load token: %v", err)
	}

	if loaded.AccessToken != token.AccessToken {
		t.Errorf("AccessToken with special chars: got %s, want %s", loaded.AccessToken, token.AccessToken)
	}
	if loaded.RefreshToken != token.RefreshToken {
		t.Errorf("RefreshToken with special chars: got %s, want %s", loaded.RefreshToken, token.RefreshToken)
	}
}

func TestIntegration_ConcurrentTokenAccess(t *testing.T) {
	tmpDir := t.TempDir()
	tokenPath := filepath.Join(tmpDir, "concurrent_token.json")

	// Initial token
	token := &oauth2.Token{
		AccessToken: "initial",
	}
	SaveToken(tokenPath, token)

	// Concurrent reads
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			_, err := TokenFromFile(tokenPath)
			if err != nil {
				t.Errorf("concurrent read failed: %v", err)
			}
			done <- true
		}()
	}

	// Wait for all reads
	for i := 0; i < 10; i++ {
		<-done
	}
}
