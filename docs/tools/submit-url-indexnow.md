---
description: Reference for the submit_url_indexnow MCP tool -- parameters, response format, and example prompts for pinging the IndexNow protocol for instant (re)indexing on Bing.
---

# submit_url_indexnow

Ping the [IndexNow](https://www.indexnow.org) protocol to request instant (re)indexing of one or
more URLs. This is a separate, simpler protocol from the classic Bing Webmaster API -- Bing (and
other participating search engines) treat it as the current recommended way to signal new or
changed content quickly.

---

## Parameters

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `host` | string | Yes | The domain the URLs belong to (e.g. `www.example.com`) |
| `url_list` | string[] | Yes | One or more URLs to submit (all must belong to `host`) |
| `key` | string | No | IndexNow key to use for this call. If omitted, uses the configured `BING_INDEXNOW_KEY` |
| `key_location` | string | No | Custom URL where your key file is hosted, if not at the default `https://<host>/<key>.txt` |

---

## Response

```json
{
  "host": "www.example.com",
  "urlCount": 2,
  "submitted": true,
  "requestedAt": "2026-02-21T19:00:00Z"
}
```

On failure:

```json
{
  "error": "IndexNow request failed with HTTP 403: key not valid for host www.example.com"
}
```

---

## Example Prompts

> "Ping IndexNow for my two updated blog posts."

---

## Notes

- Requires a **separate credential** from the Bing Webmaster API key -- see
  [Configuration](../configuration.md#indexnow-key-optional) and
  [Getting Started](../getting-started.md#step-2-optional-get-an-indexnow-key).
- Your key must be hosted as a plain-text file at your site root
  (`https://<host>/<key>.txt`) containing just the key string, unless you provide `key_location`.
- Unlike the classic submission tools, IndexNow isn't subject to Bing's per-site submission quota.
- A successful ping doesn't guarantee indexing -- it just requests expedited crawling.
