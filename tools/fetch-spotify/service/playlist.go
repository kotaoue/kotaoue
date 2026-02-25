package service

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kotaoue/kotaoue/tools/fetch-spotify/repository"
)

// RunFetchPlaylist parses flags and fetches the Spotify playlist tracks.
func RunFetchPlaylist(args []string) error {
	fs := flag.NewFlagSet("fetch-playlist", flag.ExitOnError)
	playlistID := fs.String("playlist-id", "3aARAs2A4PgdkgYzcyYPgI", "Spotify playlist ID")
	output := fs.String("output", "playlist.json", "Output file path for playlist.json")
	if err := fs.Parse(args); err != nil {
		return err
	}

	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		return fmt.Errorf("SPOTIFY_CLIENT_ID and SPOTIFY_CLIENT_SECRET environment variables must be set")
	}

	return fetchAndSavePlaylist(clientID, clientSecret, *playlistID, *output)
}

func fetchAndSavePlaylist(clientID, clientSecret, playlistID, outputFile string) error {
	accessToken, err := repository.FetchAccessToken(clientID, clientSecret)
	if err != nil {
		return fmt.Errorf("failed to fetch access token: %w", err)
	}

	tracks, err := repository.FetchPlaylistTracks(accessToken, playlistID)
	if err != nil {
		return fmt.Errorf("failed to fetch playlist tracks: %w", err)
	}

	jsonData, err := json.MarshalIndent(tracks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(outputFile, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	log.Printf("✓ Success! Playlist saved to %s", outputFile)
	log.Printf("✓ Total tracks: %d", len(tracks))

	return nil
}
