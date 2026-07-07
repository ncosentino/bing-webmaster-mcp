---
description: Reference for the get_query_page_stats MCP tool -- parameters, response format, and example prompts for finding which pages rank for a specific query in Bing.
---

# get_query_page_stats

Get the pages ranking for a specific search query.

---

## Parameters

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `site_url` | string | Yes | The URL of the site |
| `query` | string | Yes | The search query to get page statistics for |

---

## Response

```json
{
  "siteUrl": "https://www.example.com/",
  "query": "blazor dependency injection",
  "rowCount": 3,
  "rows": [
    {
      "page": "https://www.example.com/blog/my-post",
      "date": "2026-02-19T00:00:00Z",
      "clicks": 14,
      "impressions": 225,
      "avgClickPosition": 2,
      "avgImpressionPosition": 4
    },
    {
      "page": "https://www.example.com/blog/my-post",
      "date": "2026-02-20T00:00:00Z",
      "clicks": 12,
      "impressions": 210,
      "avgClickPosition": 2,
      "avgImpressionPosition": 4
    },
    {
      "page": "https://www.example.com/docs/di-guide",
      "date": "2026-02-20T00:00:00Z",
      "clicks": 4,
      "impressions": 88,
      "avgClickPosition": 6,
      "avgImpressionPosition": 8
    }
  ],
  "queriedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "Which of my pages rank for 'blazor dependency injection' on Bing?"

---

## Notes

- Complements [`get_page_query_stats`](get-page-query-stats.md), which answers the reverse
  question (which queries drive traffic to a given page).
- For a deeper breakdown of a single query+page combination, see
  [`get_query_page_detail_stats`](get-query-page-detail-stats.md).
- Rows are per page *per day*, so the same page can appear multiple times across the window --
  the example above is truncated to 3 rows for readability.
- There is no date range parameter -- Bing returns a fixed window server-side and this tool
  cannot request a specific period.
