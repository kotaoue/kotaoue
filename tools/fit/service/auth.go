package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const tokenURL = "https://oauth2.googleapis.com/token"

// tokenResponse holds the OAuth2 token refresh response.
type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// refreshAccessToken exchanges a refresh token for an access token.
func refreshAccessToken(clientID, clientSecret, refreshToken string) (string, error) {
	form := url.Values{}
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)
	form.Set("refresh_token", refreshToken)
	form.Set("grant_type", "refresh_token")

	resp, err := http.PostForm(tokenURL, form) // #nosec G107 -- URL is a fixed constant
	if err != nil {
		return "", fmt.Errorf("failed to request token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read token response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token endpoint returned status %d: %s", resp.StatusCode, body)
	}

	var tr tokenResponse
	if err := json.Unmarshal(body, &tr); err != nil {
		return "", fmt.Errorf("failed to parse token response: %w", err)
	}

	return tr.AccessToken, nil
}
