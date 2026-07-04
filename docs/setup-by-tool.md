---
description: Configure the Bing Webmaster Tools MCP server for Claude, Claude Desktop, GitHub Copilot, Cursor, VS Code, and Visual Studio with exact JSON config snippets.
---

# Setup by Tool

Replace `/path/to/binary` with the actual path to your downloaded binary. Replace
`your-api-key-here` with your [Bing Webmaster API key](getting-started.md).

---

## Claude Code / GitHub Copilot CLI

```json
{
  "mcpServers": {
    "bing-webmaster": {
      "type": "stdio",
      "command": "/path/to/bwt-mcp-go-linux-amd64",
      "args": [],
      "env": {
        "BING_WEBMASTER_API_KEY": "your-api-key-here"
      }
    }
  }
}
```

---

## Claude Desktop

Edit `~/Library/Application Support/Claude/claude_desktop_config.json` (macOS) or
`%APPDATA%\Claude\claude_desktop_config.json` (Windows):

```json
{
  "mcpServers": {
    "bing-webmaster": {
      "command": "/path/to/bwt-mcp-go-darwin-arm64",
      "env": {
        "BING_WEBMASTER_API_KEY": "your-api-key-here"
      }
    }
  }
}
```

---

## Cursor

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

---

## VS Code with GitHub Copilot

```json
{
  "mcp": {
    "servers": {
      "bing-webmaster": {
        "type": "stdio",
        "command": "/path/to/bwt-mcp-go-linux-amd64",
        "env": {
          "BING_WEBMASTER_API_KEY": "your-api-key-here"
        }
      }
    }
  }
}
```

---

## Visual Studio

```json
{
  "bing-webmaster": {
    "command": "C:\\path\\to\\bwt-mcp-csharp-win-x64.exe",
    "env": {
      "BING_WEBMASTER_API_KEY": "your-api-key-here"
    }
  }
}
```

---

## Adding the Optional IndexNow Key

Add `BING_INDEXNOW_KEY` alongside `BING_WEBMASTER_API_KEY` in any of the configs above to enable
the `submit_url_indexnow` tool:

```json
{
  "env": {
    "BING_WEBMASTER_API_KEY": "your-api-key-here",
    "BING_INDEXNOW_KEY": "your-indexnow-key-here"
  }
}
```

---

## Using a CLI Argument Instead of an Environment Variable

```json
{
  "command": "/path/to/binary",
  "args": ["--api-key", "your-api-key-here"]
}
```

---

## Troubleshooting

**Server exits immediately on startup:** No Bing Webmaster API key was found in any of the CLI
argument, environment variable, or `.env` file sources. See [Configuration](configuration.md).

**`submit_url_indexnow` returns an error but everything else works:** The IndexNow key isn't
configured -- it's optional and independent from the Bing Webmaster API key. See
[Configuration](configuration.md#indexnow-key-optional).
