---
description: Reference for the list_sites MCP tool -- response format and example prompts for listing sites in your Bing Webmaster Tools account.
---

# list_sites

List all sites in your Bing Webmaster Tools account.

---

## Parameters

No parameters required.

---

## Response

```json
{
  "sites": [
    {
      "url": "https://www.example.com/",
      "isVerified": true,
      "dnsVerificationCode": "abc123...",
      "authenticationCode": "def456..."
    }
  ],
  "queriedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "What sites do I have in my Bing Webmaster account?"

> "Is my site verified in Bing Webmaster Tools?"

---

## Notes

- Use this tool first to confirm the exact site URL to pass to every other tool.
- A site can exist in your account without being verified (`isVerified: false`) -- see
  [`verify_site`](verify-site.md).
