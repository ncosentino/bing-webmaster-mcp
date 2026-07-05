package bingwebmaster

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const httpTimeout = 30 * time.Second

// apiBaseURL is the base URL for the Bing Webmaster JSON API.
// It is a variable so tests can override it to point at a local test server.
var apiBaseURL = "https://ssl.bing.com/webmaster/api.svc/json"

type apiRequestError struct {
	StatusCode int
	Body       string
}

func (e *apiRequestError) Error() string {
	return fmt.Sprintf("Bing Webmaster API returned HTTP %d: %s", e.StatusCode, e.Body)
}

// Client calls the Bing Webmaster Tools API.
type Client struct {
	httpClient *http.Client
	apiKey     string
}

// NewClient creates a Client authenticated with the provided API key.
func NewClient(apiKey string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: httpTimeout},
		apiKey:     apiKey,
	}
}

// ListSites returns all Bing Webmaster sites accessible to the API key.
func (c *Client) ListSites(ctx context.Context) (*siteList, error) {
	var raw []rawSite
	if err := c.get(ctx, "GetUserSites", nil, &raw); err != nil {
		return nil, err
	}

	sites := make([]site, len(raw))
	for i, item := range raw {
		sites[i] = site{
			SiteURL:             item.URL,
			IsVerified:          item.IsVerified,
			DNSVerificationCode: item.DNSVerificationCode,
			AuthenticationCode:  item.AuthenticationCode,
		}
	}

	return &siteList{Sites: sites, QueriedAt: time.Now().UTC()}, nil
}

// AddSite adds a site to Bing Webmaster Tools.
func (c *Client) AddSite(ctx context.Context, siteURL string) (*addSiteResult, error) {
	// Bing's AddSite "d" payload is not a reliable success indicator -- empirically it returns
	// "d":null both for a brand-new site and a no-op repeat of an already-added site. Success is
	// defined as "the HTTP call completed without error"; confirm the actual outcome via
	// ListSites if you need to know whether the site is present.
	if err := c.postCommand(ctx, "AddSite", map[string]string{"siteUrl": siteURL}); err != nil {
		return nil, err
	}

	return &addSiteResult{SiteURL: siteURL, Success: true, RequestedAt: time.Now().UTC()}, nil
}

// VerifySite verifies a site in Bing Webmaster Tools.
func (c *Client) VerifySite(ctx context.Context, siteURL string) (*verifySiteResult, error) {
	verified, err := c.postQuery(ctx, "VerifySite", map[string]string{"siteUrl": siteURL})
	if err != nil {
		return nil, err
	}

	return &verifySiteResult{
		SiteURL:     siteURL,
		Verified:    verified,
		RequestedAt: time.Now().UTC(),
	}, nil
}

// ListSitemaps returns submitted sitemaps for a site.
func (c *Client) ListSitemaps(ctx context.Context, siteURL string) (*sitemapList, error) {
	var raw []rawFeed
	if err := c.get(ctx, "GetFeeds", map[string]string{"siteUrl": siteURL}, &raw); err != nil {
		return nil, err
	}

	return &sitemapList{
		SiteURL:   siteURL,
		Sitemaps:  mapFeeds(raw),
		QueriedAt: time.Now().UTC(),
	}, nil
}

// GetSitemapDetails returns details for a submitted sitemap.
//
// Bing's real response shape for this endpoint is unconfirmed by Microsoft's public docs.
// Live testing against a real account showed it can return either a single feed object or an
// array (e.g. for sitemap index files with multiple constituent feeds) -- this handles both.
func (c *Client) GetSitemapDetails(ctx context.Context, siteURL string, feedURL string) (*sitemapDetailResult, error) {
	var rawEnvelope json.RawMessage
	if err := c.get(ctx, "GetFeedDetails", map[string]string{"siteUrl": siteURL, "feedUrl": feedURL}, &rawEnvelope); err != nil {
		return nil, err
	}

	sitemap, err := parseFeedDetailsPayload(rawEnvelope)
	if err != nil {
		return nil, fmt.Errorf("parsing GetFeedDetails payload: %w", err)
	}

	return &sitemapDetailResult{
		SiteURL:   siteURL,
		FeedURL:   feedURL,
		Sitemap:   sitemap,
		QueriedAt: time.Now().UTC(),
	}, nil
}

// parseFeedDetailsPayload accepts either a single feed object or an array of feed objects
// (using the first element), returning nil if the payload is empty or an empty array.
func parseFeedDetailsPayload(payload json.RawMessage) (*feed, error) {
	trimmed := bytes.TrimSpace(payload)
	if len(trimmed) == 0 || string(trimmed) == "null" {
		return nil, nil
	}

	if trimmed[0] == '[' {
		var rawFeeds []rawFeed
		if err := json.Unmarshal(trimmed, &rawFeeds); err != nil {
			return nil, err
		}
		if len(rawFeeds) == 0 {
			return nil, nil
		}
		mapped := mapFeed(rawFeeds[0])
		return &mapped, nil
	}

	var raw rawFeed
	if err := json.Unmarshal(trimmed, &raw); err != nil {
		return nil, err
	}
	mapped := mapFeed(raw)
	return &mapped, nil
}

// SubmitSitemap submits a sitemap feed to Bing Webmaster Tools.
func (c *Client) SubmitSitemap(ctx context.Context, siteURL string, feedURL string) (*submitSitemapResult, error) {
	if err := c.postCommand(ctx, "SubmitFeed", map[string]string{"siteUrl": siteURL, "feedUrl": feedURL}); err != nil {
		return nil, err
	}

	return &submitSitemapResult{
		SiteURL:     siteURL,
		FeedURL:     feedURL,
		Success:     true,
		SubmittedAt: time.Now().UTC(),
	}, nil
}

// SubmitURL submits a single URL to Bing Webmaster Tools.
func (c *Client) SubmitURL(ctx context.Context, siteURL string, submittedURL string) (*submitURLResult, error) {
	if err := c.postCommand(ctx, "SubmitUrl", map[string]string{"siteUrl": siteURL, "url": submittedURL}); err != nil {
		return nil, err
	}

	return &submitURLResult{
		SiteURL:     siteURL,
		URL:         submittedURL,
		Success:     true,
		SubmittedAt: time.Now().UTC(),
	}, nil
}

// SubmitURLBatch submits up to 500 URLs to Bing Webmaster Tools.
func (c *Client) SubmitURLBatch(ctx context.Context, siteURL string, urlList []string) (*submitURLBatchResult, error) {
	if len(urlList) > 500 {
		return nil, fmt.Errorf("urlList contains %d URLs; Bing SubmitUrlBatch supports at most 500", len(urlList))
	}

	if err := c.postCommand(ctx, "SubmitUrlBatch", map[string]any{"siteUrl": siteURL, "urlList": urlList}); err != nil {
		return nil, err
	}

	return &submitURLBatchResult{
		SiteURL:        siteURL,
		URLList:        urlList,
		SubmittedCount: len(urlList),
		Success:        true,
		SubmittedAt:    time.Now().UTC(),
	}, nil
}

// GetURLSubmissionQuota returns the URL submission quotas for a site.
func (c *Client) GetURLSubmissionQuota(ctx context.Context, siteURL string) (*urlSubmissionQuotaResult, error) {
	var raw rawQuota
	if err := c.get(ctx, "GetUrlSubmissionQuota", map[string]string{"siteUrl": siteURL}, &raw); err != nil {
		return nil, err
	}

	return &urlSubmissionQuotaResult{
		SiteURL:      siteURL,
		DailyQuota:   raw.DailyQuota,
		MonthlyQuota: raw.MonthlyQuota,
		QueriedAt:    time.Now().UTC(),
	}, nil
}

// GetCrawlIssues returns crawl issues for a site.
func (c *Client) GetCrawlIssues(ctx context.Context, siteURL string) (*crawlIssuesResult, error) {
	var raw []rawCrawlIssue
	if err := c.get(ctx, "GetCrawlIssues", map[string]string{"siteUrl": siteURL}, &raw); err != nil {
		return nil, err
	}

	issues := make([]crawlIssue, len(raw))
	for i, item := range raw {
		issues[i] = crawlIssue{
			URL:      item.URL,
			HTTPCode: item.HTTPCode,
			Issues:   decodeCrawlIssueFlags(item.Issues),
			InLinks:  item.InLinks,
		}
	}

	return &crawlIssuesResult{SiteURL: siteURL, Issues: issues, QueriedAt: time.Now().UTC()}, nil
}

// GetCrawlStats returns crawl statistics for a site.
func (c *Client) GetCrawlStats(ctx context.Context, siteURL string) (*crawlStatsResult, error) {
	var raw []rawCrawlStat
	if err := c.get(ctx, "GetCrawlStats", map[string]string{"siteUrl": siteURL}, &raw); err != nil {
		return nil, err
	}

	stats := make([]crawlStat, len(raw))
	for i, item := range raw {
		stats[i] = crawlStat{
			Date:               timePointer(item.Date),
			CrawledPages:       item.CrawledPages,
			CrawlErrors:        item.CrawlErrors,
			InIndex:            item.InIndex,
			InLinks:            item.InLinks,
			Code2xx:            item.Code2xx,
			Code301:            item.Code301,
			Code302:            item.Code302,
			Code4xx:            item.Code4xx,
			Code5xx:            item.Code5xx,
			AllOtherCodes:      item.AllOtherCodes,
			BlockedByRobotsTxt: item.BlockedByRobotsTxt,
			ContainsMalware:    item.ContainsMalware,
		}
	}

	return &crawlStatsResult{SiteURL: siteURL, RowCount: len(stats), Rows: stats, QueriedAt: time.Now().UTC()}, nil
}

// GetURLInfo returns index status for a single URL.
func (c *Client) GetURLInfo(ctx context.Context, siteURL string, requestedURL string) (*urlInfoResult, error) {
	var raw rawURLInfo
	if err := c.get(ctx, "GetUrlInfo", map[string]string{"siteUrl": siteURL, "url": requestedURL}, &raw); err != nil {
		return nil, err
	}

	return &urlInfoResult{
		SiteURL:            siteURL,
		URL:                raw.URL,
		IsPage:             raw.IsPage,
		HTTPStatus:         raw.HTTPStatus,
		DocumentSize:       raw.DocumentSize,
		AnchorCount:        raw.AnchorCount,
		DiscoveryDate:      timePointer(raw.DiscoveryDate),
		LastCrawledDate:    timePointer(raw.LastCrawledDate),
		TotalChildURLCount: raw.TotalChildURLCount,
		QueriedAt:          time.Now().UTC(),
	}, nil
}

// GetURLTrafficInfo returns traffic metrics for a single URL.
func (c *Client) GetURLTrafficInfo(ctx context.Context, siteURL string, requestedURL string) (*urlTrafficInfoResult, error) {
	var raw rawURLTrafficInfo
	if err := c.get(ctx, "GetUrlTrafficInfo", map[string]string{"siteUrl": siteURL, "url": requestedURL}, &raw); err != nil {
		return nil, err
	}

	return &urlTrafficInfoResult{
		SiteURL:     siteURL,
		URL:         raw.URL,
		IsPage:      raw.IsPage,
		Clicks:      raw.Clicks,
		Impressions: raw.Impressions,
		QueriedAt:   time.Now().UTC(),
	}, nil
}

// GetURLLinks returns inbound links for a URL.
func (c *Client) GetURLLinks(ctx context.Context, siteURL string, link string, page int) (*urlLinksResult, error) {
	var raw rawURLLinkDetails
	if err := c.get(ctx, "GetUrlLinks", map[string]string{
		"siteUrl": siteURL,
		"link":    link,
		"page":    strconv.Itoa(page),
	}, &raw); err != nil {
		return nil, err
	}

	details := make([]urlLinkDetail, len(raw.Details))
	for i, item := range raw.Details {
		details[i] = urlLinkDetail{AnchorText: item.AnchorText, URL: item.URL}
	}

	return &urlLinksResult{
		SiteURL:    siteURL,
		Link:       link,
		Page:       page,
		Details:    details,
		TotalPages: raw.TotalPages,
		QueriedAt:  time.Now().UTC(),
	}, nil
}

// GetLinkCounts returns pages with inbound link counts.
func (c *Client) GetLinkCounts(ctx context.Context, siteURL string, page int) (*linkCountsResult, error) {
	var raw rawLinkCounts
	if err := c.get(ctx, "GetLinkCounts", map[string]string{
		"siteUrl": siteURL,
		"page":    strconv.Itoa(page),
	}, &raw); err != nil {
		return nil, err
	}

	links := make([]linkCount, len(raw.Links))
	for i, item := range raw.Links {
		links[i] = linkCount{Count: item.Count, URL: item.URL}
	}

	return &linkCountsResult{
		SiteURL:    siteURL,
		Page:       page,
		Links:      links,
		TotalPages: raw.TotalPages,
		QueriedAt:  time.Now().UTC(),
	}, nil
}

// GetRankAndTrafficStats returns clicks and impressions over time.
func (c *Client) GetRankAndTrafficStats(ctx context.Context, siteURL string) (*rankAndTrafficStatsResult, error) {
	var raw []rawRankTrafficStat
	if err := c.get(ctx, "GetRankAndTrafficStats", map[string]string{"siteUrl": siteURL}, &raw); err != nil {
		return nil, err
	}

	stats := make([]rankTrafficStat, len(raw))
	for i, item := range raw {
		stats[i] = rankTrafficStat{Date: timePointer(item.Date), Clicks: item.Clicks, Impressions: item.Impressions}
	}

	return &rankAndTrafficStatsResult{SiteURL: siteURL, RowCount: len(stats), Rows: stats, QueriedAt: time.Now().UTC()}, nil
}

// GetQueryStats returns top query statistics for a site.
func (c *Client) GetQueryStats(ctx context.Context, siteURL string) (*queryStatsResult, error) {
	stats, err := c.getQueryStats(ctx, "GetQueryStats", map[string]string{"siteUrl": siteURL})
	if err != nil {
		return nil, err
	}

	return &queryStatsResult{SiteURL: siteURL, RowCount: len(stats), Rows: stats, QueriedAt: time.Now().UTC()}, nil
}

// GetPageStats returns top page statistics for a site.
func (c *Client) GetPageStats(ctx context.Context, siteURL string) (*pageStatsResult, error) {
	stats, err := c.getPageStats(ctx, "GetPageStats", map[string]string{"siteUrl": siteURL})
	if err != nil {
		return nil, err
	}

	return &pageStatsResult{SiteURL: siteURL, RowCount: len(stats), Rows: stats, QueriedAt: time.Now().UTC()}, nil
}

// GetPageQueryStats returns queries for a specific page.
func (c *Client) GetPageQueryStats(ctx context.Context, siteURL string, page string) (*pageQueryStatsResult, error) {
	stats, err := c.getQueryStats(ctx, "GetPageQueryStats", map[string]string{"siteUrl": siteURL, "page": page})
	if err != nil {
		return nil, err
	}

	return &pageQueryStatsResult{SiteURL: siteURL, Page: page, RowCount: len(stats), Rows: stats, QueriedAt: time.Now().UTC()}, nil
}

// GetQueryPageStats returns pages for a specific query.
func (c *Client) GetQueryPageStats(ctx context.Context, siteURL string, query string) (*queryPageStatsResult, error) {
	stats, err := c.getPageStats(ctx, "GetQueryPageStats", map[string]string{"siteUrl": siteURL, "query": query})
	if err != nil {
		return nil, err
	}

	return &queryPageStatsResult{SiteURL: siteURL, Query: query, RowCount: len(stats), Rows: stats, QueriedAt: time.Now().UTC()}, nil
}

// GetKeywordStats returns market-wide keyword statistics.
func (c *Client) GetKeywordStats(ctx context.Context, query string, country string, language string) (*keywordStatsResult, error) {
	var raw []rawKeywordStat
	if err := c.get(ctx, "GetKeywordStats", map[string]string{"q": query, "country": country, "language": language}, &raw); err != nil {
		return nil, err
	}

	stats := make([]keywordStat, len(raw))
	for i, item := range raw {
		stats[i] = keywordStat{
			Query:            item.Query,
			Date:             timePointer(item.Date),
			Impressions:      item.Impressions,
			BroadImpressions: item.BroadImpressions,
		}
	}

	return &keywordStatsResult{Query: query, Country: country, Language: language, RowCount: len(stats), Rows: stats, QueriedAt: time.Now().UTC()}, nil
}

func (c *Client) getQueryStats(ctx context.Context, method string, params map[string]string) ([]queryStat, error) {
	var raw []rawQueryStat
	if err := c.get(ctx, method, params, &raw); err != nil {
		return nil, err
	}

	stats := make([]queryStat, len(raw))
	for i, item := range raw {
		stats[i] = queryStat{
			Query:                 item.Query,
			Date:                  timePointer(item.Date),
			Clicks:                item.Clicks,
			Impressions:           item.Impressions,
			AvgClickPosition:      item.AvgClickPosition,
			AvgImpressionPosition: item.AvgImpressionPosition,
		}
	}

	return stats, nil
}

func (c *Client) getPageStats(ctx context.Context, method string, params map[string]string) ([]pageStat, error) {
	var raw []rawQueryStat
	if err := c.get(ctx, method, params, &raw); err != nil {
		return nil, err
	}

	stats := make([]pageStat, len(raw))
	for i, item := range raw {
		stats[i] = pageStat{
			Page:                  item.Query,
			Date:                  timePointer(item.Date),
			Clicks:                item.Clicks,
			Impressions:           item.Impressions,
			AvgClickPosition:      item.AvgClickPosition,
			AvgImpressionPosition: item.AvgImpressionPosition,
		}
	}

	return stats, nil
}

func (c *Client) get(ctx context.Context, methodName string, params map[string]string, dest any) error {
	endpoint, err := c.buildURL(methodName, params)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.do(req, dest)
}

// postQuery issues a POST for a command whose "d" payload is a genuine, meaningful boolean
// answer (e.g. VerifySite: "was this site actually verified?").
func (c *Client) postQuery(ctx context.Context, methodName string, body any) (bool, error) {
	var answer bool
	if err := c.post(ctx, methodName, body, &answer); err != nil {
		return false, err
	}

	return answer, nil
}

// postCommand issues a POST for a fire-and-forget command whose "d" payload is not a reliable
// success indicator (confirmed empirically: Bing's AddSite endpoint returns "d":null both for a
// brand-new site and a no-op repeat of an already-added site -- there is no real boolean to
// read). The "d" payload is read into a RawMessage and discarded; success means the HTTP call
// completed without the client throwing an error.
func (c *Client) postCommand(ctx context.Context, methodName string, body any) error {
	var discarded json.RawMessage
	return c.post(ctx, methodName, body, &discarded)
}

func (c *Client) post(ctx context.Context, methodName string, payload any, dest any) error {
	endpoint, err := c.buildURL(methodName, nil)
	if err != nil {
		return err
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshalling request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	return c.do(req, dest)
}

func (c *Client) do(req *http.Request, dest any) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return &apiRequestError{StatusCode: resp.StatusCode, Body: truncate(string(body), 300)}
	}

	envelope := apiEnvelope[json.RawMessage]{}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return fmt.Errorf("parsing API response envelope: %w", err)
	}
	if err := json.Unmarshal(envelope.D, dest); err != nil {
		return fmt.Errorf("parsing API response payload: %w", err)
	}

	return nil
}

func (c *Client) buildURL(methodName string, params map[string]string) (string, error) {
	endpoint, err := url.Parse(apiBaseURL + "/" + methodName)
	if err != nil {
		return "", fmt.Errorf("parsing API base URL: %w", err)
	}

	query := endpoint.Query()
	query.Set("apikey", c.apiKey)
	for key, value := range params {
		query.Set(key, value)
	}
	endpoint.RawQuery = query.Encode()

	return endpoint.String(), nil
}

func mapFeeds(raw []rawFeed) []feed {
	feeds := make([]feed, len(raw))
	for i, item := range raw {
		feeds[i] = mapFeed(item)
	}
	return feeds
}

func mapFeed(item rawFeed) feed {
	return feed{
		URL:         item.URL,
		Type:        item.Type,
		Compressed:  item.Compressed,
		FileSize:    item.FileSize,
		LastCrawled: timePointer(item.LastCrawled),
		Submitted:   timePointer(item.Submitted),
		Status:      item.Status,
		URLCount:    item.URLCount,
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
