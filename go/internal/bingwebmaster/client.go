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

// SetBaseURL overrides the Bing Webmaster API base URL. Intended for pointing
// the compiled binary at a local mock server during end-to-end testing; empty
// values are ignored so the real API remains the default in production.
func SetBaseURL(url string) {
	if url != "" {
		apiBaseURL = url
	}
}

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

// RemoveSite removes a site from Bing Webmaster Tools.
func (c *Client) RemoveSite(ctx context.Context, siteURL string) (*removeSiteResult, error) {
	if err := c.postCommand(ctx, "RemoveSite", map[string]string{"siteUrl": siteURL}); err != nil {
		return nil, err
	}

	return &removeSiteResult{SiteURL: siteURL, Success: true, RequestedAt: time.Now().UTC()}, nil
}

// GetSiteRoles returns delegated roles for a site.
func (c *Client) GetSiteRoles(ctx context.Context, siteURL string, includeAllSubdomains bool) (*siteRolesResult, error) {
	var raw []rawSiteRole
	if err := c.get(ctx, "GetSiteRoles", map[string]string{
		"siteUrl":              siteURL,
		"includeAllSubdomains": strconv.FormatBool(includeAllSubdomains),
	}, &raw); err != nil {
		return nil, err
	}

	rows := make([]siteRole, len(raw))
	for i, item := range raw {
		rows[i] = siteRole{
			Email:                   item.Email,
			Role:                    decodeSiteRole(item.Role),
			Site:                    item.Site,
			VerificationSite:        item.VerificationSite,
			Expired:                 item.Expired,
			DelegatorEmail:          item.DelegatorEmail,
			DelegatedCode:           item.DelegatedCode,
			DelegatedCodeOwnerEmail: item.DelegatedCodeOwnerEmail,
			Date:                    timePointer(item.Date),
		}
	}

	return &siteRolesResult{
		SiteURL:              siteURL,
		IncludeAllSubdomains: includeAllSubdomains,
		RowCount:             len(rows),
		Rows:                 rows,
		QueriedAt:            time.Now().UTC(),
	}, nil
}

// AddSiteRole adds a delegated role for a site.
func (c *Client) AddSiteRole(ctx context.Context, siteURL string, delegatedURL string, userEmail string, authenticationCode string, isAdministrator bool, isReadOnly bool) (*addSiteRoleResult, error) {
	payload := rawAddSiteRoleRequest{
		SiteURL:            siteURL,
		DelegatedURL:       delegatedURL,
		UserEmail:          userEmail,
		AuthenticationCode: authenticationCode,
		IsAdministrator:    isAdministrator,
		IsReadOnly:         isReadOnly,
	}
	if err := c.postCommand(ctx, "AddSiteRoles", payload); err != nil {
		return nil, err
	}

	return &addSiteRoleResult{
		SiteURL:         siteURL,
		DelegatedURL:    delegatedURL,
		UserEmail:       userEmail,
		IsAdministrator: isAdministrator,
		IsReadOnly:      isReadOnly,
		Success:         true,
		RequestedAt:     time.Now().UTC(),
	}, nil
}

// RemoveSiteRole removes a delegated role for a site.
func (c *Client) RemoveSiteRole(ctx context.Context, siteURL string, email string, role string) (*removeSiteRoleResult, error) {
	roleValue, err := encodeSiteRole(role)
	if err != nil {
		return nil, err
	}

	// Bing's exact required field set for RemoveSiteRole matching is not confirmed publicly. Mirror
	// the known nested shape and populate the minimum identity fields we have from the tool input:
	// current Date, Email, Role, Site, and VerificationSite. The ancillary delegated-code/email
	// fields are left empty/omitted until live verification proves they are required.
	payload := rawSiteRoleCommandRequest{
		SiteURL: siteURL,
		SiteRole: rawSiteRoleCommand{
			Date:             formatDotNetDate(time.Now().UTC()),
			Email:            email,
			Role:             roleValue,
			Site:             siteURL,
			VerificationSite: siteURL,
		},
	}
	if err := c.postCommand(ctx, "RemoveSiteRole", payload); err != nil {
		return nil, err
	}

	return &removeSiteRoleResult{
		SiteURL:     siteURL,
		Email:       email,
		Role:        role,
		Success:     true,
		RequestedAt: time.Now().UTC(),
	}, nil
}

// GetBlockedURLs returns blocked URL removals for a site.
func (c *Client) GetBlockedURLs(ctx context.Context, siteURL string) (*blockedURLsResult, error) {
	var raw []rawBlockedURL
	if err := c.get(ctx, "GetBlockedUrls", map[string]string{"siteUrl": siteURL}, &raw); err != nil {
		return nil, err
	}

	rows := make([]blockedURL, len(raw))
	for i, item := range raw {
		rows[i] = blockedURL{
			URL:         item.URL,
			EntityType:  decodeBlockedURLEntityType(item.EntityType),
			RequestType: decodeBlockedURLRequestType(item.RequestType),
			Date:        timePointer(item.Date),
		}
	}

	return &blockedURLsResult{SiteURL: siteURL, RowCount: len(rows), Rows: rows, QueriedAt: time.Now().UTC()}, nil
}

// AddBlockedURL adds a blocked URL removal request for a site.
func (c *Client) AddBlockedURL(ctx context.Context, siteURL string, blockedURLValue string, entityType string, requestType string) (*addBlockedURLResult, error) {
	if entityType == "" {
		entityType = "Page"
	}
	if requestType == "" {
		requestType = "CacheOnly"
	}

	entityTypeValue, err := encodeBlockedURLEntityType(entityType)
	if err != nil {
		return nil, err
	}
	requestTypeValue, err := encodeBlockedURLRequestType(requestType)
	if err != nil {
		return nil, err
	}

	payload := rawBlockedURLCommandRequest{
		SiteURL: siteURL,
		BlockedURL: rawBlockedURLCommand{
			Date:        formatDotNetDate(time.Now().UTC()),
			EntityType:  entityTypeValue,
			RequestType: requestTypeValue,
			URL:         blockedURLValue,
		},
	}
	if err := c.postCommand(ctx, "AddBlockedUrl", payload); err != nil {
		return nil, err
	}

	return &addBlockedURLResult{
		SiteURL:     siteURL,
		URL:         blockedURLValue,
		EntityType:  entityType,
		RequestType: requestType,
		Success:     true,
		RequestedAt: time.Now().UTC(),
	}, nil
}

// RemoveBlockedURL removes a blocked URL removal request for a site.
func (c *Client) RemoveBlockedURL(ctx context.Context, siteURL string, blockedURLValue string, entityType string, requestType string) (*removeBlockedURLResult, error) {
	if entityType == "" {
		entityType = "Page"
	}
	if requestType == "" {
		requestType = "FullRemoval"
	}

	entityTypeValue, err := encodeBlockedURLEntityType(entityType)
	if err != nil {
		return nil, err
	}
	requestTypeValue, err := encodeBlockedURLRequestType(requestType)
	if err != nil {
		return nil, err
	}

	payload := rawBlockedURLCommandRequest{
		SiteURL: siteURL,
		BlockedURL: rawBlockedURLCommand{
			Date:        formatDotNetDate(time.Now().UTC()),
			EntityType:  entityTypeValue,
			RequestType: requestTypeValue,
			URL:         blockedURLValue,
		},
	}
	if err := c.postCommand(ctx, "RemoveBlockedUrl", payload); err != nil {
		return nil, err
	}

	return &removeBlockedURLResult{
		SiteURL:     siteURL,
		URL:         blockedURLValue,
		EntityType:  entityType,
		RequestType: requestType,
		Success:     true,
		RequestedAt: time.Now().UTC(),
	}, nil
}

// GetQueryParameters returns the query parameter normalization settings for a site.
func (c *Client) GetQueryParameters(ctx context.Context, siteURL string) (*queryParametersResult, error) {
	var raw []rawQueryParameter
	if err := c.get(ctx, "GetQueryParameters", map[string]string{"siteUrl": siteURL}, &raw); err != nil {
		return nil, err
	}

	parameters := make([]queryParameter, len(raw))
	for i, item := range raw {
		parameters[i] = queryParameter{
			Parameter: item.Parameter,
			IsEnabled: item.IsEnabled,
			Source:    item.Source,
			Date:      timePointer(item.Date),
		}
	}

	return &queryParametersResult{
		SiteURL:    siteURL,
		RowCount:   len(parameters),
		Parameters: parameters,
		QueriedAt:  time.Now().UTC(),
	}, nil
}

// AddQueryParameter adds a query parameter normalization setting for a site.
func (c *Client) AddQueryParameter(ctx context.Context, siteURL string, queryParameter string) (*addQueryParameterResult, error) {
	payload := rawQueryParameterCommandRequest{
		SiteURL:        siteURL,
		QueryParameter: queryParameter,
	}
	if err := c.postCommand(ctx, "AddQueryParameter", payload); err != nil {
		return nil, err
	}

	return &addQueryParameterResult{
		SiteURL:        siteURL,
		QueryParameter: queryParameter,
		Success:        true,
		RequestedAt:    time.Now().UTC(),
	}, nil
}

// RemoveQueryParameter removes a query parameter normalization setting for a site.
func (c *Client) RemoveQueryParameter(ctx context.Context, siteURL string, queryParameter string) (*removeQueryParameterResult, error) {
	payload := rawQueryParameterCommandRequest{
		SiteURL:        siteURL,
		QueryParameter: queryParameter,
	}
	if err := c.postCommand(ctx, "RemoveQueryParameter", payload); err != nil {
		return nil, err
	}

	return &removeQueryParameterResult{
		SiteURL:        siteURL,
		QueryParameter: queryParameter,
		Success:        true,
		RequestedAt:    time.Now().UTC(),
	}, nil
}

// EnableDisableQueryParameter enables or disables a query parameter normalization setting.
func (c *Client) EnableDisableQueryParameter(ctx context.Context, siteURL string, queryParameter string, isEnabled bool) (*enableDisableQueryParameterResult, error) {
	payload := rawEnableDisableQueryParameterRequest{
		SiteURL:        siteURL,
		QueryParameter: queryParameter,
		IsEnabled:      isEnabled,
	}
	if err := c.postCommand(ctx, "EnableDisableQueryParameter", payload); err != nil {
		return nil, err
	}

	return &enableDisableQueryParameterResult{
		SiteURL:        siteURL,
		QueryParameter: queryParameter,
		IsEnabled:      isEnabled,
		Success:        true,
		RequestedAt:    time.Now().UTC(),
	}, nil
}

// GetCountryRegionSettings returns geo-targeting settings for a site.
func (c *Client) GetCountryRegionSettings(ctx context.Context, siteURL string) (*countryRegionSettingsResult, error) {
	var raw []rawCountryRegionSettings
	if err := c.get(ctx, "GetCountryRegionSettings", map[string]string{"siteUrl": siteURL}, &raw); err != nil {
		return nil, err
	}

	settings := make([]countryRegionSetting, len(raw))
	for i, item := range raw {
		settings[i] = countryRegionSetting{
			TwoLetterIsoCountryCode: item.TwoLetterIsoCountryCode,
			SettingsType:            decodeCountryRegionSettingsType(item.Type),
			URL:                     item.URL,
			Date:                    timePointer(item.Date),
		}
	}

	return &countryRegionSettingsResult{
		SiteURL:   siteURL,
		RowCount:  len(settings),
		Settings:  settings,
		QueriedAt: time.Now().UTC(),
	}, nil
}

// AddCountryRegionSettings adds a geo-targeting setting for a site.
func (c *Client) AddCountryRegionSettings(ctx context.Context, siteURL string, twoLetterIsoCountryCode string, settingsType string, settingsURL string) (*addCountryRegionSettingsResult, error) {
	settingsTypeValue, err := encodeCountryRegionSettingsType(settingsType)
	if err != nil {
		return nil, err
	}

	payload := rawCountryRegionSettingsCommandRequest{
		SiteURL: siteURL,
		Settings: rawCountryRegionSettingsCommand{
			Date:                    formatDotNetDate(time.Now().UTC()),
			TwoLetterIsoCountryCode: twoLetterIsoCountryCode,
			Type:                    settingsTypeValue,
			URL:                     settingsURL,
		},
	}
	if err := c.postCommand(ctx, "AddCountryRegionSettings", payload); err != nil {
		return nil, err
	}

	return &addCountryRegionSettingsResult{
		SiteURL:                 siteURL,
		TwoLetterIsoCountryCode: twoLetterIsoCountryCode,
		SettingsType:            settingsType,
		URL:                     settingsURL,
		Success:                 true,
		RequestedAt:             time.Now().UTC(),
	}, nil
}

// RemoveCountryRegionSettings removes a geo-targeting setting for a site.
func (c *Client) RemoveCountryRegionSettings(ctx context.Context, siteURL string, twoLetterIsoCountryCode string, settingsType string, settingsURL string) (*removeCountryRegionSettingsResult, error) {
	settingsTypeValue, err := encodeCountryRegionSettingsType(settingsType)
	if err != nil {
		return nil, err
	}

	payload := rawCountryRegionSettingsCommandRequest{
		SiteURL: siteURL,
		Settings: rawCountryRegionSettingsCommand{
			Date:                    formatDotNetDate(time.Now().UTC()),
			TwoLetterIsoCountryCode: twoLetterIsoCountryCode,
			Type:                    settingsTypeValue,
			URL:                     settingsURL,
		},
	}
	if err := c.postCommand(ctx, "RemoveCountryRegionSettings", payload); err != nil {
		return nil, err
	}

	return &removeCountryRegionSettingsResult{
		SiteURL:                 siteURL,
		TwoLetterIsoCountryCode: twoLetterIsoCountryCode,
		SettingsType:            settingsType,
		URL:                     settingsURL,
		Success:                 true,
		RequestedAt:             time.Now().UTC(),
	}, nil
}

// GetConnectedPages returns connected pages for a site.
func (c *Client) GetConnectedPages(ctx context.Context, siteURL string) (*connectedPagesResult, error) {
	var raw []rawConnectedPage
	if err := c.get(ctx, "GetConnectedPages", map[string]string{"siteUrl": siteURL}, &raw); err != nil {
		return nil, err
	}

	pages := make([]connectedPage, len(raw))
	for i, item := range raw {
		pages[i] = connectedPage{
			URL:                      item.URL,
			IsVerified:               item.IsVerified,
			RequestedMasterSite:      item.RequestedMasterSite,
			ActualMasterSite:         item.ActualMasterSite,
			HTTPStatusCode:           item.HTTPStatusCode,
			Market:                   item.Market,
			IsBlocked:                item.IsBlocked,
			LastSuccessfullyVerified: timePointerOrNilMinDate(item.LastSuccessfullyVerified),
		}
	}

	return &connectedPagesResult{
		SiteURL:   siteURL,
		RowCount:  len(pages),
		Pages:     pages,
		QueriedAt: time.Now().UTC(),
	}, nil
}

// AddConnectedPage adds a connected page to a site.
func (c *Client) AddConnectedPage(ctx context.Context, siteURL string, masterURL string) (*addConnectedPageResult, error) {
	payload := rawAddConnectedPageRequest{
		SiteURL:   siteURL,
		MasterURL: masterURL,
	}
	if err := c.postCommand(ctx, "AddConnectedPage", payload); err != nil {
		return nil, err
	}

	return &addConnectedPageResult{
		SiteURL:     siteURL,
		MasterURL:   masterURL,
		Success:     true,
		RequestedAt: time.Now().UTC(),
	}, nil
}

// GetActivePagePreviewBlocks returns active page preview blocks for a site.
func (c *Client) GetActivePagePreviewBlocks(ctx context.Context, siteURL string) (*activePagePreviewBlocksResult, error) {
	var raw []rawPagePreview
	if err := c.get(ctx, "GetActivePagePreviewBlocks", map[string]string{"siteUrl": siteURL}, &raw); err != nil {
		return nil, err
	}

	blocks := make([]pagePreviewBlock, len(raw))
	for i, item := range raw {
		blocks[i] = pagePreviewBlock{
			URL:         item.URL,
			BlockReason: decodePagePreviewBlockReason(item.BlockReason),
			SubmitDate:  timePointer(item.SubmitDate),
		}
	}

	return &activePagePreviewBlocksResult{
		SiteURL:   siteURL,
		RowCount:  len(blocks),
		Blocks:    blocks,
		QueriedAt: time.Now().UTC(),
	}, nil
}

// AddPagePreviewBlock adds a page preview block for a site URL.
func (c *Client) AddPagePreviewBlock(ctx context.Context, siteURL string, pageURL string, reason string) (*addPagePreviewBlockResult, error) {
	reasonValue, err := encodePagePreviewBlockReason(reason)
	if err != nil {
		return nil, err
	}

	payload := rawAddPagePreviewBlockRequest{
		SiteURL: siteURL,
		URL:     pageURL,
		Reason:  reasonValue,
	}
	if err := c.postCommand(ctx, "AddPagePreviewBlock", payload); err != nil {
		return nil, err
	}

	return &addPagePreviewBlockResult{
		SiteURL:     siteURL,
		URL:         pageURL,
		Reason:      reason,
		Success:     true,
		RequestedAt: time.Now().UTC(),
	}, nil
}

// RemovePagePreviewBlock removes a page preview block for a site URL.
func (c *Client) RemovePagePreviewBlock(ctx context.Context, siteURL string, pageURL string) (*removePagePreviewBlockResult, error) {
	if err := c.postCommand(ctx, "RemovePagePreviewBlock", map[string]string{"siteUrl": siteURL, "url": pageURL}); err != nil {
		return nil, err
	}

	return &removePagePreviewBlockResult{
		SiteURL:     siteURL,
		URL:         pageURL,
		Success:     true,
		RequestedAt: time.Now().UTC(),
	}, nil
}

// GetQueryPageDetailStats returns daily stats for a specific query/page pair.
func (c *Client) GetQueryPageDetailStats(ctx context.Context, siteURL string, query string, page string) (*queryPageDetailStatsResult, error) {
	var raw []rawDetailedQueryStat
	if err := c.get(ctx, "GetQueryPageDetailStats", map[string]string{
		"siteUrl": siteURL,
		"query":   query,
		"page":    page,
	}, &raw); err != nil {
		return nil, err
	}

	rows := make([]detailedQueryStat, len(raw))
	for i, item := range raw {
		rows[i] = detailedQueryStat{
			Date:        timePointer(item.Date),
			Clicks:      item.Clicks,
			Impressions: item.Impressions,
			Position:    item.Position,
		}
	}

	return &queryPageDetailStatsResult{
		SiteURL:   siteURL,
		Query:     query,
		Page:      page,
		RowCount:  len(rows),
		Rows:      rows,
		QueriedAt: time.Now().UTC(),
	}, nil
}

// GetQueryTrafficStats returns daily clicks and impressions for a query.
func (c *Client) GetQueryTrafficStats(ctx context.Context, siteURL string, query string) (*queryTrafficStatsResult, error) {
	var raw []rawRankTrafficStat
	if err := c.get(ctx, "GetQueryTrafficStats", map[string]string{"siteUrl": siteURL, "query": query}, &raw); err != nil {
		return nil, err
	}

	rows := make([]rankTrafficStat, len(raw))
	for i, item := range raw {
		rows[i] = rankTrafficStat{Date: timePointer(item.Date), Clicks: item.Clicks, Impressions: item.Impressions}
	}

	return &queryTrafficStatsResult{SiteURL: siteURL, Query: query, RowCount: len(rows), Rows: rows, QueriedAt: time.Now().UTC()}, nil
}

// GetKeyword returns market-wide impressions for a keyword over a date range.
func (c *Client) GetKeyword(ctx context.Context, query string, country string, language string, startDate string, endDate string) (*keywordResult, error) {
	var raw rawKeywordLookup
	if err := c.get(ctx, "GetKeyword", map[string]string{
		"q":         query,
		"country":   country,
		"language":  language,
		"startDate": startDate,
		"endDate":   endDate,
	}, &raw); err != nil {
		return nil, err
	}

	found := raw.Query != ""
	return &keywordResult{
		Query:            query,
		Country:          country,
		Language:         language,
		StartDate:        startDate,
		EndDate:          endDate,
		Found:            found,
		Impressions:      raw.Impressions,
		BroadImpressions: raw.BroadImpressions,
		QueriedAt:        time.Now().UTC(),
	}, nil
}

// GetRelatedKeywords returns market-wide related keyword impressions over a date range.
func (c *Client) GetRelatedKeywords(ctx context.Context, query string, country string, language string, startDate string, endDate string) (*relatedKeywordsResult, error) {
	var raw []rawKeywordLookup
	if err := c.get(ctx, "GetRelatedKeywords", map[string]string{
		"q":         query,
		"country":   country,
		"language":  language,
		"startDate": startDate,
		"endDate":   endDate,
	}, &raw); err != nil {
		return nil, err
	}

	rows := make([]relatedKeyword, len(raw))
	for i, item := range raw {
		rows[i] = relatedKeyword{
			Query:            item.Query,
			Impressions:      item.Impressions,
			BroadImpressions: item.BroadImpressions,
		}
	}

	return &relatedKeywordsResult{
		Query:     query,
		Country:   country,
		Language:  language,
		StartDate: startDate,
		EndDate:   endDate,
		RowCount:  len(rows),
		Rows:      rows,
		QueriedAt: time.Now().UTC(),
	}, nil
}

// GetChildrenURLInfo returns child URL crawl information for a parent URL.
func (c *Client) GetChildrenURLInfo(ctx context.Context, siteURL string, requestedURL string, page int, crawlDateFilter string, discoveredDateFilter string, docFlagsFilter string, httpCodeFilter string) (*childrenURLInfoResult, error) {
	if crawlDateFilter == "" {
		crawlDateFilter = "Any"
	}
	if discoveredDateFilter == "" {
		discoveredDateFilter = "Any"
	}
	if docFlagsFilter == "" {
		docFlagsFilter = "Any"
	}
	if httpCodeFilter == "" {
		httpCodeFilter = "Any"
	}

	crawlDateFilterValue, err := encodeCrawlDateFilter(crawlDateFilter)
	if err != nil {
		return nil, err
	}
	discoveredDateFilterValue, err := encodeDiscoveredDateFilter(discoveredDateFilter)
	if err != nil {
		return nil, err
	}
	docFlagsFilterValue, err := encodeDocFlagsFilter(docFlagsFilter)
	if err != nil {
		return nil, err
	}
	httpCodeFilterValue, err := encodeHTTPCodeFilter(httpCodeFilter)
	if err != nil {
		return nil, err
	}

	payload := rawChildrenURLInfoRequest{
		SiteURL: siteURL,
		URL:     requestedURL,
		Page:    page,
		FilterProperties: rawFilterProperties{
			CrawlDateFilter:      crawlDateFilterValue,
			DiscoveredDateFilter: discoveredDateFilterValue,
			DocFlagsFilters:      docFlagsFilterValue,
			HTTPCodeFilters:      httpCodeFilterValue,
		},
	}

	var raw []rawURLInfo
	if err := c.post(ctx, "GetChildrenUrlInfo", payload, &raw); err != nil {
		return nil, err
	}

	rows := make([]childURLInfo, len(raw))
	for i, item := range raw {
		rows[i] = childURLInfo{
			URL:                item.URL,
			IsPage:             item.IsPage,
			HTTPStatus:         item.HTTPStatus,
			DocumentSize:       item.DocumentSize,
			AnchorCount:        item.AnchorCount,
			DiscoveryDate:      timePointer(item.DiscoveryDate),
			LastCrawledDate:    timePointer(item.LastCrawledDate),
			TotalChildURLCount: item.TotalChildURLCount,
		}
	}

	return &childrenURLInfoResult{
		SiteURL:              siteURL,
		URL:                  requestedURL,
		Page:                 page,
		CrawlDateFilter:      crawlDateFilter,
		DiscoveredDateFilter: discoveredDateFilter,
		DocFlagsFilter:       docFlagsFilter,
		HTTPCodeFilter:       httpCodeFilter,
		RowCount:             len(rows),
		Rows:                 rows,
		QueriedAt:            time.Now().UTC(),
	}, nil
}

// GetChildrenURLTrafficInfo returns child URL traffic information for a parent URL.
func (c *Client) GetChildrenURLTrafficInfo(ctx context.Context, siteURL string, requestedURL string, page int) (*childrenURLTrafficInfoResult, error) {
	var raw []rawURLTrafficInfo
	if err := c.get(ctx, "GetChildrenUrlTrafficInfo", map[string]string{
		"siteUrl": siteURL,
		"url":     requestedURL,
		"page":    strconv.Itoa(page),
	}, &raw); err != nil {
		return nil, err
	}

	rows := make([]childURLTrafficInfo, len(raw))
	for i, item := range raw {
		rows[i] = childURLTrafficInfo{
			URL:         item.URL,
			IsPage:      item.IsPage,
			Clicks:      item.Clicks,
			Impressions: item.Impressions,
		}
	}

	return &childrenURLTrafficInfoResult{
		SiteURL:   siteURL,
		URL:       requestedURL,
		Page:      page,
		RowCount:  len(rows),
		Rows:      rows,
		QueriedAt: time.Now().UTC(),
	}, nil
}

// FetchURL requests a fresh Bing fetch for a URL.
func (c *Client) FetchURL(ctx context.Context, siteURL string, requestedURL string) (*fetchURLResult, error) {
	if err := c.postCommand(ctx, "FetchUrl", map[string]string{"siteUrl": siteURL, "url": requestedURL}); err != nil {
		return nil, err
	}

	return &fetchURLResult{SiteURL: siteURL, URL: requestedURL, Success: true, RequestedAt: time.Now().UTC()}, nil
}

// ListFetchedURLs returns previously fetched URLs for a site.
func (c *Client) ListFetchedURLs(ctx context.Context, siteURL string) (*fetchedURLsResult, error) {
	var raw []rawFetchedURL
	if err := c.get(ctx, "GetFetchedUrls", map[string]string{"siteUrl": siteURL}, &raw); err != nil {
		return nil, err
	}

	rows := make([]fetchedURL, len(raw))
	for i, item := range raw {
		rows[i] = fetchedURL{
			URL:     item.URL,
			Date:    timePointer(item.Date),
			Fetched: item.Fetched,
			Expired: item.Expired,
		}
	}

	return &fetchedURLsResult{SiteURL: siteURL, RowCount: len(rows), Rows: rows, QueriedAt: time.Now().UTC()}, nil
}

// GetFetchedURLDetails returns the stored fetch response for a URL.
func (c *Client) GetFetchedURLDetails(ctx context.Context, siteURL string, requestedURL string) (*fetchedURLDetailsResult, error) {
	var raw rawFetchedURLDetails
	if err := c.get(ctx, "GetFetchedUrlDetails", map[string]string{"siteUrl": siteURL, "url": requestedURL}, &raw); err != nil {
		return nil, err
	}

	return &fetchedURLDetailsResult{
		SiteURL:   siteURL,
		URL:       raw.URL,
		Date:      timePointer(raw.Date),
		Status:    raw.Status,
		Headers:   raw.Headers,
		Document:  raw.Document,
		QueriedAt: time.Now().UTC(),
	}, nil
}

// RemoveSitemap removes a sitemap from Bing Webmaster Tools.
func (c *Client) RemoveSitemap(ctx context.Context, siteURL string, feedURL string) (*removeSitemapResult, error) {
	if err := c.postCommand(ctx, "RemoveFeed", map[string]string{"siteUrl": siteURL, "feedUrl": feedURL}); err != nil {
		return nil, err
	}

	return &removeSitemapResult{SiteURL: siteURL, FeedURL: feedURL, Success: true, RequestedAt: time.Now().UTC()}, nil
}

// GetSiteMoves returns site move settings for a site.
func (c *Client) GetSiteMoves(ctx context.Context, siteURL string) (*siteMovesResult, error) {
	var raw []rawSiteMoveSettings
	if err := c.get(ctx, "GetSiteMoves", map[string]string{"siteUrl": siteURL}, &raw); err != nil {
		return nil, err
	}

	rows := make([]siteMove, len(raw))
	for i, item := range raw {
		rows[i] = siteMove{
			Date:      timePointer(item.Date),
			MoveScope: decodeMoveScope(item.MoveScope),
			MoveType:  decodeMoveType(item.MoveType),
			SourceURL: item.SourceURL,
			TargetURL: item.TargetURL,
		}
	}

	return &siteMovesResult{SiteURL: siteURL, RowCount: len(rows), Rows: rows, QueriedAt: time.Now().UTC()}, nil
}

// SubmitSiteMove submits a site move request.
func (c *Client) SubmitSiteMove(ctx context.Context, siteURL string, sourceURL string, targetURL string, moveType string, moveScope string) (*submitSiteMoveResult, error) {
	if moveType == "" {
		moveType = "Local"
	}
	if moveScope == "" {
		moveScope = "Domain"
	}

	moveTypeValue, err := encodeMoveType(moveType)
	if err != nil {
		return nil, err
	}
	moveScopeValue, err := encodeMoveScope(moveScope)
	if err != nil {
		return nil, err
	}

	payload := rawSiteMoveCommandRequest{
		SiteURL: siteURL,
		Settings: rawSiteMoveCommand{
			Date:      formatDotNetDate(time.Now().UTC()),
			MoveScope: moveScopeValue,
			MoveType:  moveTypeValue,
			SourceURL: sourceURL,
			TargetURL: targetURL,
		},
	}
	if err := c.postCommand(ctx, "SubmitSiteMove", payload); err != nil {
		return nil, err
	}

	return &submitSiteMoveResult{
		SiteURL:     siteURL,
		SourceURL:   sourceURL,
		TargetURL:   targetURL,
		MoveType:    moveType,
		MoveScope:   moveScope,
		Success:     true,
		RequestedAt: time.Now().UTC(),
	}, nil
}

// SubmitContent submits cached content and structured data for a URL.
func (c *Client) SubmitContent(ctx context.Context, siteURL string, requestedURL string, httpMessage string, structuredData string, dynamicServing string) (*submitContentResult, error) {
	if dynamicServing == "" {
		dynamicServing = "None"
	}

	dynamicServingValue, err := encodeDynamicServing(dynamicServing)
	if err != nil {
		return nil, err
	}

	payload := rawSubmitContentRequest{
		SiteURL:        siteURL,
		URL:            requestedURL,
		HTTPMessage:    httpMessage,
		StructuredData: structuredData,
		DynamicServing: dynamicServingValue,
	}
	if err := c.postCommand(ctx, "SubmitContent", payload); err != nil {
		return nil, err
	}

	return &submitContentResult{
		SiteURL:        siteURL,
		URL:            requestedURL,
		DynamicServing: dynamicServing,
		Success:        true,
		RequestedAt:    time.Now().UTC(),
	}, nil
}

// GetContentSubmissionQuota returns content submission quotas for a site.
func (c *Client) GetContentSubmissionQuota(ctx context.Context, siteURL string) (*contentSubmissionQuotaResult, error) {
	var raw rawQuota
	if err := c.get(ctx, "GetContentSubmissionQuota", map[string]string{"siteUrl": siteURL}, &raw); err != nil {
		return nil, err
	}

	return &contentSubmissionQuotaResult{
		SiteURL:      siteURL,
		DailyQuota:   raw.DailyQuota,
		MonthlyQuota: raw.MonthlyQuota,
		QueriedAt:    time.Now().UTC(),
	}, nil
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
