---
description: Reference for the submit_sitemap MCP tool -- parameters, response format, and example prompts for submitting a sitemap to Bing Webmaster Tools.
---

# submit_sitemap

Submit a sitemap URL to Bing Webmaster Tools for a site.

---

## Parameters

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `site_url` | string | Yes | The URL of the site |
| `feed_url` | string | Yes | The full URL of the sitemap to submit (e.g. `https://www.example.com/sitemap.xml`) |

---

## Response

```json
{
  "siteUrl": "https://www.example.com/",
  "feedUrl": "https://www.example.com/sitemap.xml",
  "submitted": true,
  "requestedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "Submit my sitemap at https://www.example.com/sitemap.xml to Bing."

---

## Notes

- The sitemap must already be publicly accessible at the given URL before submitting.
- Use [`list_sitemaps`](list-sitemaps.md) afterward to confirm it was accepted.
