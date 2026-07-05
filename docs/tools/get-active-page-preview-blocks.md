---
description: Reference for the get_active_page_preview_blocks MCP tool -- parameters, response format, and example prompts for listing page preview blocks in Bing Webmaster Tools.
---

# get_active_page_preview_blocks

List pages currently blocked from showing a preview (e.g. a thumbnail or content snippet) in Bing
search results, along with the reason each was blocked.

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
  "blocks": [
    {
      "url": "https://www.example.com/private-preview",
      "blockReason": "Other",
      "submitDate": "2026-02-21T19:00:00Z"
    }
  ],
  "queriedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "Which of my pages have preview blocks active in Bing search results, and why?"

---

## Notes

- `blockReason` is one of `AdultContent`, `Copyright`, `IllegalContent`, or `Other`.
- Bing's real record also includes an internal action/refresh-reason state used to track review
  status; those aren't surfaced here since their exact meaning isn't independently confirmed.
- Use [`add_page_preview_block`](add-page-preview-block.md) /
  [`remove_page_preview_block`](remove-page-preview-block.md) to manage this list.
