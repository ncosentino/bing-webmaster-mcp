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
  "rowCount": 3,
  "rows": [
    {
      "date": "2026-02-18T00:00:00Z",
      "clicks": 15,
      "impressions": 312
    },
    {
      "date": "2026-02-19T00:00:00Z",
      "clicks": 21,
      "impressions": 365
    },
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
- Each row carries its own `date`, but there is no request parameter to filter or narrow that
  range -- Bing returns its full available daily history for the query and you filter client-side.
  The example above is truncated to 3 rows for readability; a real response typically spans
  several weeks.
