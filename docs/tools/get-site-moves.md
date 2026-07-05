---
description: Reference for the get_site_moves MCP tool -- parameters, response format, and example prompts for checking site move history in Bing Webmaster Tools.
---

# get_site_moves

List site move requests submitted for a site.

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
  "rows": [
    {
      "date": "2026-01-15T00:00:00Z",
      "moveScope": "Domain",
      "moveType": "Local",
      "sourceUrl": "https://old.example.com/",
      "targetUrl": "https://www.example.com/"
    }
  ],
  "queriedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "Have I submitted any site moves for this domain?"

---

## Notes

- `moveScope` is `Domain`, `Host`, or `Directory`. `moveType` is `Local` (within Bing's index
  scope) or `Global` (across markets).
- Use [`submit_site_move`](submit-site-move.md) to request a new one.
- **Live-testing note:** as of testing, this endpoint returned an HTML 404 error against a real,
  verified site (confirmed via a raw HTTP call, not a client-side bug) -- Bing may have
  deprecated or renamed this endpoint since it was documented. The tool handles this gracefully
  (a clean error, not a crash), but a genuine success response hasn't been observed. See the
  [Roadmap](../roadmap.md) for details.
