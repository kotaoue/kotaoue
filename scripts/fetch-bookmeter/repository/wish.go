package repository

import (
	"fmt"
	"log"

	"github.com/kotaoue/kotaoue/scripts/fetch-bookmeter/entity"
)

// FetchWishList fetches and parses all pages of the wish list from Bookmeter
func FetchWishList(userID string) ([]entity.Book, error) {
	var allBooks []entity.Book
	no := 1

	for page := 1; ; page++ {
		url := wishListURL(userID, page)
		htmlContent, err := fetchHTML(url, 3)
		if err != nil {
			log.Printf("Stopping wish list pagination at page %d: %v", page, err)
			break
		}

		books := parseBooks(htmlContent)
		if len(books) == 0 {
			break
		}

		for _, b := range books {
			b.No = no
			allBooks = append(allBooks, b)
			no++
		}
		log.Printf("Parsed %d books from wish list page %d", len(books), page)
	}

	log.Printf("Parsed %d books total from wish list", len(allBooks))
	return allBooks, nil
}

func wishListURL(userID string, page int) string {
	return fmt.Sprintf("https://bookmeter.com/users/%s/books/wish?page=%d", userID, page)
}
