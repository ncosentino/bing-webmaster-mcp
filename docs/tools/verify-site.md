---
description: Reference for the verify_site MCP tool -- parameters, response format, and example prompts for verifying site ownership in Bing Webmaster Tools.
---

# verify_site

Attempt to verify ownership of a site already added to your Bing Webmaster Tools account.

---

## Parameters

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `site_url` | string | Yes | The URL of the site to verify |

---

## Response

```json
{
  "siteUrl": "https://www.example.com/",
  "verified": true,
  "checkedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "Verify ownership of https://www.example.com/ in Bing Webmaster Tools."

> "Did my Bing site verification go through?"

---

## Notes

- Verification only succeeds if you've already completed one of Bing's verification methods
  outside this tool (DNS TXT record, meta tag, XML file upload, etc.) -- this tool just triggers
  Bing's check, it doesn't perform the verification itself.
- `verified: false` means the check ran but didn't find proof of ownership yet -- it's not an
  error.
