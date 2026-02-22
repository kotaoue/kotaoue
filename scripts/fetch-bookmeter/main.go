package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/kotaoue/kotaoue/scripts/fetch-bookmeter/service"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	userID := flag.String("user-id", "104", "Bookmeter user ID")
	output := flag.String("output", "wish.json", "Output file path for wish.json")
	flag.Parse()

	if err := service.FetchAndSaveWishList(*userID, *output); err != nil {
		return fmt.Errorf("failed to fetch and save wish list: %w", err)
	}
	return nil
}
