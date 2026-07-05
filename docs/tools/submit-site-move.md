---
description: Reference for the submit_site_move MCP tool -- parameters, response format, and example prompts for notifying Bing of a site migration.
---

# submit_site_move

Notify Bing that a site is moving from one URL to another (domain migration, protocol change, etc.).

---

## Parameters

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `site_url` | string | Yes | -- | The URL of the site (as registered in Bing Webmaster Tools) |
| `source_url` | string | Yes | -- | The old URL moving away from |
| `target_url` | string | Yes | -- | The new URL moving to |
| `move_type` | string | No | `Local` | `Local` or `Global` |
| `move_scope` | string | No | `Domain` | `Domain`, `Host`, or `Directory` |

---

## Response

```json
{
  "siteUrl": "https://old.example.com/",
  "sourceUrl": "https://old.example.com/",
  "targetUrl": "https://www.example.com/",
  "moveType": "Local",
  "moveScope": "Domain",
  "success": true,
  "requestedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "Tell Bing my site is moving from old.example.com to www.example.com."

---

## Notes

- **Consequential, not easily reversible.** This tells Bing to treat one site as superseded by
  another, which can affect ranking signal transfer and how both URLs are treated going forward.
  Both `source_url` and `target_url` should already be verified sites in your account, and proper
  redirects should already be in place before submitting.
- Use [`get_site_moves`](get-site-moves.md) to check status afterward -- note that endpoint
  returned an HTML 404 in live testing (see the [Roadmap](../roadmap.md)), so status may not be
  retrievable even after a successful submission.
- **Not live-tested** against a real account given the consequential/hard-to-reverse nature of
  this operation -- covered by unit tests only.
