---
description: Reference for the get_page_stats MCP tool -- parameters, response format, and example prompts for top pages by clicks and impressions in Bing.
---

# get_page_stats

Get traffic statistics for your site's top pages.

---

## Parameters

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `site_url` | string | Yes | The URL of the site |

---

## Response

```json
{
  "siteUrl": "https://www.example.com/",
  "rowCount": 3,
  "rows": [
    {
      "page": "https://www.example.com/blog/my-post",
      "date": "2026-02-19T00:00:00Z",
      "clicks": 58,
      "impressions": 1150,
      "avgClickPosition": 3,
      "avgImpressionPosition": 5
    },
    {
      "page": "https://www.example.com/blog/my-post",
      "date": "2026-02-20T00:00:00Z",
      "clicks": 61,
      "impressions": 1204,
      "avgClickPosition": 3,
      "avgImpressionPosition": 5
    },
    {
      "page": "https://www.example.com/blog/other-post",
      "date": "2026-02-20T00:00:00Z",
      "clicks": 34,
      "impressions": 812,
      "avgClickPosition": 4,
      "avgImpressionPosition": 6
    }
  ],
  "queriedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "What are my top-performing pages on Bing this month?"

---

## Notes

- Bing's underlying wire format reuses the same shape as [`get_query_stats`](get-query-stats.md)
  (its `Query` field literally holds a page URL for this endpoint) -- this server exposes it as
  `page` here so the output isn't confusing.
- Use [`get_page_query_stats`](get-page-query-stats.md) to see which queries drive traffic to a
  specific page from this list.
- Rows are per page *per day*, so the same page can appear multiple times across the window --
  the example above is truncated to 3 rows for readability.
- There is no date range parameter -- Bing returns a fixed window server-side and this tool
  cannot request a specific period.
