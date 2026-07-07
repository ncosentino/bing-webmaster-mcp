---
description: Reference for the get_keyword_stats MCP tool -- parameters, response format, and example prompts for historical keyword impression statistics in Bing (market-wide, not site-specific).
---

# get_keyword_stats

Get historical impression statistics for a keyword, market-wide (across all of Bing, not limited
to your site).

---

## Parameters

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `query` | string | Yes | The keyword to look up |
| `country` | string | Yes | Country code (e.g. `us`, `gb`) |
| `language` | string | Yes | Language-country code (e.g. `en-US`, `en-GB`) |

---

## Response

```json
{
  "query": "blazor dependency injection",
  "country": "us",
  "language": "en-US",
  "rowCount": 3,
  "rows": [
    {
      "date": "2026-02-18T00:00:00Z",
      "impressions": 812,
      "broadImpressions": 1985
    },
    {
      "date": "2026-02-19T00:00:00Z",
      "impressions": 935,
      "broadImpressions": 2260
    },
    {
      "date": "2026-02-20T00:00:00Z",
      "impressions": 890,
      "broadImpressions": 2140
    }
  ],
  "queriedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "How many Bing searches happen for 'blazor dependency injection' in the US?"

---

## Notes

- **Unlike every other tool in this list, this one does not take a `site_url` parameter** --
  keyword statistics are market-wide, reflecting overall Bing search volume for that term rather
  than data specific to any of your verified sites.
- `broadImpressions` includes broad/related-match impressions in addition to the exact query.
- Each row carries its own `date`, but there is no request parameter to filter or narrow that
  range -- Bing returns its full available daily history and you filter client-side. The example
  above is truncated to 3 rows for readability; a real response typically spans several weeks.
- For your own site's actual ranking/traffic on a query, use
  [`get_query_stats`](get-query-stats.md) or [`get_query_page_stats`](get-query-page-stats.md)
  instead.
- Unlike [`get_keyword`](get-keyword.md), this endpoint has no `start_date`/`end_date` parameter --
  Bing returns a fixed history window and there is no way to request a specific period.
