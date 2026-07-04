---
description: Planned Phase 2 and Phase 3 tools for the Bing Webmaster Tools MCP server, beyond the current 22-tool MVP.
---

# Roadmap

The current release covers a curated MVP of 22 tools (21 classic Bing Webmaster API operations
plus the IndexNow bonus tool). The full Bing Webmaster API surface is much larger (~50 operations)
-- this page tracks what's deliberately deferred and why.

---

## Phase 2 (planned)

Deeper coverage and write/destructive operations:

- `remove_site` -- destructive, requires explicit confirmation semantics
- Site role delegation: `get_site_roles`, `add_site_role`, `remove_site_role`
- Blocked URLs: `get_blocked_urls`, `add_blocked_url`, `remove_blocked_url`
- Deeper query/keyword detail stats: `get_query_page_detail_stats`, `get_query_traffic_stats`,
  `get_keyword`, `get_related_keywords`
- Directory-level inspection: `get_children_url_info`, `get_children_url_traffic_info`
- Fetch-as-Bing: `fetch_url`, `list_fetched_urls`, `get_fetched_url_details`
- `remove_sitemap`
- Site moves: `get_site_moves`, `submit_site_move`
- Content Submission API: `submit_content`, `get_content_submission_quota`
- Full OAuth 2.0 authentication support alongside the API key (Bing's 3-legged web consent flow),
  for scenarios beyond a single-operator personal-use setup

---

## Phase 3 (planned)

Niche/administrative operations, used far less frequently:

- URL normalization (query parameters): `get_query_parameters`, `add_query_parameter`,
  `enable_disable_query_parameter`, `remove_query_parameter`
- Geo-targeting: `get_country_region_settings`, `add_country_region_settings`,
  `remove_country_region_settings`
- Connected pages: `get_connected_pages`, `add_connected_page`
- Page preview blocks: `get_active_page_preview_blocks`, `add_page_preview_block`,
  `remove_page_preview_block`

---

## Excluded

The deep-link method family (`GetDeepLink`, `GetDeepLinkAlgoUrls`, `GetDeepLinkBlocks`,
`AddDeepLinkBlock`, `RemoveDeepLinkBlock`, `UpdateDeepLink`) is excluded entirely. Several of these
are explicitly marked **Obsolete** in Microsoft's own API reference, and the rest belong to the
same deprecated feature area.

---

Have a use case that needs something from Phase 2 or 3 sooner? 
[Open an issue](https://github.com/ncosentino/bing-webmaster-mcp/issues) describing what you're
trying to do.
