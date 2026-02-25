package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/kotaoue/kotaoue/tools/fetch-spotify/entity"
)

const spotifyAPIBase = "https://api.spotify.com/v1"

type tokenResponse struct {
	AccessToken string `json:"access_token"`
}

type playlistTracksResponse struct {
	Items []struct {
		AddedAt string `json:"added_at"`
		Track   struct {
			Name    string `json:"name"`
			Artists []struct {
				Name string `json:"name"`
			} `json:"artists"`
			Album struct {
				Images []struct {
					URL string `json:"url"`
				} `json:"images"`
			} `json:"album"`
			ExternalURLs struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
		} `json:"track"`
	} `json:"items"`
	Next string `json:"next"`
}

// FetchAccessToken retrieves a Spotify API access token using client credentials.
func FetchAccessToken(clientID, clientSecret string) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest(http.MethodPost, "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(clientID, clientSecret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to request token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code from token endpoint: %d", resp.StatusCode)
	}

	var tr tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	return tr.AccessToken, nil
}

// FetchPlaylistTracks fetches all tracks from the given Spotify playlist ID.
func FetchPlaylistTracks(accessToken, playlistID string) ([]entity.Track, error) {
	var tracks []entity.Track
	no := 1
	nextURL := fmt.Sprintf("%s/playlists/%s/tracks?limit=100", spotifyAPIBase, playlistID)

	for nextURL != "" {
		log.Printf("Fetching playlist tracks from %s", nextURL)

		body, err := fetchJSON(accessToken, nextURL)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch playlist tracks: %w", err)
		}

		var resp playlistTracksResponse
		if err := json.Unmarshal(body, &resp); err != nil {
			return nil, fmt.Errorf("failed to parse playlist tracks response: %w", err)
		}

		for _, item := range resp.Items {
			t := item.Track
			if t.Name == "" {
				continue
			}

			artists := make([]string, 0, len(t.Artists))
			for _, a := range t.Artists {
				artists = append(artists, a.Name)
			}

			thumb := ""
			if len(t.Album.Images) > 0 {
				thumb = t.Album.Images[0].URL
			}

			tracks = append(tracks, entity.Track{
				No:     no,
				Title:  t.Name,
				URL:    t.ExternalURLs.Spotify,
				Artist: strings.Join(artists, ", "),
				Thumb:  thumb,
				Date:   item.AddedAt,
			})
			no++
		}

		log.Printf("Fetched %d tracks so far", len(tracks))
		nextURL = resp.Next
	}

	log.Printf("Total tracks fetched: %d", len(tracks))
	return tracks, nil
}

func fetchJSON(accessToken, rawURL string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}
