---
description: Reference for the get_keyword MCP tool -- parameters, response format, and example prompts for market-wide keyword impressions on Bing for a specific period.
---

# get_keyword

Get market-wide impression data for a single keyword over a specific date range.

---

## Parameters

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `query` | string | Yes | The keyword to look up |
| `country` | string | Yes | Country code (e.g. `us`, `gb`) |
| `language` | string | Yes | Language-country code (e.g. `en-US`, `en-GB`) |
| `start_date` | string | Yes | Start date in `YYYY-MM-DD` format |
| `end_date` | string | Yes | End date in `YYYY-MM-DD` format |

---

## Response

```json
{
  "query": "blazor dependency injection",
  "country": "us",
  "language": "en-US",
  "startDate": "2026-01-01",
  "endDate": "2026-01-31",
  "found": true,
  "impressions": 2140,
  "broadImpressions": 5830,
  "queriedAt": "2026-02-21T19:00:00Z"
}
```

If Bing has no data for the keyword/period/market combination:

```json
{
  "query": "an extremely obscure phrase",
  "country": "us",
  "language": "en-US",
  "startDate": "2026-01-01",
  "endDate": "2026-01-31",
  "found": false,
  "impressions": 0,
  "broadImpressions": 0,
  "queriedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "How many Bing searches happened for 'blazor dependency injection' in the US last month?"

---

## Notes

- **Unlike most tools, this one does not take a `site_url`** -- it reflects overall Bing search
  volume for the term, not your site's specific ranking or traffic.
- `found: false` is not an error -- it just means no data exists for that exact combination.
- For a summed set of related terms instead of one exact keyword, see
  [`get_related_keywords`](get-related-keywords.md). For your own site's actual performance on
  this keyword, see [`get_query_stats`](get-query-stats.md) or
  [`get_query_traffic_stats`](get-query-traffic-stats.md).
