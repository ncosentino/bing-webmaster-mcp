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
        var success = await PostEnvelopeAsync(
            "AddSite",
            new SiteUrlRequest { SiteUrl = normalizedSiteUrl },
            BingWebmasterJsonContext.Default.SiteUrlRequest,
            cancellationToken).ConfigureAwait(false);

        return new AddSiteResponse(normalizedSiteUrl, success, DateTimeOffset.UtcNow);
    }

    internal async Task<VerifySiteResponse> VerifySiteAsync(string siteUrl, CancellationToken cancellationToken = default)
    {
        var normalizedSiteUrl = RequireText(siteUrl, nameof(siteUrl));
        var success = await PostEnvelopeAsync(
            "VerifySite",
            new SiteUrlRequest { SiteUrl = normalizedSiteUrl },
            BingWebmasterJsonContext.Default.SiteUrlRequest,
            cancellationToken).ConfigureAwait(false);

        return new VerifySiteResponse(normalizedSiteUrl, success, DateTimeOffset.UtcNow);
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
        var success = await PostEnvelopeAsync(
            "SubmitFeed",
            new SiteAndFeedRequest
            {
                SiteUrl = normalizedSiteUrl,
                FeedUrl = normalizedFeedUrl
            },
            BingWebmasterJsonContext.Default.SiteAndFeedRequest,
            cancellationToken).ConfigureAwait(false);

        return new SubmitSitemapResponse(normalizedSiteUrl, normalizedFeedUrl, success, DateTimeOffset.UtcNow);
    }

    internal async Task<SubmitUrlResponse> SubmitUrlAsync(
        string siteUrl,
        string url,
        CancellationToken cancellationToken = default)
    {
        var normalizedSiteUrl = RequireText(siteUrl, nameof(siteUrl));
        var normalizedUrl = RequireText(url, nameof(url));
        var success = await PostEnvelopeAsync(
            "SubmitUrl",
            new SiteAndUrlRequest
            {
                SiteUrl = normalizedSiteUrl,
                Url = normalizedUrl
            },
            BingWebmasterJsonContext.Default.SiteAndUrlRequest,
            cancellationToken).ConfigureAwait(false);

        return new SubmitUrlResponse(normalizedSiteUrl, normalizedUrl, success, DateTimeOffset.UtcNow);
    }

    internal async Task<SubmitUrlBatchResponse> SubmitUrlBatchAsync(
        string siteUrl,
        IReadOnlyList<string> urlList,
        CancellationToken cancellationToken = default)
    {
        var normalizedSiteUrl = RequireText(siteUrl, nameof(siteUrl));
        var normalizedUrls = NormalizeUrlList(urlList, 500);
        var success = await PostEnvelopeAsync(
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
            success,
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

    private async Task<bool> PostEnvelopeAsync<TBody>(
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
}
