---
description: Reference for the get_connected_pages MCP tool -- parameters, response format, and example prompts for listing connected pages in Bing Webmaster Tools.
---

# get_connected_pages

List pages Bing considers "connected" to a site -- other domains or pages declared (via
[`add_connected_page`](add-connected-page.md)) as syndicating or mirroring this site's content,
along with Bing's verification status for that connection.

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
  "pages": [
    {
      "url": "https://mirror.example.net/",
      "isVerified": true,
      "requestedMasterSite": "https://www.example.com/",
      "actualMasterSite": "https://www.example.com/",
      "httpStatusCode": 200,
      "market": "en-US",
      "isBlocked": false,
      "lastSuccessfullyVerified": "2026-02-21T19:00:00Z"
    }
  ],
  "queriedAt": "2026-02-21T19:00:00Z"
}
```

---

## Notes

Bing's real underlying record for a connected page has more internal bookkeeping fields (an app
ID/name, consecutive-failure counters, a deep-launch-support flag, etc.). Only the fields likely
useful to a caller are surfaced above; `requestedMasterSite` vs `actualMasterSite` differing can
indicate the connection didn't verify the way it was configured, and `httpStatusCode` reflects the
last HTTP status Bing observed when checking the connection.

---

## Example Prompts

> "What pages does Bing consider connected to my site, and are they verified?"

---

## See Also

- [`add_connected_page`](add-connected-page.md) to declare a new connection. There is no
  `remove_connected_page` in Bing's API.
