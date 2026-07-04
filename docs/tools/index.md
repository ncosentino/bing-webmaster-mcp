---
description: Overview of the 22 MCP tools exposed by the Bing Webmaster Tools MCP server, grouped by area, with common parameter and error-handling notes.
---

# MCP Tools

22 tools are exposed by this MCP server -- 21 cover the classic Bing Webmaster API, plus one bonus
IndexNow tool. All tools work identically across the Go and C# implementations. See the
[Roadmap](../roadmap.md) for what's planned beyond this set.

| Tool | Area | Description |
|------|------|-------------|
| [`list_sites`](list-sites.md) | Sites | List all sites in your Bing Webmaster account |
| [`add_site`](add-site.md) | Sites | Add a new site |
| [`verify_site`](verify-site.md) | Sites | Attempt to verify ownership of a site |
| [`list_sitemaps`](list-sitemaps.md) | Sitemaps | List submitted sitemaps for a site |
| [`get_sitemap_details`](get-sitemap-details.md) | Sitemaps | Get details for one submitted sitemap |
| [`submit_sitemap`](submit-sitemap.md) | Sitemaps | Submit a sitemap URL |
| [`submit_url`](submit-url.md) | URL Submission | Submit a single URL for indexing |
| [`submit_url_batch`](submit-url-batch.md) | URL Submission | Submit up to 500 URLs in one call |
| [`submit_url_indexnow`](submit-url-indexnow.md) | URL Submission | Ping the IndexNow protocol for instant (re)indexing |
| [`get_url_submission_quota`](get-url-submission-quota.md) | URL Submission | Check remaining daily/monthly submission quota |
| [`get_crawl_issues`](get-crawl-issues.md) | Crawling | List URLs with crawl issues |
| [`get_crawl_stats`](get-crawl-stats.md) | Crawling | Daily crawl statistics |
| [`get_url_info`](get-url-info.md) | URL & Index Inspection | Index status/details for a single URL |
| [`get_url_traffic_info`](get-url-traffic-info.md) | URL & Index Inspection | Clicks/impressions for a single URL |
| [`get_url_links`](get-url-links.md) | URL & Index Inspection | Inbound links to a specific URL |
| [`get_link_counts`](get-link-counts.md) | URL & Index Inspection | Pages with inbound links |
| [`get_rank_and_traffic_stats`](get-rank-and-traffic-stats.md) | Search Analytics | Daily clicks/impressions over time |
| [`get_query_stats`](get-query-stats.md) | Search Analytics | Top queries by clicks/impressions |
| [`get_page_stats`](get-page-stats.md) | Search Analytics | Top pages by clicks/impressions |
| [`get_page_query_stats`](get-page-query-stats.md) | Search Analytics | Queries driving traffic to one page |
| [`get_query_page_stats`](get-query-page-stats.md) | Search Analytics | Pages ranking for one query |
| [`get_keyword_stats`](get-keyword-stats.md) | Search Analytics | Historical keyword impression stats (market-wide) |

---

## Common Notes

### `site_url`

Most tools take a `site_url` parameter -- this must be a site already added and (for most
operations) verified in your Bing Webmaster account. Use [`list_sites`](list-sites.md) to see the
exact URLs your API key has access to.

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
