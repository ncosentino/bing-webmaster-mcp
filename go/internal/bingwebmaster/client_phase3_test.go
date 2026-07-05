package bingwebmaster

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetQueryParameters_UsesSiteURLAndMapsRows(t *testing.T) {
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
		_, _ = w.Write([]byte(`{"d":[{"Date":"/Date(1732612952000+0000)/","IsEnabled":true,"Parameter":"utm_campaign","Source":0}]}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.GetQueryParameters(context.Background(), "https://example.test")
	if err != nil {
		t.Fatalf("GetQueryParameters error = %v", err)
	}

	if method != http.MethodGet {
		t.Fatalf("method = %q, want GET", method)
	}
	if capturedPath != "/GetQueryParameters" {
		t.Fatalf("path = %q, want %q", capturedPath, "/GetQueryParameters")
	}
	if result.RowCount != 1 {
		t.Fatalf("RowCount = %d, want 1", result.RowCount)
	}
	if result.Parameters[0].Parameter != "utm_campaign" {
		t.Fatalf("Parameter = %q, want %q", result.Parameters[0].Parameter, "utm_campaign")
	}
	if !result.Parameters[0].IsEnabled {
		t.Fatal("expected IsEnabled true")
	}
	if result.Parameters[0].Source != 0 {
		t.Fatalf("Source = %d, want 0", result.Parameters[0].Source)
	}
	if result.Parameters[0].Date == nil {
		t.Fatal("expected Date")
	}
}

func TestAddQueryParameter_SendsPOSTBodyAndPath(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	var method string
	var contentType string
	var capturedPath string
	var body rawQueryParameterCommandRequest
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

	result, err := client.AddQueryParameter(context.Background(), "https://example.test", "utm_campaign")
	if err != nil {
		t.Fatalf("AddQueryParameter error = %v", err)
	}

	if method != http.MethodPost {
		t.Fatalf("method = %q, want POST", method)
	}
	if capturedPath != "/AddQueryParameter" {
		t.Fatalf("path = %q, want %q", capturedPath, "/AddQueryParameter")
	}
	if contentType != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", contentType)
	}
	if body.SiteURL != "https://example.test" {
		t.Fatalf("SiteURL = %q, want %q", body.SiteURL, "https://example.test")
	}
	if body.QueryParameter != "utm_campaign" {
		t.Fatalf("QueryParameter = %q, want %q", body.QueryParameter, "utm_campaign")
	}
	if result.QueryParameter != "utm_campaign" || !result.Success {
		t.Fatalf("result = %#v, want echoed query parameter and success", result)
	}
}

func TestRemoveQueryParameter_SendsPOSTBodyAndPath(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	var method string
	var capturedPath string
	var body rawQueryParameterCommandRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
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

	result, err := client.RemoveQueryParameter(context.Background(), "https://example.test", "utm_campaign")
	if err != nil {
		t.Fatalf("RemoveQueryParameter error = %v", err)
	}

	if method != http.MethodPost {
		t.Fatalf("method = %q, want POST", method)
	}
	if capturedPath != "/RemoveQueryParameter" {
		t.Fatalf("path = %q, want %q", capturedPath, "/RemoveQueryParameter")
	}
	if body.QueryParameter != "utm_campaign" {
		t.Fatalf("QueryParameter = %q, want %q", body.QueryParameter, "utm_campaign")
	}
	if !result.Success {
		t.Fatal("expected Success true")
	}
}

func TestEnableDisableQueryParameter_SendsPOSTBodyAndEchoesState(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	var method string
	var capturedPath string
	var body rawEnableDisableQueryParameterRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
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

	result, err := client.EnableDisableQueryParameter(context.Background(), "https://example.test", "utm_campaign", false)
	if err != nil {
		t.Fatalf("EnableDisableQueryParameter error = %v", err)
	}

	if method != http.MethodPost {
		t.Fatalf("method = %q, want POST", method)
	}
	if capturedPath != "/EnableDisableQueryParameter" {
		t.Fatalf("path = %q, want %q", capturedPath, "/EnableDisableQueryParameter")
	}
	if body.IsEnabled {
		t.Fatal("expected IsEnabled false")
	}
	if result.IsEnabled {
		t.Fatal("expected result IsEnabled false")
	}
}

func TestGetCountryRegionSettings_UsesSiteURLAndDecodesType(t *testing.T) {
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
		_, _ = w.Write([]byte(`{"d":[{"Date":"/Date(1732612952000+0000)/","TwoLetterIsoCountryCode":"us","Type":3,"Url":"https://blog.example.test/"}]}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.GetCountryRegionSettings(context.Background(), "https://example.test")
	if err != nil {
		t.Fatalf("GetCountryRegionSettings error = %v", err)
	}

	if method != http.MethodGet {
		t.Fatalf("method = %q, want GET", method)
	}
	if capturedPath != "/GetCountryRegionSettings" {
		t.Fatalf("path = %q, want %q", capturedPath, "/GetCountryRegionSettings")
	}
	if result.RowCount != 1 {
		t.Fatalf("RowCount = %d, want 1", result.RowCount)
	}
	if result.Settings[0].SettingsType != "Subdomain" {
		t.Fatalf("SettingsType = %q, want %q", result.Settings[0].SettingsType, "Subdomain")
	}
	if result.Settings[0].Date == nil {
		t.Fatal("expected Date")
	}
}

func TestAddCountryRegionSettings_UsesNestedBody(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	var method string
	var capturedPath string
	var body rawCountryRegionSettingsCommandRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
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

	result, err := client.AddCountryRegionSettings(context.Background(), "https://example.test", "us", "Domain", "https://example.test/")
	if err != nil {
		t.Fatalf("AddCountryRegionSettings error = %v", err)
	}

	if method != http.MethodPost {
		t.Fatalf("method = %q, want POST", method)
	}
	if capturedPath != "/AddCountryRegionSettings" {
		t.Fatalf("path = %q, want %q", capturedPath, "/AddCountryRegionSettings")
	}
	if body.Settings.Type != 2 {
		t.Fatalf("Type = %d, want 2", body.Settings.Type)
	}
	if body.Settings.TwoLetterIsoCountryCode != "us" {
		t.Fatalf("TwoLetterIsoCountryCode = %q, want %q", body.Settings.TwoLetterIsoCountryCode, "us")
	}
	if _, err := parseDotNetDate(body.Settings.Date); err != nil {
		t.Fatalf("Date = %q, parseDotNetDate error = %v", body.Settings.Date, err)
	}
	if result.SettingsType != "Domain" || !result.Success {
		t.Fatalf("result = %#v, want Domain/success", result)
	}
}

func TestRemoveCountryRegionSettings_UsesNestedBody(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	var method string
	var capturedPath string
	var body rawCountryRegionSettingsCommandRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
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

	result, err := client.RemoveCountryRegionSettings(context.Background(), "https://example.test", "us", "Directory", "https://example.test/store/")
	if err != nil {
		t.Fatalf("RemoveCountryRegionSettings error = %v", err)
	}

	if method != http.MethodPost {
		t.Fatalf("method = %q, want POST", method)
	}
	if capturedPath != "/RemoveCountryRegionSettings" {
		t.Fatalf("path = %q, want %q", capturedPath, "/RemoveCountryRegionSettings")
	}
	if body.Settings.Type != 1 {
		t.Fatalf("Type = %d, want 1", body.Settings.Type)
	}
	if body.Settings.URL != "https://example.test/store/" {
		t.Fatalf("URL = %q, want %q", body.Settings.URL, "https://example.test/store/")
	}
	if !result.Success {
		t.Fatal("expected Success true")
	}
}

func TestGetConnectedPages_UsesSiteURLAndMapsSubset(t *testing.T) {
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
		_, _ = w.Write([]byte(`{"d":[{"ActualMasterSite":"https://m.example.test/actual","HttpStatusCode":200,"IsBlocked":false,"IsVerified":true,"LastSuccessfullyVerified":"/Date(1732612952000+0000)/","Market":"en-US","RequestedMasterSite":"https://m.example.test/requested","Url":"https://example.test/page"}]}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.GetConnectedPages(context.Background(), "https://example.test")
	if err != nil {
		t.Fatalf("GetConnectedPages error = %v", err)
	}

	if method != http.MethodGet {
		t.Fatalf("method = %q, want GET", method)
	}
	if capturedPath != "/GetConnectedPages" {
		t.Fatalf("path = %q, want %q", capturedPath, "/GetConnectedPages")
	}
	if result.RowCount != 1 {
		t.Fatalf("RowCount = %d, want 1", result.RowCount)
	}
	if result.Pages[0].RequestedMasterSite != "https://m.example.test/requested" {
		t.Fatalf("RequestedMasterSite = %q, want %q", result.Pages[0].RequestedMasterSite, "https://m.example.test/requested")
	}
	if result.Pages[0].LastSuccessfullyVerified == nil {
		t.Fatal("expected LastSuccessfullyVerified")
	}
}

func TestGetConnectedPages_TreatsMinDateAsNil(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	var capturedPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"d":[{"ActualMasterSite":"","HttpStatusCode":404,"IsBlocked":true,"IsVerified":false,"LastSuccessfullyVerified":"/Date(-62135596800000+0000)/","Market":"en-GB","RequestedMasterSite":"https://m.example.test/requested","Url":"https://example.test/missing"}]}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.GetConnectedPages(context.Background(), "https://example.test")
	if err != nil {
		t.Fatalf("GetConnectedPages error = %v", err)
	}

	if capturedPath != "/GetConnectedPages" {
		t.Fatalf("path = %q, want %q", capturedPath, "/GetConnectedPages")
	}
	if result.Pages[0].LastSuccessfullyVerified != nil {
		t.Fatalf("LastSuccessfullyVerified = %v, want nil", result.Pages[0].LastSuccessfullyVerified)
	}
}

func TestAddConnectedPage_SendsPOSTBodyAndPath(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	var method string
	var capturedPath string
	var body rawAddConnectedPageRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
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

	result, err := client.AddConnectedPage(context.Background(), "https://example.test", "http://example.com/some-path")
	if err != nil {
		t.Fatalf("AddConnectedPage error = %v", err)
	}

	if method != http.MethodPost {
		t.Fatalf("method = %q, want POST", method)
	}
	if capturedPath != "/AddConnectedPage" {
		t.Fatalf("path = %q, want %q", capturedPath, "/AddConnectedPage")
	}
	if body.MasterURL != "http://example.com/some-path" {
		t.Fatalf("MasterURL = %q, want %q", body.MasterURL, "http://example.com/some-path")
	}
	if !result.Success {
		t.Fatal("expected Success true")
	}
}

func TestGetActivePagePreviewBlocks_UsesSiteURLAndDecodesReason(t *testing.T) {
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
		_, _ = w.Write([]byte(`{"d":[{"Action":0,"BlockReason":4,"Reason":"4","RefreshReason":0,"SiteUrl":"https://example.test","SubmitDate":"/Date(1732612952000+0000)/","Url":"https://example.test/preview","UserId":"user-id"}]}`))
	}))
	defer srv.Close()

	apiBaseURL = srv.URL
	client := &Client{httpClient: srv.Client(), apiKey: "test-api-key"}

	result, err := client.GetActivePagePreviewBlocks(context.Background(), "https://example.test")
	if err != nil {
		t.Fatalf("GetActivePagePreviewBlocks error = %v", err)
	}

	if method != http.MethodGet {
		t.Fatalf("method = %q, want GET", method)
	}
	if capturedPath != "/GetActivePagePreviewBlocks" {
		t.Fatalf("path = %q, want %q", capturedPath, "/GetActivePagePreviewBlocks")
	}
	if result.RowCount != 1 {
		t.Fatalf("RowCount = %d, want 1", result.RowCount)
	}
	if result.Blocks[0].BlockReason != "Other" {
		t.Fatalf("BlockReason = %q, want %q", result.Blocks[0].BlockReason, "Other")
	}
	if result.Blocks[0].SubmitDate == nil {
		t.Fatal("expected SubmitDate")
	}
}

func TestAddPagePreviewBlock_SendsPOSTBodyAndPath(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	var method string
	var capturedPath string
	var body rawAddPagePreviewBlockRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
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

	result, err := client.AddPagePreviewBlock(context.Background(), "https://example.test", "https://example.test/preview", "Other")
	if err != nil {
		t.Fatalf("AddPagePreviewBlock error = %v", err)
	}

	if method != http.MethodPost {
		t.Fatalf("method = %q, want POST", method)
	}
	if capturedPath != "/AddPagePreviewBlock" {
		t.Fatalf("path = %q, want %q", capturedPath, "/AddPagePreviewBlock")
	}
	if body.Reason != 4 {
		t.Fatalf("Reason = %d, want 4", body.Reason)
	}
	if result.Reason != "Other" || !result.Success {
		t.Fatalf("result = %#v, want Other/success", result)
	}
}

func TestRemovePagePreviewBlock_SendsPOSTBodyAndPath(t *testing.T) {
	previousBaseURL := apiBaseURL
	t.Cleanup(func() { apiBaseURL = previousBaseURL })

	var method string
	var capturedPath string
	var body map[string]string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
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

	result, err := client.RemovePagePreviewBlock(context.Background(), "https://example.test", "https://example.test/preview")
	if err != nil {
		t.Fatalf("RemovePagePreviewBlock error = %v", err)
	}

	if method != http.MethodPost {
		t.Fatalf("method = %q, want POST", method)
	}
	if capturedPath != "/RemovePagePreviewBlock" {
		t.Fatalf("path = %q, want %q", capturedPath, "/RemovePagePreviewBlock")
	}
	if body["siteUrl"] != "https://example.test" {
		t.Fatalf("siteUrl = %q, want %q", body["siteUrl"], "https://example.test")
	}
	if body["url"] != "https://example.test/preview" {
		t.Fatalf("url = %q, want %q", body["url"], "https://example.test/preview")
	}
	if !result.Success {
		t.Fatal("expected Success true")
	}
}
