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
  "rowCount": 3,
  "rows": [
    {
      "query": "blazor dependency injection",
      "date": "2026-02-19T00:00:00Z",
      "clicks": 21,
      "impressions": 365,
      "avgClickPosition": 4,
      "avgImpressionPosition": 6
    },
    {
      "query": "blazor dependency injection",
      "date": "2026-02-20T00:00:00Z",
      "clicks": 18,
      "impressions": 340,
      "avgClickPosition": 4,
      "avgImpressionPosition": 6
    },
    {
      "query": "asp.net minimal apis",
      "date": "2026-02-20T00:00:00Z",
      "clicks": 15,
      "impressions": 298,
      "avgClickPosition": 3,
      "avgImpressionPosition": 5
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
- Rows are per query *per day*, so the same query can appear multiple times across the window --
  the example above is truncated to 3 rows for readability.
- There is no date range parameter -- Bing returns a fixed window server-side and this tool
  cannot request a specific period.
