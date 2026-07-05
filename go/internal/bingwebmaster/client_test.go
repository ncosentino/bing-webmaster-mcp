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

func TestAddSite_TreatsNullPayloadAsSuccess_MatchingRealBingBehavior(t *testing.T) {
	// Regression test: real Bing "d":null payload (observed live for both a fresh add and a
	// no-op repeat of an already-added site) previously got silently coerced to Success=false
	// because Go's json.Unmarshal leaves a bool destination unchanged (at its zero value) when
	// given JSON null. AddSite has no reliable boolean signal -- success means the HTTP call
	// completed without error.
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
		_, _ = w.Write([]byte(`{"d":null}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.AddSite(context.Background(), "https://example.test")
	if err != nil {
		t.Fatalf("AddSite error = %v", err)
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
	if result.SiteURL != "https://example.test" {
		t.Fatalf("SiteURL = %q, want %q", result.SiteURL, "https://example.test")
	}
	if !result.Success {
		t.Fatal("expected Success true")
	}
}

func TestVerifySite_SendsPOSTBody(t *testing.T) {
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

	result, err := client.VerifySite(context.Background(), "https://example.test")
	if err != nil {
		t.Fatalf("VerifySite error = %v", err)
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
	if !result.Verified {
		t.Fatal("expected Verified true")
	}
}

func TestListSitemaps_UsesSiteURLAndParsesFeed(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	var method string
	var capturedPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		capturedPath = r.URL.Path
		if got := r.URL.Query().Get("siteUrl"); got != "https://example.test" {
			t.Fatalf("siteUrl query param = %q, want %q", got, "https://example.test")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"d":[{"Url":"https://example.test/sitemap.xml","Type":"Sitemap","Compressed":true,"FileSize":123,"LastCrawled":"/Date(1732612952000+0000)/","Submitted":"/Date(1732616552000+0000)/","Status":"Success","UrlCount":12}]}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.ListSitemaps(context.Background(), "https://example.test")
	if err != nil {
		t.Fatalf("ListSitemaps error = %v", err)
	}

	if method != http.MethodGet {
		t.Fatalf("method = %q, want GET", method)
	}
	if capturedPath != "/GetFeeds" {
		t.Fatalf("path = %q, want %q", capturedPath, "/GetFeeds")
	}
	if len(result.Sitemaps) != 1 {
		t.Fatalf("len(Sitemaps) = %d, want 1", len(result.Sitemaps))
	}
	if result.Sitemaps[0].URL != "https://example.test/sitemap.xml" {
		t.Fatalf("URL = %q, want %q", result.Sitemaps[0].URL, "https://example.test/sitemap.xml")
	}
	if result.Sitemaps[0].Submitted == nil {
		t.Fatal("expected Submitted date")
	}
}

func TestGetSitemapDetails_UsesFeedURLAndParsesFeed(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	var method string
	var capturedPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		capturedPath = r.URL.Path
		query := r.URL.Query()
		if got := query.Get("siteUrl"); got != "https://example.test" {
			t.Fatalf("siteUrl query param = %q, want %q", got, "https://example.test")
		}
		if got := query.Get("feedUrl"); got != "https://example.test/sitemap.xml" {
			t.Fatalf("feedUrl query param = %q, want %q", got, "https://example.test/sitemap.xml")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"d":{"Url":"https://example.test/sitemap.xml","Type":"Sitemap","Compressed":false,"FileSize":456,"LastCrawled":"/Date(1732612952000+0000)/","Submitted":"/Date(1732616552000+0000)/","Status":"Pending","UrlCount":34}}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.GetSitemapDetails(context.Background(), "https://example.test", "https://example.test/sitemap.xml")
	if err != nil {
		t.Fatalf("GetSitemapDetails error = %v", err)
	}

	if method != http.MethodGet {
		t.Fatalf("method = %q, want GET", method)
	}
	if capturedPath != "/GetFeedDetails" {
		t.Fatalf("path = %q, want %q", capturedPath, "/GetFeedDetails")
	}
	if result.FeedURL != "https://example.test/sitemap.xml" {
		t.Fatalf("FeedURL = %q, want %q", result.FeedURL, "https://example.test/sitemap.xml")
	}
	if result.Sitemap.LastCrawled == nil {
		t.Fatal("expected LastCrawled date")
	}
	if result.Sitemap.URLCount != 34 {
		t.Fatalf("URLCount = %d, want 34", result.Sitemap.URLCount)
	}
}

func TestGetSitemapDetails_BestEffortHandlesSparseResponse(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %q, want GET", r.Method)
		}
		if got := r.URL.Query().Get("feedUrl"); got != "https://example.test/sitemap.xml" {
			t.Fatalf("feedUrl query param = %q, want %q", got, "https://example.test/sitemap.xml")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"d":{"Url":"https://example.test/sitemap.xml"}}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.GetSitemapDetails(context.Background(), "https://example.test", "https://example.test/sitemap.xml")
	if err != nil {
		t.Fatalf("GetSitemapDetails error = %v", err)
	}

	if result.Sitemap.URL != "https://example.test/sitemap.xml" {
		t.Fatalf("URL = %q, want %q", result.Sitemap.URL, "https://example.test/sitemap.xml")
	}
	if result.Sitemap.LastCrawled != nil {
		t.Fatalf("LastCrawled = %v, want nil", result.Sitemap.LastCrawled)
	}
	if result.Sitemap.URLCount != 0 {
		t.Fatalf("URLCount = %d, want 0", result.Sitemap.URLCount)
	}
}

// TestGetSitemapDetails_HandlesArrayResponse is a regression test for a real bug found via
// live end-to-end testing against a real Bing Webmaster account: Bing sometimes returns an
// array of feed objects from GetFeedDetails (e.g. for sitemap index files) rather than a single
// object, which previously crashed with "cannot unmarshal array into Go value of type
// bingwebmaster.rawFeed".
func TestGetSitemapDetails_HandlesArrayResponse(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"d":[{"Url":"https://example.test/sitemap-a.xml","Type":"Sitemap","Compressed":false,"FileSize":100,"LastCrawled":"/Date(1732612952000+0000)/","Submitted":"/Date(1732526552000+0000)/","Status":"Success","UrlCount":10},{"Url":"https://example.test/sitemap-b.xml","Type":"Sitemap","Compressed":false,"FileSize":200,"LastCrawled":"/Date(1732612952000+0000)/","Submitted":"/Date(1732526552000+0000)/","Status":"Ignored","UrlCount":20}]}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.GetSitemapDetails(context.Background(), "https://example.test", "https://example.test/sitemap-index.xml")
	if err != nil {
		t.Fatalf("GetSitemapDetails error = %v", err)
	}
	if result.Sitemap == nil {
		t.Fatal("expected non-nil Sitemap from array response")
	}
	if result.Sitemap.URL != "https://example.test/sitemap-a.xml" {
		t.Fatalf("URL = %q, want first array item %q", result.Sitemap.URL, "https://example.test/sitemap-a.xml")
	}
	if result.Sitemap.URLCount != 10 {
		t.Fatalf("URLCount = %d, want 10", result.Sitemap.URLCount)
	}
}

// TestGetSitemapDetails_HandlesEmptyArrayResponse confirms an empty array (no matching feed)
// maps to a nil Sitemap rather than erroring.
func TestGetSitemapDetails_HandlesEmptyArrayResponse(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"d":[]}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.GetSitemapDetails(context.Background(), "https://example.test", "https://example.test/sitemap.xml")
	if err != nil {
		t.Fatalf("GetSitemapDetails error = %v", err)
	}
	if result.Sitemap != nil {
		t.Fatalf("Sitemap = %+v, want nil", result.Sitemap)
	}
}

func TestSubmitSitemap_SendsPOSTBody(t *testing.T) {
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

	result, err := client.SubmitSitemap(context.Background(), "https://example.test", "https://example.test/sitemap.xml")
	if err != nil {
		t.Fatalf("SubmitSitemap error = %v", err)
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
	if body["feedUrl"] != "https://example.test/sitemap.xml" {
		t.Fatalf("feedUrl = %q, want %q", body["feedUrl"], "https://example.test/sitemap.xml")
	}
	if !result.Success {
		t.Fatal("expected Success true")
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
	if len(result.Rows) != 1 {
		t.Fatalf("len(Rows) = %d, want 1", len(result.Rows))
	}
	if result.Rows[0].Date == nil {
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
	if len(result.Rows) != 1 {
		t.Fatalf("len(Rows) = %d, want 1", len(result.Rows))
	}
	got := result.Rows[0].Date
	if got == nil {
		t.Fatal("expected date, got nil")
	}
	want := time.UnixMilli(1732612952000).UTC()
	if !got.Equal(want) {
		t.Fatalf("Date = %s, want %s", got.Format(time.RFC3339), want.Format(time.RFC3339))
	}
}

func TestGetURLSubmissionQuota_UsesSiteURLAndMapsQuotas(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	var method string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		if got := r.URL.Query().Get("siteUrl"); got != "https://example.test" {
			t.Fatalf("siteUrl query param = %q, want %q", got, "https://example.test")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"d":{"DailyQuota":10,"MonthlyQuota":300}}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.GetURLSubmissionQuota(context.Background(), "https://example.test")
	if err != nil {
		t.Fatalf("GetURLSubmissionQuota error = %v", err)
	}

	if method != http.MethodGet {
		t.Fatalf("method = %q, want GET", method)
	}
	if result.DailyQuota != 10 {
		t.Fatalf("DailyQuota = %d, want 10", result.DailyQuota)
	}
	if result.MonthlyQuota != 300 {
		t.Fatalf("MonthlyQuota = %d, want 300", result.MonthlyQuota)
	}
}

func TestGetURLInfo_UsesURLParameterAndParsesDates(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	var method string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		query := r.URL.Query()
		if got := query.Get("siteUrl"); got != "https://example.test" {
			t.Fatalf("siteUrl query param = %q, want %q", got, "https://example.test")
		}
		if got := query.Get("url"); got != "https://example.test/page" {
			t.Fatalf("url query param = %q, want %q", got, "https://example.test/page")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"d":{"Url":"https://example.test/page","IsPage":true,"HttpStatus":200,"DocumentSize":5120,"AnchorCount":8,"DiscoveryDate":"/Date(1732612952000+0000)/","LastCrawledDate":"/Date(1732616552000+0000)/","TotalChildUrlCount":3}}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.GetURLInfo(context.Background(), "https://example.test", "https://example.test/page")
	if err != nil {
		t.Fatalf("GetURLInfo error = %v", err)
	}

	if method != http.MethodGet {
		t.Fatalf("method = %q, want GET", method)
	}
	if result.LastCrawledDate == nil {
		t.Fatal("expected LastCrawledDate")
	}
	if result.HTTPStatus != 200 {
		t.Fatalf("HTTPStatus = %d, want 200", result.HTTPStatus)
	}
	if result.TotalChildURLCount != 3 {
		t.Fatalf("TotalChildURLCount = %d, want 3", result.TotalChildURLCount)
	}
}

func TestGetURLTrafficInfo_UsesURLParameterAndMapsMetrics(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	var method string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		query := r.URL.Query()
		if got := query.Get("siteUrl"); got != "https://example.test" {
			t.Fatalf("siteUrl query param = %q, want %q", got, "https://example.test")
		}
		if got := query.Get("url"); got != "https://example.test/page" {
			t.Fatalf("url query param = %q, want %q", got, "https://example.test/page")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"d":{"Url":"https://example.test/page","IsPage":true,"Clicks":25,"Impressions":400}}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.GetURLTrafficInfo(context.Background(), "https://example.test", "https://example.test/page")
	if err != nil {
		t.Fatalf("GetURLTrafficInfo error = %v", err)
	}

	if method != http.MethodGet {
		t.Fatalf("method = %q, want GET", method)
	}
	if result.Clicks != 25 {
		t.Fatalf("Clicks = %d, want 25", result.Clicks)
	}
	if result.Impressions != 400 {
		t.Fatalf("Impressions = %d, want 400", result.Impressions)
	}
}

func TestGetLinkCounts_UsesPageParameterAndParsesCounts(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	var method string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		query := r.URL.Query()
		if got := query.Get("siteUrl"); got != "https://example.test" {
			t.Fatalf("siteUrl query param = %q, want %q", got, "https://example.test")
		}
		if got := query.Get("page"); got != "2" {
			t.Fatalf("page query param = %q, want 2", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"d":{"Links":[{"Count":7,"Url":"https://referrer.test/page"}],"TotalPages":4}}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.GetLinkCounts(context.Background(), "https://example.test", 2)
	if err != nil {
		t.Fatalf("GetLinkCounts error = %v", err)
	}

	if method != http.MethodGet {
		t.Fatalf("method = %q, want GET", method)
	}
	if result.TotalPages != 4 {
		t.Fatalf("TotalPages = %d, want 4", result.TotalPages)
	}
	if len(result.Links) != 1 {
		t.Fatalf("len(Links) = %d, want 1", len(result.Links))
	}
	if result.Links[0].Count != 7 {
		t.Fatalf("Count = %d, want 7", result.Links[0].Count)
	}
}

func TestGetRankAndTrafficStats_ParsesDotNetDate(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	var method string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		if got := r.URL.Query().Get("siteUrl"); got != "https://example.test" {
			t.Fatalf("siteUrl query param = %q, want %q", got, "https://example.test")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"d":[{"Date":"/Date(1732612952000+0000)/","Clicks":12,"Impressions":300}]}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.GetRankAndTrafficStats(context.Background(), "https://example.test")
	if err != nil {
		t.Fatalf("GetRankAndTrafficStats error = %v", err)
	}

	if method != http.MethodGet {
		t.Fatalf("method = %q, want GET", method)
	}
	if len(result.Rows) != 1 {
		t.Fatalf("len(Rows) = %d, want 1", len(result.Rows))
	}
	if result.Rows[0].Date == nil {
		t.Fatal("expected Date")
	}
	if result.Rows[0].Impressions != 300 {
		t.Fatalf("Impressions = %d, want 300", result.Rows[0].Impressions)
	}
}

func TestGetQueryStats_UsesSiteURLAndMapsQueryStats(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	var method string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		if got := r.URL.Query().Get("siteUrl"); got != "https://example.test" {
			t.Fatalf("siteUrl query param = %q, want %q", got, "https://example.test")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"d":[{"Query":"bing api","Date":"/Date(1732612952000+0000)/","Clicks":9,"Impressions":100,"AvgClickPosition":2,"AvgImpressionPosition":4}]}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.GetQueryStats(context.Background(), "https://example.test")
	if err != nil {
		t.Fatalf("GetQueryStats error = %v", err)
	}

	if method != http.MethodGet {
		t.Fatalf("method = %q, want GET", method)
	}
	if len(result.Rows) != 1 {
		t.Fatalf("len(Rows) = %d, want 1", len(result.Rows))
	}
	if result.Rows[0].Query != "bing api" {
		t.Fatalf("Query = %q, want %q", result.Rows[0].Query, "bing api")
	}
	if result.Rows[0].Date == nil {
		t.Fatal("expected Date")
	}
}

func TestGetPageStats_UsesQueryPayloadFieldAsPage(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	var method string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		if got := r.URL.Query().Get("siteUrl"); got != "https://example.test" {
			t.Fatalf("siteUrl query param = %q, want %q", got, "https://example.test")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"d":[{"Query":"https://example.test/page","Date":"/Date(1732612952000+0000)/","Clicks":15,"Impressions":220,"AvgClickPosition":1,"AvgImpressionPosition":3}]}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.GetPageStats(context.Background(), "https://example.test")
	if err != nil {
		t.Fatalf("GetPageStats error = %v", err)
	}

	if method != http.MethodGet {
		t.Fatalf("method = %q, want GET", method)
	}
	if len(result.Rows) != 1 {
		t.Fatalf("len(Rows) = %d, want 1", len(result.Rows))
	}
	if result.Rows[0].Page != "https://example.test/page" {
		t.Fatalf("Page = %q, want %q", result.Rows[0].Page, "https://example.test/page")
	}
	if result.Rows[0].Date == nil {
		t.Fatal("expected Date")
	}
}

func TestGetPageQueryStats_UsesPageParameterAndMapsQueries(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	var method string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		query := r.URL.Query()
		if got := query.Get("siteUrl"); got != "https://example.test" {
			t.Fatalf("siteUrl query param = %q, want %q", got, "https://example.test")
		}
		if got := query.Get("page"); got != "https://example.test/page" {
			t.Fatalf("page query param = %q, want %q", got, "https://example.test/page")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"d":[{"Query":"bing api","Date":"/Date(1732612952000+0000)/","Clicks":4,"Impressions":50,"AvgClickPosition":6,"AvgImpressionPosition":8}]}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.GetPageQueryStats(context.Background(), "https://example.test", "https://example.test/page")
	if err != nil {
		t.Fatalf("GetPageQueryStats error = %v", err)
	}

	if method != http.MethodGet {
		t.Fatalf("method = %q, want GET", method)
	}
	if result.Page != "https://example.test/page" {
		t.Fatalf("Page = %q, want %q", result.Page, "https://example.test/page")
	}
	if len(result.Rows) != 1 {
		t.Fatalf("len(Rows) = %d, want 1", len(result.Rows))
	}
	if result.Rows[0].Query != "bing api" {
		t.Fatalf("Query = %q, want %q", result.Rows[0].Query, "bing api")
	}
}

func TestGetQueryPageStats_UsesQueryParameterAndMapsPages(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	var method string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		query := r.URL.Query()
		if got := query.Get("siteUrl"); got != "https://example.test" {
			t.Fatalf("siteUrl query param = %q, want %q", got, "https://example.test")
		}
		if got := query.Get("query"); got != "bing api" {
			t.Fatalf("query query param = %q, want %q", got, "bing api")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"d":[{"Query":"https://example.test/page","Date":"/Date(1732612952000+0000)/","Clicks":11,"Impressions":60,"AvgClickPosition":5,"AvgImpressionPosition":7}]}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.GetQueryPageStats(context.Background(), "https://example.test", "bing api")
	if err != nil {
		t.Fatalf("GetQueryPageStats error = %v", err)
	}

	if method != http.MethodGet {
		t.Fatalf("method = %q, want GET", method)
	}
	if result.Query != "bing api" {
		t.Fatalf("Query = %q, want %q", result.Query, "bing api")
	}
	if len(result.Rows) != 1 {
		t.Fatalf("len(Rows) = %d, want 1", len(result.Rows))
	}
	if result.Rows[0].Page != "https://example.test/page" {
		t.Fatalf("Page = %q, want %q", result.Rows[0].Page, "https://example.test/page")
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
