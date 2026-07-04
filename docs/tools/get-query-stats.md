---
description: Reference for the get_query_stats MCP tool -- parameters, response format, and example prompts for top search queries by clicks and impressions in Bing.
---

# get_query_stats

Get traffic statistics for your site's top search queries.

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
      "query": "blazor dependency injection",
      "date": "2026-02-20T00:00:00Z",
      "clicks": 18,
      "impressions": 340,
      "avgClickPosition": 4,
      "avgImpressionPosition": 6
    }
  ],
  "queriedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "What are my top search queries on Bing this month?"

> "Which queries have high impressions but low clicks -- possible CTR opportunities?"

---

## Notes

- Use [`get_query_page_stats`](get-query-page-stats.md) to see which of your pages rank for a
  specific query from this list.
- `avgClickPosition`/`avgImpressionPosition` are integers (Bing does not report fractional
  average position the way Google Search Console does).
