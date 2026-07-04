---
description: Zero-dependency MCP server for Bing Webmaster Tools. Pre-built native binaries expose site management, sitemaps, URL submission, crawl diagnostics, and search analytics directly to AI assistants like Claude, GitHub Copilot, and Cursor.
---

# Bing Webmaster Tools MCP

> **Zero-dependency MCP server for Bing Webmaster Tools.**
> Pre-built native binaries for Linux, macOS, and Windows. No Node.js. No Python. No .NET runtime. No Go toolchain. Download one binary and configure your AI tool.

Expose real Bing search and indexing data directly to AI assistants via the
[Model Context Protocol (MCP)](https://modelcontextprotocol.io). Ask your AI to check crawl
issues, submit URLs for indexing, inspect a page's index status, or pull query and keyword
performance -- all grounded in real Bing Webmaster Tools data for your verified sites.

---

## Why This Exists

AI assistants are powerful at analyzing SEO strategy -- but they need real data. This MCP server
bridges your AI tool to Bing Webmaster Tools, giving it:

- **Site management** -- list, add, and verify sites in your Bing Webmaster account
- **URL submission & indexing** -- submit single/batch URLs, plus IndexNow for instant (re)indexing
- **Crawl diagnostics** -- crawl issues and crawl statistics for every verified site
- **URL & index inspection** -- index status, traffic, and inbound links for any URL
- **Search analytics** -- clicks, impressions, and position by query, by page, and over time
- **Keyword research** -- historical keyword impression stats

With this MCP server configured, you can ask your AI: *"Which URLs on my site have crawl issues
right now, and what's my overall click/impression trend over the last month?"* and get a real
data-backed answer.

---

## Quick Overview

22 MCP tools are exposed across six areas -- see the [MCP Tools reference](tools/index.md) for the
full list with parameters.

| Area | Example tools |
|------|----------------|
| Sites | `list_sites`, `add_site`, `verify_site` |
| Sitemaps | `list_sitemaps`, `submit_sitemap` |
| URL submission & indexing | `submit_url`, `submit_url_batch`, `submit_url_indexnow` |
| Crawling | `get_crawl_issues`, `get_crawl_stats` |
| URL & index inspection | `get_url_info`, `get_url_traffic_info`, `get_url_links` |
| Search analytics | `get_rank_and_traffic_stats`, `get_query_stats`, `get_keyword_stats` |

---

## Get Started

**[→ Getting Started](getting-started.md)** -- two steps: generate an API key, add it to your AI
tool config.

---

## About

Built by **[Nick Cosentino](https://www.devleader.ca)** (Dev Leader) -- a software engineer and
content creator covering .NET, C#, and software architecture. Available in both Go and C# (Native
AOT) with zero runtime dependencies.

- Blog: [devleader.ca](https://www.devleader.ca)
- GitHub: [ncosentino/bing-webmaster-mcp](https://github.com/ncosentino/bing-webmaster-mcp)
