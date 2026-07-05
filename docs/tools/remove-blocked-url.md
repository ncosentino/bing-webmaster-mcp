---
description: Reference for the remove_blocked_url MCP tool -- parameters, response format, and example prompts for unblocking a page or directory in Bing Webmaster Tools.
---

# remove_blocked_url

Unblock a previously blocked page or directory.

---

## Parameters

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `site_url` | string | Yes | -- | The URL of the site |
| `url` | string | Yes | -- | The page or directory URL to unblock |
| `entity_type` | string | No | `Page` | `Page` or `Directory` -- should match the original block |
| `request_type` | string | No | `FullRemoval` | `CacheOnly` or `FullRemoval` -- should match the original block |

---

## Response

```json
{
  "siteUrl": "https://www.example.com/",
  "url": "https://www.example.com/staging/",
  "entityType": "Directory",
  "requestType": "FullRemoval",
  "success": true,
  "requestedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "Unblock https://www.example.com/staging/ on Bing."

---

## Notes

- Use [`get_blocked_urls`](get-blocked-urls.md) first to confirm the exact `entity_type` and
  `request_type` the block was originally created with.
- **For `Directory` blocks, use the exact URL `get_blocked_urls` reports** (Bing normalizes
  directory URLs with a trailing slash on storage) -- passing the original un-normalized URL will
  report success but silently not match the stored block.
