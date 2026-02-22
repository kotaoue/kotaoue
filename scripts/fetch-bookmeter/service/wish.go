package service

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/kotaoue/kotaoue/scripts/fetch-bookmeter/repository"
)

// FetchAndSaveWishList fetches wish list and saves it to the given outputFile path
func FetchAndSaveWishList(userID, outputFile string) error {
	books, err := repository.FetchWishList(userID)
	if err != nil {
		return fmt.Errorf("failed to fetch wish list: %w", err)
	}

	jsonData, err := json.MarshalIndent(books, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(outputFile, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	log.Printf("✓ Success! Wish list saved to %s", outputFile)
	log.Printf("✓ Total wish books: %d", len(books))

	return nil
}
