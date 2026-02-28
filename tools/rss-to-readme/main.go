package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode/utf8"
)

const (
	startMarker = "<!-- BLOG_ENTRIES_START -->"
	endMarker   = "<!-- BLOG_ENTRIES_END -->"
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

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	rssFile := flag.String("rss", "combined_feed.xml", "Path to the combined RSS XML file")
	readmeFile := flag.String("readme", "README.md", "Path to README.md")
	flag.Parse()

	data, err := os.ReadFile(*rssFile)
	if err != nil {
		return fmt.Errorf("failed to read RSS file: %w", err)
	}

	var feed rssFeed
	if err := xml.Unmarshal(data, &feed); err != nil {
		return fmt.Errorf("failed to parse RSS XML: %w", err)
	}

	markdown := buildMarkdown(feed.Channel.Items)

	content, err := os.ReadFile(*readmeFile)
	if err != nil {
		return fmt.Errorf("failed to read README: %w", err)
	}

	newContent, err := replaceBetweenMarkers(string(content), startMarker, endMarker, markdown)
	if err != nil {
		return fmt.Errorf("failed to replace content: %w", err)
	}

	if err := os.WriteFile(*readmeFile, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write README: %w", err)
	}

	log.Printf("Updated %s with %d blog entries", *readmeFile, len(feed.Channel.Items))
	return nil
}

// buildMarkdown groups items by source (derived from the "[hostname] title" prefix)
// and renders them as grouped markdown lists.
func buildMarkdown(items []rssItem) string {
	type entry struct {
		title string
		link  string
	}

	var order []string
	grouped := make(map[string][]entry)

	for _, item := range items {
		source, title := parseTitle(item.Title)
		if _, seen := grouped[source]; !seen {
			order = append(order, source)
		}
		grouped[source] = append(grouped[source], entry{title: title, link: item.Link})
	}

	var sb strings.Builder
	sb.WriteString("\n")
	for _, src := range order {
		label := src
		r, size := utf8.DecodeRuneInString(src)
		if r != utf8.RuneError && size > 0 {
			label = strings.ToUpper(string(r)) + src[size:]
		}
		sb.WriteString(fmt.Sprintf("#### %s\n\n", mdEscape(label)))
		for _, e := range grouped[src] {
			sb.WriteString(fmt.Sprintf("- [%s](%s)\n", mdEscape(e.title), e.link))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// mdEscape escapes characters that have special meaning inside markdown link text.
func mdEscape(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, "[", `\[`)
	s = strings.ReplaceAll(s, "]", `\]`)
	return s
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
