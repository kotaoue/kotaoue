package service

import (
	"testing"
	"time"

	"github.com/kotaoue/kotaoue/tools/fetch-copilot/entity"
)

func TestBuildSummary(t *testing.T) {
	jst := time.FixedZone("JST", 9*60*60)
	date := time.Date(2025, 3, 10, 0, 0, 0, 0, jst)

	tests := []struct {
		name    string
		metrics []entity.DailyMetrics
		want    string
	}{
		{
			name:    "no metrics",
			metrics: []entity.DailyMetrics{},
			want:    "\n3月10日のCopilot使用なし\n",
		},
		{
			name: "no chats",
			metrics: []entity.DailyMetrics{
				{Date: "2025-03-10", TotalChats: 0},
			},
			want: "\n3月10日のCopilot使用なし\n",
		},
		{
			name: "chats only",
			metrics: []entity.DailyMetrics{
				{Date: "2025-03-10", TotalChats: 5},
			},
			want: "\n**3月10日のCopilotとのやりとり**\n\n- チャット: 5回\n\n",
		},
		{
			name: "chats with insertions and copies",
			metrics: []entity.DailyMetrics{
				{
					Date:                     "2025-03-10",
					TotalChats:               10,
					TotalChatInsertionEvents: 3,
					TotalChatCopyEvents:      2,
				},
			},
			want: "\n**3月10日のCopilotとのやりとり**\n\n- チャット: 10回\n- コード挿入: 3回\n- コピー: 2回\n\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := buildSummary(tc.metrics, date)
			if got != tc.want {
				t.Fatalf("buildSummary() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestReplaceBetweenMarkers(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		start       string
		end         string
		replacement string
		want        string
		wantErr     bool
	}{
		{
			name:        "basic replacement",
			content:     "prefix<!-- START -->old content<!-- END -->suffix",
			start:       "<!-- START -->",
			end:         "<!-- END -->",
			replacement: "new content",
			want:        "prefix<!-- START -->new content<!-- END -->suffix",
		},
		{
			name:    "missing start marker",
			content: "no markers here<!-- END -->",
			start:   "<!-- START -->",
			end:     "<!-- END -->",
			wantErr: true,
		},
		{
			name:    "missing end marker",
			content: "<!-- START -->no end marker",
			start:   "<!-- START -->",
			end:     "<!-- END -->",
			wantErr: true,
		},
		{
			name:    "markers in wrong order",
			content: "<!-- END --><!-- START -->",
			start:   "<!-- START -->",
			end:     "<!-- END -->",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := replaceBetweenMarkers(tc.content, tc.start, tc.end, tc.replacement)
			if (err != nil) != tc.wantErr {
				t.Fatalf("replaceBetweenMarkers() error = %v, wantErr %v", err, tc.wantErr)
			}
			if !tc.wantErr && got != tc.want {
				t.Fatalf("replaceBetweenMarkers() = %q, want %q", got, tc.want)
			}
		})
	}
}
