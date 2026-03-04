package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const tokenEndpointURL = "https://oauth2.googleapis.com/token"

// tokenResponse holds the OAuth2 token refresh response.
type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// tokenErrorResponse holds an OAuth2 error response.
type tokenErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// ErrTokenExpired is returned when the refresh token has been expired or revoked.
var ErrTokenExpired = fmt.Errorf("refresh token has expired or been revoked; run tools/prepareFit to obtain a new token and update GOOGLE_FIT_REFRESH_TOKEN in GitHub Secrets")

// refreshAccessToken exchanges a refresh token for an access token.
func refreshAccessToken(clientID, clientSecret, refreshToken string) (string, error) {
	return refreshAccessTokenWithURL(tokenEndpointURL, clientID, clientSecret, refreshToken)
}

// refreshAccessTokenWithURL is the internal implementation used by refreshAccessToken and tests.
func refreshAccessTokenWithURL(tokenURL, clientID, clientSecret, refreshToken string) (string, error) {
	form := url.Values{}
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)
	form.Set("refresh_token", refreshToken)
	form.Set("grant_type", "refresh_token")

	resp, err := http.PostForm(tokenURL, form) // #nosec G107 -- tokenURL is either the fixed constant or a test server URL
	if err != nil {
		return "", fmt.Errorf("failed to request token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read token response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp tokenErrorResponse
		if json.Unmarshal(body, &errResp) == nil && errResp.Error == "invalid_grant" {
			return "", ErrTokenExpired
		}
		return "", fmt.Errorf("token endpoint returned status %d: %s", resp.StatusCode, body)
	}

	var tr tokenResponse
	if err := json.Unmarshal(body, &tr); err != nil {
		return "", fmt.Errorf("failed to parse token response: %w", err)
	}

	return tr.AccessToken, nil
}
