package service

import (
	"encoding/json"
	"flag"
	"fmt"
	"html"
	"log"
	"os"
	"strings"

	"github.com/kotaoue/kotaoue/tools/fetch-blog-entries/entity"
)

const (
	startMarker = "<!-- BLOG_ENTRIES_START -->"
	endMarker   = "<!-- BLOG_ENTRIES_END -->"
)

// RunUpdateReadme parses flags and updates README.md with recent blog entries.
func RunUpdateReadme(args []string) error {
	fs := flag.NewFlagSet("update-readme", flag.ExitOnError)
	entriesFile := fs.String("entries-file", "blog-entries.json", "Path to blog-entries.json")
	readmeFile := fs.String("readme", "README.md", "Path to README.md")
	if err := fs.Parse(args); err != nil {
		return err
	}
	return updateReadme(*entriesFile, *readmeFile)
}

func updateReadme(entriesFile, readmeFile string) error {
	data, err := os.ReadFile(entriesFile)
	if err != nil {
		return fmt.Errorf("failed to read entries file: %w", err)
	}

	var entries []entity.Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return fmt.Errorf("failed to parse blog-entries.json: %w", err)
	}

	if len(entries) == 0 {
		log.Println("blog-entries.json is empty, skipping update")
		return nil
	}

	entriesMarkdown := buildEntriesMarkdown(entries)

	content, err := os.ReadFile(readmeFile)
	if err != nil {
		return fmt.Errorf("failed to read README file: %w", err)
	}

	newContent, err := replaceBetweenMarkers(string(content), startMarker, endMarker, entriesMarkdown)
	if err != nil {
		return fmt.Errorf("failed to replace content: %w", err)
	}

	if err := os.WriteFile(readmeFile, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write README file: %w", err)
	}

	log.Printf("Updated %s with blog entries", readmeFile)
	return nil
}

func buildEntriesMarkdown(entries []entity.Entry) string {
	// Group by source while preserving order of first appearance.
	var order []string
	grouped := make(map[string][]entity.Entry)
	feedURLs := make(map[string]string)

	for _, e := range entries {
		if _, seen := grouped[e.Source]; !seen {
			order = append(order, e.Source)
			feedURLs[e.Source] = e.FeedURL
		}
		grouped[e.Source] = append(grouped[e.Source], e)
	}

	var sb strings.Builder
	sb.WriteString("\n")

	for _, src := range order {
		srcEntries := grouped[src]
		label := src
		if len(src) > 0 {
			label = strings.ToUpper(src[:1]) + src[1:]
		}
		sb.WriteString(fmt.Sprintf("#### [%s](%s)\n\n", html.EscapeString(label), html.EscapeString(feedURLs[src])))
		for _, e := range srcEntries {
			sb.WriteString(fmt.Sprintf("- [%s](%s)\n", html.EscapeString(e.Title), html.EscapeString(e.URL)))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

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
