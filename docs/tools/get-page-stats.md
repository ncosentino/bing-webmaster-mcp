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
  "stats": [
    {
      "page": "https://www.example.com/blog/my-post",
      "date": "2026-02-20T00:00:00Z",
      "clicks": 61,
      "impressions": 1204,
      "avgClickPosition": 3,
      "avgImpressionPosition": 5
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
