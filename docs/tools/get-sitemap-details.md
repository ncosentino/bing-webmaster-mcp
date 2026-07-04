---
description: Reference for the get_sitemap_details MCP tool -- parameters, response format, and example prompts for inspecting a specific submitted sitemap in Bing Webmaster Tools.
---

# get_sitemap_details

Get details for a specific sitemap already submitted to Bing Webmaster Tools.

---

## Parameters

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `site_url` | string | Yes | The URL of the site |
| `feed_url` | string | Yes | The URL of the submitted sitemap feed |

---

## Response

```json
{
  "siteUrl": "https://www.example.com/",
  "feedUrl": "https://www.example.com/sitemap-index.xml",
  "feed": {
    "url": "https://www.example.com/sitemap-index.xml",
    "type": "Sitemap",
    "status": "Ok",
    "urlCount": 3
  },
  "queriedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "Give me the details on my sitemap index file."

---

## Notes

- Most useful for sitemap **index** files, which reference multiple constituent sitemaps --
  [`list_sitemaps`](list-sitemaps.md) alone may not show the full breakdown.
- Bing's exact response shape for this endpoint isn't published in detail; the server parses it
  defensively and may not surface every possible field.
