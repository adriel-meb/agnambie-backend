# Bible Brain API — Complete Developer Reference

> **Source**: Synthesized from two Postman collections:
> - `[M] DBP API Reference v4` (collection `48203400-3c1ed348-0c2c-433b-b683-22048465ab1e`)
> - `[Bible Brain] Example Workflows` (collection `48203400-c2767a33-1191-4641-a4ee-4cfa054f4ba1`)
>
> **Last updated**: 2026-06-13

---

## Table of Contents

1. [Overview](#1-overview)
2. [Authentication](#2-authentication)
3. [Global Conventions](#3-global-conventions)
4. [Endpoints — Alphabets & Number Systems](#4-endpoints--alphabets--number-systems)
5. [Endpoints — Languages](#5-endpoints--languages)
6. [Endpoints — Countries](#6-endpoints--countries)
7. [Endpoints — Bibles](#7-endpoints--bibles)
8. [Endpoints — Filesets & Content](#8-endpoints--filesets--content)
9. [Endpoints — Audio Timestamps](#9-endpoints--audio-timestamps)
10. [Endpoints — Search](#10-endpoints--search)
11. [Endpoints — Downloads](#11-endpoints--downloads)
12. [Endpoints — OpenAPI Schema](#12-endpoints--openapi-schema)
13. [Workflows](#13-workflows)
14. [Media Type Reference](#14-media-type-reference)
15. [Error Codes & Responses](#15-error-codes--responses)
16. [Rate Limits & Caching](#16-rate-limits--caching)

---

## 1. Overview

**Bible Brain** (also known as the Digital Bible Platform, DBP) is a fast, free API providing programmatic access to thousands of Bible translations in text, audio, dramatized audio, and video formats across hundreds of languages.

| Property | Value |
|---|---|
| **Base URL** | `https://4.dbt.io/api` |
| **Alternate Base URL** | `https://b4.dbt.io/api` |
| **API Version** | `4` (always pass `v=4` query parameter) |
| **Response Format** | JSON (`Content-Type: application/json; charset=utf-8`) |
| **Protocol** | HTTPS only |
| **Support Email** | support@digitalbibleplatform.com |

> **Note on URLs**: The collections reference both `https://4.dbt.io/api` and `https://b4.dbt.io/api`. Both point to the same v4 API. Use `https://4.dbt.io/api` for production.

---

## 2. Authentication

### API Key (Query Parameter)

All requests require an API key passed as a query parameter.

| Parameter | Location | Type | Required | Description |
|---|---|---|---|---|
| `key` | Query string | `string` | **Yes** | Your API key (UUID format) |

**Example:**
```
GET https://4.dbt.io/api/languages?v=4&key=YOUR_API_KEY
```

### Key Types

| Key Type | Description |
|---|---|
| **Personal/Production key** | Obtained via the Bible Brain developer portal. Required for production. |
| **Public demo key** | `1462b719-42d8-0874-7c50-905063472458` — rate-limited to **1,000 requests/month**. For exploration only, never for production. |
| **Public domain key** | `f4cdf23a-22c3-66c9-cc4f-05dc711b41c6` — grants access to public-domain content only. |

### CORS

The API allows cross-origin requests from any origin:

```
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: HEAD, GET, POST, PUT, PATCH, DELETE
Access-Control-Allow-Headers: *
```

---

## 3. Global Conventions

### API Version Parameter

Every request **must** include:

```
?v=4
```

This selects the v4 API. Omitting it may result in unexpected behavior.

### Pagination

All list endpoints that support pagination return a `meta.pagination` object alongside `data`:

```json
{
  "data": [...],
  "meta": {
    "pagination": {
      "total": 1798,
      "count": 50,
      "per_page": 50,
      "current_page": 1,
      "total_pages": 36,
      "links": {
        "next": "https://4.dbt.io/api/languages?page=2",
        "previous": "https://4.dbt.io/api/languages?page=1"
      }
    }
  }
}
```

| Field | Type | Description |
|---|---|---|
| `total` | integer | Total number of records matching the query |
| `count` | integer | Number of records returned in this page |
| `per_page` | integer | Records per page (default varies by endpoint; commonly 50 for languages, 25 for bibles) |
| `current_page` | integer | The current page number (1-indexed) |
| `total_pages` | integer | Total number of pages |
| `links.next` | string | URL for the next page (absent on last page) |
| `links.previous` | string | URL for the previous page (absent on first page) |

### Pagination Parameters

| Parameter | Type | Default | Description |
|---|---|---|---|
| `page` | integer | `1` | Page number to retrieve |
| `limit` | integer | varies | Number of results per page |

### Localization (`l10n`)

Many endpoints accept an `l10n` parameter. When set to a valid three-letter ISO language code (e.g., `spa` for Spanish), the API returns localized names where translations exist.

### Search Text Normalization

The Bible search and language/country search endpoints normalize search text:
- Extra whitespace is collapsed
- Multiple `+` separators are normalized (e.g., `holy     ++bible` → `holy bible`)
- Partial word matching is supported on some endpoints

---

## 4. Endpoints — Alphabets & Number Systems

### 4.1 `GET /alphabets` — List All Alphabets

Returns a list of the world's known writing scripts. Useful for font loading and script-aware rendering.

**URL**: `GET https://4.dbt.io/api/alphabets?v=4&key={API_KEY}`

#### Query Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `v` | string | **Yes** | API version. Always `4`. |
| `key` | string | **Yes** | Your API key. |

#### Example Request

```
GET https://4.dbt.io/api/alphabets?v=4&key=1462b719-42d8-0874-7c50-905063472458
```

#### Example Response (`200 OK`)

```json
{
  "data": [
    {
      "name": "Unified Canadian Aboriginal",
      "script": "Cans",
      "family": "American",
      "type": "abugida",
      "direction": "ltr"
    },
    {
      "name": "Latin",
      "script": "Latn",
      "family": "European",
      "type": "alphabet",
      "direction": "ltr"
    }
  ]
}
```

#### Response Fields

| Field | Type | Description |
|---|---|---|
| `name` | string | Human-readable name of the script |
| `script` | string | ISO 15924 script code (used as the identifier) |
| `family` | string | Script family (e.g., `American`, `Semitic`, `European`) |
| `type` | string | Script type: `alphabet`, `abugida`, `abjad`, `syllabary`, `logographic` |
| `direction` | string | Text direction: `ltr` (left-to-right) or `rtl` (right-to-left) |

---

### 4.2 `GET /alphabets/{script_id}` — Alphabet Details

Returns detailed information on a single alphabet, including any fonts and the languages/Bibles using it.

**URL**: `GET https://4.dbt.io/api/alphabets/{script_id}?v=4&key={API_KEY}`

#### Path Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `script_id` | string | **Yes** | The ISO 15924 script code (e.g., `Cans`, `Latn`, `Arab`) |

#### Query Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `v` | string | **Yes** | API version. Always `4`. |
| `key` | string | **Yes** | Your API key. |

#### Example Request

```
GET https://4.dbt.io/api/alphabets/Arab?v=4&key=1462b719-42d8-0874-7c50-905063472458
```

#### Example Response (`200 OK`)

```json
{
  "data": {
    "name": "Unified Canadian Aboriginal",
    "script": "Cans",
    "family": "American",
    "type": "abugida",
    "direction": "ltr",
    "fonts": [
      {
        "id": 7,
        "fontName": "Noto Naskh Arabic",
        "fontFileName": "NotoNaskhArabic-Regular",
        "fontWeight": 400,
        "copyright": "Creative Commons",
        "url": "https://cdn.example.com/resources/fonts/NotoNaskhArabic-Regular.ttf",
        "notes": "notes specific to this font",
        "italic": false
      }
    ],
    "languages": [
      {
        "id": 6411,
        "glotto_id": "stan1288",
        "iso": "spa",
        "name": "Spanish",
        "maps": "Andorra and France",
        "development": "Fully developed. Bible: 1553-2000.",
        "use": "60,000,000 L2 speakers.",
        "location": "Central America, South America",
        "area": "Central, south; Canary Islands...",
        "population": 24900,
        "population_notes": "Population total all countries: 334,800,758.",
        "notes": "The Aragonese dialect of Spanish is different from Aragonese language [arg]. Christian.",
        "typology": "SVO,prepositions,genitives,relatives after noun heads...",
        "description": "language description",
        "status": "6a",
        "country_id": "ES"
      }
    ],
    "bibles": [
      {
        "id": "ENGESV",
        "language_id": 6411,
        "date": 1,
        "scope": "NTPOTP",
        "derived": "English New Revised Standard Version",
        "copyright": "© 1999 Bible Society of Ghana",
        "versification": "protestant",
        "created_at": "2018-02-12 13:32:23",
        "updated_at": "2018-02-12 13:32:23"
      }
    ]
  }
}
```

#### Notes

- The `fonts` array may be empty if no font information is available for the script.
- Font URLs point to CDN-hosted TTF files you can link or embed in your app.
- Some filesets may not render correctly without the fonts listed here.

---

### 4.3 `GET /numbers` — List All Number Systems

Returns all alphabets/scripts that have a custom (non-western-arabic) numeral system.

**URL**: `GET https://4.dbt.io/api/numbers?v=4&key={API_KEY}`

#### Query Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `v` | string | **Yes** | API version. Always `4`. |
| `key` | string | **Yes** | Your API key. |

#### Example Response (`200 OK`)

```json
{
  "data": [
    {
      "id": "bengali",
      "description": "description for bengali",
      "notes": "notes for bengali"
    },
    {
      "id": "devanagari",
      "description": "description for devanagari",
      "notes": "notes for devanagari"
    }
  ]
}
```

---

### 4.4 `GET /numbers/{id}` — Number System Details

Returns the full numeral set for a single numeral system.

**URL**: `GET https://4.dbt.io/api/numbers/{id}?v=4&key={API_KEY}`

#### Path Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `id` | string | **Yes** | The numeral system ID (e.g., `western-arabic`, `bengali`, `devanagari`) |

#### Example Response (`200 OK`)

```json
{
  "data": {
    "id": "western-arabic",
    "value": 9,
    "glyph": "੧",
    "numeral_written": "est"
  }
}
```

---

### 4.5 `GET /numbers/range` — Vernacular Number Range

Returns the vernacular (script-native) glyph for each integer in the given range.

**URL**: `GET https://4.dbt.io/api/numbers/range?script_id={ID}&start={N}&end={M}&v=4&key={API_KEY}`

#### Query Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `script_id` | string | **Yes** | The numeral system ID (e.g., `western-arabic`, `bengali`) |
| `start` | integer | **Yes** | Start of the range (inclusive) |
| `end` | integer | **Yes** | End of the range (inclusive) |
| `v` | string | **Yes** | API version. Always `4`. |
| `key` | string | **Yes** | Your API key. |

#### Example Request

```
GET https://4.dbt.io/api/numbers/range?script_id=bengali&start=1&end=10&v=4&key=1462b719-42d8-0874-7c50-905063472458
```

#### Example Response (`200 OK`)

```json
{
  "data": [
    { "numeral": "1", "numeral_vernacular": "১" },
    { "numeral": "2", "numeral_vernacular": "২" },
    { "numeral": "3", "numeral_vernacular": "৩" }
  ]
}
```

#### Notes

- Use this endpoint to display chapter/verse numbers in the script of the target Bible.
- The `numeral` field is the western-arabic digit; `numeral_vernacular` is the script-native glyph.

---

## 5. Endpoints — Languages

### 5.1 `GET /languages` — List Languages

Returns a paginated list of languages that have content in the Bible Brain database. Supports numerous filters.

**URL**: `GET https://4.dbt.io/api/languages?v=4&key={API_KEY}`

#### Query Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `v` | string | **Yes** | API version. Always `4`. |
| `key` | string | **Yes** | Your API key. |
| `country` | string | No | Filter by ISO 3166 country code (e.g., `US`, `AD`). Returns languages spoken in that country. |
| `iso` | string | No | Filter to an exact ISO 639-3 language code (e.g., `eng`, `spa`, `aaa`). Returns 0 or 1 results. |
| `language_code` | string | No | Alias for `iso`. |
| `language_name` | string | No | Filter by language name (partial match supported on search endpoint; exact on list). |
| `include_translations` | boolean | No | If `true`, includes the ISO language IDs for available translations. |
| `include_alt_names` | boolean | No | If `true`, includes alternative names for the language. |
| `l10n` | string | No | Three-letter ISO code to localize returned names (e.g., `spa` → Spanish names). |
| `media` | string | No | Filter to languages that have content of this media type. Values: `audio`, `audio_drama`, `video`, `text`, `video_stream`, `audio_stream`, `audio_drama_stream`. |
| `set_type_code` | string | No | Filter by a specific fileset type code (e.g., `video_stream`, `audio_drama`). Equivalent to `media` but uses the internal code directly. |
| `page` | integer | No | Page number (default: `1`). |
| `limit` | integer | No | Results per page (default: `50`, max observed: `150`). |

#### Example Request — Page 1 (all languages)

```
GET https://4.dbt.io/api/languages?v=4&key=1462b719-42d8-0874-7c50-905063472458
```

#### Example Request — Languages with audio in United States

```
GET https://4.dbt.io/api/languages?v=4&country=US&media=audio&page=1&key=1462b719-42d8-0874-7c50-905063472458
```

#### Example Request — Language by ISO code

```
GET https://4.dbt.io/api/languages?v=4&iso=spa&limit=150&key=1462b719-42d8-0874-7c50-905063472458
```

#### Example Response (`200 OK`)

```json
{
  "data": [
    {
      "id": 1,
      "glotto_id": "aari1239",
      "iso": "aiw",
      "name": "Aari",
      "autonym": "Aari",
      "bibles": 1,
      "filesets": 4
    },
    {
      "id": 9,
      "glotto_id": "abau1245",
      "iso": "aau",
      "name": "Abau",
      "autonym": "Abau",
      "bibles": 1,
      "filesets": 2
    }
  ],
  "meta": {
    "pagination": {
      "total": 1798,
      "count": 50,
      "per_page": 50,
      "current_page": 1,
      "total_pages": 36,
      "links": {
        "next": "https://4.dbt.io/api/languages?page=2"
      }
    }
  }
}
```

#### Response Fields (`data[]`)

| Field | Type | Description |
|---|---|---|
| `id` | integer | Internal language ID |
| `glotto_id` | string | Glottolog identifier |
| `iso` | string | ISO 639-3 language code |
| `name` | string | English name of the language |
| `autonym` | string | Native name of the language in its own script |
| `bibles` | integer | Count of Bible translations available |
| `filesets` | integer | Count of filesets available (text, audio, video, etc.) |
| `country_population` | integer | Population count for this language in the filtered country (present when `country` filter is used) |

---

### 5.2 `GET /languages/{id}` — Language Details

Returns detailed metadata for a single language by its internal numeric ID.

**URL**: `GET https://4.dbt.io/api/languages/{id}?v=4&key={API_KEY}`

#### Path Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `id` | integer | **Yes** | The internal language ID (e.g., `6411` for Spanish) |

#### Query Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `v` | string | **Yes** | API version. Always `4`. |
| `key` | string | **Yes** | Your API key. |

#### Example Request

```
GET https://4.dbt.io/api/languages/6411?v=4&key=1462b719-42d8-0874-7c50-905063472458
```

#### Example Response (`200 OK`)

```json
{
  "data": {
    "id": 6411,
    "glotto_id": "stan1288",
    "iso": "spa",
    "name": "Spanish",
    "maps": "Andorra and France",
    "development": "Fully developed. Bible: 1553-2000.",
    "use": "60,000,000 L2 speakers.",
    "location": "Central America, South America, Europe",
    "area": "Central, south; Canary Islands. Also in Andorra...",
    "population": 24900,
    "population_notes": "Population total all countries: 334,800,758.",
    "notes": "The Aragonese dialect of Spanish is different from Aragonese language [arg]. Christian.",
    "typology": "SVO,prepositions,genitives,relatives after noun heads,articles,numerals before noun heads...",
    "description": "language description",
    "status": "6a",
    "country_id": "ES"
  }
}
```

#### Response Fields

| Field | Type | Description |
|---|---|---|
| `id` | integer | Internal language ID |
| `glotto_id` | string | Glottolog identifier |
| `iso` | string | ISO 639-3 code |
| `name` | string | English language name |
| `maps` | string | Geographic map references |
| `development` | string | Development status and history |
| `use` | string | Usage notes (L2 speakers, etc.) |
| `location` | string | Textual description of where spoken |
| `area` | string | Geographic distribution details |
| `population` | integer | Estimated speaker population |
| `population_notes` | string | Additional population data notes |
| `notes` | string | Additional linguistic notes |
| `typology` | string | Linguistic typology (word order, etc.) |
| `description` | string | General description |
| `status` | string | Endangerment status code (EGIDS scale, e.g., `6a` = vigorous) |
| `country_id` | string | Primary country ISO code |

---

### 5.3 `GET /languages/search/{search_text}` — Language Search

Full-text search for languages by name (English or native). Supports partial matching.

**URL**: `GET https://4.dbt.io/api/languages/search/{search_text}?v=4&key={API_KEY}`

#### Path Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `search_text` | string | **Yes** | The search query (partial or full name). Even a single letter is valid. |

#### Query Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `v` | string | **Yes** | API version. Always `4`. |
| `key` | string | **Yes** | Your API key. |
| `media` | string | No | Filter results to languages with content of this media type (e.g., `audio`, `video`). |
| `set_type_code` | string | No | Filter to languages with a specific fileset type code (e.g., `video_stream`). |
| `page` | integer | No | Page number (default: `1`). Default page size: `15`. |

#### Example Requests

```
# Search for languages matching "al"
GET https://4.dbt.io/api/languages/search/al?v=4&key=1462b719-42d8-0874-7c50-905063472458

# Search for languages matching "spanish" with audio content
GET https://4.dbt.io/api/languages/search/spanish?v=4&media=audio&key=1462b719-42d8-0874-7c50-905063472458

# Search for languages with video, matching "al"
GET https://4.dbt.io/api/languages/search/al?v=4&set_type_code=video_stream&key=1462b719-42d8-0874-7c50-905063472458
```

#### Notes

- The default `per_page` for search endpoints is `15`.
- Single-character searches return broad results; multi-character searches narrow them.
- Useful for autocomplete/typeahead UIs.

---

## 6. Endpoints — Countries

### 6.1 `GET /countries` — List Countries

Returns a list of countries, optionally with languages spoken in each.

**URL**: `GET https://4.dbt.io/api/countries?v=4&key={API_KEY}`

#### Query Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `v` | string | **Yes** | API version. Always `4`. |
| `key` | string | **Yes** | Your API key. |
| `l10n` | string | No | Localize country names in the language matching this ISO code. |
| `include_languages` | boolean | No | If `true`, includes the major languages spoken in each country. |

#### Example Request

```
GET https://4.dbt.io/api/countries?l10n=spa&include_languages=true&v=4&key=1462b719-42d8-0874-7c50-905063472458
```

#### Example Response (`200 OK`)

```json
{
  "data": [
    {
      "name": "Andorra",
      "continent_code": "EU",
      "languages": {}
    },
    {
      "name": "United States",
      "continent_code": "NA",
      "languages": {
        "eng": { "name": "English", "population": 280000000 }
      }
    }
  ]
}
```

---

### 6.2 `GET /countries/{id}` — Country Details

Returns detailed information on a single country.

**URL**: `GET https://4.dbt.io/api/countries/{id}?v=4&key={API_KEY}`

#### Path Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `id` | string | **Yes** | ISO 3166-1 alpha-2 country code (e.g., `AD`, `US`, `GB`) |

#### Query Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `v` | string | **Yes** | API version. Always `4`. |
| `key` | string | **Yes** | Your API key. |

#### Example Request

```
GET https://4.dbt.io/api/countries/AD?v=4&key=1462b719-42d8-0874-7c50-905063472458
```

#### Example Response (`200 OK`)

```json
{
  "data": {
    "name": "Andorra",
    "continent_code": "EU",
    "languages": {}
  }
}
```

---

### 6.3 `GET /countries/search/{search_text}` — Country Search

Full-text search for countries by name.

**URL**: `GET https://4.dbt.io/api/countries/search/{search_text}?v=4&key={API_KEY}`

#### Path Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `search_text` | string | **Yes** | Search text (partial or full country name, including multi-word like `burkina faso`) |

#### Query Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `v` | string | **Yes** | API version. Always `4`. |
| `key` | string | **Yes** | Your API key. |
| `page` | integer | No | Page number. Default page size: `15`. |

#### Example Requests

```
# One-letter search
GET https://4.dbt.io/api/countries/search/b?v=4&key=1462b719-42d8-0874-7c50-905063472458

# Two-letter search (returns ~Bulgaria, Burundi, etc.)
GET https://4.dbt.io/api/countries/search/bu?v=4&key=1462b719-42d8-0874-7c50-905063472458

# Full name search
GET https://4.dbt.io/api/countries/search/burkina faso?v=4&key=1462b719-42d8-0874-7c50-905063472458
```

#### Example Response (`200 OK`)

```json
{
  "data": [
    {
      "name": "Bulgaria",
      "continent_code": "EU",
      "codes": { "iso_a2": "BG", "iso_a3": "BGR" }
    }
  ],
  "meta": {
    "pagination": {
      "total": 1,
      "per_page": 15,
      "current_page": 1,
      "total_pages": 1
    }
  }
}
```

---

## 7. Endpoints — Bibles

### 7.1 `GET /bibles` — List Bibles

The primary Bible listing endpoint. Returns paginated Bible records with their associated filesets. Only returns Bibles accessible via your API key.

**URL**: `GET https://4.dbt.io/api/bibles?v=4&key={API_KEY}`

#### Query Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `v` | string | **Yes** | API version. Always `4`. |
| `key` | string | **Yes** | Your API key. |
| `language_code` | string | No | ISO 639-3 language code to filter by (e.g., `FRA` for French, `eng` for English, or internal language ID like `6411`). |
| `country_code` | string | No | ISO 3166 country code to filter by (returns Bibles spoken in that country). |
| `media` | string | No | Filter to Bibles that have at least one fileset of this media type. Values: `audio`, `audio_drama`, `text`, `text_plain`, `text_format`, `video_stream`, `audio_stream`, `audio_drama_stream`. |
| `media_exclude` | string | No | Exclude Bibles whose filesets include this media type. |
| `page` | integer | No | Page number (default: `1`). |
| `limit` | integer | No | Results per page (default: `50`). |

#### Example Requests

```
# All bibles (paginated)
GET https://4.dbt.io/api/bibles?v=4&key=1462b719-42d8-0874-7c50-905063472458

# French bibles with audio
GET https://4.dbt.io/api/bibles?language_code=FRA&media=audio&page=1&limit=25&v=4&key=1462b719-42d8-0874-7c50-905063472458

# Bibles with dramatized audio (exclude non-drama audio)
GET https://4.dbt.io/api/bibles?language_code=spa&media=audio_drama&page=1&limit=25&v=4&key=1462b719-42d8-0874-7c50-905063472458

# Bibles with video streaming
GET https://4.dbt.io/api/bibles?language_code=eng&media=video_stream&v=4&key=1462b719-42d8-0874-7c50-905063472458
```

#### Example Response (`200 OK`)

```json
{
  "data": [
    {
      "id": "ENGESV",
      "abbr": "ENGESV",
      "name": "English Standard Version",
      "vname": "English Standard Version",
      "language": "English",
      "language_id": 6414,
      "language_autonym": "English",
      "language_altNames": "English",
      "iso": "eng",
      "date": 2001,
      "filesets": {
        "dbp-prod": [
          {
            "id": "ENGESVN1DA",
            "set_type_code": "audio_drama",
            "set_size_code": "NT"
          },
          {
            "id": "ENGESPFLS",
            "set_type_code": "text_plain",
            "set_size_code": "NTPOTP"
          }
        ]
      }
    }
  ],
  "meta": {
    "pagination": {
      "total": 1801,
      "count": 25,
      "per_page": 25,
      "current_page": 1,
      "total_pages": 73
    }
  }
}
```

#### Response Fields (`data[]`)

| Field | Type | Description |
|---|---|---|
| `id` / `abbr` | string | Bible abbreviation ID (e.g., `ENGESV`, `SPARVR`) |
| `name` | string | English name of the Bible translation |
| `vname` | string | Vernacular (native-language) name of the translation |
| `language` | string | English name of the language |
| `language_id` | integer | Internal language ID |
| `language_autonym` | string | Native name of the language |
| `language_altNames` | string | Alternative language names |
| `iso` | string | ISO 639-3 language code |
| `date` | integer | Year of publication/translation |
| `filesets` | object | Keyed by bucket name (`dbp-prod`); each value is an array of fileset objects |
| `filesets[].id` | string | The fileset ID (used to access content) |
| `filesets[].set_type_code` | string | Media type code (see [Media Type Reference](#14-media-type-reference)) |
| `filesets[].set_size_code` | string | Coverage code: `NT` (New Testament), `OT` (Old Testament), `NTPOTP` (NT + portions of OT), `C` (complete), etc. |

#### Notes

- The `id` / `abbr` field identifies the Bible (e.g., `ENGESV`). Each Bible can have multiple filesets of different types.
- The default page size is **50** for this endpoint (may vary per API key configuration).
- Pagination total is approximately 1,800+ Bibles.

---

### 7.2 `GET /bibles/{id}` — Bible Details

Returns detailed metadata for a single Bible, including publisher, provider, links, and all associated filesets.

**URL**: `GET https://4.dbt.io/api/bibles/{id}?v=4&key={API_KEY}`

#### Path Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `id` | string | **Yes** | The Bible ID / abbreviation (e.g., `ENGESV`, `RUNDPI`) |

#### Query Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `v` | string | **Yes** | API version. Always `4`. |
| `key` | string | **Yes** | Your API key. |

#### Example Request

```
GET https://4.dbt.io/api/bibles/ENGESV?v=4&key=1462b719-42d8-0874-7c50-905063472458
```

#### Example Response (`200 OK`)

```json
[
  {
    "id": "ENGESV",
    "alphabet": "Latn",
    "mark": "© 2001 Crossway Bibles",
    "name": "English Standard Version",
    "description": "The English Standard Version (ESV) is an essentially literal translation...",
    "vname": "English Standard Version",
    "vdescription": "...",
    "publishers": {
      "id": 1,
      "slug": "crossway",
      "abbreviation": "ESV",
      "notes": "Publisher notes",
      "primaryColor": "#004b85",
      "secondaryColor": "#e2383f",
      "inactive": false,
      "url_facebook": "https://facebook.com/crossway",
      "url_website": "https://www.esv.org",
      "url_donate": "https://www.crossway.org/donate",
      "url_twitter": "https://twitter.com/crossway",
      "address": "1300 Crescent Street",
      "address2": "",
      "city": "Wheaton",
      "state": "IL",
      "country": "USA",
      "zip": "60187",
      "phone": "630-682-4300",
      "email": "info@crossway.org"
    },
    "providers": {
      "id": 1,
      "slug": "american-bible-society",
      "abbreviation": "ABS",
      "notes": "Provider notes",
      "primaryColor": "#004b85",
      "secondaryColor": "#e2383f",
      "inactive": false,
      "url_website": "https://www.americanbible.org"
    },
    "language": "English",
    "iso": "eng",
    "date": 2001,
    "country": "United States",
    "books": "MAT",
    "links": {
      "id": "link-id",
      "bible_id": "ENGESV",
      "type": "web",
      "url": "http://bibles.org/versions/ENGESV",
      "title": "BibleSearch",
      "organization_id": 1,
      "created_at": "2018-01-01",
      "updated_at": "2018-01-01",
      "provider": "American Bible Society"
    },
    "filesets": {
      "id": "ENGESVN1DA",
      "set_type_code": "audio_drama",
      "set_size_code": "NT"
    }
  }
]
```

> **Note**: This endpoint returns an array (not an object with a `data` wrapper).

#### Response Fields

| Field | Type | Description |
|---|---|---|
| `id` | string | Bible ID |
| `alphabet` | string | Script code used by this translation |
| `mark` | string | Copyright mark/symbol |
| `name` | string | English translation name |
| `description` | string | English description of the translation |
| `vname` | string | Vernacular translation name |
| `vdescription` | string | Vernacular description |
| `publishers` | object | Publisher organization details |
| `providers` | object | Content provider details |
| `language` | string | Language name |
| `iso` | string | ISO 639-3 language code |
| `date` | integer | Publication year |
| `country` | string | Primary country of use |
| `books` | string | Books available (USFM book ID or list) |
| `links` | object | External links to Bible resources |
| `filesets` | object | Associated fileset metadata |

---

### 7.3 `GET /bibles/{id}/book` — Bible Books

Returns translated book names and structural information (chapters, verses) for the given Bible.

**URL**: `GET https://4.dbt.io/api/bibles/{id}/book?v=4&key={API_KEY}`

#### Path Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `id` | string | **Yes** | The Bible ID (e.g., `ENGESV`) |

#### Query Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `v` | string | **Yes** | API version. Always `4`. |
| `key` | string | **Yes** | Your API key. |
| `book_id` | string | No | Filter to a specific book by USFM ID (e.g., `MAT`, `GEN`, `REV`). |
| `verify_content` | boolean | No | If `true`, only returns books that actually have content available. |
| `verse_count` | boolean | No | If `true`, includes the verse count for each chapter. |

#### Example Requests

```
# All books for ENGESV
GET https://4.dbt.io/api/bibles/ENGESV/book?v=4&key=1462b719-42d8-0874-7c50-905063472458

# Only Matthew, with verse counts, verified content
GET https://4.dbt.io/api/bibles/ENGESV/book?book_id=MAT&verify_content=true&verse_count=true&v=4&key=1462b719-42d8-0874-7c50-905063472458
```

#### Example Response (`200 OK`)

```json
{
  "data": [
    {
      "book_id": "MAT",
      "book_id_usfx": "Mt",
      "book_id_osis": "Matt",
      "name": "Matthew",
      "testament": "NT",
      "testament_order": 1,
      "book_order": 40,
      "book_group": "Gospels",
      "chapters": [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28]
    }
  ]
}
```

#### Response Fields

| Field | Type | Description |
|---|---|---|
| `book_id` | string | USFM book ID (e.g., `GEN`, `MAT`, `REV`) |
| `book_id_usfx` | string | USFX book ID |
| `book_id_osis` | string | OSIS book ID |
| `name` | string | Localized book name (in the Bible's language) |
| `testament` | string | `OT` (Old Testament) or `NT` (New Testament) |
| `testament_order` | integer | Order within the testament |
| `book_order` | integer | Global book order |
| `book_group` | string | Canonical group (e.g., `Gospels`, `Pauline Epistles`, `Major Prophets`) |
| `chapters` | array | List of chapter numbers available |

#### Notes

- The actual list of books may differ from fileset to fileset within the same Bible. For example, a KJV fileset may include deuterocanonical books not present in a sibling fileset.
- Use `verify_content=true` to ensure only books with actual content are returned for a given fileset.

---

### 7.4 `GET /bibles/defaults/types` — Default Bible for a Language

Returns the recommended default Bible(s) for a given language, grouped by media type.

**URL**: `GET https://4.dbt.io/api/bibles/defaults/types?language_code={CODE}&v=4&key={API_KEY}`

#### Query Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `v` | string | **Yes** | API version. Always `4`. |
| `key` | string | **Yes** | Your API key. |
| `language_code` | string | No | ISO 639-3 or internal language code to filter by (e.g., `AR`, `en`, `spa`). |

#### Example Request

```
GET https://4.dbt.io/api/bibles/defaults/types?language_code=en&v=4&key=1462b719-42d8-0874-7c50-905063472458
```

#### Example Response (`200 OK`)

```json
{
  "en": {
    "video": "ENGESV",
    "audio": "ENGESV"
  }
}
```

#### Notes

- The response is keyed by language code, then by media type.
- Use this to automatically pre-select the best Bible for a user's language.

---

### 7.5 `GET /bibles/{bible_id}/copyright` — Bible Copyright

Returns the copyright information for all filesets of a Bible. **This endpoint must be called** and the copyright displayed to end users as a legal requirement.

**URL**: `GET https://4.dbt.io/api/bibles/{bible_id}/copyright?v=4&key={API_KEY}`

#### Path Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `bible_id` | string | **Yes** | The Bible ID (e.g., `ENGESV`) |

#### Query Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `v` | string | **Yes** | API version. Always `4`. |
| `key` | string | **Yes** | Your API key. |
| `iso` | string | No | ISO 639-3 language code to localize the organization translations returned. |

#### Example Request

```
GET https://4.dbt.io/api/bibles/ENGESV/copyright?iso=spa&v=4&key=1462b719-42d8-0874-7c50-905063472458
```

#### Example Response (`200 OK`)

```json
[
  {
    "id": "ENGESVN1DA",
    "type": "audio_drama",
    "size": "NT",
    "copyright": {
      "copyright_date": "2001",
      "copyright": "© 2001 Crossway Bibles",
      "copyright_description": "The Holy Bible, English Standard Version. ESV® Text Edition: 2016...",
      "open_access": 0
    }
  },
  {
    "id": "ENGESV",
    "type": "text_plain",
    "size": "NTPOTP",
    "copyright": {
      "copyright_date": "2014",
      "copyright": "© Ethnos360",
      "copyright_description": "© Ethnos360",
      "open_access": 1
    }
  }
]
```

#### Response Fields

| Field | Type | Description |
|---|---|---|
| `id` | string | Fileset ID |
| `type` | string | Media type code |
| `size` | string | Coverage size code |
| `copyright.copyright_date` | string | Year of copyright |
| `copyright.copyright` | string | Copyright notice string |
| `copyright.copyright_description` | string | Full copyright description |
| `copyright.open_access` | integer | `1` = open access/public domain, `0` = restricted |

> **Legal note**: Displaying copyright information is **required** for all content accessed via the API. Always call this endpoint and surface the result to your users before presenting Bible content.

---

### 7.6 `GET /bibles/search/{search_text}` — Bible Search

Full-text search for Bibles by name or abbreviation.

**URL**: `GET https://4.dbt.io/api/bibles/search/{search_text}?v=4&key={API_KEY}`

#### Path Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `search_text` | string | **Yes** | Search term (partial name, abbreviation, even a single letter). Multi-word queries with `++` or extra spaces are normalized. |

#### Query Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `v` | string | **Yes** | API version. Always `4`. |
| `key` | string | **Yes** | Your API key. |
| `page` | integer | No | Page number. Default page size: `15`. |

#### Example Requests

```
# Single-letter search
GET https://4.dbt.io/api/bibles/search/v?v=4&key=1462b719-42d8-0874-7c50-905063472458

# Three-letter search (finds Reina Valera)
GET https://4.dbt.io/api/bibles/search/val?v=4&key=1462b719-42d8-0874-7c50-905063472458

# Multi-word search
GET https://4.dbt.io/api/bibles/search/holy     ++bible?v=4&key=1462b719-42d8-0874-7c50-905063472458

# Named example (Van Dyke Arabic Bible)
GET https://4.dbt.io/api/bibles/search/van    +++dyke?v=4&page=1&key=1462b719-42d8-0874-7c50-905063472458
```

#### Example Response (`200 OK`)

```json
{
  "data": [
    {
      "abbr": "SPARVR",
      "name": "Reina Valera 1995",
      "language_id": 6411
    }
  ],
  "meta": {
    "pagination": {
      "total": 8,
      "per_page": 15,
      "current_page": 1,
      "total_pages": 1
    }
  }
}
```

---

## 8. Endpoints — Filesets & Content

### 8.1 `GET /bibles/filesets/media/types` — Available Media Types

Returns a simple key-value map of all fileset type codes and their human-readable labels.

**URL**: `GET https://4.dbt.io/api/bibles/filesets/media/types?v=4&key={API_KEY}`

#### Query Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `v` | string | **Yes** | API version. Always `4`. |
| `key` | string | **Yes** | Your API key. |

#### Example Response (`200 OK`)

```json
{
  "audio_drama": "Dramatized Audio",
  "audio": "Audio",
  "text_plain": "Plain Text",
  "text_format": "Formatted Text",
  "video_stream": "Video",
  "audio_stream": "Audio HLS Stream",
  "audio_drama_stream": "Dramatized Audio HLS Stream"
}
```

---

### 8.2 `GET /bibles/filesets/{fileset_id}/{book}/{chapter}` — Fileset Content

Returns the actual content (text, audio URLs, or video URLs) for a given fileset, book, and chapter.

**URL**: `GET https://4.dbt.io/api/bibles/filesets/{fileset_id}/{book}/{chapter}?v=4&key={API_KEY}`

#### Path Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `fileset_id` | string | **Yes** | The fileset ID (e.g., `ENGKJVN1SA`, `ENGESVN1DA`, `BEQDPIN1DA`). Obtained from `/bibles` response. |
| `book` | string | **Yes** | USFM book ID (e.g., `MAT`, `GEN`, `REV`, `PSA`). |
| `chapter` | integer | **Yes** | Chapter number. |

#### Query Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `v` | string | **Yes** | API version. Always `4`. |
| `key` | string | **Yes** | Your API key. |
| `verse_start` | integer | No | Filter to start from this verse number (inclusive). |
| `verse_end` | integer | No | Filter to end at this verse number (inclusive). |

#### Example Requests

```
# Text content (plain text)
GET https://4.dbt.io/api/bibles/filesets/ENGKJVN1ET/MAT/1?v=4&key=1462b719-42d8-0874-7c50-905063472458

# Audio content (standard audio)
GET https://4.dbt.io/api/bibles/filesets/ENGESVN1SA/MAT/1?v=4&key=1462b719-42d8-0874-7c50-905063472458

# Dramatized audio (DA)
GET https://4.dbt.io/api/bibles/filesets/BEQDPIN1DA/MAT/1?v=4&key=1462b719-42d8-0874-7c50-905063472458

# Video content (DV)
GET https://4.dbt.io/api/bibles/filesets/ENGJESN2DV/MAT/1?v=4&key=1462b719-42d8-0874-7c50-905063472458

# Text with thumbnail
GET https://4.dbt.io/api/bibles/filesets/ENGKJVN1ET/MAT/1?v=4&key=1462b719-42d8-0874-7c50-905063472458

# Specific verse range
GET https://4.dbt.io/api/bibles/filesets/ENGKJVN1ET/MAT/5?verse_start=3&verse_end=12&v=4&key=1462b719-42d8-0874-7c50-905063472458
```

#### Response Shape (Text — `text_plain` / `text_format`)

```json
{
  "data": [
    {
      "book_id": "MAT",
      "book_name": "Matthew",
      "book_name_alt": "Matt",
      "chapter": 1,
      "chapter_alt": "1",
      "verse_start": 1,
      "verse_start_alt": "1",
      "verse_end": 1,
      "verse_end_alt": "1",
      "verse_text": "The book of the generation of Jesus Christ, the son of David, the son of Abraham."
    }
  ]
}
```

#### Response Shape (USX text — `text_usx`)

Returns XML-format text with rich markup:
```json
{
  "data": [
    {
      "book_id": "MAT",
      "chapter": 1,
      "content": "<usx version=\"3.0\">...</usx>"
    }
  ]
}
```

#### Response Shape (JSON text — `text_json`)

Returns structured JSON with verse tokens and formatting marks.

#### Response Shape (Audio — `audio` / `audio_drama`)

```json
{
  "data": [
    {
      "book_id": "MAT",
      "book_name": "Matthew",
      "chapter_start": 1,
      "chapter_end": null,
      "verse_start": null,
      "verse_end": null,
      "timestamp": null,
      "path": "https://cdn.dbt.io/audio/ENGESVN1SA/B02___01_Matthew_____ENGESVN1SA.mp3",
      "duration": 287.25,
      "filesize_in_bytes": 4598956,
      "thumbnail": null
    }
  ]
}
```

#### Response Shape (Audio Stream — `audio_stream` / `audio_drama_stream`)

Returns HLS playlist URLs:
```json
{
  "data": [
    {
      "book_id": "MAT",
      "chapter_start": 1,
      "path": "https://cdn.dbt.io/audio/ENGESVN1SA/B02___01_Matthew_____ENGESVN1SA/playlist.m3u8"
    }
  ]
}
```

#### Response Shape (Video Stream — `video_stream`)

```json
{
  "data": [
    {
      "book_id": "MAT",
      "chapter_start": 1,
      "path": "https://cdn.dbt.io/video/ENGJESN2DV/B02___01_Matthew/master.m3u8",
      "thumbnail": "https://cdn.dbt.io/video/ENGJESN2DV/B02___01_Matthew/thumbnail.jpg"
    }
  ]
}
```

#### Fileset Type Code Reference for this Endpoint

| `set_type_code` | Content Type | Notes |
|---|---|---|
| `text_plain` | Plain text verses | Simple verse-by-verse text |
| `text_format` | Formatted text | Paragraphed, formatted text |
| `text_usx` | USX XML | Rich USFM XML format |
| `text_json` | JSON structured text | JSON with tokens and formatting |
| `audio` | Standard audio MP3/HLS | Narrated, single voice |
| `audio_drama` | Dramatized audio MP3/HLS | Multi-voice dramatization |
| `audio_stream` | Audio HLS stream | Adaptive streaming (m3u8) |
| `audio_drama_stream` | Dramatized audio HLS stream | Dramatized + adaptive streaming |
| `video_stream` | Video HLS stream | Video Bible (e.g., film Jesus) |

---

### 8.3 `GET /bibles/filesets/{fileset_id}/copyright` — Fileset Copyright

Returns the copyright for a specific fileset (not the whole Bible).

**URL**: `GET https://4.dbt.io/api/bibles/filesets/{fileset_id}/copyright?v=4&key={API_KEY}`

#### Path Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `fileset_id` | string | **Yes** | The fileset ID (e.g., `ENGESVN1DA`) |

#### Query Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `v` | string | **Yes** | API version. Always `4`. |
| `key` | string | **Yes** | Your API key. |

---

## 9. Endpoints — Audio Timestamps

Audio timestamps enable karaoke-style synchronized text+audio, audio-based verse search, and jump-to-verse navigation.

### 9.1 `GET /timestamps` — Filesets with Timestamps

Returns a list of all fileset IDs that have verse-level timing metadata available.

**URL**: `GET https://4.dbt.io/api/timestamps?v=4&key={API_KEY}`

#### Query Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `v` | string | **Yes** | API version. Always `4`. |
| `key` | string | **Yes** | Your API key. |

#### Example Response (`200 OK`)

```json
[
  { "fileset_id": "ENGESVN1DA" },
  { "fileset_id": "ENGKJVN1SA" },
  { "fileset_id": "SPARVRDRA" }
]
```

#### Notes

- Only a subset of filesets have timestamp data — check this list before querying timestamps for a specific fileset.

---

### 9.2 `GET /timestamps/{fileset_id}/{book}/{chapter}` — Chapter Timestamps

Returns verse-level timing information (offset in seconds) for a specific chapter of an audio fileset.

**URL**: `GET https://4.dbt.io/api/timestamps/{fileset_id}/{book}/{chapter}?v=4&key={API_KEY}`

#### Path Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `fileset_id` | string | **Yes** | Audio fileset ID (must be in the `/timestamps` list). |
| `book` | string | **Yes** | USFM book ID (e.g., `MAT`). |
| `chapter` | integer | **Yes** | Chapter number. |

#### Query Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `v` | string | **Yes** | API version. Always `4`. |
| `key` | string | **Yes** | Your API key. |

#### Example Request

```
GET https://4.dbt.io/api/timestamps/ENGKJVN1DA/MAT/4?v=4&key=1462b719-42d8-0874-7c50-905063472458
```

#### Example Response (`200 OK`)

```json
{
  "data": [
    {
      "book": "MAT",
      "chapter": "4",
      "verse_start": "1",
      "timestamp": 0.0
    },
    {
      "book": "MAT",
      "chapter": "4",
      "verse_start": "2",
      "timestamp": 8.43
    },
    {
      "book": "MAT",
      "chapter": "4",
      "verse_start": "3",
      "timestamp": 10.19
    }
  ]
}
```

#### Response Fields

| Field | Type | Description |
|---|---|---|
| `book` | string | USFM book ID |
| `chapter` | string | Chapter number |
| `verse_start` | string | Starting verse number for this timestamp |
| `timestamp` | float | Time offset in **seconds** from the start of the audio file |

#### Use Cases

- **Karaoke / sync text-audio**: Highlight the current verse as audio plays.
- **Jump to verse**: Seek the audio player to `timestamp` seconds when a user taps a verse.
- **Audio search**: Scan timestamps to find which offset a searched passage appears at.

---

## 10. Endpoints — Search

### 10.1 `GET /search` — Bible Full-Text Search

Searches the text content of a Bible fileset for a word or phrase.

**URL**: `GET https://4.dbt.io/api/search?query={TERM}&fileset_id={FILESET}&v=4&key={API_KEY}`

#### Query Parameters

| Parameter | Type | Required | Description |
|---|---|---|---|
| `v` | string | **Yes** | API version. Always `4`. |
| `key` | string | **Yes** | Your API key. |
| `query` | string | **Yes** | The word or phrase to search for (e.g., `Jesus`, `love`, `in the beginning`). |
| `fileset_id` | string | **Yes** | The Bible **text** fileset ID to search within (e.g., `ENGESV`). |
| `limit` | integer | No | Number of results per page (default: `25`, example: `15`). |
| `page` | integer | No | Page number (default: `1`). |
| `books` | string | No | Comma-separated USFM book IDs to restrict the search scope (e.g., `GEN,EXO,MAT`). |

#### Example Requests

```
# Search for "Jesus" in ESV
GET https://4.dbt.io/api/search?query=Jesus&fileset_id=ENGESV&limit=15&page=1&v=4&key=1462b719-42d8-0874-7c50-905063472458

# Search for "love" only in Matthew
GET https://4.dbt.io/api/search?query=love&fileset_id=ENGESV&books=MAT&v=4&key=1462b719-42d8-0874-7c50-905063472458

# Search across Genesis, Exodus, Matthew
GET https://4.dbt.io/api/search?query=covenant&fileset_id=ENGESV&books=GEN,EXO,MAT&v=4&key=1462b719-42d8-0874-7c50-905063472458
```

#### Example Response (`200 OK`)

```json
{
  "verses": {
    "data": [
      {
        "book_id": "MAT",
        "book_name": "Matthew",
        "book_name_alt": "Matt",
        "chapter": "4",
        "chapter_alt": "4",
        "verse_start": "4",
        "verse_start_alt": "4",
        "verse_end": "4",
        "verse_end_alt": "4",
        "verse_text": "But he answered, 'It is written, Man shall not live by bread alone...'"
      }
    ]
  },
  "meta": {
    "pagination": {
      "total": 258,
      "count": 15,
      "per_page": 15,
      "current_page": 1,
      "total_pages": 18
    }
  }
}
```

#### Response Fields (`verses.data[]`)

| Field | Type | Description |
|---|---|---|
| `book_id` | string | USFM book ID |
| `book_name` | string | Localized book name |
| `book_name_alt` | string | Alternative book name/abbreviation |
| `chapter` | string | Chapter number |
| `chapter_alt` | string | Chapter number (vernacular) |
| `verse_start` | string | Starting verse number |
| `verse_start_alt` | string | Starting verse number (vernacular) |
| `verse_end` | string | Ending verse number |
| `verse_end_alt` | string | Ending verse number (vernacular) |
| `verse_text` | string | The matching verse text |

---

## 11. Endpoints — Downloads

The download endpoints provide access to downloadable content packages, as opposed to the streaming/fileset content endpoints.

### 11.1 `GET /bibles/filesets` — List Downloadable Content

> **Note**: The exact download list endpoint path varies by usage context. In the Example Workflows collection, it is referenced as a "Content available for download" section.

The following are documented download access patterns from the collection:

#### Download a Whole Fileset (USX)

```
GET https://4.dbt.io/api/bibles/filesets/{fileset_id}?v=4&key={API_KEY}
```

Example:
```
GET https://4.dbt.io/api/bibles/filesets/CJOWBT?v=4&key=1462b719-42d8-0874-7c50-905063472458
```

#### Download a Specific Book from a Fileset (DA audio)

```
GET https://4.dbt.io/api/bibles/filesets/{fileset_id}/{book_id}?v=4&key={API_KEY}
```

Example:
```
GET https://4.dbt.io/api/bibles/filesets/BMQBSMN1DA/MAT?v=4&key=1462b719-42d8-0874-7c50-905063472458
```

#### Download a Specific Chapter (USX or audio)

```
GET https://4.dbt.io/api/bibles/filesets/{fileset_id}/{book_id}/{chapter}?v=4&key={API_KEY}
```

Example:
```
GET https://4.dbt.io/api/bibles/filesets/CJOWBT/MAT/1?v=4&key=1462b719-42d8-0874-7c50-905063472458
```

---

## 12. Endpoints — OpenAPI Schema

### 12.1 `GET /openapi` — Retrieve OpenAPI 3 Schema

Returns the machine-readable OpenAPI 3.0 specification for the entire Bible Brain API.

**URL**: `GET https://4.dbt.io/api/openapi?v=4&key={API_KEY}`

#### Example Request

```
GET https://4.dbt.io/api/openapi?v=4&key=1462b719-42d8-0874-7c50-905063472458
```

Use this to auto-generate client SDKs or import into API tools like Swagger UI.

---

## 13. Workflows

These are step-by-step walkthroughs derived from the `[Bible Brain] Example Workflows` collection. Each workflow answers a real developer question.

---

### Workflow 1 — Discover What Languages Are Available

**Goal**: Understand what languages the API has content for.

**Step 1 — List the first page of all languages**

```
GET https://4.dbt.io/api/languages?v=4&key={API_KEY}
```

- Default page size is **50** languages per page.
- The response includes `meta.pagination.total` (typically 1,780+), `total_pages`, and a `links.next` URL.

**Step 2 — Navigate to the next page**

```
GET https://4.dbt.io/api/languages?v=4&page=2&key={API_KEY}
```

- The second page includes `links.previous` and `links.next`.

**Step 3 — Filter languages with audio content**

```
GET https://4.dbt.io/api/languages?v=4&media=audio&key={API_KEY}
```

**Step 4 — Filter languages with video content**

```
GET https://4.dbt.io/api/languages?v=4&media=video&key={API_KEY}
```

**Step 5 — Look up a language by ISO code**

```
GET https://4.dbt.io/api/languages?v=4&iso=spa&limit=150&key={API_KEY}
```

Returns exactly 1 result for a valid ISO code.

**Step 6 — Filter by country**

```
GET https://4.dbt.io/api/languages?v=4&include_alt_names=true&country=US&page=1&key={API_KEY}
```

---

### Workflow 2 — Search for a Language by Name

**Goal**: Help users find a language via a search field (typeahead).

**Step 1 — Search by partial name (one letter)**

```
GET https://4.dbt.io/api/languages/search/a?v=4&key={API_KEY}
```

Returns up to 15 results per page.

**Step 2 — Search with media filter (languages with audio)**

```
GET https://4.dbt.io/api/languages/search/al?v=4&media=audio&key={API_KEY}
```

**Step 3 — Search for languages with a specific media type code**

```
GET https://4.dbt.io/api/languages/search/al?v=4&set_type_code=video_stream&key={API_KEY}
```

---

### Workflow 3 — Find the Right Bible for a Given Language

**Goal**: Given a language, find all available Bible translations and their media types.

**Step 1 — List all Bibles (paginated)**

```
GET https://4.dbt.io/api/bibles?key={API_KEY}
```

Note: The `v=4` can be passed as a header `v: 4` in some client setups, or as a query param.

**Step 2 — Filter by country code**

```
GET https://4.dbt.io/api/bibles?country_code=BR&v=4&key={API_KEY}
```

**Step 3 — Filter by language code**

```
GET https://4.dbt.io/api/bibles?language_code=eng&v=4&key={API_KEY}
```

**Step 4 — Filter by language code + media type (audio)**

```
GET https://4.dbt.io/api/bibles?language_code=eng&media=audio&v=4&key={API_KEY}
```

**Step 5 — Filter by language code + media type (dramatized audio)**

```
GET https://4.dbt.io/api/bibles?language_code=spa&media=audio_drama_stream&v=4&key={API_KEY}
```

**Step 6 — Filter by language code + media type (video)**

```
GET https://4.dbt.io/api/bibles?language_code=eng&media=video_stream&v=4&key={API_KEY}
```

**Step 7 — Get full details on a specific Bible**

```
GET https://4.dbt.io/api/bibles/RUNDPI?v=4&key={API_KEY}
```

Confirms publisher, provider, available books, and all filesets.

**Step 8 — Get books available in that Bible**

```
GET https://4.dbt.io/api/bibles/ENGESV/book?v=4&key={API_KEY}
```

---

### Workflow 4 — Fetch Copyright Before Displaying Content (Required)

**Goal**: Retrieve and display required legal copyright information.

**Step 1 — Get copyright for a Bible**

```
GET https://4.dbt.io/api/bibles/ENGESV/copyright?v=4&key={API_KEY}
```

Response returns per-fileset copyright. Display `copyright.copyright` to users.

**Step 2 — Get copyright for a specific fileset**

```
GET https://4.dbt.io/api/bibles/filesets/ENGESVN1DA/copyright?v=4&key={API_KEY}
```

> **Important**: Displaying copyright is a **legal requirement** enforced by content publishers. Always surface this before showing Bible content.

---

### Workflow 5 — Read Bible Text Content

**Goal**: Retrieve verse text for a specific fileset, book, and chapter.

**Step 1 — Identify the text fileset ID from the Bible listing**

Call `GET /bibles?language_code=eng&media=text_plain&v=4&key={API_KEY}` and find a fileset with `set_type_code: "text_plain"`.

**Step 2 — Get verses for a chapter**

```
GET https://4.dbt.io/api/bibles/filesets/ENGKJV/MAT/1?v=4&key={API_KEY}
```

**Step 3 — Get a specific verse range**

```
GET https://4.dbt.io/api/bibles/filesets/ENGKJV/MAT/5?verse_start=3&verse_end=12&v=4&key={API_KEY}
```

Returns only the Beatitudes (Matthew 5:3–12).

**Step 4 — Get USX (rich XML) format**

```
GET https://4.dbt.io/api/bibles/filesets/ENGESVN1USX/MAT/1?v=4&key={API_KEY}
```

**Step 5 — Get JSON structured format**

```
GET https://4.dbt.io/api/bibles/filesets/ENGESVN1JSON/MAT/1?v=4&key={API_KEY}
```

---

### Workflow 6 — Stream Audio (HLS) for a Chapter

**Goal**: Stream audio Bible content using HTTP Live Streaming (HLS/m3u8).

**Step 1 — Identify an audio stream fileset**

From `/bibles`, find a fileset with `set_type_code: "audio_stream"` or `"audio_drama_stream"`. Note the fileset ID (e.g., `ENGESVN1SA`).

**Step 2 — Get the master playlist URL for a chapter**

```
GET https://4.dbt.io/api/bibles/filesets/ENGESVN1SA/MAT/1?v=4&key={API_KEY}
```

The `path` field in the response is the `.m3u8` master playlist URL.

**Step 3 — Fetch the master playlist**

```
GET {path}/playlist.m3u8
```

Example: `https://cdn.dbt.io/audio/ENGESVN1SA/B02___01_Matthew_____ENGESVN1SA/playlist.m3u8`

This returns an HLS master playlist listing available bitrates.

**Step 4 — Request a specific bitrate sub-playlist**

```
GET {cdn_base}/{fileset_id}/{file_id}/{bitrate}kbps.m3u8
```

Example: `https://cdn.dbt.io/audio/ENGESVN1SA/B02___01_Matthew_____ENGESVN1SA/64kbps.m3u8`

**Step 5 — Pass the final m3u8 URL to an HLS player**

Use a native HLS player (iOS `AVPlayer`, Android `ExoPlayer`, or `hls.js` in web browsers) to play the stream.

---

### Workflow 7 — Stream Video (HLS) for a Chapter

**Goal**: Stream video Bible content using HLS.

**Step 1 — Identify a video stream fileset**

From `/bibles`, find a fileset with `set_type_code: "video_stream"`.

**Step 2 — Retrieve metadata and master playlist URL for the fileset/book/chapter**

```
GET https://4.dbt.io/api/bibles/filesets/{video_fileset_id}/{book}/{chapter}?v=4&key={API_KEY}
```

The `path` field contains the master `.m3u8` URL.

**Step 3 — Request the master playlist for the desired section**

```
GET {path}
```

This returns an HLS master playlist with multiple resolution options (e.g., 240p, 360p, 720p).

**Step 4 — Request the playlist for the desired resolution**

```
GET {playlist_url_for_selected_resolution}
```

This returns the segment playlist (`.ts` chunk references).

**Step 5 — Feed to an HLS video player**

Pass the resolution-specific `.m3u8` to `<video>` with an HLS player library (e.g., Video.js + hls.js).

---

### Workflow 8 — Synchronized Audio+Text (Timestamps / Karaoke)

**Goal**: Display highlighted verse text that follows along with audio playback.

**Step 1 — Check if the fileset has timestamps**

```
GET https://4.dbt.io/api/timestamps?v=4&key={API_KEY}
```

Look for your fileset ID (e.g., `ENGESVN1DA`) in the returned list.

**Step 2 — Get timestamps for the chapter**

```
GET https://4.dbt.io/api/timestamps/ENGESVN1DA/MAT/4?v=4&key={API_KEY}
```

Returns per-verse `timestamp` values in seconds.

**Step 3 — Get the audio for the same chapter**

```
GET https://4.dbt.io/api/bibles/filesets/ENGESVN1DA/MAT/4?v=4&key={API_KEY}
```

Extract the audio `path` URL.

**Step 4 — Get the text for the same chapter**

```
GET https://4.dbt.io/api/bibles/filesets/{text_fileset_id}/MAT/4?v=4&key={API_KEY}
```

**Step 5 — Implement synchronized highlighting**

On `timeupdate` events from the audio player, compare `currentTime` (seconds) against the timestamp array to determine the active verse, and highlight it in the text display.

---

### Workflow 9 — Discover and Display Country/Language Data

**Goal**: Build a country → language → Bible browsing UI.

**Step 1 — List countries**

```
GET https://4.dbt.io/api/countries?v=4&key={API_KEY}
```

**Step 2 — Search for a country**

```
GET https://4.dbt.io/api/countries/search/burkina faso?v=4&key={API_KEY}
```

**Step 3 — Get country details**

```
GET https://4.dbt.io/api/countries/US?v=4&key={API_KEY}
```

**Step 4 — List languages for that country**

```
GET https://4.dbt.io/api/languages?country=US&v=4&key={API_KEY}
```

**Step 5 — For each language, list Bibles**

```
GET https://4.dbt.io/api/bibles?language_code=eng&v=4&key={API_KEY}
```

---

### Workflow 10 — Script/Alphabet Rendering Support

**Goal**: Determine the correct font for a Bible's script and load it.

**Step 1 — Get the Bible's alphabet code from the Bible detail**

```
GET https://4.dbt.io/api/bibles/ARBNKJVS?v=4&key={API_KEY}
```

Inspect the `alphabet` field (e.g., `Arab`).

**Step 2 — Get font info for that alphabet**

```
GET https://4.dbt.io/api/alphabets/Arab?v=4&key={API_KEY}
```

The `fonts` array contains CDN URLs for TTF files.

**Step 3 — Get vernacular numerals for chapter/verse display**

```
GET https://4.dbt.io/api/numbers/range?script_id=arabic-indic&start=1&end=176&v=4&key={API_KEY}
```

Use `numeral_vernacular` to display chapter and verse numbers in the correct script.

---

### Workflow 11 — Full Bible Browsing App Flow (End-to-End)

This represents the complete recommended flow for building a Bible reader app.

**Step 1** — Let the user search for their language:
```
GET /languages/search/{query}?v=4&key={API_KEY}
```

**Step 2** — Show them Bibles available in that language:
```
GET /bibles?language_code={iso}&v=4&key={API_KEY}
```

**Step 3** — User selects a Bible. Get its details and available books:
```
GET /bibles/{bible_id}?v=4&key={API_KEY}
GET /bibles/{bible_id}/book?verify_content=true&v=4&key={API_KEY}
```

**Step 4** — Get and display copyright (required before any content):
```
GET /bibles/{bible_id}/copyright?v=4&key={API_KEY}
```

**Step 5** — User selects a book and chapter. Get the content:
```
GET /bibles/filesets/{fileset_id}/{book}/{chapter}?v=4&key={API_KEY}
```

**Step 6** — (If audio) Check for timestamp availability:
```
GET /timestamps?v=4&key={API_KEY}
# If fileset is listed:
GET /timestamps/{fileset_id}/{book}/{chapter}?v=4&key={API_KEY}
```

**Step 7** — (If audio stream) Pass the `path` m3u8 URL to an HLS player.

**Step 8** — (If user searches for a verse) Call search:
```
GET /search?query={term}&fileset_id={text_fileset}&v=4&key={API_KEY}
```

---

## 14. Media Type Reference

| Code | Human Label | Description |
|---|---|---|
| `audio` | Audio | Standard narrated audio files (MP3) |
| `audio_drama` | Dramatized Audio | Multi-voice dramatized audio (MP3) |
| `audio_stream` | Audio HLS Stream | Adaptive bitrate audio via HLS (.m3u8) |
| `audio_drama_stream` | Dramatized Audio HLS Stream | Dramatized audio via HLS (.m3u8) |
| `text_plain` | Plain Text | Verse-by-verse plain text |
| `text_format` | Formatted Text | Paragraphed, formatted text |
| `text_usx` | USX XML | USFM XML format for rich text rendering |
| `text_json` | JSON Text | Structured JSON tokens |
| `video_stream` | Video | Bible video via HLS (.m3u8) |

### Fileset Size Codes

| Code | Meaning |
|---|---|
| `NT` | New Testament only |
| `OT` | Old Testament only |
| `NTPOTP` | New Testament + Portions of Old Testament |
| `C` | Complete Bible (OT + NT) |
| `P` | Portions |

---

## 15. Error Codes & Responses

The API returns standard HTTP status codes. Error bodies follow this pattern:

```json
{
  "error": {
    "message": "Human-readable error description",
    "status": 401
  }
}
```

| HTTP Status | Meaning | Common Cause |
|---|---|---|
| `200 OK` | Success | Request fulfilled normally |
| `400 Bad Request` | Invalid parameters | Missing required parameter, malformed query |
| `401 Unauthorized` | Authentication failed | Missing or invalid API key |
| `403 Forbidden` | Access denied | Your API key does not grant access to this resource |
| `404 Not Found` | Resource not found | Invalid fileset ID, Bible ID, book ID, or language ID |
| `429 Too Many Requests` | Rate limit exceeded | Too many requests within the rate limit window |
| `500 Internal Server Error` | Server error | Backend error; retry with exponential backoff |

### Rate Limiting

Rate limit information is returned in response headers:

| Header | Description |
|---|---|
| `X-RateLimit-Limit` | Maximum requests allowed per window (observed values: `1500`, `5000`) |
| `X-RateLimit-Remaining` | Remaining requests in the current window |

> The demo key (`1462b719-42d8-0874-7c50-905063472458`) is limited to **1,000 requests/month**.
> Production keys have higher limits (observed: 5,000 requests per rolling window in older versions).

---

## 16. Rate Limits & Caching

### Caching

The API returns:
```
Cache-Control: no-cache, private
```

This means responses should not be cached at the proxy level. However, you are encouraged to implement your own application-level caching for stable data (languages, countries, Bible metadata) to avoid redundant requests and stay within rate limits.

### Recommended Caching Strategy

| Resource | Recommended TTL |
|---|---|
| `/languages` | 24 hours |
| `/countries` | 24 hours |
| `/alphabets` | 24 hours |
| `/bibles` list | 1 hour |
| `/bibles/{id}` | 1 hour |
| `/bibles/{id}/book` | 1 hour |
| `/bibles/{id}/copyright` | 24 hours |
| `/bibles/filesets/{id}/{book}/{chapter}` (audio/video URLs) | 15–60 minutes (CDN URLs may expire) |
| `/timestamps` list | 24 hours |
| `/timestamps/{fileset}/{book}/{chapter}` | 24 hours |
| `/search` results | Do not cache (user-driven, dynamic) |

### Security Headers

Responses include:
```
X-Frame-Options: SAMEORIGIN
X-XSS-Protection: 1; mode=block
X-Content-Type-Options: nosniff
```

---

*This reference was generated from the `[M] DBP API Reference v4` and `[Bible Brain] Example Workflows` Postman collections (workspace `c68e2c58-dde7-4a09-8c74-68f54243c023`), analyzed on 2026-06-13.*
