---
description: Reference for the add_connected_page MCP tool -- parameters, response format, and example prompts for declaring a connected page in Bing Webmaster Tools.
---

# add_connected_page

Declare that a page is connected to (syndicates or mirrors content from) your site, identified by
a master URL. Used for cross-domain content syndication scenarios so Bing understands the
relationship rather than treating the pages as unrelated duplicates.

---

## Parameters

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `site_url` | string | Yes | The URL of the site |
| `master_url` | string | Yes | The master URL this site's content is connected to |

---

## Response

```json
{
  "siteUrl": "https://www.example.com/",
  "masterUrl": "https://mirror.example.net/",
  "success": true,
  "requestedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "Tell Bing that https://www.example.com/ is connected to https://mirror.example.net/."

---

## Notes

- There is no `remove_connected_page` tool -- Bing's API doesn't expose a removal method for this
  relationship.
- Use [`get_connected_pages`](get-connected-pages.md) afterward to confirm the connection and check
  its verification status.
