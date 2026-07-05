"""
Master manifest of all 43 MCP tools for the automated E2E suite.

Each entry defines:
  - minimal_args: only the required parameters. This exercises MCP argument
    binding for every omitted optional parameter -- exactly the layer unit
    tests cannot reach, and where two real bugs were found this session: a
    wrong Bing endpoint name (add_site_role -> "AddSiteRole" instead of
    "AddSiteRoles"), and a missing default value that only broke once MCP
    itself tried to bind a missing argument (C#'s include_all_subdomains /
    is_administrator / is_read_only, and later Go's is_read_only).
  - full_args: every parameter explicitly supplied, exercising the non-default
    code path too. None where a tool has no optional parameters.
  - bing_method: the Bing Webmaster JSON API method name this tool must
    invoke, used to assert the request actually reached the right endpoint.
    None for submit_url_indexnow, which calls the separate IndexNow endpoint.
"""

from dataclasses import dataclass
from typing import Optional


@dataclass(frozen=True)
class ToolCase:
    tool: str
    minimal_args: dict
    full_args: Optional[dict]
    bing_method: Optional[str]


TOOL_CASES: list[ToolCase] = [
    ToolCase("list_sites", {}, None, "GetUserSites"),
    ToolCase("add_site", {"site_url": "https://example.test/"}, None, "AddSite"),
    ToolCase("remove_site", {"site_url": "https://example.test/"}, None, "RemoveSite"),
    ToolCase("verify_site", {"site_url": "https://example.test/"}, None, "VerifySite"),
    ToolCase(
        "get_site_roles",
        {"site_url": "https://example.test/"},
        {"site_url": "https://example.test/", "include_all_subdomains": True},
        "GetSiteRoles",
    ),
    ToolCase(
        "add_site_role",
        {
            "site_url": "https://example.test/",
            "delegated_url": "https://example.test/",
            "user_email": "someone@example.test",
            "authentication_code": "abc123",
        },
        {
            "site_url": "https://example.test/",
            "delegated_url": "https://example.test/",
            "user_email": "someone@example.test",
            "authentication_code": "abc123",
            "is_administrator": True,
            "is_read_only": False,
        },
        "AddSiteRoles",
    ),
    ToolCase(
        "remove_site_role",
        {"site_url": "https://example.test/", "email": "someone@example.test", "role": "ReadOnly"},
        None,
        "RemoveSiteRole",
    ),
    ToolCase("list_sitemaps", {"site_url": "https://example.test/"}, None, "GetFeeds"),
    ToolCase(
        "get_sitemap_details",
        {"site_url": "https://example.test/", "feed_url": "https://example.test/sitemap.xml"},
        None,
        "GetFeedDetails",
    ),
    ToolCase(
        "submit_sitemap",
        {"site_url": "https://example.test/", "feed_url": "https://example.test/sitemap.xml"},
        None,
        "SubmitFeed",
    ),
    ToolCase(
        "remove_sitemap",
        {"site_url": "https://example.test/", "feed_url": "https://example.test/sitemap.xml"},
        None,
        "RemoveFeed",
    ),
    ToolCase(
        "submit_url", {"site_url": "https://example.test/", "url": "https://example.test/page"}, None, "SubmitUrl"
    ),
    ToolCase(
        "submit_url_batch",
        {"site_url": "https://example.test/", "url_list": ["https://example.test/a"]},
        None,
        "SubmitUrlBatch",
    ),
    ToolCase(
        "submit_url_indexnow",
        {"host": "example.test", "url_list": ["https://example.test/a"], "key": "deadbeef"},
        {
            "host": "example.test",
            "url_list": ["https://example.test/a"],
            "key": "deadbeef",
            "key_location": "https://example.test/deadbeef.txt",
        },
        None,
    ),
    ToolCase("get_url_submission_quota", {"site_url": "https://example.test/"}, None, "GetUrlSubmissionQuota"),
    ToolCase(
        "submit_content",
        {
            "site_url": "https://example.test/",
            "url": "https://example.test/page",
            "http_message": "aHR0cA==",
            "structured_data": "ZGF0YQ==",
        },
        {
            "site_url": "https://example.test/",
            "url": "https://example.test/page",
            "http_message": "aHR0cA==",
            "structured_data": "ZGF0YQ==",
            "dynamic_serving": "Mobile",
        },
        "SubmitContent",
    ),
    ToolCase(
        "get_content_submission_quota", {"site_url": "https://example.test/"}, None, "GetContentSubmissionQuota"
    ),
    ToolCase("get_crawl_issues", {"site_url": "https://example.test/"}, None, "GetCrawlIssues"),
    ToolCase("get_crawl_stats", {"site_url": "https://example.test/"}, None, "GetCrawlStats"),
    ToolCase(
        "get_url_info", {"site_url": "https://example.test/", "url": "https://example.test/page"}, None, "GetUrlInfo"
    ),
    ToolCase(
        "get_url_traffic_info",
        {"site_url": "https://example.test/", "url": "https://example.test/page"},
        None,
        "GetUrlTrafficInfo",
    ),
    ToolCase(
        "get_url_links",
        {"site_url": "https://example.test/", "link": "https://example.test/page"},
        {"site_url": "https://example.test/", "link": "https://example.test/page", "page": 2},
        "GetUrlLinks",
    ),
    ToolCase(
        "get_link_counts",
        {"site_url": "https://example.test/"},
        {"site_url": "https://example.test/", "page": 2},
        "GetLinkCounts",
    ),
    ToolCase(
        "get_children_url_info",
        {"site_url": "https://example.test/", "url": "https://example.test/dir/"},
        {
            "site_url": "https://example.test/",
            "url": "https://example.test/dir/",
            "page": 1,
            "crawl_date_filter": "LastWeek",
            "discovered_date_filter": "LastMonth",
            "doc_flags_filter": "IsMalware",
            "http_code_filter": "Code4xx",
        },
        "GetChildrenUrlInfo",
    ),
    ToolCase(
        "get_children_url_traffic_info",
        {"site_url": "https://example.test/", "url": "https://example.test/dir/"},
        {"site_url": "https://example.test/", "url": "https://example.test/dir/", "page": 1},
        "GetChildrenUrlTrafficInfo",
    ),
    ToolCase("get_blocked_urls", {"site_url": "https://example.test/"}, None, "GetBlockedUrls"),
    ToolCase(
        "add_blocked_url",
        {"site_url": "https://example.test/", "url": "https://example.test/blocked"},
        {
            "site_url": "https://example.test/",
            "url": "https://example.test/blocked/",
            "entity_type": "Directory",
            "request_type": "FullRemoval",
        },
        "AddBlockedUrl",
    ),
    ToolCase(
        "remove_blocked_url",
        {"site_url": "https://example.test/", "url": "https://example.test/blocked"},
        {
            "site_url": "https://example.test/",
            "url": "https://example.test/blocked/",
            "entity_type": "Directory",
            "request_type": "CacheOnly",
        },
        "RemoveBlockedUrl",
    ),
    ToolCase(
        "fetch_url", {"site_url": "https://example.test/", "url": "https://example.test/page"}, None, "FetchUrl"
    ),
    ToolCase("list_fetched_urls", {"site_url": "https://example.test/"}, None, "GetFetchedUrls"),
    ToolCase(
        "get_fetched_url_details",
        {"site_url": "https://example.test/", "url": "https://example.test/page"},
        None,
        "GetFetchedUrlDetails",
    ),
    ToolCase("get_site_moves", {"site_url": "https://example.test/"}, None, "GetSiteMoves"),
    ToolCase(
        "submit_site_move",
        {
            "site_url": "https://example.test/",
            "source_url": "https://example.test/",
            "target_url": "https://example-new.test/",
        },
        {
            "site_url": "https://example.test/",
            "source_url": "https://example.test/",
            "target_url": "https://example-new.test/",
            "move_type": "Global",
            "move_scope": "Host",
        },
        "SubmitSiteMove",
    ),
    ToolCase("get_rank_and_traffic_stats", {"site_url": "https://example.test/"}, None, "GetRankAndTrafficStats"),
    ToolCase("get_query_stats", {"site_url": "https://example.test/"}, None, "GetQueryStats"),
    ToolCase("get_page_stats", {"site_url": "https://example.test/"}, None, "GetPageStats"),
    ToolCase(
        "get_page_query_stats",
        {"site_url": "https://example.test/", "page": "https://example.test/page"},
        None,
        "GetPageQueryStats",
    ),
    ToolCase(
        "get_query_page_stats",
        {"site_url": "https://example.test/", "query": "mock query"},
        None,
        "GetQueryPageStats",
    ),
    ToolCase(
        "get_query_page_detail_stats",
        {"site_url": "https://example.test/", "query": "mock query", "page": "https://example.test/page"},
        None,
        "GetQueryPageDetailStats",
    ),
    ToolCase(
        "get_query_traffic_stats",
        {"site_url": "https://example.test/", "query": "mock query"},
        None,
        "GetQueryTrafficStats",
    ),
    ToolCase(
        "get_keyword_stats", {"query": "mock keyword", "country": "US", "language": "en-US"}, None, "GetKeywordStats"
    ),
    ToolCase(
        "get_keyword",
        {
            "query": "mock keyword",
            "country": "US",
            "language": "en-US",
            "start_date": "2024-01-01",
            "end_date": "2024-01-31",
        },
        None,
        "GetKeyword",
    ),
    ToolCase(
        "get_related_keywords",
        {
            "query": "mock keyword",
            "country": "US",
            "language": "en-US",
            "start_date": "2024-01-01",
            "end_date": "2024-01-31",
        },
        None,
        "GetRelatedKeywords",
    ),
    # Phase 3 -- URL normalization (query parameters).
    ToolCase("get_query_parameters", {"site_url": "https://example.test/"}, None, "GetQueryParameters"),
    ToolCase(
        "add_query_parameter",
        {"site_url": "https://example.test/", "query_parameter": "utm_campaign"},
        None,
        "AddQueryParameter",
    ),
    ToolCase(
        "remove_query_parameter",
        {"site_url": "https://example.test/", "query_parameter": "utm_campaign"},
        None,
        "RemoveQueryParameter",
    ),
    ToolCase(
        "enable_disable_query_parameter",
        {"site_url": "https://example.test/", "query_parameter": "utm_campaign", "is_enabled": False},
        None,
        "EnableDisableQueryParameter",
    ),
    # Phase 3 -- geo-targeting (country/region settings).
    ToolCase("get_country_region_settings", {"site_url": "https://example.test/"}, None, "GetCountryRegionSettings"),
    ToolCase(
        "add_country_region_settings",
        {
            "site_url": "https://example.test/",
            "two_letter_iso_country_code": "us",
            "settings_type": "Domain",
            "url": "https://example.test/",
        },
        None,
        "AddCountryRegionSettings",
    ),
    ToolCase(
        "remove_country_region_settings",
        {
            "site_url": "https://example.test/",
            "two_letter_iso_country_code": "us",
            "settings_type": "Domain",
            "url": "https://example.test/",
        },
        None,
        "RemoveCountryRegionSettings",
    ),
    # Phase 3 -- connected pages.
    ToolCase("get_connected_pages", {"site_url": "https://example.test/"}, None, "GetConnectedPages"),
    ToolCase(
        "add_connected_page",
        {"site_url": "https://example.test/", "master_url": "https://master.example.test/"},
        None,
        "AddConnectedPage",
    ),
    # Phase 3 -- page preview blocks.
    ToolCase(
        "get_active_page_preview_blocks", {"site_url": "https://example.test/"}, None, "GetActivePagePreviewBlocks"
    ),
    ToolCase(
        "add_page_preview_block",
        {"site_url": "https://example.test/", "url": "https://example.test/blocked-preview", "reason": "Other"},
        None,
        "AddPagePreviewBlock",
    ),
    ToolCase(
        "remove_page_preview_block",
        {"site_url": "https://example.test/", "url": "https://example.test/blocked-preview"},
        None,
        "RemovePagePreviewBlock",
    ),
]

_manifest_names = {c.tool for c in TOOL_CASES}


def assert_manifest_matches_tool_list(tool_names) -> None:
    """Sanity check the manifest covers exactly the tools the server advertises --
    catches a tool being added to one language and forgotten in the other, or this
    manifest going stale as new tools are added in future phases."""
    tool_name_set = set(tool_names)
    missing = tool_name_set - _manifest_names
    extra = _manifest_names - tool_name_set
    assert not missing, f"tools/list advertises tools missing from TOOL_CASES: {sorted(missing)}"
    assert not extra, f"TOOL_CASES references tools not in tools/list: {sorted(extra)}"
