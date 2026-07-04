---
description: Reference for the submit_url_batch MCP tool -- parameters, response format, and example prompts for submitting up to 500 URLs to Bing for indexing in one call.
---

# submit_url_batch

Submit multiple URLs to Bing for indexing in a single request.

---

## Parameters

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `site_url` | string | Yes | The URL of the site |
| `url_list` | string[] | Yes | The URLs to submit -- maximum 500 per call |

---

## Response

```json
{
  "siteUrl": "https://www.example.com/",
  "urlCount": 3,
  "submitted": true,
  "requestedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "Submit these 5 URLs to Bing for indexing: ..."

---

## Notes

- The 500-URL limit is enforced client-side -- exceeding it returns an error before any request is
  sent to Bing.
- Subject to your site's remaining daily/monthly submission quota -- check
  [`get_url_submission_quota`](get-url-submission-quota.md) first for large batches.
- All URLs in the batch must belong to `site_url`.
