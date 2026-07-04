---
description: Reference for the get_rank_and_traffic_stats MCP tool -- parameters, response format, and example prompts for daily clicks and impressions over time in Bing.
---

# get_rank_and_traffic_stats

Get the headline daily clicks and impressions trend for a site.

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
      "date": "2026-02-20T00:00:00Z",
      "clicks": 214,
      "impressions": 5430
    }
  ],
  "queriedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "What's my overall Bing clicks and impressions trend over the last month?"

> "Has my organic traffic from Bing gone up or down recently?"

---

## Notes

- This is the site-wide daily total -- for a breakdown by query or page, see
  [`get_query_stats`](get-query-stats.md) and [`get_page_stats`](get-page-stats.md).
- One row per day.
