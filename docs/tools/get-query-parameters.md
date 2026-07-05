---
description: Reference for the get_query_parameters MCP tool -- parameters, response format, and example prompts for listing URL normalization query parameters in Bing Webmaster Tools.
---

# get_query_parameters

List the URL normalization query parameters configured for a site -- parameters Bing should
ignore (or specifically consider) when deciding whether two URLs point at the same content.

---

## Parameters

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `site_url` | string | Yes | The URL of the site |

---

## Response

```json
{
  "siteUrl": "https://www.example.com/",
  "rowCount": 1,
  "parameters": [
    {
      "parameter": "utm_campaign",
      "isEnabled": true,
      "source": 0,
      "date": "2026-02-21T19:00:00Z"
    }
  ],
  "queriedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "What URL query parameters has Bing been told to ignore on my site?"

---

## Notes

- `source` is an undocumented raw integer from Bing's API (its exact meaning isn't published);
  it's passed through as-is rather than guessed at.
- Use [`add_query_parameter`](add-query-parameter.md) / [`remove_query_parameter`](remove-query-parameter.md)
  to manage this list, and [`enable_disable_query_parameter`](enable-disable-query-parameter.md)
  to toggle an existing entry without removing it.
