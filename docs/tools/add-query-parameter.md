---
description: Reference for the add_query_parameter MCP tool -- parameters, response format, and example prompts for adding a URL normalization query parameter in Bing Webmaster Tools.
---

# add_query_parameter

Add a URL normalization query parameter for a site -- tells Bing this parameter doesn't change
page content (e.g. a tracking or sorting parameter), so URLs differing only by it are treated as
the same page.

---

## Parameters

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `site_url` | string | Yes | The URL of the site |
| `query_parameter` | string | Yes | The query parameter name to add, for example `utm_campaign` |

---

## Response

```json
{
  "siteUrl": "https://www.example.com/",
  "queryParameter": "utm_campaign",
  "success": true,
  "requestedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "Tell Bing to ignore the utm_campaign query parameter on my site."

---

## Notes

- **Reversible** via [`remove_query_parameter`](remove-query-parameter.md).
- Use [`get_query_parameters`](get-query-parameters.md) afterward to confirm it was added.
- To temporarily disable a parameter without removing it, use
  [`enable_disable_query_parameter`](enable-disable-query-parameter.md).
