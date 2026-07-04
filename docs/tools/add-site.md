---
description: Reference for the add_site MCP tool -- parameters, response format, and example prompts for adding a new site to your Bing Webmaster Tools account.
---

# add_site

Add a new site to your Bing Webmaster Tools account.

---

## Parameters

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `site_url` | string | Yes | The URL of the site to add (e.g. `https://www.example.com/`) |

---

## Response

```json
{
  "siteUrl": "https://www.example.com/",
  "added": true,
  "requestedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "Add https://www.example.com/ to my Bing Webmaster account."

---

## Notes

- Adding a site does not verify it -- call [`verify_site`](verify-site.md) next, after completing
  whichever verification method Bing offers for the site (DNS record, meta tag, XML file, etc.,
  configured outside this MCP server).
- Use [`list_sites`](list-sites.md) afterward to confirm the site was added.
