// Package bingwebmaster provides models for the Bing Webmaster Tools API.
package bingwebmaster

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var dotNetDatePattern = regexp.MustCompile(`^/Date\((-?\d+)(?:[+-]\d{4})?\)/$`)

type wireTime struct {
	time.Time
}

func (t *wireTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("unmarshalling .NET date string: %w", err)
	}
	if raw == "" {
		return nil
	}

	parsed, err := parseDotNetDate(raw)
	if err != nil {
		return err
	}
	t.Time = parsed
	return nil
}

func parseDotNetDate(value string) (time.Time, error) {
	matches := dotNetDatePattern.FindStringSubmatch(value)
	if matches == nil {
		return time.Time{}, fmt.Errorf("invalid .NET date %q", value)
	}

	milliseconds, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("parsing .NET date milliseconds: %w", err)
	}

	return time.UnixMilli(milliseconds).UTC(), nil
}

func timePointer(value wireTime) *time.Time {
	if value.Time.IsZero() {
		return nil
	}

	t := value.Time.UTC()
	return &t
}

func formatDotNetDate(value time.Time) string {
	return fmt.Sprintf("/Date(%d+0000)/", value.UTC().UnixMilli())
}

type siteList struct {
	Sites     []site    `json:"sites"`
	QueriedAt time.Time `json:"queriedAt"`
}

type site struct {
	SiteURL             string `json:"siteUrl"`
	IsVerified          bool   `json:"isVerified"`
	DNSVerificationCode string `json:"dnsVerificationCode,omitempty"`
	AuthenticationCode  string `json:"authenticationCode,omitempty"`
}

type addSiteResult struct {
	SiteURL     string    `json:"siteUrl"`
	Success     bool      `json:"success"`
	RequestedAt time.Time `json:"requestedAt"`
}

type verifySiteResult struct {
	SiteURL     string    `json:"siteUrl"`
	Verified    bool      `json:"verified"`
	RequestedAt time.Time `json:"requestedAt"`
}

type sitemapList struct {
	SiteURL   string    `json:"siteUrl"`
	Sitemaps  []feed    `json:"sitemaps"`
	QueriedAt time.Time `json:"queriedAt"`
}

type feed struct {
	URL         string     `json:"url"`
	Type        string     `json:"type,omitempty"`
	Compressed  bool       `json:"compressed"`
	FileSize    int        `json:"fileSize"`
	LastCrawled *time.Time `json:"lastCrawled,omitempty"`
	Submitted   *time.Time `json:"submitted,omitempty"`
	Status      string     `json:"status,omitempty"`
	URLCount    int        `json:"urlCount"`
}

type sitemapDetailResult struct {
	SiteURL   string    `json:"siteUrl"`
	FeedURL   string    `json:"feedUrl"`
	Sitemap   *feed     `json:"sitemap,omitempty"`
	QueriedAt time.Time `json:"queriedAt"`
}

type submitSitemapResult struct {
	SiteURL     string    `json:"siteUrl"`
	FeedURL     string    `json:"feedUrl"`
	Success     bool      `json:"success"`
	SubmittedAt time.Time `json:"submittedAt"`
}

type submitURLResult struct {
	SiteURL     string    `json:"siteUrl"`
	URL         string    `json:"url"`
	Success     bool      `json:"success"`
	SubmittedAt time.Time `json:"submittedAt"`
}

type submitURLBatchResult struct {
	SiteURL        string    `json:"siteUrl"`
	URLList        []string  `json:"urlList"`
	SubmittedCount int       `json:"submittedCount"`
	Success        bool      `json:"success"`
	SubmittedAt    time.Time `json:"submittedAt"`
}

type urlSubmissionQuotaResult struct {
	SiteURL      string    `json:"siteUrl"`
	DailyQuota   int       `json:"dailyQuota"`
	MonthlyQuota int       `json:"monthlyQuota"`
	QueriedAt    time.Time `json:"queriedAt"`
}

type crawlIssuesResult struct {
	SiteURL   string       `json:"siteUrl"`
	Issues    []crawlIssue `json:"issues"`
	QueriedAt time.Time    `json:"queriedAt"`
}

type crawlIssue struct {
	URL      string   `json:"url"`
	HTTPCode int      `json:"httpCode"`
	Issues   []string `json:"issues"`
	InLinks  int      `json:"inLinks"`
}

type crawlStatsResult struct {
	SiteURL   string      `json:"siteUrl"`
	RowCount  int         `json:"rowCount"`
	Rows      []crawlStat `json:"rows"`
	QueriedAt time.Time   `json:"queriedAt"`
}

type crawlStat struct {
	Date               *time.Time `json:"date,omitempty"`
	CrawledPages       int        `json:"crawledPages"`
	CrawlErrors        int        `json:"crawlErrors"`
	InIndex            int        `json:"inIndex"`
	InLinks            int        `json:"inLinks"`
	Code2xx            int        `json:"code2xx"`
	Code301            int        `json:"code301"`
	Code302            int        `json:"code302"`
	Code4xx            int        `json:"code4xx"`
	Code5xx            int        `json:"code5xx"`
	AllOtherCodes      int        `json:"allOtherCodes"`
	BlockedByRobotsTxt int        `json:"blockedByRobotsTxt"`
	ContainsMalware    int        `json:"containsMalware"`
}

type urlInfoResult struct {
	SiteURL            string     `json:"siteUrl"`
	URL                string     `json:"url"`
	IsPage             bool       `json:"isPage"`
	HTTPStatus         int        `json:"httpStatus"`
	DocumentSize       int        `json:"documentSize"`
	AnchorCount        int        `json:"anchorCount"`
	DiscoveryDate      *time.Time `json:"discoveryDate,omitempty"`
	LastCrawledDate    *time.Time `json:"lastCrawledDate,omitempty"`
	TotalChildURLCount int        `json:"totalChildUrlCount"`
	QueriedAt          time.Time  `json:"queriedAt"`
}

type urlTrafficInfoResult struct {
	SiteURL     string    `json:"siteUrl"`
	URL         string    `json:"url"`
	IsPage      bool      `json:"isPage"`
	Clicks      int       `json:"clicks"`
	Impressions int       `json:"impressions"`
	QueriedAt   time.Time `json:"queriedAt"`
}

type urlLinksResult struct {
	SiteURL    string          `json:"siteUrl"`
	Link       string          `json:"link"`
	Page       int             `json:"page"`
	Details    []urlLinkDetail `json:"details"`
	TotalPages int             `json:"totalPages"`
	QueriedAt  time.Time       `json:"queriedAt"`
}

type urlLinkDetail struct {
	AnchorText string `json:"anchorText,omitempty"`
	URL        string `json:"url"`
}

type linkCountsResult struct {
	SiteURL    string      `json:"siteUrl"`
	Page       int         `json:"page"`
	Links      []linkCount `json:"links"`
	TotalPages int         `json:"totalPages"`
	QueriedAt  time.Time   `json:"queriedAt"`
}

type linkCount struct {
	Count int    `json:"count"`
	URL   string `json:"url"`
}

type rankAndTrafficStatsResult struct {
	SiteURL   string            `json:"siteUrl"`
	RowCount  int               `json:"rowCount"`
	Rows      []rankTrafficStat `json:"rows"`
	QueriedAt time.Time         `json:"queriedAt"`
}

type rankTrafficStat struct {
	Date        *time.Time `json:"date,omitempty"`
	Clicks      int        `json:"clicks"`
	Impressions int        `json:"impressions"`
}

type queryStatsResult struct {
	SiteURL   string      `json:"siteUrl"`
	RowCount  int         `json:"rowCount"`
	Rows      []queryStat `json:"rows"`
	QueriedAt time.Time   `json:"queriedAt"`
}

type queryStat struct {
	Query                 string     `json:"query"`
	Date                  *time.Time `json:"date,omitempty"`
	Clicks                int        `json:"clicks"`
	Impressions           int        `json:"impressions"`
	AvgClickPosition      int        `json:"avgClickPosition"`
	AvgImpressionPosition int        `json:"avgImpressionPosition"`
}

type pageStatsResult struct {
	SiteURL   string     `json:"siteUrl"`
	RowCount  int        `json:"rowCount"`
	Rows      []pageStat `json:"rows"`
	QueriedAt time.Time  `json:"queriedAt"`
}

type pageStat struct {
	Page                  string     `json:"page"`
	Date                  *time.Time `json:"date,omitempty"`
	Clicks                int        `json:"clicks"`
	Impressions           int        `json:"impressions"`
	AvgClickPosition      int        `json:"avgClickPosition"`
	AvgImpressionPosition int        `json:"avgImpressionPosition"`
}

type pageQueryStatsResult struct {
	SiteURL   string      `json:"siteUrl"`
	Page      string      `json:"page"`
	RowCount  int         `json:"rowCount"`
	Rows      []queryStat `json:"rows"`
	QueriedAt time.Time   `json:"queriedAt"`
}

type queryPageStatsResult struct {
	SiteURL   string     `json:"siteUrl"`
	Query     string     `json:"query"`
	RowCount  int        `json:"rowCount"`
	Rows      []pageStat `json:"rows"`
	QueriedAt time.Time  `json:"queriedAt"`
}

type keywordStatsResult struct {
	Query     string        `json:"query"`
	Country   string        `json:"country"`
	Language  string        `json:"language"`
	RowCount  int           `json:"rowCount"`
	Rows      []keywordStat `json:"rows"`
	QueriedAt time.Time     `json:"queriedAt"`
}

type keywordStat struct {
	Query            string     `json:"query"`
	Date             *time.Time `json:"date,omitempty"`
	Impressions      int        `json:"impressions"`
	BroadImpressions int        `json:"broadImpressions"`
}

type removeSiteResult struct {
	SiteURL     string    `json:"siteUrl"`
	Success     bool      `json:"success"`
	RequestedAt time.Time `json:"requestedAt"`
}

type siteRolesResult struct {
	SiteURL              string     `json:"siteUrl"`
	IncludeAllSubdomains bool       `json:"includeAllSubdomains"`
	RowCount             int        `json:"rowCount"`
	Rows                 []siteRole `json:"rows"`
	QueriedAt            time.Time  `json:"queriedAt"`
}

type siteRole struct {
	Email                   string     `json:"email"`
	Role                    string     `json:"role"`
	Site                    string     `json:"site"`
	VerificationSite        string     `json:"verificationSite"`
	Expired                 bool       `json:"expired"`
	DelegatorEmail          string     `json:"delegatorEmail,omitempty"`
	DelegatedCode           string     `json:"delegatedCode,omitempty"`
	DelegatedCodeOwnerEmail string     `json:"delegatedCodeOwnerEmail,omitempty"`
	Date                    *time.Time `json:"date,omitempty"`
}

type addSiteRoleResult struct {
	SiteURL         string    `json:"siteUrl"`
	DelegatedURL    string    `json:"delegatedUrl"`
	UserEmail       string    `json:"userEmail"`
	IsAdministrator bool      `json:"isAdministrator"`
	IsReadOnly      bool      `json:"isReadOnly"`
	Success         bool      `json:"success"`
	RequestedAt     time.Time `json:"requestedAt"`
}

type removeSiteRoleResult struct {
	SiteURL     string    `json:"siteUrl"`
	Email       string    `json:"email"`
	Role        string    `json:"role"`
	Success     bool      `json:"success"`
	RequestedAt time.Time `json:"requestedAt"`
}

type blockedURLsResult struct {
	SiteURL   string       `json:"siteUrl"`
	RowCount  int          `json:"rowCount"`
	Rows      []blockedURL `json:"rows"`
	QueriedAt time.Time    `json:"queriedAt"`
}

type blockedURL struct {
	URL         string     `json:"url"`
	EntityType  string     `json:"entityType"`
	RequestType string     `json:"requestType"`
	Date        *time.Time `json:"date,omitempty"`
}

type addBlockedURLResult struct {
	SiteURL     string    `json:"siteUrl"`
	URL         string    `json:"url"`
	EntityType  string    `json:"entityType"`
	RequestType string    `json:"requestType"`
	Success     bool      `json:"success"`
	RequestedAt time.Time `json:"requestedAt"`
}

type removeBlockedURLResult struct {
	SiteURL     string    `json:"siteUrl"`
	URL         string    `json:"url"`
	EntityType  string    `json:"entityType"`
	RequestType string    `json:"requestType"`
	Success     bool      `json:"success"`
	RequestedAt time.Time `json:"requestedAt"`
}

type queryPageDetailStatsResult struct {
	SiteURL   string              `json:"siteUrl"`
	Query     string              `json:"query"`
	Page      string              `json:"page"`
	RowCount  int                 `json:"rowCount"`
	Rows      []detailedQueryStat `json:"rows"`
	QueriedAt time.Time           `json:"queriedAt"`
}

type detailedQueryStat struct {
	Date        *time.Time `json:"date,omitempty"`
	Clicks      int        `json:"clicks"`
	Impressions int        `json:"impressions"`
	Position    int        `json:"position"`
}

type queryTrafficStatsResult struct {
	SiteURL   string            `json:"siteUrl"`
	Query     string            `json:"query"`
	RowCount  int               `json:"rowCount"`
	Rows      []rankTrafficStat `json:"rows"`
	QueriedAt time.Time         `json:"queriedAt"`
}

type keywordResult struct {
	Query            string    `json:"query"`
	Country          string    `json:"country"`
	Language         string    `json:"language"`
	StartDate        string    `json:"startDate"`
	EndDate          string    `json:"endDate"`
	Found            bool      `json:"found"`
	Impressions      int       `json:"impressions"`
	BroadImpressions int       `json:"broadImpressions"`
	QueriedAt        time.Time `json:"queriedAt"`
}

type relatedKeywordsResult struct {
	Query     string           `json:"query"`
	Country   string           `json:"country"`
	Language  string           `json:"language"`
	StartDate string           `json:"startDate"`
	EndDate   string           `json:"endDate"`
	RowCount  int              `json:"rowCount"`
	Rows      []relatedKeyword `json:"rows"`
	QueriedAt time.Time        `json:"queriedAt"`
}

type relatedKeyword struct {
	Query            string `json:"query"`
	Impressions      int    `json:"impressions"`
	BroadImpressions int    `json:"broadImpressions"`
}

type childrenURLInfoResult struct {
	SiteURL              string         `json:"siteUrl"`
	URL                  string         `json:"url"`
	Page                 int            `json:"page"`
	CrawlDateFilter      string         `json:"crawlDateFilter"`
	DiscoveredDateFilter string         `json:"discoveredDateFilter"`
	DocFlagsFilter       string         `json:"docFlagsFilter"`
	HTTPCodeFilter       string         `json:"httpCodeFilter"`
	RowCount             int            `json:"rowCount"`
	Rows                 []childURLInfo `json:"rows"`
	QueriedAt            time.Time      `json:"queriedAt"`
}

type childURLInfo struct {
	URL                string     `json:"url"`
	IsPage             bool       `json:"isPage"`
	HTTPStatus         int        `json:"httpStatus"`
	DocumentSize       int        `json:"documentSize"`
	AnchorCount        int        `json:"anchorCount"`
	DiscoveryDate      *time.Time `json:"discoveryDate,omitempty"`
	LastCrawledDate    *time.Time `json:"lastCrawledDate,omitempty"`
	TotalChildURLCount int        `json:"totalChildUrlCount"`
}

type childrenURLTrafficInfoResult struct {
	SiteURL   string                `json:"siteUrl"`
	URL       string                `json:"url"`
	Page      int                   `json:"page"`
	RowCount  int                   `json:"rowCount"`
	Rows      []childURLTrafficInfo `json:"rows"`
	QueriedAt time.Time             `json:"queriedAt"`
}

type childURLTrafficInfo struct {
	URL         string `json:"url"`
	IsPage      bool   `json:"isPage"`
	Clicks      int    `json:"clicks"`
	Impressions int    `json:"impressions"`
}

type fetchURLResult struct {
	SiteURL     string    `json:"siteUrl"`
	URL         string    `json:"url"`
	Success     bool      `json:"success"`
	RequestedAt time.Time `json:"requestedAt"`
}

type fetchedURLsResult struct {
	SiteURL   string       `json:"siteUrl"`
	RowCount  int          `json:"rowCount"`
	Rows      []fetchedURL `json:"rows"`
	QueriedAt time.Time    `json:"queriedAt"`
}

type fetchedURL struct {
	URL     string     `json:"url"`
	Date    *time.Time `json:"date,omitempty"`
	Fetched bool       `json:"fetched"`
	Expired bool       `json:"expired"`
}

type fetchedURLDetailsResult struct {
	SiteURL   string     `json:"siteUrl"`
	URL       string     `json:"url"`
	Date      *time.Time `json:"date,omitempty"`
	Status    string     `json:"status,omitempty"`
	Headers   string     `json:"headers,omitempty"`
	Document  string     `json:"document,omitempty"`
	QueriedAt time.Time  `json:"queriedAt"`
}

type removeSitemapResult struct {
	SiteURL     string    `json:"siteUrl"`
	FeedURL     string    `json:"feedUrl"`
	Success     bool      `json:"success"`
	RequestedAt time.Time `json:"requestedAt"`
}

type siteMovesResult struct {
	SiteURL   string     `json:"siteUrl"`
	RowCount  int        `json:"rowCount"`
	Rows      []siteMove `json:"rows"`
	QueriedAt time.Time  `json:"queriedAt"`
}

type siteMove struct {
	Date      *time.Time `json:"date,omitempty"`
	MoveScope string     `json:"moveScope"`
	MoveType  string     `json:"moveType"`
	SourceURL string     `json:"sourceUrl"`
	TargetURL string     `json:"targetUrl"`
}

type submitSiteMoveResult struct {
	SiteURL     string    `json:"siteUrl"`
	SourceURL   string    `json:"sourceUrl"`
	TargetURL   string    `json:"targetUrl"`
	MoveType    string    `json:"moveType"`
	MoveScope   string    `json:"moveScope"`
	Success     bool      `json:"success"`
	RequestedAt time.Time `json:"requestedAt"`
}

type submitContentResult struct {
	SiteURL        string    `json:"siteUrl"`
	URL            string    `json:"url"`
	DynamicServing string    `json:"dynamicServing"`
	Success        bool      `json:"success"`
	RequestedAt    time.Time `json:"requestedAt"`
}

type contentSubmissionQuotaResult struct {
	SiteURL      string    `json:"siteUrl"`
	DailyQuota   int       `json:"dailyQuota"`
	MonthlyQuota int       `json:"monthlyQuota"`
	QueriedAt    time.Time `json:"queriedAt"`
}

type apiEnvelope[T any] struct {
	D T `json:"d"`
}

type rawSite struct {
	URL                 string `json:"Url"`
	IsVerified          bool   `json:"IsVerified"`
	DNSVerificationCode string `json:"DnsVerificationCode"`
	AuthenticationCode  string `json:"AuthenticationCode"`
}

type rawFeed struct {
	URL         string   `json:"Url"`
	Type        string   `json:"Type"`
	Compressed  bool     `json:"Compressed"`
	FileSize    int      `json:"FileSize"`
	LastCrawled wireTime `json:"LastCrawled"`
	Submitted   wireTime `json:"Submitted"`
	Status      string   `json:"Status"`
	URLCount    int      `json:"UrlCount"`
}

type rawQuota struct {
	DailyQuota   int `json:"DailyQuota"`
	MonthlyQuota int `json:"MonthlyQuota"`
}

type rawCrawlIssue struct {
	URL      string `json:"Url"`
	HTTPCode int    `json:"HttpCode"`
	Issues   int    `json:"Issues"`
	InLinks  int    `json:"InLinks"`
}

type rawCrawlStat struct {
	Date               wireTime `json:"Date"`
	CrawledPages       int      `json:"CrawledPages"`
	CrawlErrors        int      `json:"CrawlErrors"`
	InIndex            int      `json:"InIndex"`
	InLinks            int      `json:"InLinks"`
	Code2xx            int      `json:"Code2xx"`
	Code301            int      `json:"Code301"`
	Code302            int      `json:"Code302"`
	Code4xx            int      `json:"Code4xx"`
	Code5xx            int      `json:"Code5xx"`
	AllOtherCodes      int      `json:"AllOtherCodes"`
	BlockedByRobotsTxt int      `json:"BlockedByRobotsTxt"`
	ContainsMalware    int      `json:"ContainsMalware"`
}

type rawURLInfo struct {
	URL                string   `json:"Url"`
	IsPage             bool     `json:"IsPage"`
	HTTPStatus         int      `json:"HttpStatus"`
	DocumentSize       int      `json:"DocumentSize"`
	AnchorCount        int      `json:"AnchorCount"`
	DiscoveryDate      wireTime `json:"DiscoveryDate"`
	LastCrawledDate    wireTime `json:"LastCrawledDate"`
	TotalChildURLCount int      `json:"TotalChildUrlCount"`
}

type rawURLTrafficInfo struct {
	URL         string `json:"Url"`
	IsPage      bool   `json:"IsPage"`
	Clicks      int    `json:"Clicks"`
	Impressions int    `json:"Impressions"`
}

type rawURLLinkDetails struct {
	Details    []rawURLLinkDetail `json:"Details"`
	TotalPages int                `json:"TotalPages"`
}

type rawURLLinkDetail struct {
	AnchorText string `json:"AnchorText"`
	URL        string `json:"Url"`
}

type rawLinkCounts struct {
	Links      []rawLinkCount `json:"Links"`
	TotalPages int            `json:"TotalPages"`
}

type rawLinkCount struct {
	Count int    `json:"Count"`
	URL   string `json:"Url"`
}

type rawRankTrafficStat struct {
	Date        wireTime `json:"Date"`
	Clicks      int      `json:"Clicks"`
	Impressions int      `json:"Impressions"`
}

type rawQueryStat struct {
	Query                 string   `json:"Query"`
	Date                  wireTime `json:"Date"`
	Clicks                int      `json:"Clicks"`
	Impressions           int      `json:"Impressions"`
	AvgClickPosition      int      `json:"AvgClickPosition"`
	AvgImpressionPosition int      `json:"AvgImpressionPosition"`
}

type rawKeywordStat struct {
	Query            string   `json:"Query"`
	Date             wireTime `json:"Date"`
	Impressions      int      `json:"Impressions"`
	BroadImpressions int      `json:"BroadImpressions"`
}

type rawSiteRole struct {
	Date                    wireTime `json:"Date"`
	DelegatedCode           string   `json:"DelegatedCode"`
	DelegatorEmail          string   `json:"DelegatorEmail"`
	DelegatedCodeOwnerEmail string   `json:"DelegatedCodeOwnerEmail"`
	Email                   string   `json:"Email"`
	Expired                 bool     `json:"Expired"`
	Role                    int      `json:"Role"`
	Site                    string   `json:"Site"`
	VerificationSite        string   `json:"VerificationSite"`
}

type rawBlockedURL struct {
	Date        wireTime `json:"Date"`
	EntityType  int      `json:"EntityType"`
	RequestType int      `json:"RequestType"`
	URL         string   `json:"Url"`
}

type rawDetailedQueryStat struct {
	Date        wireTime `json:"Date"`
	Clicks      int      `json:"Clicks"`
	Impressions int      `json:"Impressions"`
	Position    int      `json:"Position"`
}

type rawKeywordLookup struct {
	Query            string `json:"Query"`
	BroadImpressions int    `json:"BroadImpressions"`
	Impressions      int    `json:"Impressions"`
}

type rawFetchedURL struct {
	Date    wireTime `json:"Date"`
	Expired bool     `json:"Expired"`
	Fetched bool     `json:"Fetched"`
	URL     string   `json:"Url"`
}

type rawFetchedURLDetails struct {
	Date     wireTime `json:"Date"`
	Document string   `json:"Document"`
	Headers  string   `json:"Headers"`
	Status   string   `json:"Status"`
	URL      string   `json:"Url"`
}

type rawSiteMoveSettings struct {
	Date      wireTime `json:"Date"`
	MoveScope int      `json:"MoveScope"`
	MoveType  int      `json:"MoveType"`
	SourceURL string   `json:"SourceUrl"`
	TargetURL string   `json:"TargetUrl"`
}

type rawAddSiteRoleRequest struct {
	SiteURL            string `json:"siteUrl"`
	DelegatedURL       string `json:"delegatedUrl"`
	UserEmail          string `json:"userEmail"`
	AuthenticationCode string `json:"authenticationCode"`
	IsAdministrator    bool   `json:"isAdministrator"`
	IsReadOnly         bool   `json:"isReadOnly"`
}

type rawSiteRoleCommandRequest struct {
	SiteURL  string             `json:"siteUrl"`
	SiteRole rawSiteRoleCommand `json:"siteRole"`
}

type rawSiteRoleCommand struct {
	Date                    string `json:"Date"`
	Email                   string `json:"Email"`
	Role                    int    `json:"Role"`
	Site                    string `json:"Site"`
	VerificationSite        string `json:"VerificationSite"`
	DelegatedCode           string `json:"DelegatedCode,omitempty"`
	DelegatorEmail          string `json:"DelegatorEmail,omitempty"`
	DelegatedCodeOwnerEmail string `json:"DelegatedCodeOwnerEmail,omitempty"`
}

type rawBlockedURLCommandRequest struct {
	SiteURL    string               `json:"siteUrl"`
	BlockedURL rawBlockedURLCommand `json:"blockedUrl"`
}

type rawBlockedURLCommand struct {
	Date        string `json:"Date"`
	EntityType  int    `json:"EntityType"`
	RequestType int    `json:"RequestType"`
	URL         string `json:"Url"`
}

type rawChildrenURLInfoRequest struct {
	SiteURL          string              `json:"siteUrl"`
	URL              string              `json:"url"`
	Page             int                 `json:"page"`
	FilterProperties rawFilterProperties `json:"filterProperties"`
}

type rawFilterProperties struct {
	CrawlDateFilter      int `json:"CrawlDateFilter"`
	DiscoveredDateFilter int `json:"DiscoveredDateFilter"`
	DocFlagsFilters      int `json:"DocFlagsFilters"`
	HTTPCodeFilters      int `json:"HttpCodeFilters"`
}

type rawSiteMoveCommandRequest struct {
	SiteURL  string             `json:"siteUrl"`
	Settings rawSiteMoveCommand `json:"settings"`
}

type rawSiteMoveCommand struct {
	Date      string `json:"Date"`
	MoveScope int    `json:"MoveScope"`
	MoveType  int    `json:"MoveType"`
	SourceURL string `json:"SourceUrl"`
	TargetURL string `json:"TargetUrl"`
}

type rawSubmitContentRequest struct {
	SiteURL        string `json:"siteUrl"`
	URL            string `json:"url"`
	HTTPMessage    string `json:"httpMessage"`
	StructuredData string `json:"structuredData"`
	DynamicServing int    `json:"dynamicServing"`
}

type namedEnum struct {
	value int
	name  string
}

var siteRoleValues = []namedEnum{
	{value: 0, name: "Administrator"},
	{value: 1, name: "ReadOnly"},
	{value: 2, name: "ReadWrite"},
}

var blockedURLEntityTypeValues = []namedEnum{
	{value: 0, name: "Page"},
	{value: 1, name: "Directory"},
}

var blockedURLRequestTypeValues = []namedEnum{
	{value: 0, name: "CacheOnly"},
	{value: 1, name: "FullRemoval"},
}

var crawlDateFilterValues = []namedEnum{
	{value: 0, name: "Any"},
	{value: 1, name: "LastWeek"},
	{value: 2, name: "LastTwoWeeks"},
	{value: 4, name: "LastThreeWeeks"},
}

var discoveredDateFilterValues = []namedEnum{
	{value: 0, name: "Any"},
	{value: 1, name: "LastWeek"},
	{value: 2, name: "LastMonth"},
}

var docFlagsFilterValues = []namedEnum{
	{value: 0, name: "Any"},
	{value: 1, name: "IsBlockedByRobotsTxt"},
	{value: 2, name: "IsMalware"},
}

var httpCodeFilterValues = []namedEnum{
	{value: 0, name: "Any"},
	{value: 1, name: "Code2xx"},
	{value: 2, name: "Code3xx"},
	{value: 4, name: "Code301"},
	{value: 8, name: "Code302"},
	{value: 16, name: "Code4xx"},
	{value: 32, name: "Code5xx"},
	{value: 64, name: "AllOthers"},
}

var moveScopeValues = []namedEnum{
	{value: 0, name: "Domain"},
	{value: 1, name: "Host"},
	{value: 2, name: "Directory"},
}

var moveTypeValues = []namedEnum{
	{value: 0, name: "Local"},
	{value: 1, name: "Global"},
}

var dynamicServingValues = []namedEnum{
	{value: 0, name: "None"},
	{value: 1, name: "PcLaptop"},
	{value: 2, name: "Mobile"},
	{value: 3, name: "Amp"},
	{value: 4, name: "Tablet"},
	{value: 5, name: "NonVisualBrowser"},
}

func decodeSiteRole(value int) string {
	return decodeEnumValue(value, siteRoleValues)
}

func encodeSiteRole(value string) (int, error) {
	return encodeEnumValue("role", value, siteRoleValues)
}

func decodeBlockedURLEntityType(value int) string {
	return decodeEnumValue(value, blockedURLEntityTypeValues)
}

func encodeBlockedURLEntityType(value string) (int, error) {
	return encodeEnumValue("entity_type", value, blockedURLEntityTypeValues)
}

func decodeBlockedURLRequestType(value int) string {
	return decodeEnumValue(value, blockedURLRequestTypeValues)
}

func encodeBlockedURLRequestType(value string) (int, error) {
	return encodeEnumValue("request_type", value, blockedURLRequestTypeValues)
}

func decodeMoveScope(value int) string {
	return decodeEnumValue(value, moveScopeValues)
}

func encodeMoveScope(value string) (int, error) {
	return encodeEnumValue("move_scope", value, moveScopeValues)
}

func decodeMoveType(value int) string {
	return decodeEnumValue(value, moveTypeValues)
}

func encodeMoveType(value string) (int, error) {
	return encodeEnumValue("move_type", value, moveTypeValues)
}

func encodeCrawlDateFilter(value string) (int, error) {
	return encodeEnumValue("crawl_date_filter", value, crawlDateFilterValues)
}

func encodeDiscoveredDateFilter(value string) (int, error) {
	return encodeEnumValue("discovered_date_filter", value, discoveredDateFilterValues)
}

func encodeDocFlagsFilter(value string) (int, error) {
	return encodeEnumValue("doc_flags_filter", value, docFlagsFilterValues)
}

func encodeHTTPCodeFilter(value string) (int, error) {
	return encodeEnumValue("http_code_filter", value, httpCodeFilterValues)
}

func encodeDynamicServing(value string) (int, error) {
	return encodeEnumValue("dynamic_serving", value, dynamicServingValues)
}

func decodeEnumValue(value int, enums []namedEnum) string {
	for _, item := range enums {
		if item.value == value {
			return item.name
		}
	}

	return strconv.Itoa(value)
}

func encodeEnumValue(label string, value string, enums []namedEnum) (int, error) {
	for _, item := range enums {
		if item.name == value {
			return item.value, nil
		}
	}

	names := make([]string, len(enums))
	for i, item := range enums {
		names[i] = item.name
	}

	return 0, fmt.Errorf("invalid %s %q; expected one of: %s", label, value, strings.Join(names, ", "))
}

var crawlIssueFlags = []struct {
	mask int
	name string
}{
	{mask: 1, name: "Code301"},
	{mask: 2, name: "Code302"},
	{mask: 4, name: "Code4xx"},
	{mask: 8, name: "Code5xx"},
	{mask: 16, name: "BlockedByRobotsTxt"},
	{mask: 32, name: "ContainsMalware"},
	{mask: 64, name: "ImportantUrlBlockedByRobotsTxt"},
	{mask: 128, name: "DnsErrors"},
	{mask: 256, name: "TimeOutErrors"},
}

func decodeCrawlIssueFlags(value int) []string {
	if value == 0 {
		return []string{}
	}

	decoded := make([]string, 0, len(crawlIssueFlags))
	for _, flag := range crawlIssueFlags {
		if value&flag.mask != 0 {
			decoded = append(decoded, flag.name)
		}
	}

	return decoded
}
