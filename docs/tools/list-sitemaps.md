---
description: Reference for the list_sitemaps MCP tool -- parameters, response format, and example prompts for listing submitted sitemaps in Bing Webmaster Tools.
---

# list_sitemaps

List sitemaps submitted to Bing Webmaster Tools for a site.

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
  "feeds": [
    {
      "url": "https://www.example.com/sitemap.xml",
      "type": "Sitemap",
      "status": "Ok",
      "compressed": false,
      "fileSize": 4820,
      "urlCount": 137,
      "submitted": "2026-01-15T10:00:00Z",
      "lastCrawled": "2026-02-20T04:00:00Z"
    }
  ],
  "queriedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "Show me all submitted sitemaps for my site and when they were last crawled."

---

## Notes

- Use [`get_sitemap_details`](get-sitemap-details.md) for more detail on a specific sitemap
  (useful for sitemap index files with multiple constituent feeds).
- Use [`submit_sitemap`](submit-sitemap.md) to submit a new one.
