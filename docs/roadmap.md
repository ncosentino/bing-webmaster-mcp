---
description: Shipped Phase 2, and planned Phase 3/4 tools for the Bing Webmaster Tools MCP server, beyond the original 22-tool MVP.
---

# Roadmap

The current release covers 43 tools: the original 22-tool MVP plus Phase 2 (21 tools, shipped and
live-tested). The full Bing Webmaster API surface is much larger still -- this page tracks what's
deliberately deferred and why, plus known live-testing caveats for what's already shipped.

---

## Phase 2 -- shipped

Deeper coverage and write/destructive operations, live-tested against a real Bing Webmaster
account with both Go and C# confirmed at parity:

- `remove_site`, `get_site_roles`, `add_site_role`, `remove_site_role`
- Blocked URLs: `get_blocked_urls`, `add_blocked_url`, `remove_blocked_url`
- Deeper query/keyword detail stats: `get_query_page_detail_stats`, `get_query_traffic_stats`,
  `get_keyword`, `get_related_keywords`
- Directory-level inspection: `get_children_url_info`, `get_children_url_traffic_info`
- Fetch-as-Bing: `fetch_url`, `list_fetched_urls`, `get_fetched_url_details`
- `remove_sitemap`
- Site moves: `get_site_moves`, `submit_site_move`
- Content Submission API: `submit_content`, `get_content_submission_quota`

**Live-testing caveats:**

- `get_site_moves` returned an HTML 404 ("Service" error page) against a real, verified site
  during testing -- confirmed via a raw HTTP call that this is Bing's actual live behavior, not a
  client bug. This endpoint may have been deprecated or renamed since it was last documented. The
  tool still handles this gracefully (returns a clean error, doesn't crash), but a genuine success
  response for this specific endpoint hasn't been observed.
- `add_site_role` / `remove_site_role`'s write path is unit-tested but not live-tested (doing so
  safely requires delegating to a second real email address, which wasn't available in testing).
  `get_site_roles` (read) is confirmed live-working.
- `submit_site_move` is unit-tested only -- it's consequential and not easily reversible (see the
  [tool reference](tools/submit-site-move.md)), so it wasn't exercised against a real account.
- `submit_content` is unit-tested only -- constructing valid base64 HTTP content for a live test
  is complex and this is a rarely-used, advanced capability.

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

## Phase 4 (planned)

- Full OAuth 2.0 authentication support alongside the API key (Bing's 3-legged web consent flow),
  for scenarios beyond a single-operator personal-use setup. Deferred separately from Phase 2/3
  because it's a fundamentally different kind of work (an auth/consent flow, not new tools) with
  its own design and testing needs.

---

## Excluded

The deep-link method family (`GetDeepLink`, `GetDeepLinkAlgoUrls`, `GetDeepLinkBlocks`,
`AddDeepLinkBlock`, `RemoveDeepLinkBlock`, `UpdateDeepLink`) is excluded entirely. Several of these
are explicitly marked **Obsolete** in Microsoft's own API reference, and the rest belong to the
same deprecated feature area.

---

Have a use case that needs something from Phase 3 or 4 sooner, or ran into one of the caveats
above?
[Open an issue](https://github.com/ncosentino/bing-webmaster-mcp/issues) describing what you're
trying to do.
