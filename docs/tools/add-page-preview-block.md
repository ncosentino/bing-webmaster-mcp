---
description: Reference for the add_page_preview_block MCP tool -- parameters, response format, and example prompts for blocking a page preview in Bing Webmaster Tools.
---

# add_page_preview_block

Block a URL from showing a preview (thumbnail or content snippet) in Bing search results. The page
itself can still be indexed and ranked -- only the visual/content preview is suppressed.

---

## Parameters

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `site_url` | string | Yes | The URL of the site |
| `url` | string | Yes | The page URL to block from previews |
| `reason` | string | Yes | `AdultContent`, `Copyright`, `IllegalContent`, or `Other` |

---

## Response

```json
{
  "siteUrl": "https://www.example.com/",
  "url": "https://www.example.com/private-preview",
  "reason": "Other",
  "success": true,
  "requestedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "Block the page preview for https://www.example.com/private-preview -- reason: Other."

---

## Notes

- **Reversible** via [`remove_page_preview_block`](remove-page-preview-block.md).
- Unlike [`add_blocked_url`](add-blocked-url.md), `reason` has no default -- it must always be
  supplied explicitly since there's no sensible assumption for why a preview should be blocked.
- Use [`get_active_page_preview_blocks`](get-active-page-preview-blocks.md) afterward to confirm.
