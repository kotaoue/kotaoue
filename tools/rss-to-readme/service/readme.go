package service

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/kotaoue/kotaoue/tools/rss-to-readme/entity"
	"github.com/kotaoue/kotaoue/tools/rss-to-readme/repository"
)

const (
	startMarker = "<!-- BLOG_ENTRIES_START -->"
	endMarker   = "<!-- BLOG_ENTRIES_END -->"
)

// RunUpdateReadme parses flags and updates README.md with blog entries from an RSS file.
func RunUpdateReadme(args []string) error {
	fs := flag.NewFlagSet("update-readme", flag.ExitOnError)
	rssFile := fs.String("rss", "combined_feed.xml", "Path to the combined RSS XML file")
	readmeFile := fs.String("readme", "README.md", "Path to README.md")
	maxEntries := fs.Int("max", 5, "Maximum number of entries to display")
	if err := fs.Parse(args); err != nil {
		return err
	}
	return updateReadme(*rssFile, *readmeFile, *maxEntries)
}

func updateReadme(rssFile, readmeFile string, maxEntries int) error {
	data, err := os.ReadFile(rssFile)
	if err != nil {
		return fmt.Errorf("failed to read RSS file: %w", err)
	}

	entries, err := repository.ParseEntries(data)
	if err != nil {
		return err
	}

	if len(entries) == 0 {
		log.Println("RSS feed has no items, skipping update")
		return nil
	}

	markdown := buildMarkdown(entries, maxEntries)

	content, err := os.ReadFile(readmeFile)
	if err != nil {
		return fmt.Errorf("failed to read README: %w", err)
	}

	newContent, err := replaceBetweenMarkers(string(content), startMarker, endMarker, markdown)
	if err != nil {
		return fmt.Errorf("failed to replace content: %w", err)
	}

	if err := os.WriteFile(readmeFile, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write README: %w", err)
	}

	log.Printf("Updated %s with %d blog entries", readmeFile, len(entries))
	return nil
}

func buildMarkdown(entries []entity.Entry, maxEntries int) string {
	if maxEntries > 0 && len(entries) > maxEntries {
		entries = entries[:maxEntries]
	}

	var sb strings.Builder
	sb.WriteString("\n")
	for _, e := range entries {
		sb.WriteString(fmt.Sprintf("- [%s](%s)\n", mdEscape(e.Title), e.Link))
	}
	sb.WriteString("\n")

	return sb.String()
}

func mdEscape(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, "[", `\[`)
	s = strings.ReplaceAll(s, "]", `\]`)
	return s
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
