package service

import (
	"context"
	"encoding/json"
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const fitnessActivityReadScope = "https://www.googleapis.com/auth/fitness.activity.read"

// authorizedUserCredentials represents the authorized_user ADC JSON format.
type authorizedUserCredentials struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token"`
	Type         string `json:"type"`
}

// credentialsToAccessToken loads OAuth2 credentials from a JSON byte slice
// (Google's authorized_user ADC format) and returns a valid access token.
func credentialsToAccessToken(credJSON []byte) (string, error) {
	var creds authorizedUserCredentials
	if err := json.Unmarshal(credJSON, &creds); err != nil {
		return "", fmt.Errorf("failed to parse credentials: %w", err)
	}

	if creds.ClientID == "" || creds.ClientSecret == "" || creds.RefreshToken == "" {
		return "", fmt.Errorf("credentials JSON must contain client_id, client_secret, and refresh_token")
	}

	config := &oauth2.Config{
		ClientID:     creds.ClientID,
		ClientSecret: creds.ClientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{fitnessActivityReadScope},
	}

	// Providing only the RefreshToken causes TokenSource to exchange it for a
	// new access token on the first call to Token().
	token := &oauth2.Token{RefreshToken: creds.RefreshToken}
	tokenSource := config.TokenSource(context.Background(), token)

	newToken, err := tokenSource.Token()
	if err != nil {
		return "", fmt.Errorf("failed to refresh token: %w", err)
	}

	return newToken.AccessToken, nil
}
