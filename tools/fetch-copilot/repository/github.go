package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kotaoue/kotaoue/tools/fetch-copilot/entity"
)

const copilotMetricsURL = "https://api.github.com/orgs/%s/copilot/metrics"

// metricsResponse is the JSON response from the GitHub Copilot metrics API.
type metricsResponse struct {
	Date              string `json:"date"`
	TotalEngagedUsers int    `json:"total_engaged_users"`
	CopilotIDEChat    struct {
		TotalEngagedUsers int `json:"total_engaged_users"`
		Editors           []struct {
			Models []struct {
				TotalChats               int `json:"total_chats"`
				TotalChatInsertionEvents int `json:"total_chat_insertion_events"`
				TotalChatCopyEvents      int `json:"total_chat_copy_events"`
			} `json:"models"`
		} `json:"editors"`
	} `json:"copilot_ide_chat"`
}

// FetchDailyMetrics fetches Copilot usage metrics for the given org and date range.
func FetchDailyMetrics(token, org, since, until string) ([]entity.DailyMetrics, error) {
	url := fmt.Sprintf(copilotMetricsURL+"?since=%s&until=%s", org, since, until)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call GitHub API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d: %s", resp.StatusCode, body)
	}

	var responses []metricsResponse
	if err := json.Unmarshal(body, &responses); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	metrics := make([]entity.DailyMetrics, 0, len(responses))
	for _, r := range responses {
		m := entity.DailyMetrics{
			Date:              r.Date,
			TotalEngagedUsers: r.TotalEngagedUsers,
		}
		for _, editor := range r.CopilotIDEChat.Editors {
			for _, model := range editor.Models {
				m.TotalChats += model.TotalChats
				m.TotalChatInsertionEvents += model.TotalChatInsertionEvents
				m.TotalChatCopyEvents += model.TotalChatCopyEvents
			}
		}
		metrics = append(metrics, m)
	}

	return metrics, nil
}
