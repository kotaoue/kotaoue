package service

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kotaoue/kotaoue/tools/fetch-bookmeter/repository"
)

// RunFetchWish parses flags and fetches the wish list from Bookmeter
func RunFetchWish(args []string) error {
	fs := flag.NewFlagSet("fetch-wish", flag.ExitOnError)
	userID := fs.String("user-id", "104", "Bookmeter user ID")
	output := fs.String("output", "wish.json", "Output file path for wish.json")
	if err := fs.Parse(args); err != nil {
		return err
	}
	return fetchAndSaveWishList(*userID, *output)
}

func fetchAndSaveWishList(userID, outputFile string) error {
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
