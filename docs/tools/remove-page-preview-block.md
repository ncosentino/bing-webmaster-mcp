---
description: Reference for the remove_page_preview_block MCP tool -- parameters, response format, and example prompts for removing a page preview block in Bing Webmaster Tools.
---

# remove_page_preview_block

Remove a page preview block from a URL, added previously via
[`add_page_preview_block`](add-page-preview-block.md).

---

## Parameters

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `site_url` | string | Yes | The URL of the site |
| `url` | string | Yes | The page URL whose preview block should be removed |

---

## Response

```json
{
  "siteUrl": "https://www.example.com/",
  "url": "https://www.example.com/private-preview",
  "success": true,
  "requestedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "Remove the preview block from https://www.example.com/private-preview."

---

## Notes

- Unlike adding a block, no `reason` is needed to remove one.
- Use [`get_active_page_preview_blocks`](get-active-page-preview-blocks.md) afterward to confirm
  removal.
