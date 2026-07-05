---
description: Reference for the get_page_query_stats MCP tool -- parameters, response format, and example prompts for finding which queries drive traffic to a specific page in Bing.
---

# get_page_query_stats

Get the search queries driving traffic to a specific page.

---

## Parameters

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `site_url` | string | Yes | The URL of the site |
| `page` | string | Yes | The specific page URL to get query statistics for |

---

## Response

```json
{
  "siteUrl": "https://www.example.com/",
  "page": "https://www.example.com/blog/my-post",
  "rowCount": 1,
  "rows": [
    {
      "query": "blazor dependency injection",
      "date": "2026-02-20T00:00:00Z",
      "clicks": 12,
      "impressions": 210,
      "avgClickPosition": 2,
      "avgImpressionPosition": 4
    }
  ],
  "queriedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "Which queries are driving traffic to my blog post about dependency injection?"

---

## Notes

- Complements [`get_query_page_stats`](get-query-page-stats.md), which answers the reverse
  question (which pages rank for a given query).
