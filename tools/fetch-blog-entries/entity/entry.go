package entity

type Entry struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Source  string `json:"source"`
	Date    string `json:"date"`
	FeedURL string `json:"feed_url"`
}
