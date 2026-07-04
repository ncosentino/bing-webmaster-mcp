---
description: Reference for the get_crawl_issues MCP tool -- parameters, response format, and example prompts for listing URLs with crawl issues in Bing Webmaster Tools.
---

# get_crawl_issues

List URLs with crawl issues for a site -- pages Bing's crawler had trouble accessing or
processing.

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
  "issues": [
    {
      "url": "https://www.example.com/old-page",
      "httpCode": 404,
      "issues": ["Code4xx"],
      "inLinks": 3
    }
  ],
  "queriedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "What crawl issues does my site currently have?"

> "Which of my pages are returning 404s to Bing's crawler?"

---

## Notes

- `issues` decodes Bing's underlying bitflag value into readable labels (e.g. `Code301`,
  `Code4xx`, `BlockedByRobotsTxt`, `ContainsMalware`, `DnsErrors`, `TimeOutErrors`). A URL can have
  more than one issue at once.
- An empty `issues` array means no crawl issues were found -- not an error.
