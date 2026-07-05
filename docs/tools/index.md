---
description: Overview of the 55 MCP tools exposed by the Bing Webmaster Tools MCP server, grouped by area, with common parameter and error-handling notes.
---

# MCP Tools

55 tools are exposed by this MCP server -- the 22-tool MVP (21 classic Bing Webmaster API
operations plus one bonus IndexNow tool), 21 Phase 2 tools covering deeper site management,
blocked URLs, fetch diagnostics, site moves, and the Content Submission API, and 12 Phase 3 tools
covering URL normalization, geo-targeting, connected pages, and page preview blocks. All tools
work identically across the Go and C# implementations. See the [Roadmap](../roadmap.md) for what's
planned beyond this set.

| Tool | Area | Description |
|------|------|-------------|
| [`list_sites`](list-sites.md) | Sites | List all sites in your Bing Webmaster account |
| [`add_site`](add-site.md) | Sites | Add a new site |
| [`verify_site`](verify-site.md) | Sites | Attempt to verify ownership of a site |
| [`remove_site`](remove-site.md) | Sites | Remove a site (destructive) |
| [`get_site_roles`](get-site-roles.md) | Site Access | List delegated user access to a site |
| [`add_site_role`](add-site-role.md) | Site Access | Delegate site access to a user |
| [`remove_site_role`](remove-site-role.md) | Site Access | Revoke a user's delegated access |
| [`list_sitemaps`](list-sitemaps.md) | Sitemaps | List submitted sitemaps for a site |
| [`get_sitemap_details`](get-sitemap-details.md) | Sitemaps | Get details for one submitted sitemap |
| [`submit_sitemap`](submit-sitemap.md) | Sitemaps | Submit a sitemap URL |
| [`remove_sitemap`](remove-sitemap.md) | Sitemaps | Remove a submitted sitemap |
| [`submit_url`](submit-url.md) | URL Submission | Submit a single URL for indexing |
| [`submit_url_batch`](submit-url-batch.md) | URL Submission | Submit up to 500 URLs in one call |
| [`submit_url_indexnow`](submit-url-indexnow.md) | URL Submission | Ping the IndexNow protocol for instant (re)indexing |
| [`get_url_submission_quota`](get-url-submission-quota.md) | URL Submission | Check remaining daily/monthly submission quota |
| [`submit_content`](submit-content.md) | Content Submission | Submit raw content directly for a URL |
| [`get_content_submission_quota`](get-content-submission-quota.md) | Content Submission | Check remaining content submission quota |
| [`get_crawl_issues`](get-crawl-issues.md) | Crawling | List URLs with crawl issues |
| [`get_crawl_stats`](get-crawl-stats.md) | Crawling | Daily crawl statistics |
| [`get_url_info`](get-url-info.md) | URL & Index Inspection | Index status/details for a single URL |
| [`get_url_traffic_info`](get-url-traffic-info.md) | URL & Index Inspection | Clicks/impressions for a single URL |
| [`get_url_links`](get-url-links.md) | URL & Index Inspection | Inbound links to a specific URL |
| [`get_link_counts`](get-link-counts.md) | URL & Index Inspection | Pages with inbound links |
| [`get_children_url_info`](get-children-url-info.md) | Directory Inspection | Index details for every URL under a directory |
| [`get_children_url_traffic_info`](get-children-url-traffic-info.md) | Directory Inspection | Traffic for every URL under a directory |
| [`get_blocked_urls`](get-blocked-urls.md) | Blocked URLs | List blocked pages/directories |
| [`add_blocked_url`](add-blocked-url.md) | Blocked URLs | Block a page or directory |
| [`remove_blocked_url`](remove-blocked-url.md) | Blocked URLs | Unblock a page or directory |
| [`fetch_url`](fetch-url.md) | Fetch Diagnostics | Request an immediate Bing crawl of a URL |
| [`list_fetched_urls`](list-fetched-urls.md) | Fetch Diagnostics | List URLs submitted for fetch requests |
| [`get_fetched_url_details`](get-fetched-url-details.md) | Fetch Diagnostics | Full result of a fetch request |
| [`get_site_moves`](get-site-moves.md) | Site Moves | List site move requests |
| [`submit_site_move`](submit-site-move.md) | Site Moves | Notify Bing of a site migration |
| [`get_rank_and_traffic_stats`](get-rank-and-traffic-stats.md) | Search Analytics | Daily clicks/impressions over time |
| [`get_query_stats`](get-query-stats.md) | Search Analytics | Top queries by clicks/impressions |
| [`get_page_stats`](get-page-stats.md) | Search Analytics | Top pages by clicks/impressions |
| [`get_page_query_stats`](get-page-query-stats.md) | Search Analytics | Queries driving traffic to one page |
| [`get_query_page_stats`](get-query-page-stats.md) | Search Analytics | Pages ranking for one query |
| [`get_query_page_detail_stats`](get-query-page-detail-stats.md) | Search Analytics | Detailed daily stats for one query+page |
| [`get_query_traffic_stats`](get-query-traffic-stats.md) | Search Analytics | Daily stats for one query, site-wide |
| [`get_keyword_stats`](get-keyword-stats.md) | Search Analytics | Historical keyword impression stats (market-wide) |
| [`get_keyword`](get-keyword.md) | Keyword Research | Market-wide impressions for one keyword/period |
| [`get_related_keywords`](get-related-keywords.md) | Keyword Research | Discover related keywords with impression data |
| [`get_query_parameters`](get-query-parameters.md) | URL Normalization | List URL normalization query parameters |
| [`add_query_parameter`](add-query-parameter.md) | URL Normalization | Add a query parameter Bing should ignore |
| [`remove_query_parameter`](remove-query-parameter.md) | URL Normalization | Remove a query parameter |
| [`enable_disable_query_parameter`](enable-disable-query-parameter.md) | URL Normalization | Toggle a query parameter on/off |
| [`get_country_region_settings`](get-country-region-settings.md) | Geo-Targeting | List geo-targeting settings |
| [`add_country_region_settings`](add-country-region-settings.md) | Geo-Targeting | Target a page/directory/domain/subdomain at a country |
| [`remove_country_region_settings`](remove-country-region-settings.md) | Geo-Targeting | Remove a geo-targeting setting |
| [`get_connected_pages`](get-connected-pages.md) | Connected Pages | List pages connected to your site |
| [`add_connected_page`](add-connected-page.md) | Connected Pages | Declare a connected page |
| [`get_active_page_preview_blocks`](get-active-page-preview-blocks.md) | Page Preview Blocks | List active page preview blocks |
| [`add_page_preview_block`](add-page-preview-block.md) | Page Preview Blocks | Block a page's search result preview |
| [`remove_page_preview_block`](remove-page-preview-block.md) | Page Preview Blocks | Unblock a page's search result preview |

---

## Common Notes

### `site_url`

Most tools take a `site_url` parameter -- this must be a site already added and (for most
operations) verified in your Bing Webmaster account. Use [`list_sites`](list-sites.md) to see the
exact URLs your API key has access to. A few tools (`get_keyword`, `get_related_keywords`,
`get_keyword_stats`) are market-wide and don't take a `site_url` at all.

### Error Responses

All tools return a JSON error object when an exception occurs, instead of breaking the MCP
session:

```json
{
  "error": "ExceptionType: message"
}
```

### Timestamps

Read tools include a `queriedAt` timestamp (ISO-8601 UTC). Write/submission tools echo the key
inputs plus a `submittedAt`/`requestedAt` timestamp and a success indicator. All dates are
normalized to ISO-8601 -- Bing's underlying `/Date(...)/ ` wire format is never exposed to the
client.

### `get_query_stats` vs `get_page_stats`

Bing's API reuses one underlying response shape across several endpoints. `get_query_stats` and
`get_page_query_stats` return a `query` field holding a search query; `get_page_stats` and
`get_query_page_stats` return a `page` field holding a page URL, using the same underlying data
shape. This server renames the field appropriately per tool so the output is unambiguous.

### Destructive and consequential tools

Most write tools are safely reversible (`add_blocked_url` / `remove_blocked_url`, `submit_sitemap`
/ `remove_sitemap`, `add_country_region_settings` / `remove_country_region_settings`,
`add_page_preview_block` / `remove_page_preview_block`). A few are not: `remove_site` permanently
drops a site's history from your account view, `submit_site_move` tells Bing to treat one site as
superseded by another in a way that isn't easily undone, and `add_connected_page` has no matching
removal tool in Bing's API at all. Use all three with care.

