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

    [McpServerTool(Name = "remove_site")]
    [Description("Remove a site from Bing Webmaster Tools.")]
    internal Task<string> RemoveSite(
        [Description("The site URL to remove.")] string site_url,
        CancellationToken cancellationToken = default)
        => ExecuteAsync(ct => client.RemoveSiteAsync(site_url, ct), BingWebmasterJsonContext.Default.RemoveSiteResponse, cancellationToken);

    [McpServerTool(Name = "verify_site")]
    [Description("Verify a site in Bing Webmaster Tools.")]
    internal Task<string> VerifySite(
        [Description("The site URL to verify.")] string site_url,
        CancellationToken cancellationToken = default)
        => ExecuteAsync(ct => client.VerifySiteAsync(site_url, ct), BingWebmasterJsonContext.Default.VerifySiteResponse, cancellationToken);

    [McpServerTool(Name = "get_site_roles")]
    [Description("List delegated Bing Webmaster roles for a site.")]
    internal Task<string> GetSiteRoles(
        [Description("The Bing site URL.")] string site_url,
        [Description("Whether Bing should include roles inherited from subdomains.")] bool include_all_subdomains = false,
        CancellationToken cancellationToken = default)
        => ExecuteAsync(
            ct => client.GetSiteRolesAsync(site_url, include_all_subdomains, ct),
            BingWebmasterJsonContext.Default.GetSiteRolesResponse,
            cancellationToken);

    [McpServerTool(Name = "add_site_role")]
    [Description("Delegate Bing Webmaster access to another user for a site.")]
    internal Task<string> AddSiteRole(
        [Description("The Bing site URL.")] string site_url,
        [Description("The delegated site or subdomain URL.")] string delegated_url,
        [Description("The delegate's email address.")] string user_email,
        [Description("The Bing authentication code for the delegated site.")] string authentication_code,
        [Description("Grant administrator access when true.")] bool is_administrator = false,
        [Description("Grant read-only access when true.")] bool is_read_only = true,
        CancellationToken cancellationToken = default)
        => ExecuteAsync(
            ct => client.AddSiteRoleAsync(site_url, delegated_url, user_email, authentication_code, is_administrator, is_read_only, ct),
            BingWebmasterJsonContext.Default.AddSiteRoleResponse,
            cancellationToken);

    [McpServerTool(Name = "remove_site_role")]
    [Description("Remove a delegated Bing Webmaster role from a site.")]
    internal Task<string> RemoveSiteRole(
        [Description("The Bing site URL.")] string site_url,
        [Description("The delegated user's email address.")] string email,
        [Description("The role to remove: Administrator, ReadOnly, or ReadWrite.")] string role,
        CancellationToken cancellationToken = default)
        => ExecuteAsync(
            ct => client.RemoveSiteRoleAsync(site_url, email, role, ct),
            BingWebmasterJsonContext.Default.RemoveSiteRoleResponse,
            cancellationToken);

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

    [McpServerTool(Name = "remove_sitemap")]
    [Description("Remove a submitted sitemap from Bing Webmaster Tools.")]
    internal Task<string> RemoveSitemap(
        [Description("The Bing site URL.")] string site_url,
        [Description("The sitemap URL to remove.")] string feed_url,
        CancellationToken cancellationToken = default)
        => ExecuteAsync(
            ct => client.RemoveSitemapAsync(site_url, feed_url, ct),
            BingWebmasterJsonContext.Default.RemoveSitemapResponse,
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

    [McpServerTool(Name = "get_content_submission_quota")]
    [Description("Get the current Bing content submission quota for a site.")]
    internal Task<string> GetContentSubmissionQuota(
        [Description("The Bing site URL.")] string site_url,
        CancellationToken cancellationToken = default)
        => ExecuteAsync(
            ct => client.GetContentSubmissionQuotaAsync(site_url, ct),
            BingWebmasterJsonContext.Default.ContentSubmissionQuotaResponse,
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

    [McpServerTool(Name = "get_blocked_urls")]
    [Description("List blocked URL removal requests for a site.")]
    internal Task<string> GetBlockedUrls(
        [Description("The Bing site URL.")] string site_url,
        CancellationToken cancellationToken = default)
        => ExecuteAsync(
            ct => client.GetBlockedUrlsAsync(site_url, ct),
            BingWebmasterJsonContext.Default.GetBlockedUrlsResponse,
            cancellationToken);

    [McpServerTool(Name = "add_blocked_url")]
    [Description("Submit a blocked URL removal request to Bing.")]
    internal Task<string> AddBlockedUrl(
        [Description("The Bing site URL.")] string site_url,
        [Description("The page or directory URL to block.")] string url,
        [Description("The blocked entity type: Page or Directory.")] string entity_type = "Page",
        [Description("The removal type: CacheOnly or FullRemoval.")] string request_type = "CacheOnly",
        CancellationToken cancellationToken = default)
        => ExecuteAsync(
            ct => client.AddBlockedUrlAsync(site_url, url, entity_type, request_type, ct),
            BingWebmasterJsonContext.Default.AddBlockedUrlResponse,
            cancellationToken);

    [McpServerTool(Name = "remove_blocked_url")]
    [Description("Remove a blocked URL request from Bing.")]
    internal Task<string> RemoveBlockedUrl(
        [Description("The Bing site URL.")] string site_url,
        [Description("The page or directory URL to unblock.")] string url,
        [Description("The blocked entity type: Page or Directory.")] string entity_type = "Page",
        [Description("The removal type: CacheOnly or FullRemoval.")] string request_type = "FullRemoval",
        CancellationToken cancellationToken = default)
        => ExecuteAsync(
            ct => client.RemoveBlockedUrlAsync(site_url, url, entity_type, request_type, ct),
            BingWebmasterJsonContext.Default.RemoveBlockedUrlResponse,
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

    [McpServerTool(Name = "get_children_url_info")]
    [Description("List child URL index metadata for a page or directory, with optional Bing filters.")]
    internal Task<string> GetChildrenUrlInfo(
        [Description("The Bing site URL.")] string site_url,
        [Description("The parent page or directory URL to inspect.")] string url,
        [Description("Zero-based result page.")] int page = 0,
        [Description("The crawl date filter: Any, LastWeek, LastTwoWeeks, or LastThreeWeeks.")] string crawl_date_filter = "Any",
        [Description("The discovered date filter: Any, LastWeek, or LastMonth.")] string discovered_date_filter = "Any",
        [Description("The document flags filter: Any, IsBlockedByRobotsTxt, or IsMalware.")] string doc_flags_filter = "Any",
        [Description("The HTTP code filter: Any, Code2xx, Code3xx, Code301, Code302, Code4xx, Code5xx, or AllOthers.")] string http_code_filter = "Any",
        CancellationToken cancellationToken = default)
        => ExecuteAsync(
            ct => client.GetChildrenUrlInfoAsync(site_url, url, page, crawl_date_filter, discovered_date_filter, doc_flags_filter, http_code_filter, ct),
            BingWebmasterJsonContext.Default.ChildrenUrlInfoResponse,
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

    [McpServerTool(Name = "get_children_url_traffic_info")]
    [Description("List child URL clicks and impressions for a page or directory.")]
    internal Task<string> GetChildrenUrlTrafficInfo(
        [Description("The Bing site URL.")] string site_url,
        [Description("The parent page or directory URL to inspect.")] string url,
        [Description("Zero-based result page.")] int page = 0,
        CancellationToken cancellationToken = default)
        => ExecuteAsync(
            ct => client.GetChildrenUrlTrafficInfoAsync(site_url, url, page, ct),
            BingWebmasterJsonContext.Default.ChildrenUrlTrafficInfoResponse,
            cancellationToken);

    [McpServerTool(Name = "fetch_url")]
    [Description("Ask Bing to fetch a URL for diagnostic inspection.")]
    internal Task<string> FetchUrl(
        [Description("The Bing site URL.")] string site_url,
        [Description("The page URL Bing should fetch.")] string url,
        CancellationToken cancellationToken = default)
        => ExecuteAsync(
            ct => client.FetchUrlAsync(site_url, url, ct),
            BingWebmasterJsonContext.Default.FetchUrlResponse,
            cancellationToken);

    [McpServerTool(Name = "list_fetched_urls")]
    [Description("List URLs previously submitted to Bing's fetch tool.")]
    internal Task<string> ListFetchedUrls(
        [Description("The Bing site URL.")] string site_url,
        CancellationToken cancellationToken = default)
        => ExecuteAsync(
            ct => client.ListFetchedUrlsAsync(site_url, ct),
            BingWebmasterJsonContext.Default.ListFetchedUrlsResponse,
            cancellationToken);

    [McpServerTool(Name = "get_fetched_url_details")]
    [Description("Get the stored fetch result details for a specific URL.")]
    internal Task<string> GetFetchedUrlDetails(
        [Description("The Bing site URL.")] string site_url,
        [Description("The fetched page URL to inspect.")] string url,
        CancellationToken cancellationToken = default)
        => ExecuteAsync(
            ct => client.GetFetchedUrlDetailsAsync(site_url, url, ct),
            BingWebmasterJsonContext.Default.FetchedUrlDetailsResponse,
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

    [McpServerTool(Name = "get_query_traffic_stats")]
    [Description("Get Bing clicks and impressions over time for one search query on a site.")]
    internal Task<string> GetQueryTrafficStats(
        [Description("The Bing site URL.")] string site_url,
        [Description("The search query to inspect.")] string query,
        CancellationToken cancellationToken = default)
        => ExecuteAsync(
            ct => client.GetQueryTrafficStatsAsync(site_url, query, ct),
            BingWebmasterJsonContext.Default.QueryTrafficStatsResponse,
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

    [McpServerTool(Name = "get_query_page_detail_stats")]
    [Description("Get Bing daily clicks, impressions, and position for a specific query/page pair.")]
    internal Task<string> GetQueryPageDetailStats(
        [Description("The Bing site URL.")] string site_url,
        [Description("The search query to inspect.")] string query,
        [Description("The page URL to inspect.")] string page,
        CancellationToken cancellationToken = default)
        => ExecuteAsync(
            ct => client.GetQueryPageDetailStatsAsync(site_url, query, page, ct),
            BingWebmasterJsonContext.Default.QueryPageDetailStatsResponse,
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

    [McpServerTool(Name = "get_keyword")]
    [Description("Get market-wide Bing impressions for one keyword over a date range. This endpoint does not require a site URL.")]
    internal Task<string> GetKeyword(
        [Description("The keyword query text.")] string query,
        [Description("The market country, for example US.")] string country,
        [Description("The market language, for example en-US.")] string language,
        [Description("The start date in YYYY-MM-DD format.")] string start_date,
        [Description("The end date in YYYY-MM-DD format.")] string end_date,
        CancellationToken cancellationToken = default)
        => ExecuteAsync(
            ct => client.GetKeywordAsync(query, country, language, start_date, end_date, ct),
            BingWebmasterJsonContext.Default.GetKeywordResponse,
            cancellationToken);

    [McpServerTool(Name = "get_related_keywords")]
    [Description("Get related market-wide Bing keywords over a date range. This endpoint does not require a site URL.")]
    internal Task<string> GetRelatedKeywords(
        [Description("The keyword query text.")] string query,
        [Description("The market country, for example US.")] string country,
        [Description("The market language, for example en-US.")] string language,
        [Description("The start date in YYYY-MM-DD format.")] string start_date,
        [Description("The end date in YYYY-MM-DD format.")] string end_date,
        CancellationToken cancellationToken = default)
        => ExecuteAsync(
            ct => client.GetRelatedKeywordsAsync(query, country, language, start_date, end_date, ct),
            BingWebmasterJsonContext.Default.RelatedKeywordsResponse,
            cancellationToken);

    [McpServerTool(Name = "get_site_moves")]
    [Description("List Bing site move requests configured for a site.")]
    internal Task<string> GetSiteMoves(
        [Description("The Bing site URL.")] string site_url,
        CancellationToken cancellationToken = default)
        => ExecuteAsync(
            ct => client.GetSiteMovesAsync(site_url, ct),
            BingWebmasterJsonContext.Default.GetSiteMovesResponse,
            cancellationToken);

    [McpServerTool(Name = "submit_site_move")]
    [Description("Submit a site move request to Bing.")]
    internal Task<string> SubmitSiteMove(
        [Description("The Bing site URL.")] string site_url,
        [Description("The source URL being moved.")] string source_url,
        [Description("The target URL receiving the move.")] string target_url,
        [Description("The move type: Local or Global.")] string move_type = "Local",
        [Description("The move scope: Domain, Host, or Directory.")] string move_scope = "Domain",
        CancellationToken cancellationToken = default)
        => ExecuteAsync(
            ct => client.SubmitSiteMoveAsync(site_url, source_url, target_url, move_type, move_scope, ct),
            BingWebmasterJsonContext.Default.SubmitSiteMoveResponse,
            cancellationToken);

    [McpServerTool(Name = "submit_content")]
    [Description("Submit captured HTTP content to Bing's content submission API.")]
    internal Task<string> SubmitContent(
        [Description("The Bing site URL.")] string site_url,
        [Description("The URL whose content is being submitted.")] string url,
        [Description("The HTTP message payload, base64 encoded.")] string http_message,
        [Description("The structured data payload, base64 encoded.")] string structured_data,
        [Description("The dynamic serving target: None, PcLaptop, Mobile, Amp, Tablet, or NonVisualBrowser.")] string dynamic_serving = "None",
        CancellationToken cancellationToken = default)
        => ExecuteAsync(
            ct => client.SubmitContentAsync(site_url, url, http_message, structured_data, dynamic_serving, ct),
            BingWebmasterJsonContext.Default.SubmitContentResponse,
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
