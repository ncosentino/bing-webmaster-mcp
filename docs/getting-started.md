---
description: Set up the Bing Webmaster Tools MCP server in two steps -- generate an API key from Bing Webmaster Tools and add it to your AI tool configuration.
---

# Getting Started

Two steps: generate a Bing Webmaster API key, add it to your AI tool config.

---

## Step 1: Generate a Bing Webmaster API Key

1. Sign in to [Bing Webmaster Tools](https://www.bing.com/webmasters) (any Microsoft, Google, or
   Facebook account works).
2. Add and verify the site(s) you want to query, if you haven't already.
3. Click **Settings** (top right) → **API Access**.
4. If this is your first time, read and accept the terms, then click **API Key**.
5. Click **Generate API Key**.

!!! note "One key per user, not per site"
    Bing issues one API key per account. It grants access to every site already verified under
    that account -- you don't need a separate key per site.

!!! warning "Keep it secret"
    Don't share your API key with any third party. If it's ever compromised, delete it and
    generate a new one from the same API Access page.

---

## Step 2 (Optional): Get an IndexNow Key

The bonus `submit_url_indexnow` tool uses a *different* credential -- an arbitrary
[IndexNow](https://www.indexnow.org) key, not your Bing Webmaster API key. If you want instant
(re)indexing pings in addition to the classic submission tools:

1. Generate any key string yourself (a random hex/GUID string works) or use
   [Bing's IndexNow key generator](https://www.bing.com/indexnow/getstarted).
2. Host it as a text file at your site root: `https://yoursite.com/<key>.txt`, containing just the
   key string.

This step is entirely optional -- every other tool works without it.

---

## Step 3: Configure Your AI Tool

Add the server to your AI tool's MCP configuration. See [Setup by Tool](setup-by-tool.md) for
tool-specific instructions.

The minimal config pattern (replace the path with your actual binary location):

```json
{
  "mcpServers": {
    "bing-webmaster": {
      "command": "/path/to/bwt-mcp-go-linux-amd64",
      "env": {
        "BING_WEBMASTER_API_KEY": "your-api-key-here"
      }
    }
  }
}
```

Add `BING_INDEXNOW_KEY` alongside it if you generated one in Step 2.

---

## Next Steps

- [MCP Tools Reference](tools/index.md) -- full parameter documentation for all 55 tools
- [Configuration](configuration.md) -- credential resolution order and all configuration options
- [Setup by Tool](setup-by-tool.md) -- exact config snippets for Claude, Cursor, VS Code, Visual Studio
