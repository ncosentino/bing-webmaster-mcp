---
description: Reference for the get_link_counts MCP tool -- parameters, response format, and example prompts for finding which pages have inbound links in Bing.
---

# get_link_counts

Get a list of your site's pages that have inbound links, with counts.

---

## Parameters

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `site_url` | string | Yes | -- | The URL of the site |
| `page` | int | No | `0` | Page number of results (0-based) |

---

## Response

```json
{
  "siteUrl": "https://www.example.com/",
  "links": [
    {
      "url": "https://www.example.com/blog/my-post",
      "count": 14
    }
  ],
  "totalPages": 1,
  "queriedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "Which pages on my site have the most inbound links?"

---

## Notes

- Results are paginated -- check `totalPages` and increment `page` to fetch more.
- Use [`get_url_links`](get-url-links.md) to see the actual linking pages and anchor text for any
  specific URL from this list.
