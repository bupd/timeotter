// Package oauth provides OAuth2 authentication utilities for Google Calendar.
package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
)

// GetClient retrieves a token, saves it, and returns the generated HTTP client.
func GetClient(config *oauth2.Config, tokFile string) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	// tokFile := "~/.cal-token.json"
	tok, err := TokenFromFile(tokFile)
	if err != nil {
		tok = GetTokenFromWeb(config)
		SaveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// GetTokenFromWeb requests a token via browser authorization and returns it.
func GetTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// TokenFromFile retrieves a token from a local file.
func TokenFromFile(file string) (tok *oauth2.Token, err error) {
	f, err := os.Open(filepath.Clean(file))
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()
	tok = &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// SaveToken saves a token to a file path.
func SaveToken(path string, token *oauth2.Token) {
	cleanPath := filepath.Clean(path)
	fmt.Printf("Saving credential file to: %s\n", cleanPath)
	f, err := os.OpenFile(cleanPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	if err := json.NewEncoder(f).Encode(token); err != nil {
		closeErr := f.Close()
		if closeErr != nil {
			log.Fatalf("Unable to encode oauth token: %v (close error: %v)", err, closeErr)
		}
		log.Fatalf("Unable to encode oauth token: %v", err)
	}
	if err := f.Close(); err != nil {
		log.Fatalf("Unable to close token file: %v", err)
	}
}
