---
description: Reference for the remove_site MCP tool -- parameters, response format, and example prompts for removing a site from Bing Webmaster Tools.
---

# remove_site

Remove a site from your Bing Webmaster Tools account.

---

## Parameters

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `site_url` | string | Yes | The URL of the site to remove |

---

## Response

```json
{
  "siteUrl": "https://www.example.com/",
  "success": true,
  "requestedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "Remove https://old-project.example.com/ from my Bing Webmaster account."

---

## Notes

- **Destructive.** This removes all history Bing has accumulated for the site from your account
  view. Use with care -- always confirm the exact `site_url` with [`list_sites`](list-sites.md)
  first.
- The site can be re-added later with [`add_site`](add-site.md), but re-verification will be
  required and prior crawl/traffic history in the UI may not carry over identically.
