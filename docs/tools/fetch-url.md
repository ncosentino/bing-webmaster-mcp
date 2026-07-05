---
description: Reference for the fetch_url MCP tool -- parameters, response format, and example prompts for requesting an immediate Bing crawl of a URL.
---

# fetch_url

Request Bing to fetch (crawl) a specific URL immediately, similar to "Fetch as Bingbot."

---

## Parameters

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `site_url` | string | Yes | The URL of the site |
| `url` | string | Yes | The specific URL to fetch |

---

## Response

```json
{
  "siteUrl": "https://www.example.com/",
  "url": "https://www.example.com/blog/new-post",
  "success": true,
  "requestedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "Ask Bing to fetch my new blog post right now."

---

## Notes

- This requests an immediate crawl -- it's distinct from [`submit_url`](submit-url.md), which
  requests indexing. Fetching lets you see how Bing renders/reads the page without necessarily
  affecting the index.
- Use [`list_fetched_urls`](list-fetched-urls.md) and
  [`get_fetched_url_details`](get-fetched-url-details.md) to check the fetch result.
