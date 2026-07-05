---
description: Reference for the get_country_region_settings MCP tool -- parameters, response format, and example prompts for listing geo-targeting settings in Bing Webmaster Tools.
---

# get_country_region_settings

List the geo-targeting (country/region) settings configured for a site -- which pages,
directories, domains, or subdomains are targeted at a specific country or region.

---

## Parameters

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `site_url` | string | Yes | The URL of the site |

---

## Response

```json
{
  "siteUrl": "https://www.example.com/",
  "rowCount": 1,
  "settings": [
    {
      "twoLetterIsoCountryCode": "us",
      "settingsType": "Domain",
      "url": "https://www.example.com/",
      "date": "2026-02-21T19:00:00Z"
    }
  ],
  "queriedAt": "2026-02-21T19:00:00Z"
}
```

---

## Example Prompts

> "What country or region geo-targeting settings are configured for my site?"

---

## Notes

- `settingsType` is `Page`, `Directory`, `Domain`, or `Subdomain` -- the scope the country
  targeting applies to.
- Use [`add_country_region_settings`](add-country-region-settings.md) /
  [`remove_country_region_settings`](remove-country-region-settings.md) to manage this list.
