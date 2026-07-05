---
description: Reference for the remove_query_parameter MCP tool -- parameters, response format, and example prompts for removing a URL normalization query parameter in Bing Webmaster Tools.
---

# remove_query_parameter

Remove a URL normalization query parameter from a site, added previously via
[`add_query_parameter`](add-query-parameter.md).

---

## Parameters

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `site_url` | string | Yes | The URL of the site |
| `query_parameter` | string | Yes | The query parameter name to remove |

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

> "Stop telling Bing to ignore the utm_campaign query parameter on my site."

---

## Notes

- Reverses [`add_query_parameter`](add-query-parameter.md).
- Use [`get_query_parameters`](get-query-parameters.md) afterward to confirm removal.
