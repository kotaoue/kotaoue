package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	dateStartMarker  = "<!-- PEDOMETER_DATE_START -->"
	dateEndMarker    = "<!-- PEDOMETER_DATE_END -->"
	stepsStartMarker = "<!-- PEDOMETER_STEPS_START -->"
	stepsEndMarker   = "<!-- PEDOMETER_STEPS_END -->"
	fitnessURL       = "https://www.googleapis.com/fitness/v1/users/me/dataset:aggregate"
)

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

// RunUpdatePedometer fetches yesterday's step count from Google Fit and updates the given README file.
func RunUpdatePedometer(readmeFile string) error {
	credJSON := os.Getenv("GOOGLE_CLOUD_CREDENTIALS_JSON")
	if credJSON == "" {
		return fmt.Errorf("GOOGLE_CLOUD_CREDENTIALS_JSON environment variable is required")
	}

	accessToken, err := credentialsToAccessToken([]byte(credJSON))
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	steps, err := fetchYesterdaySteps(accessToken)
	if err != nil {
		return fmt.Errorf("failed to fetch steps: %w", err)
	}

	yesterday := time.Now().In(time.FixedZone("JST", 9*60*60)).AddDate(0, 0, -1)
	dateText := fmt.Sprintf("%d月%d日の歩数", yesterday.Month(), yesterday.Day())
	stepsText := fmt.Sprintf("%s歩", formatWithCommas(steps))

	content, err := os.ReadFile(readmeFile)
	if err != nil {
		return fmt.Errorf("failed to read README file: %w", err)
	}

	newContent, err := replaceBetweenMarkers(string(content), dateStartMarker, dateEndMarker, dateText)
	if err != nil {
		return fmt.Errorf("failed to replace date content: %w", err)
	}

	newContent, err = replaceBetweenMarkers(newContent, stepsStartMarker, stepsEndMarker, stepsText)
	if err != nil {
		return fmt.Errorf("failed to replace steps content: %w", err)
	}

	if err := os.WriteFile(readmeFile, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write README file: %w", err)
	}

	log.Printf("Updated %s with: %s %s", readmeFile, dateText, stepsText)
	return nil
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

// formatWithCommas formats an integer with comma separators (e.g. 12345 -> "12,345").
func formatWithCommas(n int) string {
	s := fmt.Sprintf("%d", n)
	neg := n < 0
	if neg {
		s = s[1:]
	}
	var b strings.Builder
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(c)
	}
	if neg {
		return "-" + b.String()
	}
	return b.String()
}
