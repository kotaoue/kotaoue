package service

import (
	"context"
	"fmt"

	"golang.org/x/oauth2/google"
)

const fitnessActivityReadScope = "https://www.googleapis.com/auth/fitness.activity.read"

// credentialsToAccessToken loads OAuth2 credentials from a JSON byte slice
// (Google's authorized_user ADC format) and returns a valid access token.
func credentialsToAccessToken(credJSON []byte) (string, error) {
	ctx := context.Background()
	creds, err := google.CredentialsFromJSON(ctx, credJSON, fitnessActivityReadScope)
	if err != nil {
		return "", fmt.Errorf("failed to load credentials: %w", err)
	}

	token, err := creds.TokenSource.Token()
	if err != nil {
		return "", fmt.Errorf("failed to get token: %w", err)
	}

	return token.AccessToken, nil
}
