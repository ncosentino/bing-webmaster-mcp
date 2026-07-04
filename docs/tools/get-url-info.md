---
description: Reference for the get_url_info MCP tool -- parameters, response format, and example prompts for checking a single URL's index status in Bing.
---

# get_url_info

Get detailed index information for a single URL -- whether it's indexed, its HTTP status,
discovery/crawl dates, and more.

---

## Parameters

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `site_url` | string | Yes | The URL of the site |
| `url` | string | Yes | The specific URL to inspect |

---

## Response

```json
{
  "siteUrl": "https://www.example.com/",
  "url": "https://www.example.com/blog/my-post",
  "isPage": true,
  "httpStatus": 200,
  "documentSize": 48213,
  "anchorCount": 22,
  "discoveryDate": "2025-11-02T08:00:00Z",
  "lastCrawledDate": "2026-02-19T03:00:00Z",
  "totalChildUrlCount": 0,
  "queriedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "Is https://www.example.com/blog/my-post indexed by Bing? When was it last crawled?"

---

## Notes

- `totalChildUrlCount` is non-zero for directory-style URLs with children Bing has discovered.
- Use [`get_url_traffic_info`](get-url-traffic-info.md) alongside this for click/impression data
  on the same URL.
