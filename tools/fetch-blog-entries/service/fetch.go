package service

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/kotaoue/kotaoue/tools/fetch-blog-entries/entity"
	"github.com/kotaoue/kotaoue/tools/fetch-blog-entries/repository"
)

const defaultLimit = 5

// feedList implements flag.Value for a repeatable -feed flag.
type feedList []string

func (f *feedList) String() string { return strings.Join(*f, ", ") }
func (f *feedList) Set(v string) error {
	*f = append(*f, v)
	return nil
}

// RunFetchEntries parses flags and fetches blog entries from RSS feeds.
func RunFetchEntries(args []string) error {
	fs := flag.NewFlagSet("fetch-entries", flag.ExitOnError)
	output := fs.String("output", "blog-entries.json", "Output file path for blog-entries.json")
	limit := fs.Int("limit", defaultLimit, "Max entries per source")
	var feeds feedList
	fs.Var(&feeds, "feed", "RSS feed URL to fetch (repeatable)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if len(feeds) == 0 {
		return fmt.Errorf("at least one -feed URL is required")
	}
	return fetchAndSaveEntries(*output, *limit, feeds)
}

// sourceNameFromURL derives a short source label from an RSS feed URL hostname.
// e.g. "https://zenn.dev/kotaoue/feed" -> "zenn"
func sourceNameFromURL(feedURL string) string {
	u, err := url.Parse(feedURL)
	if err != nil || u.Host == "" {
		return feedURL
	}
	host := u.Hostname()
	if idx := strings.Index(host, "."); idx > 0 {
		return host[:idx]
	}
	return host
}

func fetchAndSaveEntries(outputFile string, limit int, feedURLs []string) error {
	var allEntries []entity.Entry

	for _, feedURL := range feedURLs {
		name := sourceNameFromURL(feedURL)
		entries, err := repository.FetchEntries(feedURL, name, limit)
		if err != nil {
			log.Printf("Warning: failed to fetch entries from %s: %v", name, err)
			continue
		}
		allEntries = append(allEntries, entries...)
		log.Printf("✓ Fetched %d entries from %s", len(entries), name)
	}

	jsonData, err := json.MarshalIndent(allEntries, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(outputFile, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	log.Printf("✓ Success! Blog entries saved to %s", outputFile)
	log.Printf("✓ Total entries: %d", len(allEntries))

	return nil
}
