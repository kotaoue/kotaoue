package repository

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/kotaoue/kotaoue/tools/fetch-blog-entries/entity"
)

const httpTimeout = 30 * time.Second
const dateLength = 10 // len("2006-01-02")

type atomFeed struct {
	XMLName xml.Name    `xml:"feed"`
	Entries []atomEntry `xml:"entry"`
}

type atomEntry struct {
	Title     string     `xml:"title"`
	Links     []atomLink `xml:"link"`
	Published string     `xml:"published"`
	Updated   string     `xml:"updated"`
}

type atomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
}

type rssFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Channel rssChannel `xml:"channel"`
}

type rssChannel struct {
	Items []rssItem `xml:"item"`
}

type rssItem struct {
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	PubDate string `xml:"pubDate"`
}

// FetchEntries fetches RSS/Atom entries from the given URL and returns them as an Entry slice.
func FetchEntries(url, source string, limit int) ([]entity.Entry, error) {
	log.Printf("Fetching RSS from %s", url)

	data, err := fetchURL(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch RSS: %w", err)
	}

	entries, err := parseRSS(data, source)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSS: %w", err)
	}

	if limit > 0 && len(entries) > limit {
		entries = entries[:limit]
	}

	return entries, nil
}

func fetchURL(url string) ([]byte, error) {
	client := &http.Client{Timeout: httpTimeout}
	resp, err := client.Get(url) // #nosec G107 -- URL is from fixed configuration
	if err != nil {
		return nil, fmt.Errorf("failed to GET %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d from %s", resp.StatusCode, url)
	}

	return io.ReadAll(resp.Body)
}

func parseRSS(data []byte, source string) ([]entity.Entry, error) {
	// Try Atom first
	var atom atomFeed
	if err := xml.Unmarshal(data, &atom); err == nil && len(atom.Entries) > 0 {
		return atomToEntries(atom, source), nil
	}

	// Fall back to RSS 2.0
	var rss rssFeed
	if err := xml.Unmarshal(data, &rss); err == nil && len(rss.Channel.Items) > 0 {
		return rssToEntries(rss, source), nil
	}

	return nil, fmt.Errorf("failed to parse feed as Atom or RSS 2.0")
}

func atomToEntries(feed atomFeed, source string) []entity.Entry {
	entries := make([]entity.Entry, 0, len(feed.Entries))
	for _, e := range feed.Entries {
		url := ""
		for _, link := range e.Links {
			if link.Rel == "alternate" || link.Rel == "" {
				url = link.Href
				break
			}
		}

		date := e.Published
		if date == "" {
			date = e.Updated
		}
		if len(date) >= dateLength {
			date = date[:dateLength]
		}

		entries = append(entries, entity.Entry{
			Title:  e.Title,
			URL:    url,
			Source: source,
			Date:   date,
		})
	}
	return entries
}

func rssToEntries(feed rssFeed, source string) []entity.Entry {
	entries := make([]entity.Entry, 0, len(feed.Channel.Items))
	for _, item := range feed.Channel.Items {
		entries = append(entries, entity.Entry{
			Title:  item.Title,
			URL:    item.Link,
			Source: source,
			Date:   formatRSSDate(item.PubDate),
		})
	}
	return entries
}

func formatRSSDate(pubDate string) string {
	formats := []string{
		time.RFC1123Z,
		time.RFC1123,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, pubDate); err == nil {
			return t.Format("2006-01-02")
		}
	}
	return pubDate
}
