package service

import (
	"encoding/json"
	"flag"
	"fmt"
	"html"
	"log"
	"math/rand"
	"os"
	"strings"

	"github.com/kotaoue/kotaoue/tools/fetch-spotify/entity"
)

const (
	startMarker    = "<!-- SPOTIFY_TRACK_START -->"
	endMarker      = "<!-- SPOTIFY_TRACK_END -->"
	trackImageWidth = "128px"
)

// RunUpdateReadme parses flags and updates README.md with a random track from playlist.json
func RunUpdateReadme(args []string) error {
	fs := flag.NewFlagSet("update-readme", flag.ExitOnError)
	playlistFile := fs.String("playlist-file", "playlist.json", "Path to playlist.json")
	readmeFile := fs.String("readme", "README.md", "Path to README.md")
	if err := fs.Parse(args); err != nil {
		return err
	}
	return updateReadme(*playlistFile, *readmeFile)
}

func updateReadme(playlistFile, readmeFile string) error {
	data, err := os.ReadFile(playlistFile)
	if err != nil {
		return fmt.Errorf("failed to read playlist file: %w", err)
	}

	var tracks []entity.Track
	if err := json.Unmarshal(data, &tracks); err != nil {
		return fmt.Errorf("failed to parse playlist.json: %w", err)
	}

	if len(tracks) == 0 {
		log.Println("playlist.json is empty, skipping update")
		return nil
	}

	valid := filterValidTracks(tracks)
	if len(valid) == 0 {
		log.Println("no valid track entries found in playlist.json, skipping update")
		return nil
	}

	track := valid[rand.Intn(len(valid))]
	trackHTML := buildTrackHTML(track)

	content, err := os.ReadFile(readmeFile)
	if err != nil {
		return fmt.Errorf("failed to read README file: %w", err)
	}

	newContent, err := replaceBetweenMarkers(string(content), startMarker, endMarker, trackHTML)
	if err != nil {
		return fmt.Errorf("failed to replace content: %w", err)
	}

	if err := os.WriteFile(readmeFile, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write README file: %w", err)
	}

	log.Printf("Updated %s with: %s", readmeFile, trackHTML)
	return nil
}

func filterValidTracks(tracks []entity.Track) []entity.Track {
	var valid []entity.Track
	for _, t := range tracks {
		if t.URL != "" && t.Thumb != "" && t.Title != "" {
			valid = append(valid, t)
		}
	}
	return valid
}

func buildTrackHTML(track entity.Track) string {
	return fmt.Sprintf(
		`<a href="%s"><img src="%s" alt="%s" width="%s"><br>%s<br>%s</a>`,
		html.EscapeString(track.URL),
		html.EscapeString(track.Thumb),
		html.EscapeString(track.Artist+" - "+track.Title),
		trackImageWidth,
		html.EscapeString(track.Artist),
		html.EscapeString(track.Title),
	)
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
