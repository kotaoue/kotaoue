package entity

type Track struct {
	No     int    `json:"no"`
	Title  string `json:"title"`
	URL    string `json:"url"`
	Artist string `json:"artist"`
	Thumb  string `json:"thumb"`
	Date   string `json:"date"`
}
