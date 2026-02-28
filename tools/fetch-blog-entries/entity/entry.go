package entity

const (
	SourceZenn  = "zenn"
	SourceQiita = "qiita"
	SourceNote  = "note"
)

type Entry struct {
	Title  string `json:"title"`
	URL    string `json:"url"`
	Source string `json:"source"`
	Date   string `json:"date"`
}
