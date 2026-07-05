---
description: Reference for the submit_content MCP tool -- parameters, response format, and example prompts for submitting raw content directly to Bing for a URL.
---

# submit_content

Submit raw content (HTTP response and structured data) directly for a specific URL, bypassing a
live crawl.

---

## Parameters

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `site_url` | string | Yes | -- | The URL of the site |
| `url` | string | Yes | -- | The specific URL this content represents |
| `http_message` | string | Yes | -- | Base64-encoded raw HTTP response (headers + body) |
| `structured_data` | string | Yes | -- | Base64-encoded structured data (e.g. JSON-LD) for the page |
| `dynamic_serving` | string | No | `None` | Device target: `None`, `PcLaptop`, `Mobile`, `Amp`, `Tablet`, `NonVisualBrowser` |

---

## Response

```json
{
  "siteUrl": "https://www.example.com/",
  "url": "https://www.example.com/blog/my-post",
  "dynamicServing": "None",
  "success": true,
  "requestedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "Submit this pre-rendered HTML directly to Bing for my JS-heavy page that Bing's crawler
> struggles to render."

---

## Notes

- **Advanced/rarely needed.** Most sites should use [`submit_url`](submit-url.md) or
  [`fetch_url`](fetch-url.md) instead -- this exists for cases where you need to hand Bing exact
  content it can't otherwise obtain by crawling (e.g. content requiring authentication, or
  content that renders differently than what Bing's crawler sees).
- Check [`get_content_submission_quota`](get-content-submission-quota.md) before submitting --
  this has its own quota separate from classic URL submission.
- **Not live-tested** against a real account -- constructing valid base64-encoded HTTP content
  for a meaningful test is complex; this is covered by unit tests only.
