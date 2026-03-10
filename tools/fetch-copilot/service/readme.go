package service

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/kotaoue/kotaoue/tools/fetch-copilot/entity"
	"github.com/kotaoue/kotaoue/tools/fetch-copilot/repository"
)

const (
	startMarker = "<!-- COPILOT_SUMMARY_START -->"
	endMarker   = "<!-- COPILOT_SUMMARY_END -->"
)

// RunUpdateReadme parses flags and updates README.md with yesterday's Copilot usage summary.
func RunUpdateReadme(args []string) error {
	fs := flag.NewFlagSet("fetch-copilot", flag.ExitOnError)
	org := fs.String("org", os.Getenv("COPILOT_ORG"), "GitHub organization name")
	token := fs.String("token", os.Getenv("COPILOT_TOKEN"), "GitHub personal access token with manage_billing:copilot scope")
	readmeFile := fs.String("readme", "README.md", "Path to README.md")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if *org == "" || *token == "" {
		return fmt.Errorf("COPILOT_ORG and COPILOT_TOKEN environment variables are required")
	}

	return updateReadme(*token, *org, *readmeFile)
}

func updateReadme(token, org, readmeFile string) error {
	jst := time.FixedZone("JST", 9*60*60)
	yesterday := time.Now().In(jst).AddDate(0, 0, -1)
	date := yesterday.Format("2006-01-02")

	metrics, err := repository.FetchDailyMetrics(token, org, date, date)
	if err != nil {
		return fmt.Errorf("failed to fetch Copilot metrics: %w", err)
	}

	summary := buildSummary(metrics, yesterday)

	content, err := os.ReadFile(readmeFile)
	if err != nil {
		return fmt.Errorf("failed to read README file: %w", err)
	}

	newContent, err := replaceBetweenMarkers(string(content), startMarker, endMarker, summary)
	if err != nil {
		return fmt.Errorf("failed to replace content: %w", err)
	}

	if err := os.WriteFile(readmeFile, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write README file: %w", err)
	}

	log.Printf("Updated %s with Copilot summary for %s", readmeFile, date)
	return nil
}

// buildSummary creates a Markdown summary of the given Copilot daily metrics.
func buildSummary(metrics []entity.DailyMetrics, date time.Time) string {
	dateLabel := fmt.Sprintf("%d月%d日", date.Month(), date.Day())

	if len(metrics) == 0 {
		return fmt.Sprintf("\n%sのCopilot使用なし\n", dateLabel)
	}

	m := metrics[0]

	if m.TotalChats == 0 {
		return fmt.Sprintf("\n%sのCopilot使用なし\n", dateLabel)
	}

	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("**%sのCopilotとのやりとり**\n\n", dateLabel))
	sb.WriteString(fmt.Sprintf("- チャット: %d回\n", m.TotalChats))
	if m.TotalChatInsertionEvents > 0 {
		sb.WriteString(fmt.Sprintf("- コード挿入: %d回\n", m.TotalChatInsertionEvents))
	}
	if m.TotalChatCopyEvents > 0 {
		sb.WriteString(fmt.Sprintf("- コピー: %d回\n", m.TotalChatCopyEvents))
	}
	sb.WriteString("\n")

	return sb.String()
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
