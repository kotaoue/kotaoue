package service

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kotaoue/kotaoue/tools/fetch-blog-entries/entity"
	"github.com/kotaoue/kotaoue/tools/fetch-blog-entries/repository"
)

const defaultLimit = 5

var sources = []struct {
	Name string
	URL  string
}{
	{entity.SourceZenn, "https://zenn.dev/kotaoue/feed"},
	{entity.SourceQiita, "https://qiita.com/kotaoue/feed"},
	{entity.SourceNote, "https://note.com/kotaoue/rss"},
}

// RunFetchEntries parses flags and fetches blog entries from RSS feeds.
func RunFetchEntries(args []string) error {
	fs := flag.NewFlagSet("fetch-entries", flag.ExitOnError)
	output := fs.String("output", "blog-entries.json", "Output file path for blog-entries.json")
	limit := fs.Int("limit", defaultLimit, "Max entries per source")
	if err := fs.Parse(args); err != nil {
		return err
	}
	return fetchAndSaveEntries(*output, *limit)
}

func fetchAndSaveEntries(outputFile string, limit int) error {
	var allEntries []entity.Entry

	for _, src := range sources {
		entries, err := repository.FetchEntries(src.URL, src.Name, limit)
		if err != nil {
			log.Printf("Warning: failed to fetch entries from %s: %v", src.Name, err)
			continue
		}
		allEntries = append(allEntries, entries...)
		log.Printf("✓ Fetched %d entries from %s", len(entries), src.Name)
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
