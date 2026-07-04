package bingwebmaster

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestListSites_BuildsRequestAndMapsResponse(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	var capturedPath string
	var capturedQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		capturedQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"d":[{"Url":"https://example.test","IsVerified":true,"DnsVerificationCode":"dns-code","AuthenticationCode":"auth-code"}]}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.ListSites(context.Background())
	if err != nil {
		t.Fatalf("ListSites error = %v", err)
	}

	if capturedPath != "/GetUserSites" {
		t.Fatalf("path = %q, want %q", capturedPath, "/GetUserSites")
	}
	if !strings.Contains(capturedQuery, "apikey=test-api-key") {
		t.Fatalf("query = %q, want apikey", capturedQuery)
	}
	if len(result.Sites) != 1 {
		t.Fatalf("len(Sites) = %d, want 1", len(result.Sites))
	}
	if result.Sites[0].SiteURL != "https://example.test" {
		t.Fatalf("SiteURL = %q, want %q", result.Sites[0].SiteURL, "https://example.test")
	}
	if !result.Sites[0].IsVerified {
		t.Fatal("expected IsVerified true")
	}
}

func TestSubmitURLBatch_ValidatesMax500(t *testing.T) {
	urlList := make([]string, 501)
	for i := range urlList {
		urlList[i] = "https://example.test/page"
	}

	client := NewClient("test-api-key")
	_, err := client.SubmitURLBatch(context.Background(), "https://example.test", urlList)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "at most 500") {
		t.Fatalf("error = %q, want max-500 message", err.Error())
	}
}

func TestGetURLLinks_UsesLinkParameterAndParsesFixture(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	var capturedQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedQuery = r.URL.RawQuery
		if got := r.URL.Query().Get("link"); got != "https://example.test" {
			t.Fatalf("link query param = %q, want %q", got, "https://example.test")
		}
		if got := r.URL.Query().Get("page"); got != "0" {
			t.Fatalf("page query param = %q, want 0", got)
		}
		if got := r.URL.Query().Get("siteUrl"); got != "https://example.test" {
			t.Fatalf("siteUrl query param = %q, want %q", got, "https://example.test")
		}
		if got := r.URL.Query().Get("url"); got != "" {
			t.Fatalf("url query param should be empty, got %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"d":{"__type":"LinkDetails:#Microsoft.Bing.Webmaster.Api","Details":[],"TotalPages":0}}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.GetURLLinks(context.Background(), "https://example.test", "https://example.test", 0)
	if err != nil {
		t.Fatalf("GetURLLinks error = %v", err)
	}

	if !strings.Contains(capturedQuery, "apikey=test-api-key") {
		t.Fatalf("query = %q, want apikey", capturedQuery)
	}
	if result.TotalPages != 0 {
		t.Fatalf("TotalPages = %d, want 0", result.TotalPages)
	}
	if len(result.Details) != 0 {
		t.Fatalf("len(Details) = %d, want 0", len(result.Details))
	}
}

func TestGetCrawlIssues_ParsesEmptyArrayFixture(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"d":[]}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.GetCrawlIssues(context.Background(), "https://example.test")
	if err != nil {
		t.Fatalf("GetCrawlIssues error = %v", err)
	}
	if len(result.Issues) != 0 {
		t.Fatalf("len(Issues) = %d, want 0", len(result.Issues))
	}
}

func TestSubmitURL_SendsPOSTBody(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	var method string
	var contentType string
	var body map[string]string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		contentType = r.Header.Get("Content-Type")
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"d":true}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.SubmitURL(context.Background(), "https://example.test", "https://example.test/page")
	if err != nil {
		t.Fatalf("SubmitURL error = %v", err)
	}

	if method != http.MethodPost {
		t.Fatalf("method = %q, want POST", method)
	}
	if contentType != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", contentType)
	}
	if body["siteUrl"] != "https://example.test" {
		t.Fatalf("siteUrl = %q, want %q", body["siteUrl"], "https://example.test")
	}
	if body["url"] != "https://example.test/page" {
		t.Fatalf("url = %q, want %q", body["url"], "https://example.test/page")
	}
	if !result.Success {
		t.Fatal("expected Success true")
	}
}

func TestGetKeywordStats_UsesQParameterWithoutSiteURL(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if got := query.Get("q"); got != "golang" {
			t.Fatalf("q = %q, want %q", got, "golang")
		}
		if got := query.Get("country"); got != "US" {
			t.Fatalf("country = %q, want %q", got, "US")
		}
		if got := query.Get("language"); got != "en" {
			t.Fatalf("language = %q, want %q", got, "en")
		}
		if got := query.Get("siteUrl"); got != "" {
			t.Fatalf("siteUrl should be empty, got %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"d":[{"Query":"golang","Date":"/Date(1732612952000+0000)/","Impressions":42,"BroadImpressions":84}]}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.GetKeywordStats(context.Background(), "golang", "US", "en")
	if err != nil {
		t.Fatalf("GetKeywordStats error = %v", err)
	}
	if len(result.Stats) != 1 {
		t.Fatalf("len(Stats) = %d, want 1", len(result.Stats))
	}
	if result.Stats[0].Date == nil {
		t.Fatal("expected parsed date")
	}
}

func TestGetCrawlStats_ParsesDotNetDate(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"d":[{"Date":"/Date(1732612952000+0000)/","CrawledPages":5,"CrawlErrors":1,"InIndex":2,"InLinks":3,"Code2xx":4,"Code301":0,"Code302":0,"Code4xx":1,"Code5xx":0,"AllOtherCodes":0,"BlockedByRobotsTxt":0,"ContainsMalware":0}]}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.GetCrawlStats(context.Background(), "https://example.test")
	if err != nil {
		t.Fatalf("GetCrawlStats error = %v", err)
	}
	if len(result.Stats) != 1 {
		t.Fatalf("len(Stats) = %d, want 1", len(result.Stats))
	}
	got := result.Stats[0].Date
	if got == nil {
		t.Fatal("expected date, got nil")
	}
	want := time.UnixMilli(1732612952000).UTC()
	if !got.Equal(want) {
		t.Fatalf("Date = %s, want %s", got.Format(time.RFC3339), want.Format(time.RFC3339))
	}
}

func TestDo_Non2xxReturnsTypedErrorWithTruncatedBody(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	longBody := strings.Repeat("x", 400)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(longBody))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	_, err := client.ListSites(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var apiErr *apiRequestError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error type = %T, want *apiRequestError", err)
	}
	if apiErr.StatusCode != http.StatusBadRequest {
		t.Fatalf("StatusCode = %d, want %d", apiErr.StatusCode, http.StatusBadRequest)
	}
	if len(apiErr.Body) != 303 {
		t.Fatalf("len(Body) = %d, want 303", len(apiErr.Body))
	}
	if !strings.HasSuffix(apiErr.Body, "...") {
		t.Fatalf("Body should end with ellipsis, got %q", apiErr.Body[len(apiErr.Body)-3:])
	}
}
