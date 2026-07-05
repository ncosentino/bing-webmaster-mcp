---
description: Reference for the get_children_url_info MCP tool -- parameters, response format, and example prompts for index details across a directory of pages in Bing.
---

# get_children_url_info

Get index details for every child URL under a directory.

---

## Parameters

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `site_url` | string | Yes | -- | The URL of the site |
| `url` | string | Yes | -- | The parent/directory URL |
| `page` | int | No | `0` | Page number of results (0-based) |
| `crawl_date_filter` | string | No | `Any` | `Any`, `LastWeek`, `LastTwoWeeks`, `LastThreeWeeks` |
| `discovered_date_filter` | string | No | `Any` | `Any`, `LastWeek`, `LastMonth` |
| `doc_flags_filter` | string | No | `Any` | `Any`, `IsBlockedByRobotsTxt`, `IsMalware` |
| `http_code_filter` | string | No | `Any` | `Any`, `Code2xx`, `Code3xx`, `Code301`, `Code302`, `Code4xx`, `Code5xx`, `AllOthers` |

---

## Response

```json
{
  "siteUrl": "https://www.example.com/",
  "url": "https://www.example.com/blog/",
  "page": 0,
  "crawlDateFilter": "Any",
  "discoveredDateFilter": "Any",
  "docFlagsFilter": "Any",
  "httpCodeFilter": "Any",
  "rowCount": 1,
  "rows": [
    {
      "url": "https://www.example.com/blog/my-post",
      "isPage": true,
      "httpStatus": 200,
      "documentSize": 48213,
      "anchorCount": 22,
      "discoveryDate": "2025-11-02T08:00:00Z",
      "lastCrawledDate": "2026-02-19T03:00:00Z",
      "totalChildUrlCount": 0
    }
  ],
  "queriedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "Show me the index status of every page under /blog/ on my site."

> "Which pages under /blog/ are blocked by robots.txt?"

---

## Notes

- Results are paginated -- increment `page` for more.
- The filters let you narrow to, e.g., only recently-discovered pages or only pages currently
  returning a 4xx/5xx.
- For traffic instead of index status, see
  [`get_children_url_traffic_info`](get-children-url-traffic-info.md).
