package entity

type Book struct {
	No        int    `json:"no"`
	Title     string `json:"title"`
	URL       string `json:"url"`
	Author    string `json:"author"`
	AuthorURL string `json:"authorUrl"`
	Thumb     string `json:"thumb"`
	Date      string `json:"date"`
}
