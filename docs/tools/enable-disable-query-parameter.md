---
description: Reference for the enable_disable_query_parameter MCP tool -- parameters, response format, and example prompts for toggling a URL normalization query parameter in Bing Webmaster Tools.
---

# enable_disable_query_parameter

Enable or disable an existing URL normalization query parameter for a site, without removing it
from the list.

---

## Parameters

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `site_url` | string | Yes | The URL of the site |
| `query_parameter` | string | Yes | The query parameter name to update |
| `is_enabled` | boolean | Yes | `true` to enable, `false` to disable |

---

## Response

```json
{
  "siteUrl": "https://www.example.com/",
  "queryParameter": "utm_campaign",
  "isEnabled": false,
  "success": true,
  "requestedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "Disable the utm_campaign query parameter on my site without removing it."

---

## Notes

- Unlike most other optional parameters in this server, `is_enabled` has no default -- it must
  always be supplied explicitly since there's no sensible assumption for whether you mean to
  enable or disable.
- Use [`get_query_parameters`](get-query-parameters.md) to see current enabled/disabled state.
- Use [`remove_query_parameter`](remove-query-parameter.md) instead if you want the entry gone
  entirely rather than just disabled.
