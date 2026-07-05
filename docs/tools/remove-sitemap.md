---
description: Reference for the remove_sitemap MCP tool -- parameters, response format, and example prompts for removing a submitted sitemap from Bing Webmaster Tools.
---

# remove_sitemap

Remove a previously submitted sitemap.

---

## Parameters

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `site_url` | string | Yes | The URL of the site |
| `feed_url` | string | Yes | The URL of the sitemap to remove |

---

## Response

```json
{
  "siteUrl": "https://www.example.com/",
  "feedUrl": "https://www.example.com/old-sitemap.xml",
  "success": true,
  "requestedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "Remove the old sitemap at https://www.example.com/old-sitemap.xml from Bing."

---

## Notes

- **Reversible** -- re-submit with [`submit_sitemap`](submit-sitemap.md) if needed.
- Use [`list_sitemaps`](list-sitemaps.md) afterward to confirm removal.
