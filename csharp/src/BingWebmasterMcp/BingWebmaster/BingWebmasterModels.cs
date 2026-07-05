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

internal sealed record RemoveSiteResponse(
    string SiteUrl,
    bool Success,
    DateTimeOffset RequestedAt);

internal sealed record SiteRoleEntry(
    string Email,
    string Role,
    string Site,
    string VerificationSite,
    bool Expired,
    string? DelegatorEmail,
    string? DelegatedCode,
    string? DelegatedCodeOwnerEmail,
    DateTimeOffset Date);

internal sealed record GetSiteRolesResponse(
    string SiteUrl,
    bool IncludeAllSubdomains,
    int RowCount,
    IReadOnlyList<SiteRoleEntry> Rows,
    DateTimeOffset QueriedAt);

internal sealed record AddSiteRoleResponse(
    string SiteUrl,
    string DelegatedUrl,
    string UserEmail,
    bool IsAdministrator,
    bool IsReadOnly,
    bool Success,
    DateTimeOffset RequestedAt);

internal sealed record RemoveSiteRoleResponse(
    string SiteUrl,
    string Email,
    string Role,
    bool Success,
    DateTimeOffset RequestedAt);

internal sealed record BlockedUrlEntry(
    string Url,
    string EntityType,
    string RequestType,
    DateTimeOffset Date);

internal sealed record GetBlockedUrlsResponse(
    string SiteUrl,
    int RowCount,
    IReadOnlyList<BlockedUrlEntry> Rows,
    DateTimeOffset QueriedAt);

internal sealed record AddBlockedUrlResponse(
    string SiteUrl,
    string Url,
    string EntityType,
    string RequestType,
    bool Success,
    DateTimeOffset RequestedAt);

internal sealed record RemoveBlockedUrlResponse(
    string SiteUrl,
    string Url,
    string EntityType,
    string RequestType,
    bool Success,
    DateTimeOffset RequestedAt);

internal sealed record DetailedQueryStatEntry(
    DateTimeOffset Date,
    int Clicks,
    int Impressions,
    int Position);

internal sealed record QueryPageDetailStatsResponse(
    string SiteUrl,
    string Query,
    string Page,
    int RowCount,
    IReadOnlyList<DetailedQueryStatEntry> Rows,
    DateTimeOffset QueriedAt);

internal sealed record QueryTrafficStatsResponse(
    string SiteUrl,
    string Query,
    int RowCount,
    IReadOnlyList<RankAndTrafficStatEntry> Rows,
    DateTimeOffset QueriedAt);

internal sealed record GetKeywordResponse(
    string Query,
    string Country,
    string Language,
    string StartDate,
    string EndDate,
    bool Found,
    int Impressions,
    int BroadImpressions,
    DateTimeOffset QueriedAt);

internal sealed record RelatedKeywordEntry(
    string Query,
    int Impressions,
    int BroadImpressions);

internal sealed record RelatedKeywordsResponse(
    string Query,
    string Country,
    string Language,
    string StartDate,
    string EndDate,
    int RowCount,
    IReadOnlyList<RelatedKeywordEntry> Rows,
    DateTimeOffset QueriedAt);

internal sealed record ChildUrlInfoEntry(
    string Url,
    bool IsPage,
    int HttpStatus,
    int DocumentSize,
    int AnchorCount,
    DateTimeOffset? DiscoveryDate,
    DateTimeOffset? LastCrawledDate,
    int TotalChildUrlCount);

internal sealed record ChildrenUrlInfoResponse(
    string SiteUrl,
    string Url,
    int Page,
    string CrawlDateFilter,
    string DiscoveredDateFilter,
    string DocFlagsFilter,
    string HttpCodeFilter,
    int RowCount,
    IReadOnlyList<ChildUrlInfoEntry> Rows,
    DateTimeOffset QueriedAt);

internal sealed record ChildUrlTrafficInfoEntry(
    string Url,
    bool IsPage,
    int Clicks,
    int Impressions);

internal sealed record ChildrenUrlTrafficInfoResponse(
    string SiteUrl,
    string Url,
    int Page,
    int RowCount,
    IReadOnlyList<ChildUrlTrafficInfoEntry> Rows,
    DateTimeOffset QueriedAt);

internal sealed record FetchUrlResponse(
    string SiteUrl,
    string Url,
    bool Success,
    DateTimeOffset RequestedAt);

internal sealed record FetchedUrlEntry(
    string Url,
    DateTimeOffset Date,
    bool Fetched,
    bool Expired);

internal sealed record ListFetchedUrlsResponse(
    string SiteUrl,
    int RowCount,
    IReadOnlyList<FetchedUrlEntry> Rows,
    DateTimeOffset QueriedAt);

internal sealed record FetchedUrlDetailsResponse(
    string SiteUrl,
    string Url,
    DateTimeOffset Date,
    string Status,
    string Headers,
    string Document,
    DateTimeOffset QueriedAt);

internal sealed record RemoveSitemapResponse(
    string SiteUrl,
    string FeedUrl,
    bool Success,
    DateTimeOffset RequestedAt);

internal sealed record SiteMoveEntry(
    DateTimeOffset Date,
    string MoveScope,
    string MoveType,
    string SourceUrl,
    string TargetUrl);

internal sealed record GetSiteMovesResponse(
    string SiteUrl,
    int RowCount,
    IReadOnlyList<SiteMoveEntry> Rows,
    DateTimeOffset QueriedAt);

internal sealed record SubmitSiteMoveResponse(
    string SiteUrl,
    string SourceUrl,
    string TargetUrl,
    string MoveType,
    string MoveScope,
    bool Success,
    DateTimeOffset RequestedAt);

internal sealed record SubmitContentResponse(
    string SiteUrl,
    string Url,
    string DynamicServing,
    bool Success,
    DateTimeOffset RequestedAt);

internal sealed record ContentSubmissionQuotaResponse(
    string SiteUrl,
    int DailyQuota,
    int MonthlyQuota,
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

internal sealed class ApiSiteRole
{
    [JsonPropertyName("Date")]
    public string Date { get; set; } = string.Empty;

    [JsonPropertyName("DelegatedCode")]
    public string? DelegatedCode { get; set; }

    [JsonPropertyName("DelegatorEmail")]
    public string? DelegatorEmail { get; set; }

    [JsonPropertyName("DelegatedCodeOwnerEmail")]
    public string? DelegatedCodeOwnerEmail { get; set; }

    [JsonPropertyName("Email")]
    public string Email { get; set; } = string.Empty;

    [JsonPropertyName("Expired")]
    public bool Expired { get; set; }

    [JsonPropertyName("Role")]
    public int Role { get; set; }

    [JsonPropertyName("Site")]
    public string Site { get; set; } = string.Empty;

    [JsonPropertyName("VerificationSite")]
    public string VerificationSite { get; set; } = string.Empty;
}

internal sealed class ApiBlockedUrl
{
    [JsonPropertyName("Date")]
    public string Date { get; set; } = string.Empty;

    [JsonPropertyName("EntityType")]
    public int EntityType { get; set; }

    [JsonPropertyName("RequestType")]
    public int RequestType { get; set; }

    [JsonPropertyName("Url")]
    public string Url { get; set; } = string.Empty;
}

internal sealed class ApiDetailedQueryStat
{
    [JsonPropertyName("Date")]
    public string Date { get; set; } = string.Empty;

    [JsonPropertyName("Clicks")]
    public int Clicks { get; set; }

    [JsonPropertyName("Impressions")]
    public int Impressions { get; set; }

    [JsonPropertyName("Position")]
    public int Position { get; set; }
}

internal sealed class ApiKeywordDetails
{
    [JsonPropertyName("Query")]
    public string? Query { get; set; }

    [JsonPropertyName("BroadImpressions")]
    public int BroadImpressions { get; set; }

    [JsonPropertyName("Impressions")]
    public int Impressions { get; set; }
}

internal sealed class ApiFetchedUrl
{
    [JsonPropertyName("Date")]
    public string Date { get; set; } = string.Empty;

    [JsonPropertyName("Expired")]
    public bool Expired { get; set; }

    [JsonPropertyName("Fetched")]
    public bool Fetched { get; set; }

    [JsonPropertyName("Url")]
    public string Url { get; set; } = string.Empty;
}

internal sealed class ApiFetchedUrlDetails
{
    [JsonPropertyName("Date")]
    public string Date { get; set; } = string.Empty;

    [JsonPropertyName("Document")]
    public string Document { get; set; } = string.Empty;

    [JsonPropertyName("Headers")]
    public string Headers { get; set; } = string.Empty;

    [JsonPropertyName("Status")]
    public string Status { get; set; } = string.Empty;

    [JsonPropertyName("Url")]
    public string Url { get; set; } = string.Empty;
}

internal sealed class ApiSiteMoveSettings
{
    [JsonPropertyName("Date")]
    public string Date { get; set; } = string.Empty;

    [JsonPropertyName("MoveScope")]
    public int MoveScope { get; set; }

    [JsonPropertyName("MoveType")]
    public int MoveType { get; set; }

    [JsonPropertyName("SourceUrl")]
    public string SourceUrl { get; set; } = string.Empty;

    [JsonPropertyName("TargetUrl")]
    public string TargetUrl { get; set; } = string.Empty;
}

internal sealed class ApiFilterProperties
{
    [JsonPropertyName("CrawlDateFilter")]
    public int CrawlDateFilter { get; set; }

    [JsonPropertyName("DiscoveredDateFilter")]
    public int DiscoveredDateFilter { get; set; }

    [JsonPropertyName("DocFlagsFilters")]
    public int DocFlagsFilters { get; set; }

    [JsonPropertyName("HttpCodeFilters")]
    public int HttpCodeFilters { get; set; }
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

internal sealed class AddSiteRoleRequest
{
    [JsonPropertyName("siteUrl")]
    public string SiteUrl { get; set; } = string.Empty;

    [JsonPropertyName("delegatedUrl")]
    public string DelegatedUrl { get; set; } = string.Empty;

    [JsonPropertyName("userEmail")]
    public string UserEmail { get; set; } = string.Empty;

    [JsonPropertyName("authenticationCode")]
    public string AuthenticationCode { get; set; } = string.Empty;

    [JsonPropertyName("isAdministrator")]
    public bool IsAdministrator { get; set; }

    [JsonPropertyName("isReadOnly")]
    public bool IsReadOnly { get; set; }
}

internal sealed class RemoveSiteRoleRequest
{
    [JsonPropertyName("siteUrl")]
    public string SiteUrl { get; set; } = string.Empty;

    [JsonPropertyName("siteRole")]
    public RemoveSiteRoleItem SiteRole { get; set; } = new();
}

internal sealed class RemoveSiteRoleItem
{
    [JsonPropertyName("Date")]
    public string Date { get; set; } = string.Empty;

    [JsonPropertyName("Email")]
    public string Email { get; set; } = string.Empty;

    [JsonPropertyName("Role")]
    public int Role { get; set; }

    [JsonPropertyName("Site")]
    public string Site { get; set; } = string.Empty;

    [JsonPropertyName("VerificationSite")]
    public string VerificationSite { get; set; } = string.Empty;

    [JsonPropertyName("DelegatedCode")]
    public string? DelegatedCode { get; set; }

    [JsonPropertyName("DelegatorEmail")]
    public string? DelegatorEmail { get; set; }

    [JsonPropertyName("DelegatedCodeOwnerEmail")]
    public string? DelegatedCodeOwnerEmail { get; set; }
}

internal sealed class BlockedUrlRequest
{
    [JsonPropertyName("siteUrl")]
    public string SiteUrl { get; set; } = string.Empty;

    [JsonPropertyName("blockedUrl")]
    public ApiBlockedUrl BlockedUrl { get; set; } = new();
}

internal sealed class ChildrenUrlInfoRequest
{
    [JsonPropertyName("siteUrl")]
    public string SiteUrl { get; set; } = string.Empty;

    [JsonPropertyName("url")]
    public string Url { get; set; } = string.Empty;

    [JsonPropertyName("page")]
    public int Page { get; set; }

    [JsonPropertyName("filterProperties")]
    public ApiFilterProperties FilterProperties { get; set; } = new();
}

internal sealed class SubmitSiteMoveRequest
{
    [JsonPropertyName("siteUrl")]
    public string SiteUrl { get; set; } = string.Empty;

    [JsonPropertyName("settings")]
    public ApiSiteMoveSettings Settings { get; set; } = new();
}

internal sealed class SubmitContentRequest
{
    [JsonPropertyName("siteUrl")]
    public string SiteUrl { get; set; } = string.Empty;

    [JsonPropertyName("url")]
    public string Url { get; set; } = string.Empty;

    [JsonPropertyName("httpMessage")]
    public string HttpMessage { get; set; } = string.Empty;

    [JsonPropertyName("structuredData")]
    public string StructuredData { get; set; } = string.Empty;

    [JsonPropertyName("dynamicServing")]
    public int DynamicServing { get; set; }
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
[JsonSerializable(typeof(ApiSiteRole[]))]
[JsonSerializable(typeof(ApiBlockedUrl[]))]
[JsonSerializable(typeof(ApiDetailedQueryStat[]))]
[JsonSerializable(typeof(ApiKeywordDetails))]
[JsonSerializable(typeof(ApiKeywordDetails[]))]
[JsonSerializable(typeof(ApiUrlInfo[]))]
[JsonSerializable(typeof(ApiUrlTrafficInfo[]))]
[JsonSerializable(typeof(ApiFetchedUrl[]))]
[JsonSerializable(typeof(ApiFetchedUrlDetails))]
[JsonSerializable(typeof(ApiSiteMoveSettings[]))]
[JsonSerializable(typeof(ApiFilterProperties))]
[JsonSerializable(typeof(SiteUrlRequest))]
[JsonSerializable(typeof(SiteAndFeedRequest))]
[JsonSerializable(typeof(SiteAndUrlRequest))]
[JsonSerializable(typeof(SiteAndUrlListRequest))]
[JsonSerializable(typeof(AddSiteRoleRequest))]
[JsonSerializable(typeof(RemoveSiteRoleRequest))]
[JsonSerializable(typeof(RemoveSiteRoleItem))]
[JsonSerializable(typeof(BlockedUrlRequest))]
[JsonSerializable(typeof(ChildrenUrlInfoRequest))]
[JsonSerializable(typeof(SubmitSiteMoveRequest))]
[JsonSerializable(typeof(SubmitContentRequest))]
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
[JsonSerializable(typeof(RemoveSiteResponse))]
[JsonSerializable(typeof(GetSiteRolesResponse))]
[JsonSerializable(typeof(AddSiteRoleResponse))]
[JsonSerializable(typeof(RemoveSiteRoleResponse))]
[JsonSerializable(typeof(GetBlockedUrlsResponse))]
[JsonSerializable(typeof(AddBlockedUrlResponse))]
[JsonSerializable(typeof(RemoveBlockedUrlResponse))]
[JsonSerializable(typeof(QueryPageDetailStatsResponse))]
[JsonSerializable(typeof(QueryTrafficStatsResponse))]
[JsonSerializable(typeof(GetKeywordResponse))]
[JsonSerializable(typeof(RelatedKeywordsResponse))]
[JsonSerializable(typeof(ChildrenUrlInfoResponse))]
[JsonSerializable(typeof(ChildrenUrlTrafficInfoResponse))]
[JsonSerializable(typeof(FetchUrlResponse))]
[JsonSerializable(typeof(ListFetchedUrlsResponse))]
[JsonSerializable(typeof(FetchedUrlDetailsResponse))]
[JsonSerializable(typeof(RemoveSitemapResponse))]
[JsonSerializable(typeof(GetSiteMovesResponse))]
[JsonSerializable(typeof(SubmitSiteMoveResponse))]
[JsonSerializable(typeof(SubmitContentResponse))]
[JsonSerializable(typeof(ContentSubmissionQuotaResponse))]
[JsonSerializable(typeof(ErrorResult))]
[JsonSourceGenerationOptions(
    PropertyNamingPolicy = JsonKnownNamingPolicy.CamelCase,
    WriteIndented = false,
    DefaultIgnoreCondition = JsonIgnoreCondition.WhenWritingNull)]
internal partial class BingWebmasterJsonContext : JsonSerializerContext;
