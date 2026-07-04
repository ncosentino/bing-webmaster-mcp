---
description: Reference for the get_url_submission_quota MCP tool -- parameters, response format, and example prompts for checking your Bing URL submission quota.
---

# get_url_submission_quota

Check the remaining daily and monthly URL submission quota for a site.

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

> "How many URLs can I still submit to Bing today for my site?"

---

## Notes

- Check this before a large [`submit_url_batch`](submit-url-batch.md) call to avoid a partial
  failure from exceeding quota.
- Quota is per-site, not account-wide.
- [`submit_url_indexnow`](submit-url-indexnow.md) is not subject to this quota.
