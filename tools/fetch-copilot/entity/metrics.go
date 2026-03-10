package entity

// DailyMetrics represents a summary of Copilot usage for a single day.
type DailyMetrics struct {
	Date                     string
	TotalEngagedUsers        int
	TotalChats               int
	TotalChatInsertionEvents int
	TotalChatCopyEvents      int
}
