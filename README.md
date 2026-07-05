# Bing Webmaster Tools MCP Server -- Bing SEO Data for AI Assistants

[![Latest Release](https://img.shields.io/github/v/release/ncosentino/bing-webmaster-mcp?style=flat-square)](https://github.com/ncosentino/bing-webmaster-mcp/releases/latest)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg?style=flat-square)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.26-00ADD8?style=flat-square&logo=go)](go/go.mod)
[![.NET Version](https://img.shields.io/badge/.NET-10-512BD4?style=flat-square&logo=dotnet)](csharp/Directory.Build.props)
[![CI](https://img.shields.io/github/actions/workflow/status/ncosentino/bing-webmaster-mcp/ci.yml?label=CI&style=flat-square)](https://github.com/ncosentino/bing-webmaster-mcp/actions/workflows/ci.yml)

> **Zero-dependency MCP server for Bing Webmaster Tools.**
> Pre-built native binaries for Linux, macOS, and Windows. No Node.js. No Python. No .NET runtime. No Go toolchain. Download one binary and configure your AI tool.

Expose real Bing search and indexing data directly to AI assistants like Claude, GitHub Copilot, and Cursor via the [Model Context Protocol (MCP)](https://modelcontextprotocol.io). Ask your AI to check crawl issues, submit URLs for indexing, inspect a page's index status, or pull query and keyword performance -- all grounded in real Bing Webmaster Tools data for your verified sites.

---

## Why This Exists

AI assistants are powerful at analyzing SEO strategy -- but they need real data. This MCP server bridges your AI tool to Bing Webmaster Tools, giving it:

- **Site management** -- list, add, and verify sites in your Bing Webmaster account
- **URL submission & indexing** -- submit single/batch URLs, plus [IndexNow](https://www.indexnow.org) for instant (re)indexing
- **Crawl diagnostics** -- crawl issues (decoded into readable labels, not raw bitmasks) and crawl statistics
- **URL & index inspection** -- index status, traffic, and inbound links for any URL
- **Search analytics** -- clicks, impressions, and position by query, by page, and over time
- **Keyword research** -- historical, market-wide keyword impression stats

With this MCP server configured, you can ask your AI: _"Which URLs on my site have crawl issues right now, and what's my overall click/impression trend over the last month?"_ and get a real data-backed answer.

---

## Quick Start

**Two steps: generate an API key, download a binary, add it to your MCP config.**

### Step 1: Generate a Bing Webmaster API Key

1. Sign in to [Bing Webmaster Tools](https://www.bing.com/webmasters) (any Microsoft, Google, or Facebook account works)
2. Add and verify the site(s) you want to query, if you haven't already
3. Click **Settings** (top right) → **API Access**
4. Accept the terms if this is your first time, then click **API Key**
5. Click **Generate API Key**

> **Note:** Bing issues one API key per account -- it grants access to every site already verified under that account.

> **Optional:** The bonus `submit_url_indexnow` tool uses a *separate* [IndexNow](https://www.indexnow.org) key, not your Bing Webmaster API key. See [Getting Started](https://github.devleader.ca/bing-webmaster-mcp/getting-started/) for details -- every other tool works without it.

### Step 2: Download a Binary

Go to the [Releases page](https://github.com/ncosentino/bing-webmaster-mcp/releases/latest) and download the binary for your platform:

| Platform | Go binary | C# binary |
|----------|-----------|-----------|
| Linux x64 | `bwt-mcp-go-linux-amd64` | `bwt-mcp-csharp-linux-x64` |
| Linux arm64 | `bwt-mcp-go-linux-arm64` | `bwt-mcp-csharp-linux-arm64` |
| macOS x64 (Intel) | `bwt-mcp-go-darwin-amd64` | `bwt-mcp-csharp-osx-x64` |
| macOS arm64 (Apple Silicon) | `bwt-mcp-go-darwin-arm64` | `bwt-mcp-csharp-osx-arm64` |
| Windows x64 | `bwt-mcp-go-windows-amd64.exe` | `bwt-mcp-csharp-win-x64.exe` |
| Windows arm64 | `bwt-mcp-go-windows-arm64.exe` | `bwt-mcp-csharp-win-arm64.exe` |

On Linux/macOS, make the binary executable after downloading:

```bash
chmod +x bwt-mcp-go-linux-amd64
```

### Step 3: Add to Your AI Tool Config

See the [Setup by Tool](#setup-by-tool) section below for your specific client.

---

## Setup by Tool

Replace `/path/to/binary` with the actual path to your downloaded binary. Replace `your-api-key-here` with your Bing Webmaster API key.

### Claude Code / GitHub Copilot CLI

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

### Claude Desktop

Edit `~/Library/Application Support/Claude/claude_desktop_config.json` (macOS) or `%APPDATA%\Claude\claude_desktop_config.json` (Windows):

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

### Cursor

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

### VS Code with GitHub Copilot

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

### Visual Studio

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

### Using CLI Arguments

```json
{
  "command": "/path/to/binary",
  "args": ["--api-key", "your-api-key-here", "--indexnow-key", "your-indexnow-key-here"]
}
```

---

## Available Tools

55 MCP tools are exposed -- the 22-tool MVP (21 classic Bing Webmaster API operations plus one bonus IndexNow tool), 21 Phase 2 tools, and 12 Phase 3 tools. Full parameter documentation for every tool is on the [docs site](https://github.devleader.ca/bing-webmaster-mcp/tools/).

| Area | Tools |
|------|-------|
| Sites | `list_sites`, `add_site`, `verify_site`, `remove_site` |
| Site Access | `get_site_roles`, `add_site_role`, `remove_site_role` |
| Sitemaps | `list_sitemaps`, `get_sitemap_details`, `submit_sitemap`, `remove_sitemap` |
| URL Submission & Indexing | `submit_url`, `submit_url_batch`, `submit_url_indexnow`, `get_url_submission_quota` |
| Content Submission | `submit_content`, `get_content_submission_quota` |
| Crawling | `get_crawl_issues`, `get_crawl_stats` |
| URL & Index Inspection | `get_url_info`, `get_url_traffic_info`, `get_url_links`, `get_link_counts` |
| Directory Inspection | `get_children_url_info`, `get_children_url_traffic_info` |
| Blocked URLs | `get_blocked_urls`, `add_blocked_url`, `remove_blocked_url` |
| Fetch Diagnostics | `fetch_url`, `list_fetched_urls`, `get_fetched_url_details` |
| Site Moves | `get_site_moves`, `submit_site_move` |
| Search Analytics | `get_rank_and_traffic_stats`, `get_query_stats`, `get_page_stats`, `get_page_query_stats`, `get_query_page_stats`, `get_query_page_detail_stats`, `get_query_traffic_stats`, `get_keyword_stats` |
| Keyword Research | `get_keyword`, `get_related_keywords` |
| URL Normalization | `get_query_parameters`, `add_query_parameter`, `remove_query_parameter`, `enable_disable_query_parameter` |
| Geo-Targeting | `get_country_region_settings`, `add_country_region_settings`, `remove_country_region_settings` |
| Connected Pages | `get_connected_pages`, `add_connected_page` |
| Page Preview Blocks | `get_active_page_preview_blocks`, `add_page_preview_block`, `remove_page_preview_block` |

**Example prompts:**

> "Which URLs on my site currently have crawl issues, and what are they?"

> "What's my clicks and impressions trend on Bing over the last month?"

> "Submit these 5 new blog post URLs to Bing for indexing."

> "Which queries are driving traffic to my homepage on Bing?"

### Response Structure

Every tool normalizes Bing's raw wire format (PascalCase fields, `/Date(...)/ ` timestamps, an undocumented `d` envelope) into clean, camelCase JSON with ISO-8601 dates. For example, `get_crawl_issues` returns:

```json
{
  "siteUrl": "https://www.example.com/",
  "issues": [
    {
      "url": "https://www.example.com/old-page",
      "httpCode": 404,
      "issues": ["Code4xx"],
      "inLinks": 3
    }
  ],
  "queriedAt": "2026-02-21T19:00:00Z"
}
```

See the [MCP Tools reference](https://github.devleader.ca/bing-webmaster-mcp/tools/) for every tool's exact response shape.

---

## Configuration Reference

Two independent credentials are resolved, each with the same priority order (highest to lowest):

### Bing Webmaster API Key (required)

1. **CLI argument:** `--api-key <key>`
2. **Environment variable:** `BING_WEBMASTER_API_KEY`
3. **`.env` file** (dev convenience): `BING_WEBMASTER_API_KEY=your-api-key-here`

### IndexNow Key (optional -- only needed for `submit_url_indexnow`)

1. **CLI argument:** `--indexnow-key <key>`
2. **Environment variable:** `BING_INDEXNOW_KEY`
3. **`.env` file:** `BING_INDEXNOW_KEY=your-indexnow-key-here`

Missing the IndexNow key does not prevent the server from starting -- only `submit_url_indexnow` returns an error if it's not configured. See the [Configuration guide](https://github.devleader.ca/bing-webmaster-mcp/configuration/) for full details.

---

## Go vs C# -- Which Binary?

Both implementations expose identical tools with identical behavior.

| Aspect | Go | C# Native AOT |
|--------|----|----|
| Binary size | ~10-15 MB | ~25-40 MB |
| Startup time | ~10-50ms | ~50-100ms |
| Runtime dependency | None | None |
| Language | Go 1.26 | C# / .NET 10 |
| MCP SDK | Official `go-sdk` | Official `ModelContextProtocol` |

**Recommendation:** Both work great. Pick Go for smaller binary size, C# if you prefer the .NET ecosystem.

---

## Building from Source

### Go

Requires Go 1.26+:

```bash
cd go
go mod tidy
go build -ldflags="-s -w" -trimpath -o bwt-mcp-go .
```

Run tests:

```bash
go test ./...
```

### C# (.NET 10 SDK Required)

```bash
cd csharp

# Build (non-AOT, for development)
dotnet build BingWebmasterMcp.slnx

# Publish Native AOT
dotnet publish src/BingWebmasterMcp/BingWebmasterMcp.csproj -r linux-x64 -c Release --self-contained true

# Run tests
dotnet test BingWebmasterMcp.slnx
```

---

## Roadmap

The current release covers 55 tools: the 22-tool MVP, Phase 2 (site role delegation, blocked URLs, fetch-as-Bing, site moves, the Content Submission API), and Phase 3 (URL normalization, geo-targeting, connected pages, page preview blocks). See the [Roadmap](https://github.devleader.ca/bing-webmaster-mcp/roadmap/) for what's still planned: Phase 4 (full OAuth 2.0 alongside the API key).

---

## Related Projects

- [google-search-console-mcp](https://github.com/ncosentino/google-search-console-mcp) -- Zero-dependency MCP server for Google Search Console
- [google-psi-mcp](https://github.com/ncosentino/google-psi-mcp) -- Zero-dependency MCP server for Google PageSpeed Insights Core Web Vitals
- [google-keyword-planner-mcp](https://github.com/ncosentino/google-keyword-planner-mcp) -- Zero-dependency MCP server for Google Ads Keyword Planner

---

## About

### Nick Cosentino -- Dev Leader

This MCP server was built by **[Nick Cosentino](https://www.devleader.ca)**, a software engineer and content creator known as **Dev Leader**. Nick creates practical .NET, C#, ASP.NET Core, Blazor, and software engineering content for intermediate to advanced developers -- covering everything from performance optimization and clean architecture to real-world career advice.

This tool was born out of the same real-world SEO analysis work behind its Google-focused siblings, extended to cover Bing. It serves as a practical example of building Native AOT C# and idiomatic Go MCP servers with zero runtime dependencies against an older, quirkier SOAP-era API surface.

**Find Nick online:**

- Blog: [https://www.devleader.ca](https://www.devleader.ca)
- YouTube: [https://www.youtube.com/@devleaderca](https://www.youtube.com/@devleaderca)
- Newsletter: [https://weekly.devleader.ca](https://weekly.devleader.ca)
- LinkedIn: [https://linkedin.com/in/nickcosentino](https://linkedin.com/in/nickcosentino)
- All My Links: [https://links.devleader.ca](https://links.devleader.ca)

### BrandGhost

[BrandGhost](https://www.brandghost.ai) is a social media automation platform built by Nick that lets content creators cross-post and schedule content across all social platforms in one click. If you create content and want to spend less time on distribution and more time creating, check it out.

---

## Contributing

Contributions are welcome! Please:

1. Open an issue describing the bug or feature request before submitting a PR
2. Run `golangci-lint run` (Go) or `dotnet build` with zero warnings (C#) before submitting
3. Keep both implementations in sync -- a feature added to Go should also be added to C#, and vice versa

---

## License

MIT License -- see [LICENSE](LICENSE) for details.
