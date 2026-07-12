---
description: Bing Webmaster Tools MCP server credential configuration reference -- resolution order and environment variables for both the Bing Webmaster API key and the optional IndexNow key.
---

# Configuration

This server resolves **two independent credentials**. Only the Bing Webmaster API key is
required -- the IndexNow key is only needed for the `submit_url_indexnow` tool.

---

## Bing Webmaster API Key (required)

Priority order (highest to lowest):

### 1. CLI Argument

```bash
/path/to/bwt-mcp-go-linux-amd64 --api-key YOUR_API_KEY
```

### 2. Environment Variable

```bash
export BING_WEBMASTER_API_KEY=your-api-key-here
```

### 3. `.env` File (Lowest Priority -- Dev Convenience)

Create a `.env` file in the working directory:

```
BING_WEBMASTER_API_KEY=your-api-key-here
```

If no API key is found in any of these sources, the server logs an error to stderr and exits --
none of the 21 classic tools can function without it.

---

## IndexNow Key (optional)

Priority order (highest to lowest):

### 1. CLI Argument

```bash
/path/to/bwt-mcp-go-linux-amd64 --indexnow-key YOUR_INDEXNOW_KEY
```

### 2. Environment Variable

```bash
export BING_INDEXNOW_KEY=your-indexnow-key-here
```

### 3. `.env` File

```
BING_INDEXNOW_KEY=your-indexnow-key-here
```

If no IndexNow key is configured, the server still starts normally and all 21 classic tools work
as expected -- only `submit_url_indexnow` returns an error explaining the key is missing.

---

## Why Two Separate Credentials?

The classic Bing Webmaster API (site management, crawl diagnostics, search analytics, and classic
URL submission) and [IndexNow](https://www.indexnow.org) are two unrelated protocols:

| | Bing Webmaster API key | IndexNow key |
|---|---|---|
| Issued by | Bing Webmaster Tools, one per account | Chosen by you (any string) |
| Proves | You control the verified Bing Webmaster account | You control the target domain (via a hosted key file) |
| Used by | 21 classic tools | `submit_url_indexnow` only |

See [Getting Started](getting-started.md) for how to obtain each one.

---

## HTTP host

```bash
./bwt-mcp-go-linux-amd64 \
  --transport http \
  --listen-address 127.0.0.1 \
  --port 8083
```

- `--listen-address` overrides `MCP_LISTEN_ADDRESS`; the default is `127.0.0.1`.
- `--port` overrides `PORT`; the default is `8080`.
- MCP is served at `/mcp`.
- Health metadata is served at `/health`.

Go accepts a comma-separated `--allowed-hosts` list. C# uses the standard
semicolon-separated ASP.NET Core `AllowedHosts` setting.
