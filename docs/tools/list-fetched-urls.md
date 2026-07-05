---
description: Reference for the list_fetched_urls MCP tool -- parameters, response format, and example prompts for listing URLs submitted for Bing fetch requests.
---

# list_fetched_urls

List URLs that have been submitted via [`fetch_url`](fetch-url.md), with their status.

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
  "rows": [
    {
      "url": "https://www.example.com/blog/new-post",
      "date": "2026-02-21T18:55:00Z",
      "fetched": true,
      "expired": false
    }
  ],
  "queriedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "What URLs have I asked Bing to fetch recently, and did they complete?"

---

## Notes

- Use [`get_fetched_url_details`](get-fetched-url-details.md) for the full fetch result
  (headers, document, status) of a specific URL from this list.
