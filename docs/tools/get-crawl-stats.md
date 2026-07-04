---
description: Reference for the get_crawl_stats MCP tool -- parameters, response format, and example prompts for daily crawl statistics in Bing Webmaster Tools.
---

# get_crawl_stats

Get daily crawl statistics for a site -- crawled page counts, HTTP status code breakdowns, and
index totals over time.

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
  "stats": [
    {
      "date": "2026-02-20T00:00:00Z",
      "crawledPages": 412,
      "crawlErrors": 2,
      "inIndex": 389,
      "inLinks": 1204,
      "code2xx": 405,
      "code301": 4,
      "code302": 1,
      "code4xx": 2,
      "code5xx": 0,
      "allOtherCodes": 0,
      "blockedByRobotsTxt": 0,
      "containsMalware": 0
    }
  ],
  "queriedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "How many pages did Bing crawl on my site yesterday, and were there any errors?"

> "Show me my crawl stats trend over the last few weeks."

---

## Notes

- One row per day.
- Use alongside [`get_crawl_issues`](get-crawl-issues.md) to see the specific URLs behind any
  error counts.
