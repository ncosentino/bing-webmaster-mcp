package bingwebmaster

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRemoveSite_SendsPOSTBody(t *testing.T) {
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

	result, err := client.RemoveSite(context.Background(), "https://example.test")
	if err != nil {
		t.Fatalf("RemoveSite error = %v", err)
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
	if !result.Success {
		t.Fatal("expected Success true")
	}
}

func TestGetSiteRoles_UsesBooleanQueryAndDecodesRole(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	var method string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		query := r.URL.Query()
		if got := query.Get("siteUrl"); got != "https://example.test" {
			t.Fatalf("siteUrl query param = %q, want %q", got, "https://example.test")
		}
		if got := query.Get("includeAllSubdomains"); got != "true" {
			t.Fatalf("includeAllSubdomains query param = %q, want true", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"d":[{"Date":"/Date(1732612952000+0000)/","DelegatedCode":"delegated-code","DelegatorEmail":"owner@example.test","DelegatedCodeOwnerEmail":"code-owner@example.test","Email":"reader@example.test","Expired":false,"Role":2,"Site":"https://example.test","VerificationSite":"https://verify.example.test"}]}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.GetSiteRoles(context.Background(), "https://example.test", true)
	if err != nil {
		t.Fatalf("GetSiteRoles error = %v", err)
	}

	if method != http.MethodGet {
		t.Fatalf("method = %q, want GET", method)
	}
	if result.RowCount != 1 {
		t.Fatalf("RowCount = %d, want 1", result.RowCount)
	}
	if result.Rows[0].Role != "ReadWrite" {
		t.Fatalf("Role = %q, want %q", result.Rows[0].Role, "ReadWrite")
	}
	if result.Rows[0].Date == nil {
		t.Fatal("expected Date")
	}
}

func TestAddSiteRole_SendsPOSTBody(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	var method string
	var contentType string
	var capturedPath string
	var body rawAddSiteRoleRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		contentType = r.Header.Get("Content-Type")
		capturedPath = r.URL.Path
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"d":null}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.AddSiteRole(context.Background(), "https://example.test", "https://delegated.example.test", "reader@example.test", "auth-code", true, false)
	if err != nil {
		t.Fatalf("AddSiteRole error = %v", err)
	}

	if method != http.MethodPost {
		t.Fatalf("method = %q, want POST", method)
	}
	// Regression check: Bing's real endpoint is "AddSiteRoles" (plural) -- an earlier version of
	// this client incorrectly called "AddSiteRole" (singular), which would 404 against the live API.
	if capturedPath != "/AddSiteRoles" {
		t.Fatalf("path = %q, want %q", capturedPath, "/AddSiteRoles")
	}
	if contentType != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", contentType)
	}
	if body.DelegatedURL != "https://delegated.example.test" {
		t.Fatalf("DelegatedURL = %q, want %q", body.DelegatedURL, "https://delegated.example.test")
	}
	if !body.IsAdministrator {
		t.Fatal("expected IsAdministrator true")
	}
	if result.UserEmail != "reader@example.test" {
		t.Fatalf("UserEmail = %q, want %q", result.UserEmail, "reader@example.test")
	}
}

func TestRemoveSiteRole_ConstructsNestedSiteRoleBody(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	var method string
	var body rawSiteRoleCommandRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"d":null}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.RemoveSiteRole(context.Background(), "https://example.test", "reader@example.test", "ReadOnly")
	if err != nil {
		t.Fatalf("RemoveSiteRole error = %v", err)
	}

	if method != http.MethodPost {
		t.Fatalf("method = %q, want POST", method)
	}
	if body.SiteURL != "https://example.test" {
		t.Fatalf("SiteURL = %q, want %q", body.SiteURL, "https://example.test")
	}
	if body.SiteRole.Role != 1 {
		t.Fatalf("Role = %d, want 1", body.SiteRole.Role)
	}
	if body.SiteRole.Site != "https://example.test" {
		t.Fatalf("Site = %q, want %q", body.SiteRole.Site, "https://example.test")
	}
	if body.SiteRole.VerificationSite != "https://example.test" {
		t.Fatalf("VerificationSite = %q, want %q", body.SiteRole.VerificationSite, "https://example.test")
	}
	if _, err := parseDotNetDate(body.SiteRole.Date); err != nil {
		t.Fatalf("Date = %q, parseDotNetDate error = %v", body.SiteRole.Date, err)
	}
	if result.Role != "ReadOnly" {
		t.Fatalf("Role = %q, want %q", result.Role, "ReadOnly")
	}
}

func TestGetBlockedURLs_MapsDecodedEnums(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("siteUrl"); got != "https://example.test" {
			t.Fatalf("siteUrl query param = %q, want %q", got, "https://example.test")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"d":[{"Date":"/Date(1732612952000+0000)/","EntityType":1,"RequestType":1,"Url":"https://example.test/private/"}]}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.GetBlockedURLs(context.Background(), "https://example.test")
	if err != nil {
		t.Fatalf("GetBlockedURLs error = %v", err)
	}

	if result.RowCount != 1 {
		t.Fatalf("RowCount = %d, want 1", result.RowCount)
	}
	if result.Rows[0].EntityType != "Directory" {
		t.Fatalf("EntityType = %q, want %q", result.Rows[0].EntityType, "Directory")
	}
	if result.Rows[0].RequestType != "FullRemoval" {
		t.Fatalf("RequestType = %q, want %q", result.Rows[0].RequestType, "FullRemoval")
	}
}

func TestAddBlockedURL_UsesNestedBodyAndDefaults(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	var method string
	var body rawBlockedURLCommandRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"d":null}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.AddBlockedURL(context.Background(), "https://example.test", "https://example.test/private/", "", "")
	if err != nil {
		t.Fatalf("AddBlockedURL error = %v", err)
	}

	if method != http.MethodPost {
		t.Fatalf("method = %q, want POST", method)
	}
	if body.BlockedURL.EntityType != 0 {
		t.Fatalf("EntityType = %d, want 0", body.BlockedURL.EntityType)
	}
	if body.BlockedURL.RequestType != 0 {
		t.Fatalf("RequestType = %d, want 0", body.BlockedURL.RequestType)
	}
	if _, err := parseDotNetDate(body.BlockedURL.Date); err != nil {
		t.Fatalf("Date = %q, parseDotNetDate error = %v", body.BlockedURL.Date, err)
	}
	if result.EntityType != "Page" || result.RequestType != "CacheOnly" {
		t.Fatalf("defaults = %q/%q, want Page/CacheOnly", result.EntityType, result.RequestType)
	}
}

func TestRemoveBlockedURL_UsesNestedBodyAndDefaults(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	var body rawBlockedURLCommandRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"d":null}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.RemoveBlockedURL(context.Background(), "https://example.test", "https://example.test/private/", "", "")
	if err != nil {
		t.Fatalf("RemoveBlockedURL error = %v", err)
	}

	if body.BlockedURL.RequestType != 1 {
		t.Fatalf("RequestType = %d, want 1", body.BlockedURL.RequestType)
	}
	if result.RequestType != "FullRemoval" {
		t.Fatalf("RequestType = %q, want %q", result.RequestType, "FullRemoval")
	}
}

func TestGetQueryPageDetailStats_UsesParamsAndMapsRows(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	var method string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		query := r.URL.Query()
		if got := query.Get("query"); got != "bing api" {
			t.Fatalf("query query param = %q, want %q", got, "bing api")
		}
		if got := query.Get("page"); got != "https://example.test/page" {
			t.Fatalf("page query param = %q, want %q", got, "https://example.test/page")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"d":[{"Date":"/Date(1732612952000+0000)/","Clicks":3,"Impressions":20,"Position":7}]}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.GetQueryPageDetailStats(context.Background(), "https://example.test", "bing api", "https://example.test/page")
	if err != nil {
		t.Fatalf("GetQueryPageDetailStats error = %v", err)
	}

	if method != http.MethodGet {
		t.Fatalf("method = %q, want GET", method)
	}
	if result.RowCount != 1 || result.Rows[0].Position != 7 {
		t.Fatalf("got %+v, want one row with position 7", result.Rows)
	}
}

func TestGetQueryTrafficStats_UsesParamsAndMapsRows(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if got := query.Get("siteUrl"); got != "https://example.test" {
			t.Fatalf("siteUrl query param = %q, want %q", got, "https://example.test")
		}
		if got := query.Get("query"); got != "bing api" {
			t.Fatalf("query query param = %q, want %q", got, "bing api")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"d":[{"Date":"/Date(1732612952000+0000)/","Clicks":12,"Impressions":90}]}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.GetQueryTrafficStats(context.Background(), "https://example.test", "bing api")
	if err != nil {
		t.Fatalf("GetQueryTrafficStats error = %v", err)
	}

	if result.RowCount != 1 {
		t.Fatalf("RowCount = %d, want 1", result.RowCount)
	}
	if result.Rows[0].Clicks != 12 {
		t.Fatalf("Clicks = %d, want 12", result.Rows[0].Clicks)
	}
}

func TestGetKeyword_UsesDateRangeAndMapsFoundResult(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if got := query.Get("q"); got != "golang" {
			t.Fatalf("q = %q, want %q", got, "golang")
		}
		if got := query.Get("startDate"); got != "2024-01-01" {
			t.Fatalf("startDate = %q, want %q", got, "2024-01-01")
		}
		if got := query.Get("endDate"); got != "2024-01-31" {
			t.Fatalf("endDate = %q, want %q", got, "2024-01-31")
		}
		if got := query.Get("siteUrl"); got != "" {
			t.Fatalf("siteUrl should be empty, got %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"d":{"Query":"golang","BroadImpressions":200,"Impressions":75}}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.GetKeyword(context.Background(), "golang", "US", "en", "2024-01-01", "2024-01-31")
	if err != nil {
		t.Fatalf("GetKeyword error = %v", err)
	}

	if !result.Found {
		t.Fatal("expected Found true")
	}
	if result.Impressions != 75 || result.BroadImpressions != 200 {
		t.Fatalf("got impressions %+v, want 75/200", result)
	}
}

func TestGetKeyword_TreatsMissingQueryAsNotFound(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"d":{"BroadImpressions":0,"Impressions":0}}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.GetKeyword(context.Background(), "golang", "US", "en", "2024-01-01", "2024-01-31")
	if err != nil {
		t.Fatalf("GetKeyword error = %v", err)
	}

	if result.Found {
		t.Fatal("expected Found false")
	}
	if result.Query != "golang" {
		t.Fatalf("Query = %q, want request query", result.Query)
	}
}

func TestGetRelatedKeywords_UsesDateRangeAndMapsRows(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if got := query.Get("q"); got != "golang" {
			t.Fatalf("q = %q, want %q", got, "golang")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"d":[{"Query":"golang tutorial","BroadImpressions":300,"Impressions":120}]}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.GetRelatedKeywords(context.Background(), "golang", "US", "en", "2024-01-01", "2024-01-31")
	if err != nil {
		t.Fatalf("GetRelatedKeywords error = %v", err)
	}

	if result.RowCount != 1 {
		t.Fatalf("RowCount = %d, want 1", result.RowCount)
	}
	if result.Rows[0].Query != "golang tutorial" {
		t.Fatalf("Query = %q, want %q", result.Rows[0].Query, "golang tutorial")
	}
}

func TestGetChildrenURLInfo_SendsNestedFiltersAndMapsRows(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	var method string
	var body rawChildrenURLInfoRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"d":[{"Url":"https://example.test/child","IsPage":true,"HttpStatus":200,"DocumentSize":1024,"AnchorCount":3,"DiscoveryDate":"/Date(1732612952000+0000)/","LastCrawledDate":"/Date(1732616552000+0000)/","TotalChildUrlCount":2}]}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.GetChildrenURLInfo(context.Background(), "https://example.test", "https://example.test/parent", 4, "LastTwoWeeks", "LastMonth", "IsMalware", "Code4xx")
	if err != nil {
		t.Fatalf("GetChildrenURLInfo error = %v", err)
	}

	if method != http.MethodPost {
		t.Fatalf("method = %q, want POST", method)
	}
	if body.Page != 4 {
		t.Fatalf("Page = %d, want 4", body.Page)
	}
	if body.FilterProperties.CrawlDateFilter != 2 {
		t.Fatalf("CrawlDateFilter = %d, want 2", body.FilterProperties.CrawlDateFilter)
	}
	if body.FilterProperties.DiscoveredDateFilter != 2 {
		t.Fatalf("DiscoveredDateFilter = %d, want 2", body.FilterProperties.DiscoveredDateFilter)
	}
	if body.FilterProperties.DocFlagsFilters != 2 {
		t.Fatalf("DocFlagsFilters = %d, want 2", body.FilterProperties.DocFlagsFilters)
	}
	if body.FilterProperties.HTTPCodeFilters != 16 {
		t.Fatalf("HTTPCodeFilters = %d, want 16", body.FilterProperties.HTTPCodeFilters)
	}
	if result.RowCount != 1 || result.Rows[0].HTTPStatus != 200 {
		t.Fatalf("got %+v, want one row with HTTPStatus 200", result.Rows)
	}
}

func TestGetChildrenURLTrafficInfo_UsesParamsAndMapsRows(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if got := query.Get("page"); got != "3" {
			t.Fatalf("page query param = %q, want 3", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"d":[{"Url":"https://example.test/child","IsPage":false,"Clicks":7,"Impressions":55}]}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.GetChildrenURLTrafficInfo(context.Background(), "https://example.test", "https://example.test/parent", 3)
	if err != nil {
		t.Fatalf("GetChildrenURLTrafficInfo error = %v", err)
	}

	if result.RowCount != 1 {
		t.Fatalf("RowCount = %d, want 1", result.RowCount)
	}
	if result.Rows[0].Impressions != 55 {
		t.Fatalf("Impressions = %d, want 55", result.Rows[0].Impressions)
	}
}

func TestFetchURL_SendsPOSTBody(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	var method string
	var body map[string]string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"d":null}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.FetchURL(context.Background(), "https://example.test", "https://example.test/page")
	if err != nil {
		t.Fatalf("FetchURL error = %v", err)
	}

	if method != http.MethodPost {
		t.Fatalf("method = %q, want POST", method)
	}
	if body["url"] != "https://example.test/page" {
		t.Fatalf("url = %q, want %q", body["url"], "https://example.test/page")
	}
	if !result.Success {
		t.Fatal("expected Success true")
	}
}

func TestListFetchedURLs_MapsRows(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"d":[{"Date":"/Date(1732612952000+0000)/","Expired":false,"Fetched":true,"Url":"https://example.test/page"}]}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.ListFetchedURLs(context.Background(), "https://example.test")
	if err != nil {
		t.Fatalf("ListFetchedURLs error = %v", err)
	}

	if result.RowCount != 1 {
		t.Fatalf("RowCount = %d, want 1", result.RowCount)
	}
	if !result.Rows[0].Fetched || result.Rows[0].Date == nil {
		t.Fatalf("got %+v, want fetched row with date", result.Rows[0])
	}
}

func TestGetFetchedURLDetails_UsesParamsAndMapsObject(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if got := query.Get("url"); got != "https://example.test/page" {
			t.Fatalf("url query param = %q, want %q", got, "https://example.test/page")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"d":{"Date":"/Date(1732612952000+0000)/","Document":"PGh0bWw+PC9odG1sPg==","Headers":"HTTP/1.1 200 OK","Status":"Success","Url":"https://example.test/page"}}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.GetFetchedURLDetails(context.Background(), "https://example.test", "https://example.test/page")
	if err != nil {
		t.Fatalf("GetFetchedURLDetails error = %v", err)
	}

	if result.URL != "https://example.test/page" {
		t.Fatalf("URL = %q, want %q", result.URL, "https://example.test/page")
	}
	if result.Status != "Success" {
		t.Fatalf("Status = %q, want %q", result.Status, "Success")
	}
}

func TestRemoveSitemap_SendsPOSTBody(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	var body map[string]string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"d":null}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.RemoveSitemap(context.Background(), "https://example.test", "https://example.test/sitemap.xml")
	if err != nil {
		t.Fatalf("RemoveSitemap error = %v", err)
	}

	if body["feedUrl"] != "https://example.test/sitemap.xml" {
		t.Fatalf("feedUrl = %q, want %q", body["feedUrl"], "https://example.test/sitemap.xml")
	}
	if !result.Success {
		t.Fatal("expected Success true")
	}
}

func TestGetSiteMoves_MapsDecodedEnums(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"d":[{"Date":"/Date(1732612952000+0000)/","MoveScope":2,"MoveType":1,"SourceUrl":"https://example.test/old/","TargetUrl":"https://example.test/new/"}]}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.GetSiteMoves(context.Background(), "https://example.test")
	if err != nil {
		t.Fatalf("GetSiteMoves error = %v", err)
	}

	if result.RowCount != 1 {
		t.Fatalf("RowCount = %d, want 1", result.RowCount)
	}
	if result.Rows[0].MoveScope != "Directory" {
		t.Fatalf("MoveScope = %q, want %q", result.Rows[0].MoveScope, "Directory")
	}
	if result.Rows[0].MoveType != "Global" {
		t.Fatalf("MoveType = %q, want %q", result.Rows[0].MoveType, "Global")
	}
}

func TestSubmitSiteMove_UsesNestedBodyAndDefaults(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	var body rawSiteMoveCommandRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"d":null}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.SubmitSiteMove(context.Background(), "https://example.test", "https://example.test/old/", "https://example.test/new/", "", "")
	if err != nil {
		t.Fatalf("SubmitSiteMove error = %v", err)
	}

	if body.Settings.MoveType != 0 {
		t.Fatalf("MoveType = %d, want 0", body.Settings.MoveType)
	}
	if body.Settings.MoveScope != 0 {
		t.Fatalf("MoveScope = %d, want 0", body.Settings.MoveScope)
	}
	if _, err := parseDotNetDate(body.Settings.Date); err != nil {
		t.Fatalf("Date = %q, parseDotNetDate error = %v", body.Settings.Date, err)
	}
	if result.MoveType != "Local" || result.MoveScope != "Domain" {
		t.Fatalf("defaults = %q/%q, want Local/Domain", result.MoveType, result.MoveScope)
	}
}

func TestSubmitContent_UsesPOSTBodyAndMapsDynamicServing(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	var body rawSubmitContentRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"d":null}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.SubmitContent(context.Background(), "https://example.test", "https://example.test/page", "aHR0cA==", "eyJAY29udGV4dCI6IiJ9", "Amp")
	if err != nil {
		t.Fatalf("SubmitContent error = %v", err)
	}

	if body.DynamicServing != 3 {
		t.Fatalf("DynamicServing = %d, want 3", body.DynamicServing)
	}
	if body.HTTPMessage != "aHR0cA==" {
		t.Fatalf("HTTPMessage = %q, want %q", body.HTTPMessage, "aHR0cA==")
	}
	if result.DynamicServing != "Amp" {
		t.Fatalf("DynamicServing = %q, want %q", result.DynamicServing, "Amp")
	}
}

func TestSubmitContent_InvalidDynamicServingReturnsClientError(t *testing.T) {
	client := NewClient("test-api-key")

	_, err := client.SubmitContent(context.Background(), "https://example.test", "https://example.test/page", "aHR0cA==", "eyJAY29udGV4dCI6IiJ9", "Spaceship")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), `invalid dynamic_serving "Spaceship"`) {
		t.Fatalf("error = %q, want invalid dynamic_serving message", err.Error())
	}
}

func TestGetContentSubmissionQuota_UsesSiteURLAndMapsQuota(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	var method string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		if got := r.URL.Query().Get("siteUrl"); got != "https://example.test" {
			t.Fatalf("siteUrl query param = %q, want %q", got, "https://example.test")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"d":{"DailyQuota":25,"MonthlyQuota":500}}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.GetContentSubmissionQuota(context.Background(), "https://example.test")
	if err != nil {
		t.Fatalf("GetContentSubmissionQuota error = %v", err)
	}

	if method != http.MethodGet {
		t.Fatalf("method = %q, want GET", method)
	}
	if result.DailyQuota != 25 || result.MonthlyQuota != 500 {
		t.Fatalf("got %+v, want DailyQuota 25 and MonthlyQuota 500", result)
	}
}
