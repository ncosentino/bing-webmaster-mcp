---
description: Reference for the get_blocked_urls MCP tool -- parameters, response format, and example prompts for listing blocked pages and directories in Bing Webmaster Tools.
---

# get_blocked_urls

List pages and directories currently blocked from Bing's index for a site.

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
  "rowCount": 1,
  "rows": [
    {
      "url": "https://www.example.com/private/",
      "entityType": "Directory",
      "requestType": "FullRemoval",
      "date": "2026-01-05T00:00:00Z"
    }
  ],
  "queriedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "What pages or directories are blocked from Bing on my site?"

---

## Notes

- `entityType` is `Page` or `Directory`.
- `requestType` is `CacheOnly` (removes the cached snapshot but the page can still be indexed) or
  `FullRemoval` (removes the page entirely from Bing's index).
- Use [`add_blocked_url`](add-blocked-url.md) / [`remove_blocked_url`](remove-blocked-url.md) to
  manage this list.
