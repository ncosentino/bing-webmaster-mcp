---
description: Reference for the get_query_page_detail_stats MCP tool -- parameters, response format, and example prompts for detailed query+page statistics in Bing.
---

# get_query_page_detail_stats

Get detailed daily statistics for a specific query and page combination.

---

## Parameters

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `site_url` | string | Yes | The URL of the site |
| `query` | string | Yes | The search query |
| `page` | string | Yes | The specific page URL |

---

## Response

```json
{
  "siteUrl": "https://www.example.com/",
  "query": "blazor dependency injection",
  "page": "https://www.example.com/blog/my-post",
  "rowCount": 1,
  "rows": [
    {
      "date": "2026-02-20T00:00:00Z",
      "clicks": 3,
      "impressions": 40,
      "position": 4
    }
  ],
  "queriedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "Give me the day-by-day breakdown of how 'blazor dependency injection' performs for my blog post specifically."

---

## Notes

- More granular than [`get_query_page_stats`](get-query-page-stats.md) -- this drills into a
  single query+page pair over time, including `position` (rank), which the broader stats tools
  don't report per-row.
