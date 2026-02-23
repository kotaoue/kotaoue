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

	"github.com/kotaoue/kotaoue/tools/fetch-bookmeter/entity"
)

const (
	startMarker    = "<!-- WISH_BOOK_START -->"
	endMarker      = "<!-- WISH_BOOK_END -->"
	bookImageWidth = "128px"
)

// RunUpdateReadme parses flags and updates README.md with a random book from wish.json
func RunUpdateReadme(args []string) error {
	fs := flag.NewFlagSet("update-readme", flag.ExitOnError)
	wishFile := fs.String("wish-file", "wish.json", "Path to wish.json")
	readmeFile := fs.String("readme", "README.md", "Path to README.md")
	if err := fs.Parse(args); err != nil {
		return err
	}
	return updateReadme(*wishFile, *readmeFile)
}

func updateReadme(wishFile, readmeFile string) error {
	data, err := os.ReadFile(wishFile)
	if err != nil {
		return fmt.Errorf("failed to read wish file: %w", err)
	}

	var books []entity.Book
	if err := json.Unmarshal(data, &books); err != nil {
		return fmt.Errorf("failed to parse wish.json: %w", err)
	}

	if len(books) == 0 {
		log.Println("wish.json is empty, skipping update")
		return nil
	}

	valid := filterValidBooks(books)
	if len(valid) == 0 {
		log.Println("no valid book entries found in wish.json, skipping update")
		return nil
	}

	book := valid[rand.Intn(len(valid))]
	bookHTML := buildBookHTML(book)

	content, err := os.ReadFile(readmeFile)
	if err != nil {
		return fmt.Errorf("failed to read README file: %w", err)
	}

	newContent, err := replaceBetweenMarkers(string(content), startMarker, endMarker, bookHTML)
	if err != nil {
		return fmt.Errorf("failed to replace content: %w", err)
	}

	if err := os.WriteFile(readmeFile, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write README file: %w", err)
	}

	log.Printf("Updated %s with: %s", readmeFile, bookHTML)
	return nil
}

func filterValidBooks(books []entity.Book) []entity.Book {
	var valid []entity.Book
	for _, b := range books {
		if b.URL != "" && b.Thumb != "" && b.Title != "" {
			valid = append(valid, b)
		}
	}
	return valid
}

func buildBookHTML(book entity.Book) string {
	return fmt.Sprintf(
		`<a href="%s"><img src="%s" alt="%s" width="%s"></a>`,
		html.EscapeString(book.URL),
		html.EscapeString(book.Thumb),
		html.EscapeString(book.Title),
		bookImageWidth,
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
