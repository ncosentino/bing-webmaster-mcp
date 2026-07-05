---
description: Reference for the get_related_keywords MCP tool -- parameters, response format, and example prompts for discovering related keywords and their market-wide Bing impressions.
---

# get_related_keywords

Discover keywords related to a seed term, with market-wide Bing impression data for each.

---

## Parameters

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `query` | string | Yes | The seed keyword |
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
  "rowCount": 2,
  "rows": [
    { "query": "blazor di container", "impressions": 210, "broadImpressions": 640 },
    { "query": "blazor service lifetime", "impressions": 180, "broadImpressions": 520 }
  ],
  "queriedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "What keywords are related to 'blazor dependency injection' that people search for on Bing?"

---

## Notes

- **Unlike most tools, this one does not take a `site_url`** -- it's market-wide keyword
  discovery, not tied to your site's data.
- Useful for content ideation -- pair with [`get_keyword`](get-keyword.md) to check volume for any
  specific term this surfaces.
