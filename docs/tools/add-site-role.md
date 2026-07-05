---
description: Reference for the add_site_role MCP tool -- parameters, response format, and example prompts for delegating site access in Bing Webmaster Tools.
---

# add_site_role

Delegate access to a site to another user.

---

## Parameters

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `site_url` | string | Yes | -- | The URL of your site |
| `delegated_url` | string | Yes | -- | The URL being delegated (usually the same as `site_url`) |
| `user_email` | string | Yes | -- | The email of the user to delegate access to |
| `authentication_code` | string | Yes | -- | The site's authentication code (from [`list_sites`](list-sites.md)) |
| `is_administrator` | boolean | No | `false` | Whether to grant administrator privileges |
| `is_read_only` | boolean | No | `true` | Whether access should be read-only |

---

## Response

```json
{
  "siteUrl": "https://www.example.com/",
  "delegatedUrl": "https://www.example.com/",
  "userEmail": "colleague@example.com",
  "isAdministrator": false,
  "isReadOnly": true,
  "success": true,
  "requestedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "Give colleague@example.com read-only access to my Bing Webmaster site."

---

## Notes

- The recipient must accept the delegation from their own Bing Webmaster account before access is
  active.
- Use [`get_site_roles`](get-site-roles.md) afterward to confirm the delegation appears (it may
  show as pending until accepted).
- Use [`remove_site_role`](remove-site-role.md) to revoke access later.
- **Live-testing note:** `get_site_roles` (read) is confirmed working against a real account.
  This tool's write path is unit-tested only -- live-testing it safely requires delegating to a
  second real email address, which wasn't available during testing.
