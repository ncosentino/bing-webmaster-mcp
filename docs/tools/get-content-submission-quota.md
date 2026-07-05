---
description: Reference for the get_content_submission_quota MCP tool -- parameters, response format, and example prompts for checking Bing's Content Submission API quota.
---

# get_content_submission_quota

Check the remaining daily and monthly quota for [`submit_content`](submit-content.md).

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
  "dailyQuota": 10,
  "monthlyQuota": 300,
  "queriedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "How much Content Submission API quota do I have left today?"

---

## Notes

- This is a **separate quota** from [`get_url_submission_quota`](get-url-submission-quota.md) --
  it applies specifically to [`submit_content`](submit-content.md).
