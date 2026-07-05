---
description: Reference for the remove_country_region_settings MCP tool -- parameters, response format, and example prompts for removing a geo-targeting setting in Bing Webmaster Tools.
---

# remove_country_region_settings

Remove a geo-targeting (country/region) setting from a site, added previously via
[`add_country_region_settings`](add-country-region-settings.md).

---

## Parameters

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `site_url` | string | Yes | The URL of the site |
| `two_letter_iso_country_code` | string | Yes | Two-letter ISO country code to remove |
| `settings_type` | string | Yes | `Page`, `Directory`, `Domain`, or `Subdomain` |
| `url` | string | Yes | The page, directory, domain, or subdomain URL the setting applies to |

---

## Response

```json
{
  "siteUrl": "https://www.example.com/",
  "twoLetterIsoCountryCode": "us",
  "settingsType": "Domain",
  "url": "https://www.example.com/",
  "success": true,
  "requestedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "Remove the United States geo-targeting setting from https://www.example.com/."

---

## Notes

- All four parameters must exactly match what was originally added -- this identifies which
  setting to remove, the same way [`remove_blocked_url`](remove-blocked-url.md) needs to match
  [`add_blocked_url`](add-blocked-url.md)'s values exactly.
- Use [`get_country_region_settings`](get-country-region-settings.md) afterward to confirm removal.
