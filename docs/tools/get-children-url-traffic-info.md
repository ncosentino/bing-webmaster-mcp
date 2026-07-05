---
description: Reference for the get_children_url_traffic_info MCP tool -- parameters, response format, and example prompts for traffic across a directory of pages in Bing.
---

# get_children_url_traffic_info

Get clicks and impressions for every child URL under a directory.

---

## Parameters

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `site_url` | string | Yes | -- | The URL of the site |
| `url` | string | Yes | -- | The parent/directory URL |
| `page` | int | No | `0` | Page number of results (0-based) |

---

## Response

```json
{
  "siteUrl": "https://www.example.com/",
  "url": "https://www.example.com/blog/",
  "page": 0,
  "rowCount": 1,
  "rows": [
    {
      "url": "https://www.example.com/blog/my-post",
      "isPage": true,
      "clicks": 61,
      "impressions": 1204
    }
  ],
  "queriedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "Which pages under /blog/ get the most Bing traffic?"

---

## Notes

- Results are paginated -- increment `page` for more.
- For index status instead of traffic, see [`get_children_url_info`](get-children-url-info.md).
