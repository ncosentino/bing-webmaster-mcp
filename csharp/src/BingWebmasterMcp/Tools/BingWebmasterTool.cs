using System.ComponentModel;
using System.Text.Json;
using System.Text.Json.Serialization.Metadata;
using BingWebmasterMcp.BingWebmaster;
using BingWebmasterMcp.IndexNow;
using ModelContextProtocol.Server;

namespace BingWebmasterMcp.Tools;

/// <summary>MCP tools for Bing Webmaster Tools and IndexNow.</summary>
[McpServerToolType]
internal sealed class BingWebmasterTool(BingWebmasterClient client, IndexNowClient indexNowClient)
{
    [McpServerTool(Name = "list_sites")]
    [Description("List all Bing Webmaster Tools sites accessible to the configured API key.")]
    internal Task<string> ListSites(CancellationToken cancellationToken = default)
        => ExecuteAsync(ct => client.ListSitesAsync(ct), BingWebmasterJsonContext.Default.ListSitesResponse, cancellationToken);

    [McpServerTool(Name = "add_site")]
    [Description("Add a site to Bing Webmaster Tools.")]
    internal Task<string> AddSite(
        [Description("The site URL to add.")] string site_url,
        CancellationToken cancellationToken = default)
        => ExecuteAsync(ct => client.AddSiteAsync(site_url, ct), BingWebmasterJsonContext.Default.AddSiteResponse, cancellationToken);

    [McpServerTool(Name = "verify_site")]
    [Description("Verify a site in Bing Webmaster Tools.")]
    internal Task<string> VerifySite(
        [Description("The site URL to verify.")] string site_url,
        CancellationToken cancellationToken = default)
        => ExecuteAsync(ct => client.VerifySiteAsync(site_url, ct), BingWebmasterJsonContext.Default.VerifySiteResponse, cancellationToken);

    [McpServerTool(Name = "list_sitemaps")]
    [Description("List sitemaps submitted for a Bing Webmaster Tools site.")]
    internal Task<string> ListSitemaps(
        [Description("The Bing site URL.")] string site_url,
        CancellationToken cancellationToken = default)
        => ExecuteAsync(ct => client.ListSitemapsAsync(site_url, ct), BingWebmasterJsonContext.Default.ListSitemapsResponse, cancellationToken);

    [McpServerTool(Name = "get_sitemap_details")]
    [Description("Get details for a specific submitted sitemap.")]
    internal Task<string> GetSitemapDetails(
        [Description("The Bing site URL.")] string site_url,
        [Description("The sitemap URL to inspect.")] string feed_url,
        CancellationToken cancellationToken = default)
        => ExecuteAsync(
            ct => client.GetSitemapDetailsAsync(site_url, feed_url, ct),
            BingWebmasterJsonContext.Default.GetSitemapDetailsResponse,
            cancellationToken);

    [McpServerTool(Name = "submit_sitemap")]
    [Description("Submit a sitemap to Bing Webmaster Tools.")]
    internal Task<string> SubmitSitemap(
        [Description("The Bing site URL.")] string site_url,
        [Description("The sitemap URL to submit.")] string feed_url,
        CancellationToken cancellationToken = default)
        => ExecuteAsync(
            ct => client.SubmitSitemapAsync(site_url, feed_url, ct),
            BingWebmasterJsonContext.Default.SubmitSitemapResponse,
            cancellationToken);

    [McpServerTool(Name = "submit_url")]
    [Description("Submit a single URL to Bing Webmaster Tools.")]
    internal Task<string> SubmitUrl(
        [Description("The Bing site URL.")] string site_url,
        [Description("The page URL to submit.")] string url,
        CancellationToken cancellationToken = default)
        => ExecuteAsync(
            ct => client.SubmitUrlAsync(site_url, url, ct),
            BingWebmasterJsonContext.Default.SubmitUrlResponse,
            cancellationToken);

    [McpServerTool(Name = "submit_url_batch")]
    [Description("Submit up to 500 URLs to Bing Webmaster Tools in one request.")]
    internal Task<string> SubmitUrlBatch(
        [Description("The Bing site URL.")] string site_url,
        [Description("The URLs to submit. Maximum 500.")] string[] url_list,
        CancellationToken cancellationToken = default)
        => ExecuteAsync(
            ct => client.SubmitUrlBatchAsync(site_url, url_list, ct),
            BingWebmasterJsonContext.Default.SubmitUrlBatchResponse,
            cancellationToken);

    [McpServerTool(Name = "submit_url_indexnow")]
    [Description("Submit one or more URLs to Bing via the IndexNow batch protocol.")]
    internal Task<string> SubmitUrlIndexNow(
        [Description("The host that owns all submitted URLs, for example example.com.")] string host,
        [Description("The URLs to submit via IndexNow.")] string[] url_list,
        [Description("Optional IndexNow key override. If omitted, the configured BING_INDEXNOW_KEY is used.")] string? key = null,
        [Description("Optional absolute URL for the hosted IndexNow key file.")] string? key_location = null,
        CancellationToken cancellationToken = default)
        => ExecuteAsync(
            ct => indexNowClient.SubmitUrlIndexNowAsync(host, url_list, key, key_location, ct),
            BingWebmasterJsonContext.Default.SubmitUrlIndexNowResponse,
            cancellationToken);

    [McpServerTool(Name = "get_url_submission_quota")]
    [Description("Get the current Bing URL submission quota for a site.")]
    internal Task<string> GetUrlSubmissionQuota(
        [Description("The Bing site URL.")] string site_url,
        CancellationToken cancellationToken = default)
        => ExecuteAsync(
            ct => client.GetUrlSubmissionQuotaAsync(site_url, ct),
            BingWebmasterJsonContext.Default.UrlSubmissionQuotaResponse,
            cancellationToken);

    [McpServerTool(Name = "get_crawl_issues")]
    [Description("List current Bing crawl issues for a site.")]
    internal Task<string> GetCrawlIssues(
        [Description("The Bing site URL.")] string site_url,
        CancellationToken cancellationToken = default)
        => ExecuteAsync(
            ct => client.GetCrawlIssuesAsync(site_url, ct),
            BingWebmasterJsonContext.Default.CrawlIssuesResponse,
            cancellationToken);

    [McpServerTool(Name = "get_crawl_stats")]
    [Description("Get Bing crawl statistics over time for a site.")]
    internal Task<string> GetCrawlStats(
        [Description("The Bing site URL.")] string site_url,
        CancellationToken cancellationToken = default)
        => ExecuteAsync(
            ct => client.GetCrawlStatsAsync(site_url, ct),
            BingWebmasterJsonContext.Default.CrawlStatsResponse,
            cancellationToken);

    [McpServerTool(Name = "get_url_info")]
    [Description("Get Bing index metadata for a specific URL.")]
    internal Task<string> GetUrlInfo(
        [Description("The Bing site URL.")] string site_url,
        [Description("The page URL to inspect.")] string url,
        CancellationToken cancellationToken = default)
        => ExecuteAsync(
            ct => client.GetUrlInfoAsync(site_url, url, ct),
            BingWebmasterJsonContext.Default.UrlInfoResponse,
            cancellationToken);

    [McpServerTool(Name = "get_url_traffic_info")]
    [Description("Get Bing clicks and impressions for a specific URL.")]
    internal Task<string> GetUrlTrafficInfo(
        [Description("The Bing site URL.")] string site_url,
        [Description("The page URL to inspect.")] string url,
        CancellationToken cancellationToken = default)
        => ExecuteAsync(
            ct => client.GetUrlTrafficInfoAsync(site_url, url, ct),
            BingWebmasterJsonContext.Default.UrlTrafficInfoResponse,
            cancellationToken);

    [McpServerTool(Name = "get_url_links")]
    [Description("Get inbound link details for a specific URL.")]
    internal Task<string> GetUrlLinks(
        [Description("The Bing site URL.")] string site_url,
        [Description("The URL whose inbound links should be listed.")] string link,
        [Description("Zero-based result page.")] int page = 0,
        CancellationToken cancellationToken = default)
        => ExecuteAsync(
            ct => client.GetUrlLinksAsync(site_url, link, page, ct),
            BingWebmasterJsonContext.Default.UrlLinksResponse,
            cancellationToken);

    [McpServerTool(Name = "get_link_counts")]
    [Description("Get paged inbound link counts for pages on a site.")]
    internal Task<string> GetLinkCounts(
        [Description("The Bing site URL.")] string site_url,
        [Description("Zero-based result page.")] int page = 0,
        CancellationToken cancellationToken = default)
        => ExecuteAsync(
            ct => client.GetLinkCountsAsync(site_url, page, ct),
            BingWebmasterJsonContext.Default.LinkCountsResponse,
            cancellationToken);

    [McpServerTool(Name = "get_rank_and_traffic_stats")]
    [Description("Get Bing clicks and impressions over time for a site.")]
    internal Task<string> GetRankAndTrafficStats(
        [Description("The Bing site URL.")] string site_url,
        CancellationToken cancellationToken = default)
        => ExecuteAsync(
            ct => client.GetRankAndTrafficStatsAsync(site_url, ct),
            BingWebmasterJsonContext.Default.RankAndTrafficStatsResponse,
            cancellationToken);

    [McpServerTool(Name = "get_query_stats")]
    [Description("Get Bing search query statistics for a site.")]
    internal Task<string> GetQueryStats(
        [Description("The Bing site URL.")] string site_url,
        CancellationToken cancellationToken = default)
        => ExecuteAsync(
            ct => client.GetQueryStatsAsync(site_url, ct),
            BingWebmasterJsonContext.Default.QueryStatsResponse,
            cancellationToken);

    [McpServerTool(Name = "get_page_stats")]
    [Description("Get Bing page statistics for a site. The output uses page instead of Bing's reused Query field.")]
    internal Task<string> GetPageStats(
        [Description("The Bing site URL.")] string site_url,
        CancellationToken cancellationToken = default)
        => ExecuteAsync(
            ct => client.GetPageStatsAsync(site_url, ct),
            BingWebmasterJsonContext.Default.PageStatsResponse,
            cancellationToken);

    [McpServerTool(Name = "get_page_query_stats")]
    [Description("Get Bing search query statistics for a specific page.")]
    internal Task<string> GetPageQueryStats(
        [Description("The Bing site URL.")] string site_url,
        [Description("The page URL to inspect.")] string page,
        CancellationToken cancellationToken = default)
        => ExecuteAsync(
            ct => client.GetPageQueryStatsAsync(site_url, page, ct),
            BingWebmasterJsonContext.Default.PageQueryStatsResponse,
            cancellationToken);

    [McpServerTool(Name = "get_query_page_stats")]
    [Description("Get Bing page statistics for a specific query. The output uses page instead of Bing's reused Query field.")]
    internal Task<string> GetQueryPageStats(
        [Description("The Bing site URL.")] string site_url,
        [Description("The search query to inspect.")] string query,
        CancellationToken cancellationToken = default)
        => ExecuteAsync(
            ct => client.GetQueryPageStatsAsync(site_url, query, ct),
            BingWebmasterJsonContext.Default.QueryPageStatsResponse,
            cancellationToken);

    [McpServerTool(Name = "get_keyword_stats")]
    [Description("Get market-wide Bing keyword statistics. This endpoint does not require a site URL.")]
    internal Task<string> GetKeywordStats(
        [Description("The keyword query text.")] string query,
        [Description("The market country, for example US.")] string country,
        [Description("The market language, for example en-US.")] string language,
        CancellationToken cancellationToken = default)
        => ExecuteAsync(
            ct => client.GetKeywordStatsAsync(query, country, language, ct),
            BingWebmasterJsonContext.Default.KeywordStatsResponse,
            cancellationToken);

    private static async Task<string> ExecuteAsync<T>(
        Func<CancellationToken, Task<T>> action,
        JsonTypeInfo<T> typeInfo,
        CancellationToken cancellationToken)
    {
        try
        {
            var result = await action(cancellationToken).ConfigureAwait(false);
            return JsonSerializer.Serialize(result, typeInfo);
        }
        catch (Exception ex)
        {
            return JsonSerializer.Serialize(
                new ErrorResult($"{ex.GetType().Name}: {ex.Message}"),
                BingWebmasterJsonContext.Default.ErrorResult);
        }
    }
}
