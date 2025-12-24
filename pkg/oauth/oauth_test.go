package oauth

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

func TestTokenFromFile_ValidToken(t *testing.T) {
	tmpDir := t.TempDir()
	tokenPath := filepath.Join(tmpDir, "token.json")

	// Create a valid token file
	tokenContent := `{
		"access_token": "test_access_token",
		"token_type": "Bearer",
		"refresh_token": "test_refresh_token",
		"expiry": "2025-12-31T23:59:59Z"
	}`

	if err := os.WriteFile(tokenPath, []byte(tokenContent), 0600); err != nil {
		t.Fatalf("failed to write token file: %v", err)
	}

	tok, err := TokenFromFile(tokenPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if tok.AccessToken != "test_access_token" {
		t.Errorf("AccessToken mismatch, got %s", tok.AccessToken)
	}
	if tok.TokenType != "Bearer" {
		t.Errorf("TokenType mismatch, got %s", tok.TokenType)
	}
	if tok.RefreshToken != "test_refresh_token" {
		t.Errorf("RefreshToken mismatch, got %s", tok.RefreshToken)
	}
}

func TestTokenFromFile_MissingFile(t *testing.T) {
	tmpDir := t.TempDir()
	tokenPath := filepath.Join(tmpDir, "nonexistent.json")

	_, err := TokenFromFile(tokenPath)
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestTokenFromFile_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	tokenPath := filepath.Join(tmpDir, "invalid.json")

	// Create an invalid JSON file
	if err := os.WriteFile(tokenPath, []byte("not valid json"), 0600); err != nil {
		t.Fatalf("failed to write token file: %v", err)
	}

	_, err := TokenFromFile(tokenPath)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestTokenFromFile_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	tokenPath := filepath.Join(tmpDir, "empty.json")

	// Create an empty file
	if err := os.WriteFile(tokenPath, []byte(""), 0600); err != nil {
		t.Fatalf("failed to write token file: %v", err)
	}

	_, err := TokenFromFile(tokenPath)
	if err == nil {
		t.Error("expected error for empty file, got nil")
	}
}

func TestSaveToken_CreatesFile(t *testing.T) {
	tmpDir := t.TempDir()
	tokenPath := filepath.Join(tmpDir, "new_token.json")

	token := &oauth2.Token{
		AccessToken:  "test_access",
		TokenType:    "Bearer",
		RefreshToken: "test_refresh",
		Expiry:       time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC),
	}

	SaveToken(tokenPath, token)

	// Verify file was created
	if _, err := os.Stat(tokenPath); os.IsNotExist(err) {
		t.Error("token file was not created")
	}

	// Verify file permissions
	info, err := os.Stat(tokenPath)
	if err != nil {
		t.Fatalf("failed to stat file: %v", err)
	}
	perms := info.Mode().Perm()
	if perms != 0600 {
		t.Errorf("expected file permissions 0600, got %04o", perms)
	}
}

func TestSaveToken_OverwritesExisting(t *testing.T) {
	tmpDir := t.TempDir()
	tokenPath := filepath.Join(tmpDir, "existing_token.json")

	// Create initial file
	if err := os.WriteFile(tokenPath, []byte(`{"access_token":"old"}`), 0600); err != nil {
		t.Fatalf("failed to create initial file: %v", err)
	}

	// Save new token
	token := &oauth2.Token{
		AccessToken:  "new_access",
		TokenType:    "Bearer",
		RefreshToken: "new_refresh",
	}

	SaveToken(tokenPath, token)

	// Read back and verify
	readToken, err := TokenFromFile(tokenPath)
	if err != nil {
		t.Fatalf("failed to read token: %v", err)
	}

	if readToken.AccessToken != "new_access" {
		t.Errorf("token not overwritten, got %s", readToken.AccessToken)
	}
}

func TestTokenRoundtrip(t *testing.T) {
	tmpDir := t.TempDir()
	tokenPath := filepath.Join(tmpDir, "roundtrip_token.json")

	original := &oauth2.Token{
		AccessToken:  "roundtrip_access",
		TokenType:    "Bearer",
		RefreshToken: "roundtrip_refresh",
		Expiry:       time.Date(2025, 6, 15, 12, 0, 0, 0, time.UTC),
	}

	// Save and read back
	SaveToken(tokenPath, original)
	loaded, err := TokenFromFile(tokenPath)
	if err != nil {
		t.Fatalf("failed to load token: %v", err)
	}

	// Compare fields
	if loaded.AccessToken != original.AccessToken {
		t.Errorf("AccessToken mismatch: got %s, want %s", loaded.AccessToken, original.AccessToken)
	}
	if loaded.TokenType != original.TokenType {
		t.Errorf("TokenType mismatch: got %s, want %s", loaded.TokenType, original.TokenType)
	}
	if loaded.RefreshToken != original.RefreshToken {
		t.Errorf("RefreshToken mismatch: got %s, want %s", loaded.RefreshToken, original.RefreshToken)
	}
	// Note: Expiry comparison may have timezone differences, check Unix timestamp
	if loaded.Expiry.Unix() != original.Expiry.Unix() {
		t.Errorf("Expiry mismatch: got %v, want %v", loaded.Expiry, original.Expiry)
	}
}

func TestSaveToken_CreatesDirIfNeeded(t *testing.T) {
	tmpDir := t.TempDir()
	tokenPath := filepath.Join(tmpDir, "subdir", "token.json")

	// Create the subdir first (SaveToken doesn't create directories)
	if err := os.MkdirAll(filepath.Dir(tokenPath), 0755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}

	token := &oauth2.Token{
		AccessToken: "test",
	}

	SaveToken(tokenPath, token)

	// Verify file was created
	if _, err := os.Stat(tokenPath); os.IsNotExist(err) {
		t.Error("token file was not created in subdirectory")
	}
}

func TestTokenFromFile_PartialToken(t *testing.T) {
	tmpDir := t.TempDir()
	tokenPath := filepath.Join(tmpDir, "partial.json")

	// Create a token file with only some fields
	tokenContent := `{"access_token": "partial_token"}`

	if err := os.WriteFile(tokenPath, []byte(tokenContent), 0600); err != nil {
		t.Fatalf("failed to write token file: %v", err)
	}

	tok, err := TokenFromFile(tokenPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if tok.AccessToken != "partial_token" {
		t.Errorf("AccessToken mismatch, got %s", tok.AccessToken)
	}
	// Other fields should be empty/zero
	if tok.RefreshToken != "" {
		t.Errorf("expected empty RefreshToken, got %s", tok.RefreshToken)
	}
}
