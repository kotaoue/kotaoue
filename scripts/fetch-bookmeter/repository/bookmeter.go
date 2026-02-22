package repository

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/kotaoue/kotaoue/scripts/fetch-bookmeter/entity"
)

func fetchHTML(url string, retries int) (string, error) {
	log.Printf("Fetching HTML from %s", url)

	var lastErr error
	for i := 0; i <= retries; i++ {
		if i > 0 {
			time.Sleep(time.Duration(i) * time.Second)
			log.Printf("Retrying %s (attempt %d/%d)", url, i+1, retries)
		}
		content, err := doFetch(url)
		if err == nil {
			return content, nil
		}
		lastErr = err
	}
	return "", lastErr
}

func doFetch(url string) (string, error) {
	resp, err := http.Get(url) // #nosec G107 -- URL is constructed from user-supplied ID and a fixed pattern
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	htmlBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(htmlBytes), nil
}

func parseBooks(html string) []entity.Book {
	books := []entity.Book{}

	booksRegex := regexp.MustCompile(`<li class="group__book">(.*?)</div></li>`)
	bookMatches := booksRegex.FindAllStringSubmatch(html, -1)

	for i, match := range bookMatches {
		if len(match) > 1 {
			bookHTML := match[1]
			book := parseBook(bookHTML, i+1)
			books = append(books, book)
		}
	}

	return books
}

func parseBook(bookHTML string, no int) entity.Book {
	book := entity.Book{
		No: no,
	}

	titleRegex := regexp.MustCompile(`<div class="thumbnail__cover"><a href="(?P<url>.*?)"><img alt="(?P<title>.*?)" class`)
	if matches := titleRegex.FindStringSubmatch(bookHTML); len(matches) > 0 {
		book.URL = "https://bookmeter.com" + matches[1]
		book.Title = matches[2]
	}

	authorRegex := regexp.MustCompile(`<ul class="detail__authors"><li><a href="(?P<url>.*?)">(?P<author>.*?)</a></li></ul>`)
	if matches := authorRegex.FindStringSubmatch(bookHTML); len(matches) > 0 {
		book.AuthorURL = "https://bookmeter.com" + matches[1]
		book.Author = matches[2]
	}

	thumbRegex := regexp.MustCompile(`class="cover__image" src="(.*?)" />`)
	if matches := thumbRegex.FindStringSubmatch(bookHTML); len(matches) > 1 {
		book.Thumb = matches[1]
	}

	dateRegex := regexp.MustCompile(`<div class="detail__date">(.*?)</div>`)
	if matches := dateRegex.FindStringSubmatch(bookHTML); len(matches) > 1 {
		book.Date = strings.TrimSpace(matches[1])
	}

	return book
}
