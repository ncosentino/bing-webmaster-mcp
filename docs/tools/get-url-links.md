---
description: Reference for the get_url_links MCP tool -- parameters, response format, and example prompts for finding inbound links to a specific URL in Bing.
---

# get_url_links

Get inbound links pointing to a specific URL.

---

## Parameters

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `site_url` | string | Yes | -- | The URL of the site |
| `link` | string | Yes | -- | The specific URL to get inbound links for |
| `page` | int | No | `0` | Page number of results (0-based) |

---

## Response

```json
{
  "siteUrl": "https://www.example.com/",
  "link": "https://www.example.com/blog/my-post",
  "details": [
    {
      "url": "https://another-site.com/references",
      "anchorText": "this great article"
    }
  ],
  "totalPages": 1,
  "queriedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "Who's linking to my blog post about dependency injection?"

---

## Notes

- Results are paginated -- check `totalPages` and increment `page` to fetch more.
- Use [`get_link_counts`](get-link-counts.md) first to see which of your pages have the most
  inbound links before drilling into a specific one with this tool.
