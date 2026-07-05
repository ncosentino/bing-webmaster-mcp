using System.Globalization;
using System.Text;
using System.Text.Json;
using System.Text.Json.Serialization.Metadata;

namespace BingWebmasterMcp.BingWebmaster;

/// <summary>Thrown when the Bing Webmaster API returns a non-success status code.</summary>
internal sealed class BingWebmasterApiException : Exception
{
    internal int StatusCode { get; }

    internal BingWebmasterApiException(int statusCode, string message) : base(message)
        => StatusCode = statusCode;
}

/// <summary>Client for the Bing Webmaster JSON API.</summary>
internal sealed class BingWebmasterClient(HttpClient httpClient, string apiKey, string? baseUrlOverride = null)
{
    private const string DefaultBaseUrl = "https://ssl.bing.com/webmaster/api.svc/json";

    private readonly string _apiKey = string.IsNullOrWhiteSpace(apiKey)
        ? throw new ArgumentException("API key is required.", nameof(apiKey))
        : apiKey;
    private readonly string _baseUrl = baseUrlOverride ?? DefaultBaseUrl;
    private readonly HttpClient _httpClient = httpClient;

    internal async Task<ListSitesResponse> ListSitesAsync(CancellationToken cancellationToken = default)
    {
        var rawSites = await GetEnvelopeAsync("GetUserSites", null, BingWebmasterJsonContext.Default.ApiSiteArray, cancellationToken).ConfigureAwait(false)
            ?? [];

        var sites = rawSites
            .Select(site => new SiteInfo(
                site.Url,
                site.IsVerified,
                site.DnsVerificationCode,
                site.AuthenticationCode))
            .ToList();

        return new ListSitesResponse(sites, DateTimeOffset.UtcNow);
    }

    internal async Task<AddSiteResponse> AddSiteAsync(string siteUrl, CancellationToken cancellationToken = default)
    {
        var normalizedSiteUrl = RequireText(siteUrl, nameof(siteUrl));
        await PostCommandAsync(
            "AddSite",
            new SiteUrlRequest { SiteUrl = normalizedSiteUrl },
            BingWebmasterJsonContext.Default.SiteUrlRequest,
            cancellationToken).ConfigureAwait(false);

        return new AddSiteResponse(normalizedSiteUrl, true, DateTimeOffset.UtcNow);
    }

    internal async Task<VerifySiteResponse> VerifySiteAsync(string siteUrl, CancellationToken cancellationToken = default)
    {
        var normalizedSiteUrl = RequireText(siteUrl, nameof(siteUrl));
        var verified = await PostQueryAsync(
            "VerifySite",
            new SiteUrlRequest { SiteUrl = normalizedSiteUrl },
            BingWebmasterJsonContext.Default.SiteUrlRequest,
            cancellationToken).ConfigureAwait(false);

        return new VerifySiteResponse(normalizedSiteUrl, verified, DateTimeOffset.UtcNow);
    }

    internal async Task<ListSitemapsResponse> ListSitemapsAsync(string siteUrl, CancellationToken cancellationToken = default)
    {
        var normalizedSiteUrl = RequireText(siteUrl, nameof(siteUrl));
        var rawFeeds = await GetEnvelopeAsync(
            "GetFeeds",
            new Dictionary<string, string?>
            {
                ["siteUrl"] = normalizedSiteUrl
            },
            BingWebmasterJsonContext.Default.ApiFeedArray,
            cancellationToken).ConfigureAwait(false) ?? [];

        return new ListSitemapsResponse(
            normalizedSiteUrl,
            rawFeeds.Select(ToSitemapInfo).ToList(),
            DateTimeOffset.UtcNow);
    }

    internal async Task<GetSitemapDetailsResponse> GetSitemapDetailsAsync(
        string siteUrl,
        string feedUrl,
        CancellationToken cancellationToken = default)
    {
        var normalizedSiteUrl = RequireText(siteUrl, nameof(siteUrl));
        var normalizedFeedUrl = RequireText(feedUrl, nameof(feedUrl));

        var request = new HttpRequestMessage(
            HttpMethod.Get,
            BuildUrl(
                "GetFeedDetails",
                new Dictionary<string, string?>
                {
                    ["siteUrl"] = normalizedSiteUrl,
                    ["feedUrl"] = normalizedFeedUrl
                }));

        using var response = await _httpClient.SendAsync(request, cancellationToken).ConfigureAwait(false);
        await EnsureSuccessAsync(response, cancellationToken).ConfigureAwait(false);

        using var document = await ReadEnvelopeDocumentAsync(response, cancellationToken).ConfigureAwait(false);
        SitemapInfo? details = null;
        var payload = document.RootElement.GetProperty("d");

        // Bing's GetFeedDetails wire shape is not confirmed by a recorded example. We accept either
        // a single feed object or an array and map the first item if Bing returns a collection shape.
        if (payload.ValueKind == JsonValueKind.Object)
        {
            var rawFeed = payload.Deserialize(BingWebmasterJsonContext.Default.ApiFeed);
            if (rawFeed is not null)
                details = ToSitemapInfo(rawFeed);
        }
        else if (payload.ValueKind == JsonValueKind.Array)
        {
            var rawFeeds = payload.Deserialize(BingWebmasterJsonContext.Default.ApiFeedArray);
            if (rawFeeds is { Length: > 0 })
                details = ToSitemapInfo(rawFeeds[0]);
        }

        return new GetSitemapDetailsResponse(
            normalizedSiteUrl,
            normalizedFeedUrl,
            details,
            DateTimeOffset.UtcNow);
    }

    internal async Task<SubmitSitemapResponse> SubmitSitemapAsync(
        string siteUrl,
        string feedUrl,
        CancellationToken cancellationToken = default)
    {
        var normalizedSiteUrl = RequireText(siteUrl, nameof(siteUrl));
        var normalizedFeedUrl = RequireText(feedUrl, nameof(feedUrl));
        await PostCommandAsync(
            "SubmitFeed",
            new SiteAndFeedRequest
            {
                SiteUrl = normalizedSiteUrl,
                FeedUrl = normalizedFeedUrl
            },
            BingWebmasterJsonContext.Default.SiteAndFeedRequest,
            cancellationToken).ConfigureAwait(false);

        return new SubmitSitemapResponse(normalizedSiteUrl, normalizedFeedUrl, true, DateTimeOffset.UtcNow);
    }

    internal async Task<SubmitUrlResponse> SubmitUrlAsync(
        string siteUrl,
        string url,
        CancellationToken cancellationToken = default)
    {
        var normalizedSiteUrl = RequireText(siteUrl, nameof(siteUrl));
        var normalizedUrl = RequireText(url, nameof(url));
        await PostCommandAsync(
            "SubmitUrl",
            new SiteAndUrlRequest
            {
                SiteUrl = normalizedSiteUrl,
                Url = normalizedUrl
            },
            BingWebmasterJsonContext.Default.SiteAndUrlRequest,
            cancellationToken).ConfigureAwait(false);

        return new SubmitUrlResponse(normalizedSiteUrl, normalizedUrl, true, DateTimeOffset.UtcNow);
    }

    internal async Task<SubmitUrlBatchResponse> SubmitUrlBatchAsync(
        string siteUrl,
        IReadOnlyList<string> urlList,
        CancellationToken cancellationToken = default)
    {
        var normalizedSiteUrl = RequireText(siteUrl, nameof(siteUrl));
        var normalizedUrls = NormalizeUrlList(urlList, 500);
        await PostCommandAsync(
            "SubmitUrlBatch",
            new SiteAndUrlListRequest
            {
                SiteUrl = normalizedSiteUrl,
                UrlList = normalizedUrls
            },
            BingWebmasterJsonContext.Default.SiteAndUrlListRequest,
            cancellationToken).ConfigureAwait(false);

        return new SubmitUrlBatchResponse(
            normalizedSiteUrl,
            normalizedUrls,
            normalizedUrls.Count,
            true,
            DateTimeOffset.UtcNow);
    }

    internal async Task<UrlSubmissionQuotaResponse> GetUrlSubmissionQuotaAsync(
        string siteUrl,
        CancellationToken cancellationToken = default)
    {
        var normalizedSiteUrl = RequireText(siteUrl, nameof(siteUrl));
        var rawQuota = await GetEnvelopeAsync(
            "GetUrlSubmissionQuota",
            new Dictionary<string, string?>
            {
                ["siteUrl"] = normalizedSiteUrl
            },
            BingWebmasterJsonContext.Default.ApiUrlSubmissionQuota,
            cancellationToken).ConfigureAwait(false)
            ?? throw new InvalidOperationException("Bing returned an empty URL submission quota payload.");

        return new UrlSubmissionQuotaResponse(
            normalizedSiteUrl,
            rawQuota.DailyQuota,
            rawQuota.MonthlyQuota,
            DateTimeOffset.UtcNow);
    }

    internal async Task<CrawlIssuesResponse> GetCrawlIssuesAsync(
        string siteUrl,
        CancellationToken cancellationToken = default)
    {
        var normalizedSiteUrl = RequireText(siteUrl, nameof(siteUrl));
        var rawIssues = await GetEnvelopeAsync(
            "GetCrawlIssues",
            new Dictionary<string, string?>
            {
                ["siteUrl"] = normalizedSiteUrl
            },
            BingWebmasterJsonContext.Default.ApiUrlWithCrawlIssuesArray,
            cancellationToken).ConfigureAwait(false) ?? [];

        return new CrawlIssuesResponse(
            normalizedSiteUrl,
            rawIssues.Select(issue => new CrawlIssueEntry(
                issue.Url,
                issue.HttpCode,
                DecodeCrawlIssues(issue.Issues),
                issue.InLinks)).ToList(),
            DateTimeOffset.UtcNow);
    }

    internal async Task<CrawlStatsResponse> GetCrawlStatsAsync(
        string siteUrl,
        CancellationToken cancellationToken = default)
    {
        var normalizedSiteUrl = RequireText(siteUrl, nameof(siteUrl));
        var rawStats = await GetEnvelopeAsync(
            "GetCrawlStats",
            new Dictionary<string, string?>
            {
                ["siteUrl"] = normalizedSiteUrl
            },
            BingWebmasterJsonContext.Default.ApiCrawlStatsArray,
            cancellationToken).ConfigureAwait(false) ?? [];

        var rows = rawStats.Select(stat => new CrawlStatEntry(
            ParseRequiredDate(stat.Date, nameof(stat.Date)),
            stat.CrawledPages,
            stat.CrawlErrors,
            stat.InIndex,
            stat.InLinks,
            stat.Code2xx,
            stat.Code301,
            stat.Code302,
            stat.Code4xx,
            stat.Code5xx,
            stat.AllOtherCodes,
            stat.BlockedByRobotsTxt,
            stat.ContainsMalware)).ToList();

        return new CrawlStatsResponse(normalizedSiteUrl, rows.Count, rows, DateTimeOffset.UtcNow);
    }

    internal async Task<UrlInfoResponse> GetUrlInfoAsync(
        string siteUrl,
        string url,
        CancellationToken cancellationToken = default)
    {
        var normalizedSiteUrl = RequireText(siteUrl, nameof(siteUrl));
        var normalizedUrl = RequireText(url, nameof(url));
        var rawInfo = await GetEnvelopeAsync(
            "GetUrlInfo",
            new Dictionary<string, string?>
            {
                ["siteUrl"] = normalizedSiteUrl,
                ["url"] = normalizedUrl
            },
            BingWebmasterJsonContext.Default.ApiUrlInfo,
            cancellationToken).ConfigureAwait(false)
            ?? throw new InvalidOperationException("Bing returned an empty URL info payload.");

        return new UrlInfoResponse(
            normalizedSiteUrl,
            rawInfo.Url,
            rawInfo.IsPage,
            rawInfo.HttpStatus,
            rawInfo.DocumentSize,
            rawInfo.AnchorCount,
            ParseOptionalDate(rawInfo.DiscoveryDate),
            ParseOptionalDate(rawInfo.LastCrawledDate),
            rawInfo.TotalChildUrlCount,
            DateTimeOffset.UtcNow);
    }

    internal async Task<UrlTrafficInfoResponse> GetUrlTrafficInfoAsync(
        string siteUrl,
        string url,
        CancellationToken cancellationToken = default)
    {
        var normalizedSiteUrl = RequireText(siteUrl, nameof(siteUrl));
        var normalizedUrl = RequireText(url, nameof(url));
        var rawInfo = await GetEnvelopeAsync(
            "GetUrlTrafficInfo",
            new Dictionary<string, string?>
            {
                ["siteUrl"] = normalizedSiteUrl,
                ["url"] = normalizedUrl
            },
            BingWebmasterJsonContext.Default.ApiUrlTrafficInfo,
            cancellationToken).ConfigureAwait(false)
            ?? throw new InvalidOperationException("Bing returned an empty URL traffic info payload.");

        return new UrlTrafficInfoResponse(
            normalizedSiteUrl,
            rawInfo.Url,
            rawInfo.IsPage,
            rawInfo.Clicks,
            rawInfo.Impressions,
            DateTimeOffset.UtcNow);
    }

    internal async Task<UrlLinksResponse> GetUrlLinksAsync(
        string siteUrl,
        string link,
        int page = 0,
        CancellationToken cancellationToken = default)
    {
        var normalizedSiteUrl = RequireText(siteUrl, nameof(siteUrl));
        var normalizedLink = RequireText(link, nameof(link));
        if (page < 0)
            throw new ArgumentOutOfRangeException(nameof(page), "Page must be zero or greater.");

        var rawResponse = await GetEnvelopeAsync(
            "GetUrlLinks",
            new Dictionary<string, string?>
            {
                ["siteUrl"] = normalizedSiteUrl,
                ["link"] = normalizedLink,
                ["page"] = page.ToString(CultureInfo.InvariantCulture)
            },
            BingWebmasterJsonContext.Default.ApiLinkDetailsResponse,
            cancellationToken).ConfigureAwait(false)
            ?? new ApiLinkDetailsResponse();

        var details = (rawResponse.Details ?? [])
            .Select(detail => new LinkDetailEntry(detail.AnchorText, detail.Url))
            .ToList();

        return new UrlLinksResponse(
            normalizedSiteUrl,
            normalizedLink,
            page,
            rawResponse.TotalPages,
            details,
            DateTimeOffset.UtcNow);
    }

    internal async Task<LinkCountsResponse> GetLinkCountsAsync(
        string siteUrl,
        int page = 0,
        CancellationToken cancellationToken = default)
    {
        var normalizedSiteUrl = RequireText(siteUrl, nameof(siteUrl));
        if (page < 0)
            throw new ArgumentOutOfRangeException(nameof(page), "Page must be zero or greater.");

        var rawResponse = await GetEnvelopeAsync(
            "GetLinkCounts",
            new Dictionary<string, string?>
            {
                ["siteUrl"] = normalizedSiteUrl,
                ["page"] = page.ToString(CultureInfo.InvariantCulture)
            },
            BingWebmasterJsonContext.Default.ApiLinkCountsResponse,
            cancellationToken).ConfigureAwait(false)
            ?? new ApiLinkCountsResponse();

        var links = (rawResponse.Links ?? [])
            .Select(linkCount => new LinkCountEntry(linkCount.Count, linkCount.Url))
            .ToList();

        return new LinkCountsResponse(
            normalizedSiteUrl,
            page,
            rawResponse.TotalPages,
            links,
            DateTimeOffset.UtcNow);
    }

    internal async Task<RankAndTrafficStatsResponse> GetRankAndTrafficStatsAsync(
        string siteUrl,
        CancellationToken cancellationToken = default)
    {
        var normalizedSiteUrl = RequireText(siteUrl, nameof(siteUrl));
        var rawStats = await GetEnvelopeAsync(
            "GetRankAndTrafficStats",
            new Dictionary<string, string?>
            {
                ["siteUrl"] = normalizedSiteUrl
            },
            BingWebmasterJsonContext.Default.ApiRankAndTrafficStatArray,
            cancellationToken).ConfigureAwait(false) ?? [];

        var rows = rawStats.Select(stat => new RankAndTrafficStatEntry(
            ParseRequiredDate(stat.Date, nameof(stat.Date)),
            stat.Clicks,
            stat.Impressions)).ToList();

        return new RankAndTrafficStatsResponse(normalizedSiteUrl, rows.Count, rows, DateTimeOffset.UtcNow);
    }

    internal async Task<QueryStatsResponse> GetQueryStatsAsync(
        string siteUrl,
        CancellationToken cancellationToken = default)
    {
        var normalizedSiteUrl = RequireText(siteUrl, nameof(siteUrl));
        var rows = await GetQueryStatsRowsAsync(
            "GetQueryStats",
            new Dictionary<string, string?> { ["siteUrl"] = normalizedSiteUrl },
            cancellationToken).ConfigureAwait(false);

        return new QueryStatsResponse(normalizedSiteUrl, rows.Count, rows, DateTimeOffset.UtcNow);
    }

    internal async Task<PageStatsResponse> GetPageStatsAsync(
        string siteUrl,
        CancellationToken cancellationToken = default)
    {
        var normalizedSiteUrl = RequireText(siteUrl, nameof(siteUrl));
        var rows = await GetPageStatsRowsAsync(
            "GetPageStats",
            new Dictionary<string, string?> { ["siteUrl"] = normalizedSiteUrl },
            cancellationToken).ConfigureAwait(false);

        return new PageStatsResponse(normalizedSiteUrl, rows.Count, rows, DateTimeOffset.UtcNow);
    }

    internal async Task<PageQueryStatsResponse> GetPageQueryStatsAsync(
        string siteUrl,
        string page,
        CancellationToken cancellationToken = default)
    {
        var normalizedSiteUrl = RequireText(siteUrl, nameof(siteUrl));
        var normalizedPage = RequireText(page, nameof(page));
        var rows = await GetQueryStatsRowsAsync(
            "GetPageQueryStats",
            new Dictionary<string, string?>
            {
                ["siteUrl"] = normalizedSiteUrl,
                ["page"] = normalizedPage
            },
            cancellationToken).ConfigureAwait(false);

        return new PageQueryStatsResponse(normalizedSiteUrl, normalizedPage, rows.Count, rows, DateTimeOffset.UtcNow);
    }

    internal async Task<QueryPageStatsResponse> GetQueryPageStatsAsync(
        string siteUrl,
        string query,
        CancellationToken cancellationToken = default)
    {
        var normalizedSiteUrl = RequireText(siteUrl, nameof(siteUrl));
        var normalizedQuery = RequireText(query, nameof(query));
        var rows = await GetPageStatsRowsAsync(
            "GetQueryPageStats",
            new Dictionary<string, string?>
            {
                ["siteUrl"] = normalizedSiteUrl,
                ["query"] = normalizedQuery
            },
            cancellationToken).ConfigureAwait(false);

        return new QueryPageStatsResponse(normalizedSiteUrl, normalizedQuery, rows.Count, rows, DateTimeOffset.UtcNow);
    }

    internal async Task<KeywordStatsResponse> GetKeywordStatsAsync(
        string query,
        string country,
        string language,
        CancellationToken cancellationToken = default)
    {
        var normalizedQuery = RequireText(query, nameof(query));
        var normalizedCountry = RequireText(country, nameof(country));
        var normalizedLanguage = RequireText(language, nameof(language));
        var rawRows = await GetEnvelopeAsync(
            "GetKeywordStats",
            new Dictionary<string, string?>
            {
                ["q"] = normalizedQuery,
                ["country"] = normalizedCountry,
                ["language"] = normalizedLanguage
            },
            BingWebmasterJsonContext.Default.ApiKeywordStatArray,
            cancellationToken).ConfigureAwait(false) ?? [];

        var rows = rawRows.Select(row => new KeywordStatsEntry(
            row.Query,
            ParseRequiredDate(row.Date, nameof(row.Date)),
            row.Impressions,
            row.BroadImpressions)).ToList();

        return new KeywordStatsResponse(
            normalizedQuery,
            normalizedCountry,
            normalizedLanguage,
            rows.Count,
            rows,
            DateTimeOffset.UtcNow);
    }

    internal async Task<RemoveSiteResponse> RemoveSiteAsync(string siteUrl, CancellationToken cancellationToken = default)
    {
        var normalizedSiteUrl = RequireText(siteUrl, nameof(siteUrl));
        await PostCommandAsync(
            "RemoveSite",
            new SiteUrlRequest { SiteUrl = normalizedSiteUrl },
            BingWebmasterJsonContext.Default.SiteUrlRequest,
            cancellationToken).ConfigureAwait(false);

        return new RemoveSiteResponse(normalizedSiteUrl, true, DateTimeOffset.UtcNow);
    }

    internal async Task<GetSiteRolesResponse> GetSiteRolesAsync(
        string siteUrl,
        bool includeAllSubdomains,
        CancellationToken cancellationToken = default)
    {
        var normalizedSiteUrl = RequireText(siteUrl, nameof(siteUrl));
        var rawRoles = await GetEnvelopeAsync(
            "GetSiteRoles",
            new Dictionary<string, string?>
            {
                ["siteUrl"] = normalizedSiteUrl,
                ["includeAllSubdomains"] = ToBooleanString(includeAllSubdomains)
            },
            BingWebmasterJsonContext.Default.ApiSiteRoleArray,
            cancellationToken).ConfigureAwait(false) ?? [];

        var rows = rawRoles.Select(role => new SiteRoleEntry(
            role.Email,
            DecodeRole(role.Role),
            role.Site,
            role.VerificationSite,
            role.Expired,
            role.DelegatorEmail,
            role.DelegatedCode,
            role.DelegatedCodeOwnerEmail,
            ParseRequiredDate(role.Date, nameof(role.Date)))).ToList();

        return new GetSiteRolesResponse(
            normalizedSiteUrl,
            includeAllSubdomains,
            rows.Count,
            rows,
            DateTimeOffset.UtcNow);
    }

    internal async Task<AddSiteRoleResponse> AddSiteRoleAsync(
        string siteUrl,
        string delegatedUrl,
        string userEmail,
        string authenticationCode,
        bool isAdministrator,
        bool isReadOnly,
        CancellationToken cancellationToken = default)
    {
        var normalizedSiteUrl = RequireText(siteUrl, nameof(siteUrl));
        var normalizedDelegatedUrl = RequireText(delegatedUrl, nameof(delegatedUrl));
        var normalizedUserEmail = RequireText(userEmail, nameof(userEmail));
        var normalizedAuthenticationCode = RequireText(authenticationCode, nameof(authenticationCode));

        await PostCommandAsync(
            "AddSiteRoles",
            new AddSiteRoleRequest
            {
                SiteUrl = normalizedSiteUrl,
                DelegatedUrl = normalizedDelegatedUrl,
                UserEmail = normalizedUserEmail,
                AuthenticationCode = normalizedAuthenticationCode,
                IsAdministrator = isAdministrator,
                IsReadOnly = isReadOnly
            },
            BingWebmasterJsonContext.Default.AddSiteRoleRequest,
            cancellationToken).ConfigureAwait(false);

        return new AddSiteRoleResponse(
            normalizedSiteUrl,
            normalizedDelegatedUrl,
            normalizedUserEmail,
            isAdministrator,
            isReadOnly,
            true,
            DateTimeOffset.UtcNow);
    }

    internal async Task<RemoveSiteRoleResponse> RemoveSiteRoleAsync(
        string siteUrl,
        string email,
        string role,
        CancellationToken cancellationToken = default)
    {
        var normalizedSiteUrl = RequireText(siteUrl, nameof(siteUrl));
        var normalizedEmail = RequireText(email, nameof(email));
        var roleValue = ParseRole(role, nameof(role));
        var requestDate = DateTimeOffset.UtcNow;

        await PostCommandAsync(
            "RemoveSiteRole",
            new RemoveSiteRoleRequest
            {
                SiteUrl = normalizedSiteUrl,
                SiteRole = new RemoveSiteRoleItem
                {
                    // Bing documents a SiteRole object here but does not confirm which fields are
                    // required to identify the role for removal. We mirror the existing site + role
                    // identity fields and synthesize Date/VerificationSite from the current request
                    // so callers can use a stable remove_site_role(site_url, email, role) contract.
                    Date = BingDateParser.Format(requestDate),
                    Email = normalizedEmail,
                    Role = roleValue,
                    Site = normalizedSiteUrl,
                    VerificationSite = normalizedSiteUrl
                }
            },
            BingWebmasterJsonContext.Default.RemoveSiteRoleRequest,
            cancellationToken).ConfigureAwait(false);

        return new RemoveSiteRoleResponse(normalizedSiteUrl, normalizedEmail, DecodeRole(roleValue), true, DateTimeOffset.UtcNow);
    }

    internal async Task<GetBlockedUrlsResponse> GetBlockedUrlsAsync(
        string siteUrl,
        CancellationToken cancellationToken = default)
    {
        var normalizedSiteUrl = RequireText(siteUrl, nameof(siteUrl));
        var rawRows = await GetEnvelopeAsync(
            "GetBlockedUrls",
            new Dictionary<string, string?>
            {
                ["siteUrl"] = normalizedSiteUrl
            },
            BingWebmasterJsonContext.Default.ApiBlockedUrlArray,
            cancellationToken).ConfigureAwait(false) ?? [];

        var rows = rawRows.Select(row => new BlockedUrlEntry(
            row.Url,
            DecodeEntityType(row.EntityType),
            DecodeRequestType(row.RequestType),
            ParseRequiredDate(row.Date, nameof(row.Date)))).ToList();

        return new GetBlockedUrlsResponse(normalizedSiteUrl, rows.Count, rows, DateTimeOffset.UtcNow);
    }

    internal async Task<AddBlockedUrlResponse> AddBlockedUrlAsync(
        string siteUrl,
        string url,
        string entityType = "Page",
        string requestType = "CacheOnly",
        CancellationToken cancellationToken = default)
    {
        var normalizedSiteUrl = RequireText(siteUrl, nameof(siteUrl));
        var normalizedUrl = RequireText(url, nameof(url));
        var entityTypeValue = ParseEntityType(entityType, nameof(entityType));
        var requestTypeValue = ParseRequestType(requestType, nameof(requestType));
        var requestDate = DateTimeOffset.UtcNow;

        await PostCommandAsync(
            "AddBlockedUrl",
            new BlockedUrlRequest
            {
                SiteUrl = normalizedSiteUrl,
                BlockedUrl = new ApiBlockedUrl
                {
                    Date = BingDateParser.Format(requestDate),
                    EntityType = entityTypeValue,
                    RequestType = requestTypeValue,
                    Url = normalizedUrl
                }
            },
            BingWebmasterJsonContext.Default.BlockedUrlRequest,
            cancellationToken).ConfigureAwait(false);

        return new AddBlockedUrlResponse(
            normalizedSiteUrl,
            normalizedUrl,
            DecodeEntityType(entityTypeValue),
            DecodeRequestType(requestTypeValue),
            true,
            DateTimeOffset.UtcNow);
    }

    internal async Task<RemoveBlockedUrlResponse> RemoveBlockedUrlAsync(
        string siteUrl,
        string url,
        string entityType = "Page",
        string requestType = "FullRemoval",
        CancellationToken cancellationToken = default)
    {
        var normalizedSiteUrl = RequireText(siteUrl, nameof(siteUrl));
        var normalizedUrl = RequireText(url, nameof(url));
        var entityTypeValue = ParseEntityType(entityType, nameof(entityType));
        var requestTypeValue = ParseRequestType(requestType, nameof(requestType));
        var requestDate = DateTimeOffset.UtcNow;

        await PostCommandAsync(
            "RemoveBlockedUrl",
            new BlockedUrlRequest
            {
                SiteUrl = normalizedSiteUrl,
                BlockedUrl = new ApiBlockedUrl
                {
                    Date = BingDateParser.Format(requestDate),
                    EntityType = entityTypeValue,
                    RequestType = requestTypeValue,
                    Url = normalizedUrl
                }
            },
            BingWebmasterJsonContext.Default.BlockedUrlRequest,
            cancellationToken).ConfigureAwait(false);

        return new RemoveBlockedUrlResponse(
            normalizedSiteUrl,
            normalizedUrl,
            DecodeEntityType(entityTypeValue),
            DecodeRequestType(requestTypeValue),
            true,
            DateTimeOffset.UtcNow);
    }

    internal async Task<QueryPageDetailStatsResponse> GetQueryPageDetailStatsAsync(
        string siteUrl,
        string query,
        string page,
        CancellationToken cancellationToken = default)
    {
        var normalizedSiteUrl = RequireText(siteUrl, nameof(siteUrl));
        var normalizedQuery = RequireText(query, nameof(query));
        var normalizedPage = RequireText(page, nameof(page));
        var rawRows = await GetEnvelopeAsync(
            "GetQueryPageDetailStats",
            new Dictionary<string, string?>
            {
                ["siteUrl"] = normalizedSiteUrl,
                ["query"] = normalizedQuery,
                ["page"] = normalizedPage
            },
            BingWebmasterJsonContext.Default.ApiDetailedQueryStatArray,
            cancellationToken).ConfigureAwait(false) ?? [];

        var rows = rawRows.Select(row => new DetailedQueryStatEntry(
            ParseRequiredDate(row.Date, nameof(row.Date)),
            row.Clicks,
            row.Impressions,
            row.Position)).ToList();

        return new QueryPageDetailStatsResponse(
            normalizedSiteUrl,
            normalizedQuery,
            normalizedPage,
            rows.Count,
            rows,
            DateTimeOffset.UtcNow);
    }

    internal async Task<QueryTrafficStatsResponse> GetQueryTrafficStatsAsync(
        string siteUrl,
        string query,
        CancellationToken cancellationToken = default)
    {
        var normalizedSiteUrl = RequireText(siteUrl, nameof(siteUrl));
        var normalizedQuery = RequireText(query, nameof(query));
        var rawRows = await GetEnvelopeAsync(
            "GetQueryTrafficStats",
            new Dictionary<string, string?>
            {
                ["siteUrl"] = normalizedSiteUrl,
                ["query"] = normalizedQuery
            },
            BingWebmasterJsonContext.Default.ApiRankAndTrafficStatArray,
            cancellationToken).ConfigureAwait(false) ?? [];

        var rows = rawRows.Select(stat => new RankAndTrafficStatEntry(
            ParseRequiredDate(stat.Date, nameof(stat.Date)),
            stat.Clicks,
            stat.Impressions)).ToList();

        return new QueryTrafficStatsResponse(normalizedSiteUrl, normalizedQuery, rows.Count, rows, DateTimeOffset.UtcNow);
    }

    internal async Task<GetKeywordResponse> GetKeywordAsync(
        string query,
        string country,
        string language,
        string startDate,
        string endDate,
        CancellationToken cancellationToken = default)
    {
        var normalizedQuery = RequireText(query, nameof(query));
        var normalizedCountry = RequireText(country, nameof(country));
        var normalizedLanguage = RequireText(language, nameof(language));
        var normalizedStartDate = RequireIsoDate(startDate, nameof(startDate));
        var normalizedEndDate = RequireIsoDate(endDate, nameof(endDate));
        var rawKeyword = await GetEnvelopeAsync(
            "GetKeyword",
            new Dictionary<string, string?>
            {
                ["q"] = normalizedQuery,
                ["country"] = normalizedCountry,
                ["language"] = normalizedLanguage,
                ["startDate"] = normalizedStartDate,
                ["endDate"] = normalizedEndDate
            },
            BingWebmasterJsonContext.Default.ApiKeywordDetails,
            cancellationToken).ConfigureAwait(false);

        var found = !string.IsNullOrWhiteSpace(rawKeyword?.Query);
        return new GetKeywordResponse(
            normalizedQuery,
            normalizedCountry,
            normalizedLanguage,
            normalizedStartDate,
            normalizedEndDate,
            found,
            found ? rawKeyword!.Impressions : 0,
            found ? rawKeyword!.BroadImpressions : 0,
            DateTimeOffset.UtcNow);
    }

    internal async Task<RelatedKeywordsResponse> GetRelatedKeywordsAsync(
        string query,
        string country,
        string language,
        string startDate,
        string endDate,
        CancellationToken cancellationToken = default)
    {
        var normalizedQuery = RequireText(query, nameof(query));
        var normalizedCountry = RequireText(country, nameof(country));
        var normalizedLanguage = RequireText(language, nameof(language));
        var normalizedStartDate = RequireIsoDate(startDate, nameof(startDate));
        var normalizedEndDate = RequireIsoDate(endDate, nameof(endDate));
        var rawRows = await GetEnvelopeAsync(
            "GetRelatedKeywords",
            new Dictionary<string, string?>
            {
                ["q"] = normalizedQuery,
                ["country"] = normalizedCountry,
                ["language"] = normalizedLanguage,
                ["startDate"] = normalizedStartDate,
                ["endDate"] = normalizedEndDate
            },
            BingWebmasterJsonContext.Default.ApiKeywordDetailsArray,
            cancellationToken).ConfigureAwait(false) ?? [];

        var rows = rawRows
            .Where(row => !string.IsNullOrWhiteSpace(row.Query))
            .Select(row => new RelatedKeywordEntry(row.Query!, row.Impressions, row.BroadImpressions))
            .ToList();

        return new RelatedKeywordsResponse(
            normalizedQuery,
            normalizedCountry,
            normalizedLanguage,
            normalizedStartDate,
            normalizedEndDate,
            rows.Count,
            rows,
            DateTimeOffset.UtcNow);
    }

    internal async Task<ChildrenUrlInfoResponse> GetChildrenUrlInfoAsync(
        string siteUrl,
        string url,
        int page = 0,
        string crawlDateFilter = "Any",
        string discoveredDateFilter = "Any",
        string docFlagsFilter = "Any",
        string httpCodeFilter = "Any",
        CancellationToken cancellationToken = default)
    {
        var normalizedSiteUrl = RequireText(siteUrl, nameof(siteUrl));
        var normalizedUrl = RequireText(url, nameof(url));
        RequireNonNegative(page, nameof(page));
        var crawlDateFilterValue = ParseCrawlDateFilter(crawlDateFilter, nameof(crawlDateFilter));
        var discoveredDateFilterValue = ParseDiscoveredDateFilter(discoveredDateFilter, nameof(discoveredDateFilter));
        var docFlagsFilterValue = ParseDocFlagsFilter(docFlagsFilter, nameof(docFlagsFilter));
        var httpCodeFilterValue = ParseHttpCodeFilter(httpCodeFilter, nameof(httpCodeFilter));

        var rawRows = await PostEnvelopeAsync(
            "GetChildrenUrlInfo",
            new ChildrenUrlInfoRequest
            {
                SiteUrl = normalizedSiteUrl,
                Url = normalizedUrl,
                Page = page,
                FilterProperties = new ApiFilterProperties
                {
                    CrawlDateFilter = crawlDateFilterValue,
                    DiscoveredDateFilter = discoveredDateFilterValue,
                    DocFlagsFilters = docFlagsFilterValue,
                    HttpCodeFilters = httpCodeFilterValue
                }
            },
            BingWebmasterJsonContext.Default.ChildrenUrlInfoRequest,
            BingWebmasterJsonContext.Default.ApiUrlInfoArray,
            cancellationToken).ConfigureAwait(false) ?? [];

        var rows = rawRows.Select(ToChildUrlInfoEntry).ToList();
        return new ChildrenUrlInfoResponse(
            normalizedSiteUrl,
            normalizedUrl,
            page,
            DecodeCrawlDateFilter(crawlDateFilterValue),
            DecodeDiscoveredDateFilter(discoveredDateFilterValue),
            DecodeDocFlagsFilter(docFlagsFilterValue),
            DecodeHttpCodeFilter(httpCodeFilterValue),
            rows.Count,
            rows,
            DateTimeOffset.UtcNow);
    }

    internal async Task<ChildrenUrlTrafficInfoResponse> GetChildrenUrlTrafficInfoAsync(
        string siteUrl,
        string url,
        int page = 0,
        CancellationToken cancellationToken = default)
    {
        var normalizedSiteUrl = RequireText(siteUrl, nameof(siteUrl));
        var normalizedUrl = RequireText(url, nameof(url));
        RequireNonNegative(page, nameof(page));

        var rawRows = await GetEnvelopeAsync(
            "GetChildrenUrlTrafficInfo",
            new Dictionary<string, string?>
            {
                ["siteUrl"] = normalizedSiteUrl,
                ["url"] = normalizedUrl,
                ["page"] = page.ToString(CultureInfo.InvariantCulture)
            },
            BingWebmasterJsonContext.Default.ApiUrlTrafficInfoArray,
            cancellationToken).ConfigureAwait(false) ?? [];

        var rows = rawRows.Select(ToChildUrlTrafficInfoEntry).ToList();
        return new ChildrenUrlTrafficInfoResponse(normalizedSiteUrl, normalizedUrl, page, rows.Count, rows, DateTimeOffset.UtcNow);
    }

    internal async Task<FetchUrlResponse> FetchUrlAsync(
        string siteUrl,
        string url,
        CancellationToken cancellationToken = default)
    {
        var normalizedSiteUrl = RequireText(siteUrl, nameof(siteUrl));
        var normalizedUrl = RequireText(url, nameof(url));

        await PostCommandAsync(
            "FetchUrl",
            new SiteAndUrlRequest
            {
                SiteUrl = normalizedSiteUrl,
                Url = normalizedUrl
            },
            BingWebmasterJsonContext.Default.SiteAndUrlRequest,
            cancellationToken).ConfigureAwait(false);

        return new FetchUrlResponse(normalizedSiteUrl, normalizedUrl, true, DateTimeOffset.UtcNow);
    }

    internal async Task<ListFetchedUrlsResponse> ListFetchedUrlsAsync(
        string siteUrl,
        CancellationToken cancellationToken = default)
    {
        var normalizedSiteUrl = RequireText(siteUrl, nameof(siteUrl));
        var rawRows = await GetEnvelopeAsync(
            "GetFetchedUrls",
            new Dictionary<string, string?>
            {
                ["siteUrl"] = normalizedSiteUrl
            },
            BingWebmasterJsonContext.Default.ApiFetchedUrlArray,
            cancellationToken).ConfigureAwait(false) ?? [];

        var rows = rawRows.Select(row => new FetchedUrlEntry(
            row.Url,
            ParseRequiredDate(row.Date, nameof(row.Date)),
            row.Fetched,
            row.Expired)).ToList();

        return new ListFetchedUrlsResponse(normalizedSiteUrl, rows.Count, rows, DateTimeOffset.UtcNow);
    }

    internal async Task<FetchedUrlDetailsResponse> GetFetchedUrlDetailsAsync(
        string siteUrl,
        string url,
        CancellationToken cancellationToken = default)
    {
        var normalizedSiteUrl = RequireText(siteUrl, nameof(siteUrl));
        var normalizedUrl = RequireText(url, nameof(url));
        var rawDetails = await GetEnvelopeAsync(
            "GetFetchedUrlDetails",
            new Dictionary<string, string?>
            {
                ["siteUrl"] = normalizedSiteUrl,
                ["url"] = normalizedUrl
            },
            BingWebmasterJsonContext.Default.ApiFetchedUrlDetails,
            cancellationToken).ConfigureAwait(false)
            ?? throw new InvalidOperationException("Bing returned an empty fetched URL details payload.");

        return new FetchedUrlDetailsResponse(
            normalizedSiteUrl,
            rawDetails.Url,
            ParseRequiredDate(rawDetails.Date, nameof(rawDetails.Date)),
            rawDetails.Status,
            rawDetails.Headers,
            rawDetails.Document,
            DateTimeOffset.UtcNow);
    }

    internal async Task<RemoveSitemapResponse> RemoveSitemapAsync(
        string siteUrl,
        string feedUrl,
        CancellationToken cancellationToken = default)
    {
        var normalizedSiteUrl = RequireText(siteUrl, nameof(siteUrl));
        var normalizedFeedUrl = RequireText(feedUrl, nameof(feedUrl));
        await PostCommandAsync(
            "RemoveFeed",
            new SiteAndFeedRequest
            {
                SiteUrl = normalizedSiteUrl,
                FeedUrl = normalizedFeedUrl
            },
            BingWebmasterJsonContext.Default.SiteAndFeedRequest,
            cancellationToken).ConfigureAwait(false);

        return new RemoveSitemapResponse(normalizedSiteUrl, normalizedFeedUrl, true, DateTimeOffset.UtcNow);
    }

    internal async Task<GetSiteMovesResponse> GetSiteMovesAsync(
        string siteUrl,
        CancellationToken cancellationToken = default)
    {
        var normalizedSiteUrl = RequireText(siteUrl, nameof(siteUrl));
        var rawRows = await GetEnvelopeAsync(
            "GetSiteMoves",
            new Dictionary<string, string?>
            {
                ["siteUrl"] = normalizedSiteUrl
            },
            BingWebmasterJsonContext.Default.ApiSiteMoveSettingsArray,
            cancellationToken).ConfigureAwait(false) ?? [];

        var rows = rawRows.Select(row => new SiteMoveEntry(
            ParseRequiredDate(row.Date, nameof(row.Date)),
            DecodeMoveScope(row.MoveScope),
            DecodeMoveType(row.MoveType),
            row.SourceUrl,
            row.TargetUrl)).ToList();

        return new GetSiteMovesResponse(normalizedSiteUrl, rows.Count, rows, DateTimeOffset.UtcNow);
    }

    internal async Task<SubmitSiteMoveResponse> SubmitSiteMoveAsync(
        string siteUrl,
        string sourceUrl,
        string targetUrl,
        string moveType = "Local",
        string moveScope = "Domain",
        CancellationToken cancellationToken = default)
    {
        var normalizedSiteUrl = RequireText(siteUrl, nameof(siteUrl));
        var normalizedSourceUrl = RequireText(sourceUrl, nameof(sourceUrl));
        var normalizedTargetUrl = RequireText(targetUrl, nameof(targetUrl));
        var moveTypeValue = ParseMoveType(moveType, nameof(moveType));
        var moveScopeValue = ParseMoveScope(moveScope, nameof(moveScope));
        var requestDate = DateTimeOffset.UtcNow;

        await PostCommandAsync(
            "SubmitSiteMove",
            new SubmitSiteMoveRequest
            {
                SiteUrl = normalizedSiteUrl,
                Settings = new ApiSiteMoveSettings
                {
                    Date = BingDateParser.Format(requestDate),
                    MoveScope = moveScopeValue,
                    MoveType = moveTypeValue,
                    SourceUrl = normalizedSourceUrl,
                    TargetUrl = normalizedTargetUrl
                }
            },
            BingWebmasterJsonContext.Default.SubmitSiteMoveRequest,
            cancellationToken).ConfigureAwait(false);

        return new SubmitSiteMoveResponse(
            normalizedSiteUrl,
            normalizedSourceUrl,
            normalizedTargetUrl,
            DecodeMoveType(moveTypeValue),
            DecodeMoveScope(moveScopeValue),
            true,
            DateTimeOffset.UtcNow);
    }

    internal async Task<SubmitContentResponse> SubmitContentAsync(
        string siteUrl,
        string url,
        string httpMessage,
        string structuredData,
        string dynamicServing = "None",
        CancellationToken cancellationToken = default)
    {
        var normalizedSiteUrl = RequireText(siteUrl, nameof(siteUrl));
        var normalizedUrl = RequireText(url, nameof(url));
        var normalizedHttpMessage = RequireText(httpMessage, nameof(httpMessage));
        var normalizedStructuredData = RequireText(structuredData, nameof(structuredData));
        var dynamicServingValue = ParseDynamicServing(dynamicServing, nameof(dynamicServing));

        await PostCommandAsync(
            "SubmitContent",
            new SubmitContentRequest
            {
                SiteUrl = normalizedSiteUrl,
                Url = normalizedUrl,
                HttpMessage = normalizedHttpMessage,
                StructuredData = normalizedStructuredData,
                DynamicServing = dynamicServingValue
            },
            BingWebmasterJsonContext.Default.SubmitContentRequest,
            cancellationToken).ConfigureAwait(false);

        return new SubmitContentResponse(
            normalizedSiteUrl,
            normalizedUrl,
            DecodeDynamicServing(dynamicServingValue),
            true,
            DateTimeOffset.UtcNow);
    }

    internal async Task<ContentSubmissionQuotaResponse> GetContentSubmissionQuotaAsync(
        string siteUrl,
        CancellationToken cancellationToken = default)
    {
        var normalizedSiteUrl = RequireText(siteUrl, nameof(siteUrl));
        var rawQuota = await GetEnvelopeAsync(
            "GetContentSubmissionQuota",
            new Dictionary<string, string?>
            {
                ["siteUrl"] = normalizedSiteUrl
            },
            BingWebmasterJsonContext.Default.ApiUrlSubmissionQuota,
            cancellationToken).ConfigureAwait(false)
            ?? throw new InvalidOperationException("Bing returned an empty content submission quota payload.");

        return new ContentSubmissionQuotaResponse(
            normalizedSiteUrl,
            rawQuota.DailyQuota,
            rawQuota.MonthlyQuota,
            DateTimeOffset.UtcNow);
    }

    private async Task<List<QueryStatsEntry>> GetQueryStatsRowsAsync(
        string methodName,
        IReadOnlyDictionary<string, string?> queryParameters,
        CancellationToken cancellationToken)
    {
        var rawRows = await GetEnvelopeAsync(
            methodName,
            queryParameters,
            BingWebmasterJsonContext.Default.ApiQueryStatArray,
            cancellationToken).ConfigureAwait(false) ?? [];

        return rawRows.Select(row => new QueryStatsEntry(
            row.Query,
            ParseRequiredDate(row.Date, nameof(row.Date)),
            row.Clicks,
            row.Impressions,
            row.AvgClickPosition,
            row.AvgImpressionPosition)).ToList();
    }

    private async Task<List<PageStatsEntry>> GetPageStatsRowsAsync(
        string methodName,
        IReadOnlyDictionary<string, string?> queryParameters,
        CancellationToken cancellationToken)
    {
        var rawRows = await GetEnvelopeAsync(
            methodName,
            queryParameters,
            BingWebmasterJsonContext.Default.ApiQueryStatArray,
            cancellationToken).ConfigureAwait(false) ?? [];

        return rawRows.Select(row => new PageStatsEntry(
            row.Query,
            ParseRequiredDate(row.Date, nameof(row.Date)),
            row.Clicks,
            row.Impressions,
            row.AvgClickPosition,
            row.AvgImpressionPosition)).ToList();
    }

    private async Task<T?> GetEnvelopeAsync<T>(
        string methodName,
        IReadOnlyDictionary<string, string?>? queryParameters,
        JsonTypeInfo<T> typeInfo,
        CancellationToken cancellationToken)
    {
        var request = new HttpRequestMessage(HttpMethod.Get, BuildUrl(methodName, queryParameters));
        using var response = await _httpClient.SendAsync(request, cancellationToken).ConfigureAwait(false);
        await EnsureSuccessAsync(response, cancellationToken).ConfigureAwait(false);
        return await DeserializeEnvelopeAsync(response, typeInfo, cancellationToken).ConfigureAwait(false);
    }

    private async Task<bool> PostQueryAsync<TBody>(
        string methodName,
        TBody body,
        JsonTypeInfo<TBody> typeInfo,
        CancellationToken cancellationToken)
    {
        var jsonBody = JsonSerializer.Serialize(body, typeInfo);
        var request = new HttpRequestMessage(HttpMethod.Post, BuildUrl(methodName, null))
        {
            Content = new StringContent(jsonBody, Encoding.UTF8, "application/json")
        };

        using var response = await _httpClient.SendAsync(request, cancellationToken).ConfigureAwait(false);
        await EnsureSuccessAsync(response, cancellationToken).ConfigureAwait(false);
        return await DeserializeEnvelopeAsync(response, BingWebmasterJsonContext.Default.Boolean, cancellationToken).ConfigureAwait(false);
    }

    private async Task<TResponse?> PostEnvelopeAsync<TBody, TResponse>(
        string methodName,
        TBody body,
        JsonTypeInfo<TBody> bodyTypeInfo,
        JsonTypeInfo<TResponse> responseTypeInfo,
        CancellationToken cancellationToken)
    {
        var jsonBody = JsonSerializer.Serialize(body, bodyTypeInfo);
        var request = new HttpRequestMessage(HttpMethod.Post, BuildUrl(methodName, null))
        {
            Content = new StringContent(jsonBody, Encoding.UTF8, "application/json")
        };

        using var response = await _httpClient.SendAsync(request, cancellationToken).ConfigureAwait(false);
        await EnsureSuccessAsync(response, cancellationToken).ConfigureAwait(false);
        return await DeserializeEnvelopeAsync(response, responseTypeInfo, cancellationToken).ConfigureAwait(false);
    }

    /// <summary>
    /// Issues a POST for a fire-and-forget command whose "d" payload is not a reliable success
    /// indicator (confirmed empirically: Bing's AddSite endpoint returns "d":null both for a
    /// brand-new site and a no-op repeat of an already-added site -- there is no real boolean to
    /// read). The payload is read and discarded; success means the HTTP call completed without
    /// the client throwing.
    /// </summary>
    private async Task PostCommandAsync<TBody>(
        string methodName,
        TBody body,
        JsonTypeInfo<TBody> typeInfo,
        CancellationToken cancellationToken)
    {
        var jsonBody = JsonSerializer.Serialize(body, typeInfo);
        var request = new HttpRequestMessage(HttpMethod.Post, BuildUrl(methodName, null))
        {
            Content = new StringContent(jsonBody, Encoding.UTF8, "application/json")
        };

        using var response = await _httpClient.SendAsync(request, cancellationToken).ConfigureAwait(false);
        await EnsureSuccessAsync(response, cancellationToken).ConfigureAwait(false);
        using var document = await ReadEnvelopeDocumentAsync(response, cancellationToken).ConfigureAwait(false);
        // Intentionally discard document.RootElement.GetProperty("d") -- it carries no reliable
        // signal for this class of command.
    }

    private async Task<T?> DeserializeEnvelopeAsync<T>(
        HttpResponseMessage response,
        JsonTypeInfo<T> typeInfo,
        CancellationToken cancellationToken)
    {
        using var document = await ReadEnvelopeDocumentAsync(response, cancellationToken).ConfigureAwait(false);
        return document.RootElement.GetProperty("d").Deserialize(typeInfo);
    }

    private static async Task<JsonDocument> ReadEnvelopeDocumentAsync(
        HttpResponseMessage response,
        CancellationToken cancellationToken)
    {
        await using var stream = await response.Content.ReadAsStreamAsync(cancellationToken).ConfigureAwait(false);
        return await JsonDocument.ParseAsync(stream, cancellationToken: cancellationToken).ConfigureAwait(false);
    }

    private string BuildUrl(string methodName, IReadOnlyDictionary<string, string?>? queryParameters)
    {
        var builder = new StringBuilder($"{_baseUrl}/{methodName}?apikey={Uri.EscapeDataString(_apiKey)}");
        if (queryParameters is not null)
        {
            foreach (var pair in queryParameters)
            {
                if (string.IsNullOrWhiteSpace(pair.Value))
                    continue;

                builder.Append('&')
                    .Append(Uri.EscapeDataString(pair.Key))
                    .Append('=')
                    .Append(Uri.EscapeDataString(pair.Value));
            }
        }

        return builder.ToString();
    }

    private static async Task EnsureSuccessAsync(HttpResponseMessage response, CancellationToken cancellationToken)
    {
        if (response.IsSuccessStatusCode)
            return;

        var body = await response.Content.ReadAsStringAsync(cancellationToken).ConfigureAwait(false);
        var snippet = body.Length > 300 ? body[..300] + "..." : body;
        throw new BingWebmasterApiException(
            (int)response.StatusCode,
            $"Bing Webmaster API returned HTTP {(int)response.StatusCode} {response.StatusCode}: {snippet}");
    }

    private static SitemapInfo ToSitemapInfo(ApiFeed feed)
        => new(
            feed.Url,
            feed.Type,
            feed.Compressed,
            feed.FileSize,
            ParseOptionalDate(feed.LastCrawled),
            ParseOptionalDate(feed.Submitted),
            feed.Status,
            feed.UrlCount);

    private static ChildUrlInfoEntry ToChildUrlInfoEntry(ApiUrlInfo info)
        => new(
            info.Url,
            info.IsPage,
            info.HttpStatus,
            info.DocumentSize,
            info.AnchorCount,
            ParseOptionalDate(info.DiscoveryDate),
            ParseOptionalDate(info.LastCrawledDate),
            info.TotalChildUrlCount);

    private static ChildUrlTrafficInfoEntry ToChildUrlTrafficInfoEntry(ApiUrlTrafficInfo info)
        => new(
            info.Url,
            info.IsPage,
            info.Clicks,
            info.Impressions);

    private static DateTimeOffset ParseRequiredDate(string value, string fieldName)
        => BingDateParser.TryParse(value, out var parsed)
            ? parsed
            : throw new FormatException($"Bing date field '{fieldName}' was not in the expected /Date(...)/ format.");

    private static DateTimeOffset? ParseOptionalDate(string? value)
        => string.IsNullOrWhiteSpace(value)
            ? null
            : ParseRequiredDate(value, "date");

    private static string RequireText(string value, string paramName)
        => string.IsNullOrWhiteSpace(value)
            ? throw new ArgumentException("Value is required.", paramName)
            : value.Trim();

    private static string RequireIsoDate(string value, string paramName)
    {
        var normalized = RequireText(value, paramName);
        return DateOnly.ParseExact(normalized, "yyyy-MM-dd", CultureInfo.InvariantCulture)
            .ToString("yyyy-MM-dd", CultureInfo.InvariantCulture);
    }

    private static void RequireNonNegative(int value, string paramName)
    {
        if (value < 0)
            throw new ArgumentOutOfRangeException(paramName, "Value must be zero or greater.");
    }

    private static IReadOnlyList<string> NormalizeUrlList(IReadOnlyList<string> urlList, int maxCount)
    {
        ArgumentNullException.ThrowIfNull(urlList);

        var normalized = urlList
            .Where(url => !string.IsNullOrWhiteSpace(url))
            .Select(url => url.Trim())
            .ToList();

        if (normalized.Count == 0)
            throw new ArgumentException("At least one URL is required.", nameof(urlList));

        if (normalized.Count > maxCount)
            throw new ArgumentException($"A maximum of {maxCount} URLs is allowed.", nameof(urlList));

        return normalized;
    }

    private static IReadOnlyList<string> DecodeCrawlIssues(int issues)
    {
        if (issues == 0)
            return [];

        var decoded = new List<string>();
        AddIssueFlag(decoded, issues, 1, "Code301");
        AddIssueFlag(decoded, issues, 2, "Code302");
        AddIssueFlag(decoded, issues, 4, "Code4xx");
        AddIssueFlag(decoded, issues, 8, "Code5xx");
        AddIssueFlag(decoded, issues, 16, "BlockedByRobotsTxt");
        AddIssueFlag(decoded, issues, 32, "ContainsMalware");
        AddIssueFlag(decoded, issues, 64, "ImportantUrlBlockedByRobotsTxt");
        AddIssueFlag(decoded, issues, 128, "DnsErrors");
        AddIssueFlag(decoded, issues, 256, "TimeOutErrors");
        return decoded;
    }

    private static void AddIssueFlag(List<string> decoded, int issues, int flag, string name)
    {
        if ((issues & flag) == flag)
            decoded.Add(name);
    }

    private static string ToBooleanString(bool value)
        => value ? "true" : "false";

    private static int ParseRole(string value, string paramName)
        => ParseNamedValue(
            value,
            paramName,
            [
                ("Administrator", 0),
                ("ReadOnly", 1),
                ("ReadWrite", 2)
            ]);

    private static string DecodeRole(int value)
        => value switch
        {
            0 => "Administrator",
            1 => "ReadOnly",
            2 => "ReadWrite",
            _ => $"Unknown({value})"
        };

    private static int ParseEntityType(string value, string paramName)
        => ParseNamedValue(
            value,
            paramName,
            [
                ("Page", 0),
                ("Directory", 1)
            ]);

    private static string DecodeEntityType(int value)
        => value switch
        {
            0 => "Page",
            1 => "Directory",
            _ => $"Unknown({value})"
        };

    private static int ParseRequestType(string value, string paramName)
        => ParseNamedValue(
            value,
            paramName,
            [
                ("CacheOnly", 0),
                ("FullRemoval", 1)
            ]);

    private static string DecodeRequestType(int value)
        => value switch
        {
            0 => "CacheOnly",
            1 => "FullRemoval",
            _ => $"Unknown({value})"
        };

    private static int ParseMoveScope(string value, string paramName)
        => ParseNamedValue(
            value,
            paramName,
            [
                ("Domain", 0),
                ("Host", 1),
                ("Directory", 2)
            ]);

    private static string DecodeMoveScope(int value)
        => value switch
        {
            0 => "Domain",
            1 => "Host",
            2 => "Directory",
            _ => $"Unknown({value})"
        };

    private static int ParseMoveType(string value, string paramName)
        => ParseNamedValue(
            value,
            paramName,
            [
                ("Local", 0),
                ("Global", 1)
            ]);

    private static string DecodeMoveType(int value)
        => value switch
        {
            0 => "Local",
            1 => "Global",
            _ => $"Unknown({value})"
        };

    private static int ParseCrawlDateFilter(string value, string paramName)
        => ParseNamedValue(
            value,
            paramName,
            [
                ("Any", 0),
                ("LastWeek", 1),
                ("LastTwoWeeks", 2),
                ("LastThreeWeeks", 4)
            ]);

    private static string DecodeCrawlDateFilter(int value)
        => value switch
        {
            0 => "Any",
            1 => "LastWeek",
            2 => "LastTwoWeeks",
            4 => "LastThreeWeeks",
            _ => $"Unknown({value})"
        };

    private static int ParseDiscoveredDateFilter(string value, string paramName)
        => ParseNamedValue(
            value,
            paramName,
            [
                ("Any", 0),
                ("LastWeek", 1),
                ("LastMonth", 2)
            ]);

    private static string DecodeDiscoveredDateFilter(int value)
        => value switch
        {
            0 => "Any",
            1 => "LastWeek",
            2 => "LastMonth",
            _ => $"Unknown({value})"
        };

    private static int ParseDocFlagsFilter(string value, string paramName)
        => ParseNamedValue(
            value,
            paramName,
            [
                ("Any", 0),
                ("IsBlockedByRobotsTxt", 1),
                ("IsMalware", 2)
            ]);

    private static string DecodeDocFlagsFilter(int value)
        => value switch
        {
            0 => "Any",
            1 => "IsBlockedByRobotsTxt",
            2 => "IsMalware",
            _ => $"Unknown({value})"
        };

    private static int ParseHttpCodeFilter(string value, string paramName)
        => ParseNamedValue(
            value,
            paramName,
            [
                ("Any", 0),
                ("Code2xx", 1),
                ("Code3xx", 2),
                ("Code301", 4),
                ("Code302", 8),
                ("Code4xx", 16),
                ("Code5xx", 32),
                ("AllOthers", 64)
            ]);

    private static string DecodeHttpCodeFilter(int value)
        => value switch
        {
            0 => "Any",
            1 => "Code2xx",
            2 => "Code3xx",
            4 => "Code301",
            8 => "Code302",
            16 => "Code4xx",
            32 => "Code5xx",
            64 => "AllOthers",
            _ => $"Unknown({value})"
        };

    private static int ParseDynamicServing(string value, string paramName)
        => ParseNamedValue(
            value,
            paramName,
            [
                ("None", 0),
                ("PcLaptop", 1),
                ("Mobile", 2),
                ("Amp", 3),
                ("Tablet", 4),
                ("NonVisualBrowser", 5)
            ]);

    private static string DecodeDynamicServing(int value)
        => value switch
        {
            0 => "None",
            1 => "PcLaptop",
            2 => "Mobile",
            3 => "Amp",
            4 => "Tablet",
            5 => "NonVisualBrowser",
            _ => $"Unknown({value})"
        };

    private static int ParseNamedValue(
        string value,
        string paramName,
        ReadOnlySpan<(string Name, int Value)> supportedValues)
    {
        var normalized = RequireText(value, paramName);
        foreach (var supportedValue in supportedValues)
        {
            if (string.Equals(normalized, supportedValue.Name, StringComparison.OrdinalIgnoreCase))
                return supportedValue.Value;
        }

        throw new ArgumentException(
            $"Unsupported value '{value}'. Expected one of: {string.Join(", ", supportedValues.ToArray().Select(v => v.Name))}.",
            paramName);
    }
}
