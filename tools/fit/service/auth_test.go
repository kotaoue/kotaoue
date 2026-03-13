package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRefreshAccessToken(t *testing.T) {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		wantToken   string
		wantErrFrag string
	}{
		{
			name: "success",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				resp := tokenResponse{
					AccessToken: "test-access-token",
					TokenType:   "Bearer",
					ExpiresIn:   3600,
				}
				_ = json.NewEncoder(w).Encode(resp)
			},
			wantToken: "test-access-token",
		},
		{
			name: "non-200 status",
			handler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
			},
			wantErrFrag: "token endpoint returned status 401",
		},
		{
			name: "invalid json response",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("not-json"))
			},
			wantErrFrag: "failed to parse token response",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewServer(tc.handler)
			defer srv.Close()

			original := tokenURL
			tokenURL = srv.URL
			defer func() { tokenURL = original }()

			got, err := refreshAccessToken("client-id", "client-secret", "refresh-token")
			if tc.wantErrFrag != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tc.wantErrFrag)
				}
				if errStr := err.Error(); !strings.Contains(errStr, tc.wantErrFrag) {
					t.Fatalf("expected error containing %q, got %q", tc.wantErrFrag, errStr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.wantToken {
				t.Fatalf("refreshAccessToken() = %q, want %q", got, tc.wantToken)
			}
		})
	}
}
