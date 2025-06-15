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

## Description

This is a ready-made solution in the form of an API server that generates social media posts containing descriptions of GitHub repositories using [AI](https://mistral.ai/) 🤖 and stores them for later use. It is a standalone solution that you can manage using your own utility or by leveraging existing tools available in this organization's repositories, for example [content-maestro](https://github.com/think-root/content-maestro)

### Key Features

- RESTful API for automatic generation of repository descriptions based on GitHub trends
- RESTful API for manual generation of repository descriptions by specifying the repository URL
- RESTful API for content management
- Database storage
- Support for multiple AI providers (Mistral AI, OpenAI, OpenRouter)


### Technology Stack

- Go 1.23
- PostgreSQL 16
- Docker & Docker Compose
- Multiple AI Providers:
  - Mistral AI [API](https://docs.mistral.ai/api/)
  - OpenAI [API](https://platform.openai.com/docs/api-reference)
  - OpenRouter [API](https://openrouter.ai/docs)

## Installation

### Prerequisites

- Docker v20.10.0 or later
- Docker Compose v2.0.0 or later
- API key for at least one of the supported AI providers

## How to run

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
3. If using Mistral AI agent, created and configured a Mistral agent
4. Set up your PostgreSQL instance

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
MISTRAL_AGENT=<get agent api id https://console.mistral.ai/build/agents>

# OpenAI Provider Settings (optional)
OPENAI_TOKEN=<openai api key>

# OpenRouter Provider Settings (optional)
OPENROUTER_TOKEN=<openrouter api key>
```

### Mistral AI Agent configuration

- create mistral [agent](https://console.mistral.ai/build/agents) (model: mistral large 2.1, temperature: 0.1)
- system prompt (UA, translate yourself if necessary):

```text
Ти слухняний і корисний помічник, який суворо дотримується усіх вимог які описані нижче. Твоя основна задача — генерувати короткі описи GitHub репозиторіїв, українською мовою, з текстів, які будуть надані. Ці тексти є описами (README) GitHub репозиторіїв. При генеруванні тексту обов'язково дотримуйся таких вимог:

1. В описі має бути не більше трьох ключових функцій репозиторію.
2. Використовуй простий і зрозумілий стиль тексту без перерахувань. Інформацію про функції репозиторію вплітай у зв'язний текст.
3. Не згадуй інформацію про сумісність, платформи, авторів.
4. Не використовуй розмітку тексту, таку як HTML теги, Markdown розмітку тощо.
5. Опис має бути лаконічним і точним, обсягом від 200 до 400 символів (з урахуванням пробілів та інших символів).
6. Якщо зустрічаються технічні терміни, такі як назви мов програмування, бібліотек, команд або інструментів, видів програмування, залишай їх англійською мовою без перекладу.  
7. Перед генерацією тексту переконайся, що він повністю відповідає усім вищезазначеним вимогам.  

Далі тобі буде надано назву GitHub репозиторію та його README. Твоє завдання — створити чіткий, короткий і зрозумілий опис, який відповідає всім вищезазначеним вимогам.
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

## API

> [!IMPORTANT]
> All API requests must include an Authorization header in the following format:
> Authorization: Bearer <BEARER_TOKEN>
> 
> Rate Limit: 60 requests per minute per IP address
> All endpoints return JSON responses with appropriate HTTP status codes

### /api/manual-generate/

**Endpoint:** `/think-root/api/manual-generate/`

**Method:** `POST`

**Description:** This endpoint is used to manually generate description for a provided repository URL, and add it to the database.

**Curl Example:**

```bash
curl -X POST \
  'http://localhost:9111/think-root/api/manual-generate/' \
  -H 'Authorization: Bearer <BEARER_TOKEN>' \
  -H 'Content-Type: application/json' \
  -d '{"url": "https://github.com/example/repo"}'
```

**Request Example:**

```json
{
  "url": "https://github.com/example/repo"
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

**Description:** This endpoint is used to automatically parse trending repositories and generate description based on certain parameters. It also adds the generated posts to the database.

**Curl Example:**

```bash
curl -X POST \
  'http://localhost:9111/think-root/api/auto-generate/' \
  -H 'Authorization: Bearer <BEARER_TOKEN>' \
  -H 'Content-Type: application/json' \
  -d '{
    "max_repos": 5,
    "since": "weekly",
    "spoken_language_code": "en"
  }'
```

**Request Example:**

```json
{
  "max_repos": 5,
  "since": "weekly",
  "spoken_language_code": "en"
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

**Description:** This endpoint retrieves a list of repositories based on the provided limit, posted status, and sorting preferences. Results can be sorted by different fields and directions, with special handling for null values in publication dates.

**Curl Example:**

```bash
curl -X POST \
  'http://localhost:9111/think-root/api/get-repository/' \
  -H 'Authorization: Bearer <BEARER_TOKEN>' \
  -H 'Content-Type: application/json' \
  -d '{
    "limit": 1,
    "posted": false,
    "sort_by": "date_added",
    "sort_order": "DESC"
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

---

### /api/update-posted/

**Endpoint:** `/think-root/api/update-posted/`

**Method:** `PATCH`

**Description:** This endpoint updates the posted status of a repository identified by its URL.

**Curl Example:**

```bash
curl -X PATCH \
  'http://localhost:9111/think-root/api/update-posted/' \
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

## Contribution

### Development Setup

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
