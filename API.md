# API

> [!IMPORTANT]
> All API requests must include an Authorization header in the following format:
> Authorization: Bearer <BEARER_TOKEN>
>
> Rate Limit: Configurable via RATE_LIMIT environment variable (default 5) requests per minute per IP address
> All endpoints return JSON responses with appropriate HTTP status codes

### /api/manual-generate/

**Endpoint:** `/think-root/api/manual-generate/`

**Method:** `POST`

**Description:** This endpoint is used to manually generate description for a provided repository URL, and add it to the database. Supports multilingual text generation.

**Curl Example:**

```bash
curl -X POST \
  'http://localhost:8080/think-root/api/manual-generate/' \
  -H 'Authorization: Bearer <BEARER_TOKEN>' \
  -H 'Content-Type: application/json' \
  -d '{
    "url": "https://github.com/example/repo",
    "llm_output_language": "en,uk,fr",
    "llm_provider": "mistral_api",
    "llm_config": {
      "model": "mistral-small-latest"
    }
  }'
```

**Request Parameters:**

| Parameter               | Type    | Required | Description                                                                                                                                                                                                         |
| ----------------------- | ------- | -------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `url`                 | string  | Yes      | GitHub repository URL. Supports multiple whitespace-separated URLs to process in a single request.                                                                                                                  |
| `llm_output_language` | string  | No       | Comma-separated language codes (e.g., "en,uk,fr"). Default: "uk".                                                                                                                                                    |
| `llm_provider`        | string  | No       | LLM provider name. Values: `mistral_api` (default), `openai`, `openrouter`, `chutes`. If omitted, defaults to `mistral_api`.                                                                                         |
| `llm_config`          | object  | Yes*     | Provider configuration. **Required** for all providers to specify at least the `model`. See llm_config Structure section below.                                                                                     |
| `use_direct_url`      | boolean | No       | If true, the URL string is used directly as LLM input instead of README content.                                                                                                                                    |

\* `llm_config` is required for all providers to specify the model.

**`llm_config` Structure:**

The `llm_config` object is passed to the LLM provider's chat completion API. Common parameters:

-   `model`: (Required) The model ID (e.g., `mistral-small-latest`, `gpt-4o`, `google/gemini-2.0-flash-exp:free`, `moonshotai/Kimi-K2-Instruct-0905`).
-   `temperature`: (Optional) Sampling temperature (0.0 to 1.0).
-   `max_tokens`: (Optional) Maximum tokens to generate.
-   `top_p`: (Optional) Nucleus sampling probability.
-   `messages`: (Optional) Array of message objects to specify a custom prompt. See below.

**Custom Prompt:**

You can provide a custom system prompt via the `messages` array. The server will **append** its multilingual instructions to your prompt:

```json
"llm_config": {
  "model": "mistral-small-latest",
  "messages": [
    {
      "role": "system",
      "content": "Your custom prompt here. Describe repos in a fun, engaging way."
    }
  ]
}
```

If you don't provide a `messages` array, the server creates one with default multilingual instructions.

**Request Examples:**

1. Basic request with Mistral (default provider):

```json
{
  "url": "https://github.com/example/repo",
  "llm_config": {
    "model": "mistral-small-latest"
  }
}
```

2. Request with Chutes provider:

```json
{
  "url": "https://github.com/example/repo",
  "llm_provider": "chutes",
  "llm_output_language": "en",
  "llm_config": {
    "model": "moonshotai/Kimi-K2-Instruct-0905"
  }
}
```

3. Multilingual request with OpenRouter:

```json
{
  "url": "https://github.com/example/repo",
  "llm_output_language": "en,uk,fr",
  "llm_provider": "openrouter",
  "llm_config": {
    "model": "google/gemini-2.0-flash-exp:free"
  }
}
```

**Status Codes:**

- 200: Success
- 400: Invalid request
- 401: Unauthorized

**Response Example:**

```json
{
  "status": "ok",
  "added": ["https://github.com/example/repo"],
  "dont_added": []
}
```

### /api/auto-generate/

**Endpoint:** `/think-root/api/auto-generate/`

**Method:** `POST`

**Description:** This endpoint is used to automatically parse trending repositories and generate description based on certain parameters. It also adds the generated posts to the database. Supports multilingual text generation. Supports multiple data sources (GitHub, OssInsight).

**Curl Example:**

```bash
curl -X POST \
  'http://localhost:8080/think-root/api/auto-generate/' \
  -H 'Authorization: Bearer <BEARER_TOKEN>' \
  -H 'Content-Type: application/json' \
  -d '{
    "max_repos": 5,
    "resource": "github",
    "since": "weekly",
    "spoken_language_code": "en",
    "llm_output_language": "en,uk,fr",
    "llm_provider": "mistral_api",
    "llm_config": {
      "model": "mistral-small-latest"
    }
  }'
```

**Request Parameters:**

| Parameter                | Type    | Required | Description                                                                                                                                                                                                                    |
| ------------------------ | ------- | -------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `max_repos`            | integer | Yes      | Maximum number of repositories to process. Must be > 0.                                                                                                                                                                        |
| `resource`             | string  | No       | Data source. Values: `github` (default), `ossinsight`.                                                                                                                                                                         |
| `since`                | string  | No       | **For GitHub resource**: Time period for trending repos (`daily`, `weekly`, `monthly`).                                                                                                                                        |
| `spoken_language_code` | string  | No       | **For GitHub resource**: Spoken language filter for GitHub Trending.                                                                                                                                                           |
| `period`               | string  | No       | **For OssInsight resource**: Time period (`past_24_hours`, `past_week`, `past_month`, `past_3_months`). Default: `past_24_hours`.                                                                                              |
| `language`             | string  | No       | **For OssInsight resource**: Programming language filter (e.g., `Python`, `All`). Default: `All`.                                                                                                                              |
| `llm_output_language`  | string  | No       | Comma-separated language codes for output (e.g., `en,uk,fr`). Default: `uk`.                                                                                                                                                   |
| `llm_provider`         | string  | No       | LLM provider name. Values: `mistral_api` (default), `openai`, `openrouter`, `chutes`. If omitted, defaults to `mistral_api`.                                                                                                   |
| `llm_config`           | object  | Yes*     | Provider configuration. **Required** for all providers to specify at least the `model`.                                                                                                                                        |
| `use_direct_url`       | boolean | No       | If true, the repository URL string is used directly as LLM input instead of README content.                                                                                                                                    |

\* `llm_config` is required for all providers to specify the model.

**`llm_config` Structure:**

The `llm_config` object is passed as the JSON body to the chosen LLM provider's chat completion API (e.g., OpenAI, Mistral, OpenRouter). Common parameters include:

-   `model`: (Required) The ID of the model to use (e.g., `gpt-4o`, `mistral-large-latest`, `google/gemini-2.0-flash-exp:free`).
-   `temperature`: (Optional) Sampling temperature (0.0 to 1.0). Higher values mean more random output.
-   `max_tokens`: (Optional) The maximum number of tokens to generate.
-   `top_p`: (Optional) Nucleus sampling probability.
-   `stream`: (Optional) Boolean to stream responses (usually `false` for this API).
-   `messages`: (Optional) Array of message objects to specify a custom prompt. See below.

**Custom Prompt:**

You can provide a custom system prompt via the `messages` array. The server will **append** its multilingual instructions to your custom prompt:

```json
"llm_config": {
  "model": "google/gemini-2.0-flash-exp:free",
  "messages": [
    {
      "role": "system",
      "content": "Your custom prompt here. Be concise and creative."
    }
  ]
}
```

If you don't provide a `messages` array, the server creates one with default multilingual instructions based on `llm_output_language`.

**Request Examples:**

1. Basic request (GitHub source):

```json
{
  "max_repos": 5,
  "resource": "github",
  "since": "weekly",
  "spoken_language_code": "en",
  "llm_provider": "mistral_api",
  "llm_config": {
    "model": "mistral-tiny"
  }
}
```

2. Request with Chutes.ai provider:

```json
{
  "max_repos": 3,
  "since": "daily",
  "spoken_language_code": "en",
  "llm_provider": "chutes",
  "llm_output_language": "uk,en",
  "llm_config": {
    "model": "moonshotai/Kimi-K2-Instruct-0905",
    "temperature": 0.5,
    "max_tokens": 1024
  }
}
```

3. Multilingual request with OpenRouter:

```json
{
  "max_repos": 5,
  "since": "weekly",
  "spoken_language_code": "en",
  "llm_output_language": "en,uk,fr",
  "llm_provider": "openrouter",
  "llm_config": {
    "model": "google/gemini-2.0-flash-exp:free"
  }
}
```

4. Request using OssInsight source:

```json
{
  "max_repos": 10,
  "resource": "ossinsight",
  "period": "past_month",
  "language": "Python",
  "llm_provider": "mistral_api",
  "llm_config": {
    "model": "mistral-small"
  }
}
```

**Response Example:**

```json
{
  "status": "ok",
  "added": ["https://github.com/example/repo1", "https://github.com/example/repo2"],
  "dont_added": ["https://github.com/example/repo3"]
}
```

### /api/get-repository/

**Endpoint:** `/think-root/api/get-repository/`

**Method:** `POST`

**Description:** This endpoint retrieves a list of repositories based on the provided limit, posted status, and sorting preferences. Results can be sorted by different fields and directions, with special handling for null values in publication dates. By default, if `text_language` is omitted, the endpoint returns the raw multilingual text exactly as stored, e.g., "===(en)text===(uk)текст===". If `text_language` is provided (e.g., "en" or "uk"), the endpoint returns only that language’s text. If the requested language is not available, the response preserves the existing error handling behavior.

**Curl Example:**

```bash
curl -X POST \
  'http://localhost:8080/think-root/api/get-repository/' \
  -H 'Authorization: Bearer <BEARER_TOKEN>' \
  -H 'Content-Type: application/json' \
  -d '{
    "limit": 1,
    "posted": false,
    "sort_by": "date_added",
    "sort_order": "DESC",
    "text_language": "uk"
  }'
```

**Request Parameters:**

| Parameter         | Type    | Required | Description                                                                                                                                                                                                                                                                                                   |
| ----------------- | ------- | -------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `limit`         | integer | No       | Maximum number of repositories to return. Set to 0 to either return all records (if page and page_size are not specified) or use pagination mode (if page or page_size are specified).                                                                                                                        |
| `posted`        | boolean | No       | Filter repositories by posted status. If not specified and limit is 0, returns all records regardless of posted status.                                                                                                                                                                                       |
| `sort_by`       | string  | No       | Field to sort results by. Valid values:`id`, `date_added`, `date_posted`. Default: `date_added` for unposted repositories, `date_posted` for posted repositories. When sorting by `date_posted`, repositories without a publication date (null) will be displayed according to the sorting order. |
| `sort_order`    | string  | No       | Order of sorting. Valid values:`ASC` (ascending), `DESC` (descending). Default: `DESC`.                                                                                                                                                                                                                 |
| `page`          | integer | No       | Page number for pagination (1-based). If not specified along with page_size and limit is 0, all records will be returned without pagination.                                                                                                                                                                  |
| `page_size`     | integer | No       | Number of items per page. If not specified along with page and limit is 0, all records will be returned without pagination.                                                                                                                                                                                   |
| `text_language` | string  | No       | Optional. When omitted, raw multilingual text is returned in the original format, for example "===(en)text===(uk)text===". When provided (e.g., "en", "uk"), the API extracts and returns only the specified language’s text.                                                                                |

**Request Examples:**

1. Get all records without pagination:

```json
{
  "limit": 0,
  "posted": null,
  "sort_by": "date_added",
  "sort_order": "DESC"
}
```

2. Get records with pagination:

```json
{
  "limit": 0,
  "posted": null,
  "sort_by": "date_added",
  "sort_order": "DESC",
  "page": 1,
  "page_size": 10
}
```

3. Get limited number of records:

```json
{
  "limit": 5,
  "posted": true,
  "sort_by": "date_posted",
  "sort_order": "DESC"
}
```

4. **Example: Specifying text_language: 'uk' returns only the Ukrainian text**

```json
{
  "limit": 10,
  "text_language": "uk"
}
```

5. ** Get English text with pagination:**

```json
{
  "limit": 0,
  "page": 1,
  "page_size": 5,
  "text_language": "en"
}
```

6. **Example: Without text_language (raw multilingual text)**

```json
{
  "limit": 10
}
```

Returns raw multilingual text segments in the original format, e.g., "===(en)text===(uk)text===".

**Pagination Details:**

- When `limit` is 0:
  - If neither `page` nor `page_size` are specified, returns all matching records without pagination. In this case, `page=0`, `page_size=0`, `total_pages=1`, and `total_items` equals the count of all matching records.
  - If either `page` or `page_size` are specified, pagination mode is used.
- When `limit` > 0:
  - Pagination mode is used. If `page` < 1, it defaults to 1. If `page_size` < 1, it defaults to 10.
- Response always includes:
  - `page`: Current page number (0 when returning all records without pagination; otherwise the active page number).
  - `page_size`: Number of items per page (0 when returning all records without pagination; otherwise the active page size).
  - `total_pages`: Total number of pages (1 when returning all records without pagination).
  - `total_items`: Total number of items matching the query.

**Sorting Behavior:**

- When sorting by `date_posted`:
  - If `sort_order` = `ASC`: entries with null values are shown first, followed by dates in ascending order
  - If `sort_order` = `DESC`: entries with dates are shown in descending order first, followed by those with null values
- When sorting by `date_added` or `id`: standard ascending or descending sort
- If `sort_by` is not specified, `date_posted` is used for posted=true and `date_added` for posted=false
- If `sort_order` is not specified, `DESC` is used as default

**Response Example:**

```json
{
  "status": "ok",
  "message": "Repositories fetched successfully",
  "data": {
    "all": 50,
    "posted": 20,
    "unposted": 30,
    "items": [
      {
        "id": 1,
        "posted": false,
        "url": "https://github.com/example/repo",
        "text": "Repository description here.",
        "date_added": "2025-03-20T15:30:45Z",
        "date_posted": null
      }
    ],
    "page": 1,
    "page_size": 10,
    "total_pages": 5,
    "total_items": 50
  }
}
```

Additional Response Example (raw multilingual text):

```json
{
  "text": "===(en)An open-source project===(uk)Відкритий проект===",
  "... other fields ...": "..."
}
```

Note: This indicates raw multilingual text as stored.

### /api/update-posted/

**Endpoint:** `/think-root/api/update-posted/`

**Method:** `PATCH`

**Description:** This endpoint updates the posted status of a repository identified by its URL. Note: when the URL does not exist, the current implementation returns 500 with a generic error rather than 404. Setting `posted=true` sets `date_posted` to current time; `posted=false` clears `date_posted` per [database.UpdatePostedStatusByURL()](database/update.go:8).

**Curl Example:**

```bash
curl -X PATCH \
  'http://localhost:8080/think-root/api/update-posted/' \
  -H 'Authorization: Bearer <BEARER_TOKEN>' \
  -H 'Content-Type: application/json' \
  -d '{
    "url": "https://github.com/example/repo",
    "posted": true
  }'
```

**Request Example:**

```json
{
  "url": "https://github.com/example/repo",
  "posted": true
}
```

**Response Example:**

```json
{
  "status": "ok",
  "message": "Posted status updated successfully"
}
```

### /api/update-repository-text/

**Endpoint:** `/think-root/api/update-repository-text/`

**Method:** `PATCH`

**Description:** Updates the repository text with two modes:

- Full replace when `text_language` is omitted.
- Strict language-specific update when `text_language` is provided. If the specified language does not exist in the existing multilingual content, returns 422 Unprocessable Entity.

**Request Schema:**

- Exactly one of `id` or `url` must be provided
- `text` is required
- `text_language` is optional (language code). When provided, triggers language-specific update.
- Language code validation is performed by [language.ValidateLanguageCodes()](server/language_validator.go:26)

**Curl and JSON Examples:**

1. Full replace (no text_language):

```bash
curl -X PATCH \
  'http://localhost:8080/think-root/api/update-repository-text/' \
  -H 'Authorization: Bearer <BEARER_TOKEN>' \
  -H 'Content-Type: application/json' \
  -d '{
    "id": 1172,
    "text": "Updated text content via ID"
  }'
```

Result in DB: "Updated text content via ID"
Response includes `available_languages` (for plain: ["uk"]) and omits `updated_language`.

2. Language update on existing multilingual:

```bash
curl -X PATCH \
  'http://localhost:8080/think-root/api/update-repository-text/' \
  -H 'Authorization: Bearer <BEARER_TOKEN>' \
  -H 'Content-Type: application/json' \
  -d '{
    "id": 1172,
    "text": "Updated text content via ID",
    "text_language": "en"
  }'
```

Result in DB: only the `en` segment updated; other segments unchanged.

3. Language update on plain existing text:

```bash
curl -X PATCH \
  'http://localhost:8080/think-root/api/update-repository-text/' \
  -H 'Authorization: Bearer <BEARER_TOKEN>' \
  -H 'Content-Type: application/json' \
  -d '{
    "id": 1172,
    "text": "тут якийсь текст",
    "text_language": "uk"
  }'
```

Result in DB: "===(uk)тут якийсь текст==="

4. Error when language missing in existing multilingual:

```bash
curl -X PATCH \
  'http://localhost:8080/think-root/api/update-repository-text/' \
  -H 'Authorization: Bearer <BEARER_TOKEN>' \
  -H 'Content-Type: application/json' \
  -d '{
    "id": 1172,
    "text": "KURWA! Ja perdoly jjajajajaj.",
    "text_language": "pl"
  }'
```

Response: `422 Unprocessable Entity` with message `language 'pl' not found in existing content`.

**Request Parameters:**

| Parameter         | Type    | Required | Description                                                                             |
| ----------------- | ------- | -------- | --------------------------------------------------------------------------------------- |
| `id`            | integer | No\*     | Repository ID (positive integer)                                                        |
| `url`           | string  | No\*     | Repository URL (non-empty string)                                                       |
| `text`          | string  | Yes      | New text content (1-1000 characters, valid UTF-8)                                       |
| `text_language` | string  | No       | Optional language code. When provided, performs a language-specific update (validated). |

\*Exactly one of `id` or `url` must be provided.

**Validation Rules:**

- Exactly one identifier (`id` or `url`) must be provided
- `text` is required and non-empty
- `text` length ≤ 1000 characters
- `text` must be valid UTF-8
- `text_language` validated via [language.ValidateLanguageCodes()](server/language_validator.go:25)

**Response Fields:**

- `status` and `message`
- `data.id`, `data.url`, `data.text` (final text stored)
- `data.updated_language` present only when `text_language` is provided
- `data.available_languages` via [multilingual.GetAvailableLanguages()](server/multilingual_helper.go:156)
- `data.updated_at`

**Success Response Examples:**

Full replace (plain text):

```json
{
  "status": "ok",
  "message": "Repository text updated successfully",
  "data": {
    "id": 1172,
    "url": "https://github.com/example/repo",
    "text": "Updated text content via ID",
    "available_languages": ["uk"],
    "updated_at": "2025-06-22T15:00:00Z"
  }
}
```

Language-specific update:

```json
{
  "status": "ok",
  "message": "Repository text updated successfully",
  "data": {
    "id": 1172,
    "url": "https://github.com/example/repo",
    "text": "Updated text content via ID",
    "updated_language": "en",
    "available_languages": ["en", "uk"],
    "updated_at": "2025-06-22T15:00:00Z"
  }
}
```

**Error Response Example (missing language in existing multilingual):**

```json
{
  "status": "error",
  "message": "language 'pl' not found in existing content"
}
```

### /api/delete-repository/

**Endpoint:** `/think-root/api/delete-repository/`

**Method:** `DELETE`

**Description:** This endpoint deletes a repository from the database. The repository can be identified by either its unique ID or URL.

**Curl Examples:**

Delete by ID:

```bash
curl -X DELETE \
  'http://localhost:8080/think-root/api/delete-repository/' \
  -H 'Authorization: Bearer <BEARER_TOKEN>' \
  -H 'Content-Type: application/json' \
  -d '{
    "id": 123
  }'
```

Delete by URL:

```bash
curl -X DELETE \
  'http://localhost:8080/think-root/api/delete-repository/' \
  -H 'Authorization: Bearer <BEARER_TOKEN>' \
  -H 'Content-Type: application/json' \
  -d '{
    "url": "https://github.com/example/repo"
  }'
```

**Request Parameters:**

| Parameter | Type    | Required | Description                       |
| --------- | ------- | -------- | --------------------------------- |
| `id`    | integer | No\*     | Repository ID (positive integer)  |
| `url`   | string  | No\*     | Repository URL (non-empty string) |

\*Either `id` or `url` must be provided, but not both.

**Request Examples:**

1. Delete by ID:

```json
{
  "id": 123
}
```

2. Delete by URL:

```json
{
  "url": "https://github.com/example/awesome-project"
}
```

**Validation Rules:**

- Exactly one identifier (`id` or `url`) must be provided
- ID must be a positive integer if provided
- URL must be a non-empty string if provided

**Status Codes:**

- 200: Success - Repository deleted
- 400: Bad Request - Validation errors
- 401: Unauthorized - Invalid or missing Bearer token
- 405: Method Not Allowed - Wrong HTTP method
- 500: Internal Server Error - Database or server error (including when URL not found in update-posted)

**Success Response Example:**

```json
{
  "status": "ok",
  "message": "Repository deleted successfully"
}
```

**Error Response Examples:**

```json
{
  "status": "error",
  "message": "Either id or url must be provided"
}
```

```json
{
  "status": "error",
  "message": "Provide either id or url, not both"
}
```

```json
{
  "status": "error",
  "message": "repository with ID 123 not found"
}
```

```json
{
  "status": "error",
  "message": "repository with URL https://github.com/example/repo not found"
}
```
