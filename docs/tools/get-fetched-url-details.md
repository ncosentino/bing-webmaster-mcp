---
description: Reference for the get_fetched_url_details MCP tool -- parameters, response format, and example prompts for inspecting the full result of a Bing fetch request.
---

# get_fetched_url_details

Get the full result of a previously requested [`fetch_url`](fetch-url.md) call.

---

## Parameters

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `site_url` | string | Yes | The URL of the site |
| `url` | string | Yes | The specific URL to get fetch details for |

---

## Response

```json
{
  "siteUrl": "https://www.example.com/",
  "url": "https://www.example.com/blog/new-post",
  "date": "2026-02-21T18:56:00Z",
  "status": "Success",
  "headers": "HTTP/1.1 200 OK\nContent-Type: text/html; charset=utf-8\n...",
  "document": "<!DOCTYPE html>...",
  "queriedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "Show me exactly what Bing saw when it fetched my new blog post."

---

## Notes

- `document` can be large (the full fetched HTML) -- expect a sizable response for content-heavy
  pages.
- Useful for diagnosing rendering issues Bing's crawler may be encountering.
