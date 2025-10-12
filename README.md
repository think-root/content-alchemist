<h1 align="center">Content Alchemist</h1>

<div align="center">

![License](https://img.shields.io/github/license/think-root/content-alchemist?style=flat-square&color=blue)
[![Go Report Card](https://goreportcard.com/badge/github.com/think-root/content-alchemist?style=flat-square)](https://goreportcard.com/report/github.com/think-root/content-alchemist)
[![Go Version](https://img.shields.io/github/go-mod/go-version/think-root/content-alchemist?style=flat-square)](https://github.com/think-root/content-alchemist)
[![Deploy Status](https://img.shields.io/github/actions/workflow/status/think-root/content-maestro/deploy.yml?branch=main&label=Deploy&style=flat-square)](https://github.com/think-root/content-alchemist/actions/workflows/deploy.yml)
[![Version](https://img.shields.io/github/v/release/think-root/content-alchemist?style=flat-square)](https://github.com/think-root/content-alchemist/releases)
[![Changelog](https://img.shields.io/badge/changelog-view-blue?style=flat-square)](CHANGELOG.md)


<img src="baner.png" alt="baner">

</div>

## 📖 Description

This is a ready-made solution in the form of an API server that generates social media posts containing descriptions of GitHub repositories using AI and stores them for later use. It is a standalone solution that you can manage using your own utility or by leveraging existing tools available in this organization's repositories, for example [content-maestro](https://github.com/think-root/content-maestro)

### Key Features

- RESTful API for automatic generation of repository descriptions based on GitHub trends
- RESTful API for manual generation of repository descriptions by specifying the repository URL
- RESTful API for content management and text editing
- Generate and retrieve content in multiple languages
- Database storage with PostgreSQL
- Support for multiple AI providers (Mistral AI, OpenAI, OpenRouter, Chutes.ai)


### Technology Stack

- Go 1.23
- PostgreSQL 16
- Docker & Docker Compose
- Multiple AI Providers:
  - Mistral AI [API](https://docs.mistral.ai/api/)
  - OpenAI [API](https://platform.openai.com/docs/api-reference)
  - OpenRouter [API](https://openrouter.ai/docs)
  - Chutes.ai [API](https://chutes.ai)

## 🛠️ Installation

### Prerequisites

- Docker v20.10.0 or later
- Docker Compose v2.0.0 or later
- API key for at least one of the supported AI providers

## ▶️ How to run

### Requirements

- [docker](https://docs.docker.com/engine/install/) or/and [docker-compose](https://docs.docker.com/compose/install/)
- [mistral ai](https://console.mistral.ai/api-keys/) api key (or other supported provider)

### Clone repo

```shell
git clone https://github.com/think-root/content-alchemist.git
```

### Config

Create a **.env** file in the app root directory and ensure you have:

1. Created an account with at least one of the supported AI providers
2. Generated API key(s) for the provider(s) you plan to use
3. Set up your PostgreSQL instance

```properties
# Required for database and API protection
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_EXTERNAL_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=your_secure_password
POSTGRES_DB=think-root-db
BEARER_TOKEN=<your token for API protection>

# Mistral AI Provider Settings
MISTRAL_TOKEN=<mistral api key>

# OpenAI Provider Settings (optional)
OPENAI_TOKEN=<openai api key>

# OpenRouter Provider Settings (optional)
OPENROUTER_TOKEN=<openrouter api key>

# Chutes.ai Provider Settings (optional)
CHUTES_API_TOKEN=<chutes api key>
```


### Deploy

1. Create Docker network:

   ```bash
   docker network create think-root-network
   ```

2. Deploy PostgreSQL:

   ```bash
   docker compose -f docker-compose.db.yml up -d
   ```

3. Deploy content-alchemist:
   ```bash
   docker compose -f docker-compose.app.yml up -d
   ```

## 🌐 Multilingual Support

Content Alchemist now supports generating and retrieving repository descriptions in multiple languages. This feature allows you to:

- **Generate multilingual content**: Create repository descriptions in multiple languages simultaneously
- **Language-specific retrieval**: Get repository texts in specific languages
- **Flexible text formats**: Support for single language and multilingual text formats
- **Automatic validation**: Language codes are validated against ISO 639-1 standards

### Quick Examples

**Generate content in multiple languages:**
```bash
curl -X POST \
  'http://localhost:8080/think-root/api/manual-generate/' \
  -H 'Authorization: Bearer <BEARER_TOKEN>' \
  -H 'Content-Type: application/json' \
  -d '{
    "url": "https://github.com/facebook/react",
    "llm_output_language": "en,uk,fr"
  }'
```

**Retrieve content in Ukrainian:**
```bash
curl -X POST \
  'http://localhost:8080/think-root/api/get-repository/' \
  -H 'Authorization: Bearer <BEARER_TOKEN>' \
  -H 'Content-Type: application/json' \
  -d '{
    "limit": 5,
    "text_language": "uk"
  }'
```

---

## 🔌 API

> [!IMPORTANT]
> All API requests must include an Authorization header in the following format:
> Authorization: Bearer <BEARER_TOKEN>
> 
> Rate Limit: 60 requests per minute per IP address
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
    "llm_output_language": "en,uk,fr"
  }'
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `url` | string | Yes | GitHub repository URL |
| `llm_output_language` | string | No | Comma-separated language codes (e.g., "en,uk,fr"). Default: "uk" |

**Request Examples:**

1. Basic request (Ukrainian only):
```json
{
  "url": "https://github.com/example/repo"
}
```

2. Request with specific AI provider:
```json
{
  "url": "https://github.com/example/repo",
  "llm_provider": "chutes",
  "llm_output_language": "en"
}
```

3. Multilingual request:
```json
{
  "url": "https://github.com/example/repo",
  "llm_output_language": "en,uk,fr"
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

---

### /api/auto-generate/

**Endpoint:** `/think-root/api/auto-generate/`

**Method:** `POST`

**Description:** This endpoint is used to automatically parse trending repositories and generate description based on certain parameters. It also adds the generated posts to the database. Supports multilingual text generation.

**Curl Example:**

```bash
curl -X POST \
  'http://localhost:8080/think-root/api/auto-generate/' \
  -H 'Authorization: Bearer <BEARER_TOKEN>' \
  -H 'Content-Type: application/json' \
  -d '{
    "max_repos": 5,
    "since": "weekly",
    "spoken_language_code": "en",
    "llm_output_language": "en,uk,fr",
    "llm_provider": "chutes"
  }'
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `max_repos` | integer | No | Maximum number of repositories to process |
| `since` | string | No | Time period for trending repos ("daily", "weekly", "monthly") |
| `spoken_language_code` | string | No | Programming language filter |
| `llm_output_language` | string | No | Comma-separated language codes for output (e.g., "en,uk,fr"). Default: "uk" |

**Request Examples:**

1. Basic request:
```json
{
  "max_repos": 5,
  "since": "weekly",
  "spoken_language_code": "en"
}
```

2. Request with Chutes.ai provider:
```json
{
  "max_repos": 3,
  "since": "daily",
  "spoken_language_code": "en",
  "llm_provider": "chutes",
  "llm_output_language": "uk,en"
}
```

3. Multilingual request:
```json
{
  "max_repos": 5,
  "since": "weekly",
  "spoken_language_code": "en",
  "llm_output_language": "en,uk,fr"
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

---

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

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `limit` | integer | No | Maximum number of repositories to return. Set to 0 to either return all records (if page and page_size are not specified) or use pagination mode (if page or page_size are specified). |
| `posted` | boolean | No | Filter repositories by posted status. If not specified and limit is 0, returns all records regardless of posted status. |
| `sort_by` | string | No | Field to sort results by. Valid values: `id`, `date_added`, `date_posted`. Default: `date_added` for unposted repositories, `date_posted` for posted repositories. When sorting by `date_posted`, repositories without a publication date (null) will be displayed according to the sorting order. |
| `sort_order` | string | No | Order of sorting. Valid values: `ASC` (ascending), `DESC` (descending). Default: `DESC`. |
| `page` | integer | No | Page number for pagination (1-based). If not specified along with page_size and limit is 0, all records will be returned without pagination. |
| `page_size` | integer | No | Number of items per page. If not specified along with page and limit is 0, all records will be returned without pagination. |
| `text_language` | string | No | Optional. When omitted, raw multilingual text is returned in the original format, for example "===(en)text===(uk)text===". When provided (e.g., "en", "uk"), the API extracts and returns only the specified language’s text. |

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
  - If neither `page` nor `page_size` are specified, returns all matching records without pagination
  - If either `page` or `page_size` are specified, uses pagination mode
- When `limit` > 0:
  - Uses the specified limit with pagination
  - Default page size is 10 if not specified
- Response always includes:
  - `page`: Current page number (1 when returning all records)
  - `page_size`: Number of items per page (equal to total items when returning all records)
  - `total_pages`: Total number of pages (1 when returning all records)
  - `total_items`: Total number of items matching the query
- If `page` is less than 1, it defaults to 1
- If `page_size` is less than 1, it defaults to 10

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

---

### /api/update-posted/

**Endpoint:** `/think-root/api/update-posted/`

**Method:** `PATCH`

**Description:** This endpoint updates the posted status of a repository identified by its URL.

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

---

### /api/update-repository-text/

**Endpoint:** `/think-root/api/update-repository-text/`

**Method:** `PATCH`

**Description:** This endpoint updates the text field of a repository. The repository can be identified by either its unique ID or URL.

**Curl Examples:**

Update by ID:
```bash
curl -X PATCH \
  'http://localhost:8080/think-root/api/update-repository-text/' \
  -H 'Authorization: Bearer <BEARER_TOKEN>' \
  -H 'Content-Type: application/json' \
  -d '{
    "id": 123,
    "text": "Updated repository description"
  }'
```

Update by URL:
```bash
curl -X PATCH \
  'http://localhost:8080/think-root/api/update-repository-text/' \
  -H 'Authorization: Bearer <BEARER_TOKEN>' \
  -H 'Content-Type: application/json' \
  -d '{
    "url": "https://github.com/example/repo",
    "text": "Updated repository description"
  }'
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | integer | No* | Repository ID (positive integer) |
| `url` | string | No* | Repository URL (non-empty string) |
| `text` | string | Yes | New text content (1-1000 characters, valid UTF-8) |

*Either `id` or `url` must be provided, but not both.

**Request Examples:**

1. Update by ID:
```json
{
  "id": 123,
  "text": "Updated repository description"
}
```

2. Update by URL:
```json
{
  "url": "https://github.com/example/awesome-project",
  "text": "This is an awesome project with new features"
}
```

**Validation Rules:**
- Exactly one identifier (`id` or `url`) must be provided
- Text field is required and cannot be empty
- Text length must not exceed 1000 characters
- Text must be valid UTF-8 encoding
- ID must be a positive integer if provided
- URL must be a non-empty string if provided

**Status Codes:**
- 200: Success - Repository text updated
- 400: Bad Request - Validation errors
- 401: Unauthorized - Invalid or missing Bearer token
- 404: Not Found - Repository not found
- 500: Internal Server Error - Database or server error

**Success Response Example:**

```json
{
  "status": "ok",
  "message": "Repository text updated successfully",
  "data": {
    "id": 123,
    "url": "https://github.com/example/repo",
    "text": "Updated repository description",
    "updated_at": "2025-06-22T15:00:00Z"
  }
}
```

**Error Response Examples:**

```json
{
  "status": "error",
  "message": "Text field is required and cannot be empty"
}
```

```json
{
  "status": "error",
  "message": "Either id or url must be provided"
}
```

```json
{
  "status": "error",
  "message": "repository with ID 123 not found"
}
```

---

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
    "url": "https://github.com/example/repo",
    "llm_provider": "chutes"
  }'
```

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | integer | No* | Repository ID (positive integer) |
| `url` | string | No* | Repository URL (non-empty string) |

*Either `id` or `url` must be provided, but not both.

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
- 404: Not Found - Repository not found
- 405: Method Not Allowed - Wrong HTTP method
- 500: Internal Server Error - Database or server error

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

## 🤖 AI Provider Examples

### Using Chutes.ai Provider

The Chutes.ai provider supports the moonshotai/Kimi-K2-Instruct-0905 model and other OpenAI-compatible models. Here's how to use it:

**Environment Setup:**
```bash
export CHUTES_API_TOKEN="your_chutes_api_token_here"
```

**Basic Usage Example:**
```bash
curl -X POST \
  'http://localhost:8080/think-root/api/manual-generate/' \
  -H 'Authorization: Bearer <BEARER_TOKEN>' \
  -H 'Content-Type: application/json' \
  -d '{
    "url": "https://github.com/think-root/content-alchemist",
    "llm_provider": "chutes",
    "llm_output_language": "en"
  }'
```

**Multilingual Generation with Chutes:**
```bash
curl -X POST \
  'http://localhost:8080/think-root/api/auto-generate/' \
  -H 'Authorization: Bearer <BEARER_TOKEN>' \
  -H 'Content-Type: application/json' \
  -d '{
    "max_repos": 5,
    "since": "weekly",
    "llm_provider": "chutes",
    "llm_output_language": "en,uk,fr"
  }'
```

### Provider Comparison

| Provider | Model Examples | Best For |
|----------|---------------|----------|
| Mistral AI | mistral-large-latest | European languages, technical content |
| OpenAI | gpt-4-turbo | General purpose, creative content |
| OpenRouter | claude-3, gpt-4 | Flexibility, multiple models |
| **Chutes.ai** | moonshotai/Kimi-K2-Instruct-0905 | **New!** Advanced reasoning, multilingual support |

## 🧑‍💻 Development Setup

1. Install Go 1.23 or later
2. Install PostgreSQL 16 or later
3. Clone the repository
4. Install dependencies: `go mod download`

### Running Locally

1. Set up your .env file with PostgreSQL connection details
2. Start PostgreSQL
3. Run the server:
   ```bash
   go run ./cmd/server/main.go
   ```

### Building

```bash
go build -o content-alchemist ./cmd/server/main.go
```
