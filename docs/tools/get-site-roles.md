---
description: Reference for the get_site_roles MCP tool -- parameters, response format, and example prompts for listing delegated user access to a Bing Webmaster Tools site.
---

# get_site_roles

List users with delegated access to a site.

---

## Parameters

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `site_url` | string | Yes | -- | The URL of the site |
| `include_all_subdomains` | boolean | No | `false` | Whether to include roles delegated for all subdomains |

---

## Response

```json
{
  "siteUrl": "https://www.example.com/",
  "includeAllSubdomains": false,
  "rowCount": 1,
  "rows": [
    {
      "email": "colleague@example.com",
      "role": "ReadOnly",
      "site": "https://www.example.com/",
      "verificationSite": "https://www.example.com/",
      "expired": false,
      "delegatorEmail": "you@example.com",
      "delegatedCode": null,
      "delegatedCodeOwnerEmail": null,
      "date": "2026-01-10T00:00:00Z"
    }
  ],
  "queriedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "Who has access to my Bing Webmaster site, and what permission level do they have?"

---

## Notes

- `role` is one of `Administrator`, `ReadOnly`, or `ReadWrite`.
- Use [`add_site_role`](add-site-role.md) to delegate access, and
  [`remove_site_role`](remove-site-role.md) to revoke it.
