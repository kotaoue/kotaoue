package service

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	startMarker = "<!-- PEDOMETER_START -->"
	endMarker   = "<!-- PEDOMETER_END -->"
	tokenURL    = "https://oauth2.googleapis.com/token"
	fitnessURL  = "https://www.googleapis.com/fitness/v1/users/me/dataset:aggregate"
)

// tokenResponse holds the OAuth2 token refresh response.
type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// aggregateRequest is the request body for the Google Fit aggregate API.
type aggregateRequest struct {
	AggregateBy  []aggregateBy `json:"aggregateBy"`
	BucketByTime bucketByTime  `json:"bucketByTime"`
	StartTimeMs  int64         `json:"startTimeMillis,string"`
	EndTimeMs    int64         `json:"endTimeMillis,string"`
}

type aggregateBy struct {
	DataTypeName string `json:"dataTypeName"`
}

type bucketByTime struct {
	DurationMillis int64 `json:"durationMillis"`
}

// aggregateResponse is the response from the Google Fit aggregate API.
type aggregateResponse struct {
	Bucket []struct {
		Dataset []struct {
			Point []struct {
				Value []struct {
					IntVal int `json:"intVal"`
				} `json:"value"`
			} `json:"point"`
		} `json:"dataset"`
	} `json:"bucket"`
}

// RunUpdatePedometer fetches yesterday's step count from Google Fit and updates README.md.
func RunUpdatePedometer(args []string) error {
	fs := flag.NewFlagSet("update-pedometer", flag.ExitOnError)
	readmeFile := fs.String("readme", "README.md", "Path to README.md")
	if err := fs.Parse(args); err != nil {
		return err
	}

	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	refreshToken := os.Getenv("GOOGLE_REFRESH_TOKEN")

	if clientID == "" || clientSecret == "" || refreshToken == "" {
		return fmt.Errorf("GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET, and GOOGLE_REFRESH_TOKEN environment variables are required")
	}

	accessToken, err := refreshAccessToken(clientID, clientSecret, refreshToken)
	if err != nil {
		return fmt.Errorf("failed to refresh access token: %w", err)
	}

	steps, err := fetchYesterdaySteps(accessToken)
	if err != nil {
		return fmt.Errorf("failed to fetch steps: %w", err)
	}

	yesterday := time.Now().In(time.FixedZone("JST", 9*60*60)).AddDate(0, 0, -1)
	stepsText := fmt.Sprintf("%d月%d日の歩数: %d歩", yesterday.Month(), yesterday.Day(), steps)

	content, err := os.ReadFile(*readmeFile)
	if err != nil {
		return fmt.Errorf("failed to read README file: %w", err)
	}

	newContent, err := replaceBetweenMarkers(string(content), startMarker, endMarker, stepsText)
	if err != nil {
		return fmt.Errorf("failed to replace content: %w", err)
	}

	if err := os.WriteFile(*readmeFile, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write README file: %w", err)
	}

	log.Printf("Updated %s with: %s", *readmeFile, stepsText)
	return nil
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

// fetchYesterdaySteps returns the total step count for yesterday using the Google Fit API.
func fetchYesterdaySteps(accessToken string) (int, error) {
	jst := time.FixedZone("JST", 9*60*60)
	now := time.Now().In(jst)
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, jst)
	yesterdayStart := todayStart.AddDate(0, 0, -1)

	reqBody := aggregateRequest{
		AggregateBy:  []aggregateBy{{DataTypeName: "com.google.step_count.delta"}},
		BucketByTime: bucketByTime{DurationMillis: 86400000},
		StartTimeMs:  yesterdayStart.UnixMilli(),
		EndTimeMs:    todayStart.UnixMilli(),
	}

	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, fitnessURL, bytes.NewReader(reqJSON))
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to call fitness API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read fitness response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("fitness API returned status %d: %s", resp.StatusCode, body)
	}

	var ar aggregateResponse
	if err := json.Unmarshal(body, &ar); err != nil {
		return 0, fmt.Errorf("failed to parse fitness response: %w", err)
	}

	total := 0
	for _, bucket := range ar.Bucket {
		for _, dataset := range bucket.Dataset {
			for _, point := range dataset.Point {
				for _, v := range point.Value {
					total += v.IntVal
				}
			}
		}
	}

	return total, nil
}

// replaceBetweenMarkers replaces content between start and end markers.
func replaceBetweenMarkers(content, start, end, replacement string) (string, error) {
	startIdx := strings.Index(content, start)
	endIdx := strings.Index(content, end)
	if startIdx == -1 || endIdx == -1 {
		return "", fmt.Errorf("markers not found in content: %q, %q", start, end)
	}
	if startIdx >= endIdx {
		return "", fmt.Errorf("start marker must appear before end marker")
	}
	return content[:startIdx+len(start)] + replacement + content[endIdx:], nil
}
