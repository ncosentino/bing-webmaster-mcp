using System.Text.Json.Serialization;

namespace BingWebmasterMcp.BingWebmaster;

internal sealed record SiteInfo(
    string SiteUrl,
    bool IsVerified,
    string? DnsVerificationCode,
    string? AuthenticationCode);

internal sealed record ListSitesResponse(
    IReadOnlyList<SiteInfo> Sites,
    DateTimeOffset QueriedAt);

internal sealed record AddSiteResponse(
    string SiteUrl,
    bool Success,
    DateTimeOffset RequestedAt);

internal sealed record VerifySiteResponse(
    string SiteUrl,
    bool Verified,
    DateTimeOffset RequestedAt);

internal sealed record SitemapInfo(
    string Url,
    string? Type,
    bool Compressed,
    int FileSize,
    DateTimeOffset? LastCrawled,
    DateTimeOffset? Submitted,
    string? Status,
    int UrlCount);

internal sealed record ListSitemapsResponse(
    string SiteUrl,
    IReadOnlyList<SitemapInfo> Sitemaps,
    DateTimeOffset QueriedAt);

internal sealed record GetSitemapDetailsResponse(
    string SiteUrl,
    string FeedUrl,
    SitemapInfo? Sitemap,
    DateTimeOffset QueriedAt);

internal sealed record SubmitSitemapResponse(
    string SiteUrl,
    string FeedUrl,
    bool Success,
    DateTimeOffset SubmittedAt);

internal sealed record SubmitUrlResponse(
    string SiteUrl,
    string Url,
    bool Success,
    DateTimeOffset SubmittedAt);

internal sealed record SubmitUrlBatchResponse(
    string SiteUrl,
    IReadOnlyList<string> UrlList,
    int SubmittedCount,
    bool Success,
    DateTimeOffset SubmittedAt);

internal sealed record SubmitUrlIndexNowResponse(
    string Host,
    IReadOnlyList<string> UrlList,
    string? KeyLocation,
    bool Success,
    string KeySource,
    DateTimeOffset SubmittedAt);

internal sealed record UrlSubmissionQuotaResponse(
    string SiteUrl,
    int DailyQuota,
    int MonthlyQuota,
    DateTimeOffset QueriedAt);

internal sealed record CrawlIssueEntry(
    string Url,
    int HttpCode,
    IReadOnlyList<string> Issues,
    int InLinks);

internal sealed record CrawlIssuesResponse(
    string SiteUrl,
    IReadOnlyList<CrawlIssueEntry> Issues,
    DateTimeOffset QueriedAt);

internal sealed record CrawlStatEntry(
    DateTimeOffset Date,
    int CrawledPages,
    int CrawlErrors,
    int InIndex,
    int InLinks,
    int Code2xx,
    int Code301,
    int Code302,
    int Code4xx,
    int Code5xx,
    int AllOtherCodes,
    int BlockedByRobotsTxt,
    int ContainsMalware);

internal sealed record CrawlStatsResponse(
    string SiteUrl,
    int RowCount,
    IReadOnlyList<CrawlStatEntry> Rows,
    DateTimeOffset QueriedAt);

internal sealed record UrlInfoResponse(
    string SiteUrl,
    string Url,
    bool IsPage,
    int HttpStatus,
    int DocumentSize,
    int AnchorCount,
    DateTimeOffset? DiscoveryDate,
    DateTimeOffset? LastCrawledDate,
    int TotalChildUrlCount,
    DateTimeOffset QueriedAt);

internal sealed record UrlTrafficInfoResponse(
    string SiteUrl,
    string Url,
    bool IsPage,
    int Clicks,
    int Impressions,
    DateTimeOffset QueriedAt);

internal sealed record LinkDetailEntry(
    string? AnchorText,
    string? Url);

internal sealed record UrlLinksResponse(
    string SiteUrl,
    string Link,
    int Page,
    int TotalPages,
    IReadOnlyList<LinkDetailEntry> Details,
    DateTimeOffset QueriedAt);

internal sealed record LinkCountEntry(
    int Count,
    string? Url);

internal sealed record LinkCountsResponse(
    string SiteUrl,
    int Page,
    int TotalPages,
    IReadOnlyList<LinkCountEntry> Links,
    DateTimeOffset QueriedAt);

internal sealed record RankAndTrafficStatEntry(
    DateTimeOffset Date,
    int Clicks,
    int Impressions);

internal sealed record RankAndTrafficStatsResponse(
    string SiteUrl,
    int RowCount,
    IReadOnlyList<RankAndTrafficStatEntry> Rows,
    DateTimeOffset QueriedAt);

internal sealed record QueryStatsEntry(
    string Query,
    DateTimeOffset Date,
    int Clicks,
    int Impressions,
    int AvgClickPosition,
    int AvgImpressionPosition);

internal sealed record PageStatsEntry(
    string Page,
    DateTimeOffset Date,
    int Clicks,
    int Impressions,
    int AvgClickPosition,
    int AvgImpressionPosition);

internal sealed record QueryStatsResponse(
    string SiteUrl,
    int RowCount,
    IReadOnlyList<QueryStatsEntry> Rows,
    DateTimeOffset QueriedAt);

internal sealed record PageStatsResponse(
    string SiteUrl,
    int RowCount,
    IReadOnlyList<PageStatsEntry> Rows,
    DateTimeOffset QueriedAt);

internal sealed record PageQueryStatsResponse(
    string SiteUrl,
    string Page,
    int RowCount,
    IReadOnlyList<QueryStatsEntry> Rows,
    DateTimeOffset QueriedAt);

internal sealed record QueryPageStatsResponse(
    string SiteUrl,
    string Query,
    int RowCount,
    IReadOnlyList<PageStatsEntry> Rows,
    DateTimeOffset QueriedAt);

internal sealed record KeywordStatsEntry(
    string Query,
    DateTimeOffset Date,
    int Impressions,
    int BroadImpressions);

internal sealed record KeywordStatsResponse(
    string Query,
    string Country,
    string Language,
    int RowCount,
    IReadOnlyList<KeywordStatsEntry> Rows,
    DateTimeOffset QueriedAt);

internal sealed record ErrorResult(string Error);

internal sealed class ApiSite
{
    [JsonPropertyName("Url")]
    public string Url { get; set; } = string.Empty;

    [JsonPropertyName("IsVerified")]
    public bool IsVerified { get; set; }

    [JsonPropertyName("DnsVerificationCode")]
    public string? DnsVerificationCode { get; set; }

    [JsonPropertyName("AuthenticationCode")]
    public string? AuthenticationCode { get; set; }
}

internal sealed class ApiFeed
{
    [JsonPropertyName("Url")]
    public string Url { get; set; } = string.Empty;

    [JsonPropertyName("Type")]
    public string? Type { get; set; }

    [JsonPropertyName("Compressed")]
    public bool Compressed { get; set; }

    [JsonPropertyName("FileSize")]
    public int FileSize { get; set; }

    [JsonPropertyName("LastCrawled")]
    public string? LastCrawled { get; set; }

    [JsonPropertyName("Submitted")]
    public string? Submitted { get; set; }

    [JsonPropertyName("Status")]
    public string? Status { get; set; }

    [JsonPropertyName("UrlCount")]
    public int UrlCount { get; set; }
}

internal sealed class ApiUrlSubmissionQuota
{
    [JsonPropertyName("DailyQuota")]
    public int DailyQuota { get; set; }

    [JsonPropertyName("MonthlyQuota")]
    public int MonthlyQuota { get; set; }
}

internal sealed class ApiUrlWithCrawlIssues
{
    [JsonPropertyName("Url")]
    public string Url { get; set; } = string.Empty;

    [JsonPropertyName("HttpCode")]
    public int HttpCode { get; set; }

    [JsonPropertyName("Issues")]
    public int Issues { get; set; }

    [JsonPropertyName("InLinks")]
    public int InLinks { get; set; }
}

internal sealed class ApiCrawlStats
{
    [JsonPropertyName("Date")]
    public string Date { get; set; } = string.Empty;

    [JsonPropertyName("CrawledPages")]
    public int CrawledPages { get; set; }

    [JsonPropertyName("CrawlErrors")]
    public int CrawlErrors { get; set; }

    [JsonPropertyName("InIndex")]
    public int InIndex { get; set; }

    [JsonPropertyName("InLinks")]
    public int InLinks { get; set; }

    [JsonPropertyName("Code2xx")]
    public int Code2xx { get; set; }

    [JsonPropertyName("Code301")]
    public int Code301 { get; set; }

    [JsonPropertyName("Code302")]
    public int Code302 { get; set; }

    [JsonPropertyName("Code4xx")]
    public int Code4xx { get; set; }

    [JsonPropertyName("Code5xx")]
    public int Code5xx { get; set; }

    [JsonPropertyName("AllOtherCodes")]
    public int AllOtherCodes { get; set; }

    [JsonPropertyName("BlockedByRobotsTxt")]
    public int BlockedByRobotsTxt { get; set; }

    [JsonPropertyName("ContainsMalware")]
    public int ContainsMalware { get; set; }
}

internal sealed class ApiUrlInfo
{
    [JsonPropertyName("Url")]
    public string Url { get; set; } = string.Empty;

    [JsonPropertyName("IsPage")]
    public bool IsPage { get; set; }

    [JsonPropertyName("HttpStatus")]
    public int HttpStatus { get; set; }

    [JsonPropertyName("DocumentSize")]
    public int DocumentSize { get; set; }

    [JsonPropertyName("AnchorCount")]
    public int AnchorCount { get; set; }

    [JsonPropertyName("DiscoveryDate")]
    public string? DiscoveryDate { get; set; }

    [JsonPropertyName("LastCrawledDate")]
    public string? LastCrawledDate { get; set; }

    [JsonPropertyName("TotalChildUrlCount")]
    public int TotalChildUrlCount { get; set; }
}

internal sealed class ApiUrlTrafficInfo
{
    [JsonPropertyName("Url")]
    public string Url { get; set; } = string.Empty;

    [JsonPropertyName("IsPage")]
    public bool IsPage { get; set; }

    [JsonPropertyName("Clicks")]
    public int Clicks { get; set; }

    [JsonPropertyName("Impressions")]
    public int Impressions { get; set; }
}

internal sealed class ApiUrlLinkDetail
{
    [JsonPropertyName("AnchorText")]
    public string? AnchorText { get; set; }

    [JsonPropertyName("Url")]
    public string? Url { get; set; }
}

internal sealed class ApiLinkDetailsResponse
{
    [JsonPropertyName("Details")]
    public ApiUrlLinkDetail[]? Details { get; set; }

    [JsonPropertyName("TotalPages")]
    public int TotalPages { get; set; }
}

internal sealed class ApiLinkCount
{
    [JsonPropertyName("Count")]
    public int Count { get; set; }

    [JsonPropertyName("Url")]
    public string? Url { get; set; }
}

internal sealed class ApiLinkCountsResponse
{
    [JsonPropertyName("Links")]
    public ApiLinkCount[]? Links { get; set; }

    [JsonPropertyName("TotalPages")]
    public int TotalPages { get; set; }
}

internal sealed class ApiRankAndTrafficStat
{
    [JsonPropertyName("Date")]
    public string Date { get; set; } = string.Empty;

    [JsonPropertyName("Clicks")]
    public int Clicks { get; set; }

    [JsonPropertyName("Impressions")]
    public int Impressions { get; set; }
}

internal sealed class ApiQueryStat
{
    [JsonPropertyName("Query")]
    public string Query { get; set; } = string.Empty;

    [JsonPropertyName("Date")]
    public string Date { get; set; } = string.Empty;

    [JsonPropertyName("Clicks")]
    public int Clicks { get; set; }

    [JsonPropertyName("Impressions")]
    public int Impressions { get; set; }

    [JsonPropertyName("AvgClickPosition")]
    public int AvgClickPosition { get; set; }

    [JsonPropertyName("AvgImpressionPosition")]
    public int AvgImpressionPosition { get; set; }
}

internal sealed class ApiKeywordStat
{
    [JsonPropertyName("Query")]
    public string Query { get; set; } = string.Empty;

    [JsonPropertyName("Date")]
    public string Date { get; set; } = string.Empty;

    [JsonPropertyName("Impressions")]
    public int Impressions { get; set; }

    [JsonPropertyName("BroadImpressions")]
    public int BroadImpressions { get; set; }
}

internal sealed class SiteUrlRequest
{
    [JsonPropertyName("siteUrl")]
    public string SiteUrl { get; set; } = string.Empty;
}

internal sealed class SiteAndFeedRequest
{
    [JsonPropertyName("siteUrl")]
    public string SiteUrl { get; set; } = string.Empty;

    [JsonPropertyName("feedUrl")]
    public string FeedUrl { get; set; } = string.Empty;
}

internal sealed class SiteAndUrlRequest
{
    [JsonPropertyName("siteUrl")]
    public string SiteUrl { get; set; } = string.Empty;

    [JsonPropertyName("url")]
    public string Url { get; set; } = string.Empty;
}

internal sealed class SiteAndUrlListRequest
{
    [JsonPropertyName("siteUrl")]
    public string SiteUrl { get; set; } = string.Empty;

    [JsonPropertyName("urlList")]
    public IReadOnlyList<string> UrlList { get; set; } = [];
}

internal sealed class IndexNowRequest
{
    [JsonPropertyName("host")]
    public string Host { get; set; } = string.Empty;

    [JsonPropertyName("key")]
    public string Key { get; set; } = string.Empty;

    [JsonPropertyName("keyLocation")]
    public string? KeyLocation { get; set; }

    [JsonPropertyName("urlList")]
    public IReadOnlyList<string> UrlList { get; set; } = [];
}

[JsonSerializable(typeof(bool))]
[JsonSerializable(typeof(ApiSite[]))]
[JsonSerializable(typeof(ApiFeed))]
[JsonSerializable(typeof(ApiFeed[]))]
[JsonSerializable(typeof(ApiUrlSubmissionQuota))]
[JsonSerializable(typeof(ApiUrlWithCrawlIssues[]))]
[JsonSerializable(typeof(ApiCrawlStats[]))]
[JsonSerializable(typeof(ApiUrlInfo))]
[JsonSerializable(typeof(ApiUrlTrafficInfo))]
[JsonSerializable(typeof(ApiLinkDetailsResponse))]
[JsonSerializable(typeof(ApiLinkCountsResponse))]
[JsonSerializable(typeof(ApiRankAndTrafficStat[]))]
[JsonSerializable(typeof(ApiQueryStat[]))]
[JsonSerializable(typeof(ApiKeywordStat[]))]
[JsonSerializable(typeof(SiteUrlRequest))]
[JsonSerializable(typeof(SiteAndFeedRequest))]
[JsonSerializable(typeof(SiteAndUrlRequest))]
[JsonSerializable(typeof(SiteAndUrlListRequest))]
[JsonSerializable(typeof(IndexNowRequest))]
[JsonSerializable(typeof(ListSitesResponse))]
[JsonSerializable(typeof(AddSiteResponse))]
[JsonSerializable(typeof(VerifySiteResponse))]
[JsonSerializable(typeof(ListSitemapsResponse))]
[JsonSerializable(typeof(GetSitemapDetailsResponse))]
[JsonSerializable(typeof(SubmitSitemapResponse))]
[JsonSerializable(typeof(SubmitUrlResponse))]
[JsonSerializable(typeof(SubmitUrlBatchResponse))]
[JsonSerializable(typeof(SubmitUrlIndexNowResponse))]
[JsonSerializable(typeof(UrlSubmissionQuotaResponse))]
[JsonSerializable(typeof(CrawlIssuesResponse))]
[JsonSerializable(typeof(CrawlStatsResponse))]
[JsonSerializable(typeof(UrlInfoResponse))]
[JsonSerializable(typeof(UrlTrafficInfoResponse))]
[JsonSerializable(typeof(UrlLinksResponse))]
[JsonSerializable(typeof(LinkCountsResponse))]
[JsonSerializable(typeof(RankAndTrafficStatsResponse))]
[JsonSerializable(typeof(QueryStatsResponse))]
[JsonSerializable(typeof(PageStatsResponse))]
[JsonSerializable(typeof(PageQueryStatsResponse))]
[JsonSerializable(typeof(QueryPageStatsResponse))]
[JsonSerializable(typeof(KeywordStatsResponse))]
[JsonSerializable(typeof(ErrorResult))]
[JsonSourceGenerationOptions(
    PropertyNamingPolicy = JsonKnownNamingPolicy.CamelCase,
    WriteIndented = false,
    DefaultIgnoreCondition = JsonIgnoreCondition.WhenWritingNull)]
internal partial class BingWebmasterJsonContext : JsonSerializerContext;
