---
description: Reference for the submit_url MCP tool -- parameters, response format, and example prompts for submitting a single URL to Bing for indexing.
---

# submit_url

Submit a single URL to Bing for crawling and indexing.

---

## Parameters

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `site_url` | string | Yes | The URL of the site |
| `url` | string | Yes | The specific URL to submit |

---

## Response

```json
{
  "siteUrl": "https://www.example.com/",
  "url": "https://www.example.com/new-post",
  "submitted": true,
  "requestedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "Submit https://www.example.com/new-post to Bing for indexing."

---

## Notes

- Subject to your site's daily/monthly submission quota -- check
  [`get_url_submission_quota`](get-url-submission-quota.md) first if you're submitting frequently.
- For more than one URL, use [`submit_url_batch`](submit-url-batch.md) instead of calling this
  tool repeatedly.
- Consider [`submit_url_indexnow`](submit-url-indexnow.md) as a faster, quota-independent
  alternative once you've configured an IndexNow key.
