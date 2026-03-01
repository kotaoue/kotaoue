package repository

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/kotaoue/kotaoue/tools/rss-to-readme/entity"
)

type rssFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Channel rssChannel `xml:"channel"`
}

type rssChannel struct {
	Items []rssItem `xml:"item"`
}

type rssItem struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
}

// ParseEntries parses RSS 2.0 XML bytes and returns a slice of Entry.
func ParseEntries(data []byte) ([]entity.Entry, error) {
	var feed rssFeed
	if err := xml.Unmarshal(data, &feed); err != nil {
		return nil, fmt.Errorf("failed to parse RSS XML: %w", err)
	}

	entries := make([]entity.Entry, 0, len(feed.Channel.Items))
	for _, item := range feed.Channel.Items {
		source, title := parseTitle(item.Title)
		entries = append(entries, entity.Entry{
			Title:  title,
			Link:   item.Link,
			Source: source,
		})
	}

	return entries, nil
}

// parseTitle splits a title of the form "[hostname] actual title" into (hostname, actual title).
// If the format doesn't match, the full title is returned with an empty source.
func parseTitle(raw string) (source, title string) {
	if strings.HasPrefix(raw, "[") {
		end := strings.Index(raw, "]")
		if end > 0 {
			return raw[1:end], strings.TrimSpace(raw[end+1:])
		}
	}
	return "", raw
}

