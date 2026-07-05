---
description: Reference for the get_query_traffic_stats MCP tool -- parameters, response format, and example prompts for daily traffic statistics for a specific query in Bing.
---

# get_query_traffic_stats

Get daily clicks and impressions for a specific search query across your whole site.

---

## Parameters

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `site_url` | string | Yes | The URL of the site |
| `query` | string | Yes | The search query |

---

## Response

```json
{
  "siteUrl": "https://www.example.com/",
  "query": "blazor dependency injection",
  "rowCount": 1,
  "rows": [
    {
      "date": "2026-02-20T00:00:00Z",
      "clicks": 18,
      "impressions": 340
    }
  ],
  "queriedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "How has traffic for 'blazor dependency injection' trended over time on my site?"

---

## Notes

- Site-wide daily trend for one query (no per-page breakdown) -- use
  [`get_query_page_stats`](get-query-page-stats.md) to see which pages contribute.
