---
description: Reference for the get_url_traffic_info MCP tool -- parameters, response format, and example prompts for checking a single URL's clicks and impressions in Bing.
---

# get_url_traffic_info

Get click and impression traffic data for a single URL.

---

## Parameters

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `site_url` | string | Yes | The URL of the site |
| `url` | string | Yes | The specific URL to inspect |

---

## Response

```json
{
  "siteUrl": "https://www.example.com/",
  "url": "https://www.example.com/blog/my-post",
  "isPage": true,
  "clicks": 142,
  "impressions": 3820,
  "queriedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "How much search traffic has https://www.example.com/blog/my-post gotten from Bing?"

---

## Notes

- Pair with [`get_page_query_stats`](get-page-query-stats.md) to see which specific queries are
  driving that traffic.
- Pair with [`get_url_info`](get-url-info.md) for index status/crawl metadata on the same URL.
