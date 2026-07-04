package indexnow

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSubmitURLs_UsesDefaultKey(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	var gotMethod string
	var gotContentType string
	var gotBody submitRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotContentType = r.Header.Get("Content-Type")
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), defaultKey: "configured-key"}

	result, err := client.SubmitURLs(context.Background(), "example.com", []string{"https://example.com/a"}, "", "https://example.com/configured-key.txt")
	if err != nil {
		t.Fatalf("SubmitURLs error = %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Fatalf("method = %q, want POST", gotMethod)
	}
	if gotContentType != "application/json; charset=utf-8" {
		t.Fatalf("Content-Type = %q, want application/json; charset=utf-8", gotContentType)
	}
	if gotBody.Key != "configured-key" {
		t.Fatalf("Key = %q, want %q", gotBody.Key, "configured-key")
	}
	if !result.Success {
		t.Fatal("expected success")
	}
}

func TestSubmitURLs_KeyOverrideWins(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	var gotBody submitRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), defaultKey: "configured-key"}

	_, err := client.SubmitURLs(context.Background(), "example.com", []string{"https://example.com/a"}, "override-key", "")
	if err != nil {
		t.Fatalf("SubmitURLs error = %v", err)
	}
	if gotBody.Key != "override-key" {
		t.Fatalf("Key = %q, want %q", gotBody.Key, "override-key")
	}
}

func TestSubmitURLs_MissingKeyFails(t *testing.T) {
	client := NewClient("")
	_, err := client.SubmitURLs(context.Background(), "example.com", []string{"https://example.com/a"}, "", "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestSubmitURLs_StatusErrorIncludesReason(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), defaultKey: "configured-key"}

	_, err := client.SubmitURLs(context.Background(), "example.com", []string{"https://example.com/a"}, "", "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "IndexNow returned HTTP 403: key invalid" {
		t.Fatalf("error = %q", err.Error())
	}
}
