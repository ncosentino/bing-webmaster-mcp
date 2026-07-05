---
description: Reference for the remove_site_role MCP tool -- parameters, response format, and example prompts for revoking delegated site access in Bing Webmaster Tools.
---

# remove_site_role

Revoke a user's delegated access to a site.

---

## Parameters

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `site_url` | string | Yes | The URL of the site |
| `email` | string | Yes | The email of the user whose access should be revoked |
| `role` | string | Yes | The role to remove: `Administrator`, `ReadOnly`, or `ReadWrite` (must match what [`get_site_roles`](get-site-roles.md) reports) |

---

## Response

```json
{
  "siteUrl": "https://www.example.com/",
  "email": "colleague@example.com",
  "role": "ReadOnly",
  "success": true,
  "requestedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "Revoke colleague@example.com's access to my Bing Webmaster site."

---

## Notes

- Use [`get_site_roles`](get-site-roles.md) first to confirm the exact email and role to remove.
- Bing's exact required fields for matching a role removal aren't fully documented publicly --
  this tool submits the minimal identifying information (site, email, role). If removal doesn't
  take effect, please open an issue with the (redacted) details.
