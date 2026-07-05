---
description: Reference for the add_blocked_url MCP tool -- parameters, response format, and example prompts for blocking a page or directory in Bing Webmaster Tools.
---

# add_blocked_url

Block a page or directory from Bing's index.

---

## Parameters

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `site_url` | string | Yes | -- | The URL of the site |
| `url` | string | Yes | -- | The page or directory URL to block |
| `entity_type` | string | No | `Page` | `Page` or `Directory` |
| `request_type` | string | No | `CacheOnly` | `CacheOnly` or `FullRemoval` |

---

## Response

```json
{
  "siteUrl": "https://www.example.com/",
  "url": "https://www.example.com/staging/",
  "entityType": "Directory",
  "requestType": "CacheOnly",
  "success": true,
  "requestedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "Block https://www.example.com/staging/ from Bing's index."

---

## Notes

- **Reversible** via [`remove_blocked_url`](remove-blocked-url.md).
- Use `request_type: "FullRemoval"` for a stronger block that removes the page from the index
  entirely, not just its cached snapshot.
- Use [`get_blocked_urls`](get-blocked-urls.md) afterward to confirm the block was applied.
- **Directory blocks get normalized.** Live testing confirmed Bing appends a trailing slash to
  `Directory`-type URLs when storing the block (e.g. `.../staging` is stored as `.../staging/`).
  When removing a directory block, use the exact normalized form reported by
  [`get_blocked_urls`](get-blocked-urls.md) -- passing the original un-normalized URL to
  [`remove_blocked_url`](remove-blocked-url.md) will report success but silently not match
  anything.
