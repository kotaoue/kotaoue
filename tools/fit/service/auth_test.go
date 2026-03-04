package service

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRefreshAccessToken_InvalidGrant(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid_grant","error_description":"Token has been expired or revoked."}`)) //nolint:errcheck
	}))
	defer srv.Close()

	_, err := refreshAccessTokenWithURL(srv.URL, "id", "secret", "expired-token")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrTokenExpired) {
		t.Fatalf("expected ErrTokenExpired, got %v", err)
	}
}

func TestRefreshAccessToken_OtherError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid_client","error_description":"The OAuth client was not found."}`)) //nolint:errcheck
	}))
	defer srv.Close()

	_, err := refreshAccessTokenWithURL(srv.URL, "bad-id", "secret", "token")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if errors.Is(err, ErrTokenExpired) {
		t.Fatal("expected generic error, not ErrTokenExpired")
	}
}
