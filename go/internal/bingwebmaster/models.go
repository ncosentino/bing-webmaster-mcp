// Package bingwebmaster provides models for the Bing Webmaster Tools API.
package bingwebmaster

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
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
