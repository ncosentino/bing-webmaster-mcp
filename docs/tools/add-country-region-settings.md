---
description: Reference for the add_country_region_settings MCP tool -- parameters, response format, and example prompts for adding a geo-targeting setting in Bing Webmaster Tools.
---

# add_country_region_settings

Add a geo-targeting (country/region) setting for a site -- tells Bing that a specific page,
directory, domain, or subdomain scope targets a particular country.

---

## Parameters

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `site_url` | string | Yes | The URL of the site |
| `two_letter_iso_country_code` | string | Yes | Two-letter ISO country code, for example `us` |
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

> "Target https://www.example.com/ at the United States."

---

## Notes

- **Reversible** via [`remove_country_region_settings`](remove-country-region-settings.md), which
  must be called with the exact same `two_letter_iso_country_code`, `settings_type`, and `url` to
  identify which setting to remove.
- Use [`get_country_region_settings`](get-country-region-settings.md) afterward to confirm.
