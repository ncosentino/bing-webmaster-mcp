---
description: Fix common issues with the Bing Webmaster Tools MCP server -- missing credentials, API key errors, and quota limits.
---

# Troubleshooting

## Startup Issues

**Error at startup: "no service credentials provided" / server exits immediately**

The server couldn't find a Bing Webmaster API key in any of the three sources (CLI argument,
environment variable, `.env` file). Verify the environment variable name is exactly
`BING_WEBMASTER_API_KEY` -- it's case-sensitive. See [Configuration](configuration.md).

**HTTP client receives 404**

Use the MCP endpoint `http://127.0.0.1:8083/mcp`. The server root is not an MCP
endpoint. Service supervisors can check `http://127.0.0.1:8083/health`.

**"Invalid API key" or 401/403 errors from every tool**

1. Confirm the key was copied in full from Bing Webmaster Tools → Settings → API Access.
2. Confirm the key hasn't been deleted/regenerated since -- only one key is valid per account at a
   time.
3. Confirm the site you're querying is actually verified under the same Bing Webmaster account
   that generated the key. Use [`list_sites`](tools/list-sites.md) to see which sites your key can
   access.

---

## Quota and Rate Limits

Bing enforces daily/monthly quotas on URL submissions per site. Use
[`get_url_submission_quota`](tools/get-url-submission-quota.md) before a large batch submission to
confirm you have enough quota remaining. `submit_url_batch` accepts at most 500 URLs per call --
larger lists must be split across multiple calls (and multiple days, if quota is exhausted).

---

## IndexNow-Specific Issues

**`submit_url_indexnow` returns "key not configured"**

The IndexNow key is a separate, optional credential from the Bing Webmaster API key -- set
`BING_INDEXNOW_KEY` (or pass `--indexnow-key`, or pass an explicit `key` argument to the tool
call) to use this tool. See [Configuration](configuration.md#indexnow-key-optional).

**`submit_url_indexnow` returns HTTP 403**

Bing couldn't validate your key. Confirm the key file is hosted at
`https://yoursite.com/<key>.txt` at your site's root (not a subdirectory) and that its contents
are exactly the key string with no extra whitespace.

**`submit_url_indexnow` returns HTTP 422**

The URLs you submitted don't belong to the `host` you specified, or the key doesn't match the
host. Double-check the `host` argument matches the domain of every URL in the batch.

---

## Getting More Help

If you're still stuck, [open an issue](https://github.com/ncosentino/bing-webmaster-mcp/issues)
with the tool name, a redacted version of your config (remove your API key), and the exact error
message.
