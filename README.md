<h1 align="center">Content Alchemist</h1>

<div align="center">

![License](https://img.shields.io/github/license/think-root/content-alchemist?style=flat-square)
[![Go Version](https://img.shields.io/github/go-mod/go-version/think-root/content-alchemist)](https://github.com/think-root/content-alchemist)
[![Version](https://img.shields.io/github/v/release/think-root/content-alchemist)](https://github.com/think-root/content-alchemist/releases)
[![Changelog](https://img.shields.io/badge/changelog-view-blue)](CHANGELOG.md)
[![Deploy Status](https://github.com/think-root/content-alchemist/workflows/Deploy%20content-alchemist/badge.svg)](https://github.com/think-root/content-alchemist/actions/workflows/deploy.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/think-root/content-alchemist)](https://goreportcard.com/report/github.com/think-root/content-alchemist)

<img src="baner.png" alt="baner">

</div>

## Description

This is a ready-made solution in the form of an API server that generates social media posts containing descriptions of GitHub repositories using [AI](https://mistral.ai/) 🤖 and stores them for later use. It is a standalone solution that you can manage using your own utility or by leveraging existing tools available in this organization's repositories, for example [content-maestro](https://github.com/think-root/content-maestro)

### Key Features

- RESTful API for automatic generation of repository descriptions based on GitHub trends
- RESTful API for manual generation of repository descriptions by specifying the repository URL
- RESTful API for content management
- Database storage


### Technology Stack

- Go 1.23
- MariaDB/MySQL
- Docker & Docker Compose
- Mistral AI [API](https://docs.mistral.ai/api/)

## Installation

### Prerequisites

- Docker v20.10.0 or later
- Docker Compose v2.0.0 or later
- Mistral AI API key

## How to run

### Requirements

- [docker](https://docs.docker.com/engine/install/) or/and [docker-compose](https://docs.docker.com/compose/install/)
- [mistral ai](https://console.mistral.ai/api-keys/) api key

### Preparation

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

### Clone repo

```shell
git clone https://github.com/think-root/content-alchemist.git
```

### Config

create a **.env** file in the app root directory

Before creating the .env file, ensure you have:

1. Created a Mistral AI account
2. Generated an API key
3. Created and configured a Mistral agent
4. Set up your MariaDB instance

```properties
MISTRAL_TOKEN=<mistral api key>
MISTRAL_AGENT=<get agent api id https://console.mistral.ai/build/agents>
DB_CONNECTION=<db connection string e.g. user:password@tcp(localhost:3306)>
BEARER_TOKEN=<your token for API protection>
```

### Deploy

1. Create Docker network:

   ```bash
   docker network create think-root-network
   ```

2. Deploy MariaDB:

   ```bash
   docker run -d --name mariadb --network think-root-network -e MYSQL_ROOT_PASSWORD=your_password -p 3306:3306 mariadb:latest
   ```

3. Deploy content-alchemist:
   ```bash
   docker compose up -d
   ```

## API

```text
All API requests must include an Authorization header in the following format:
Authorization: Bearer <BEARER_TOKEN>
```

Rate Limit: 100 requests per minute per IP address
All endpoints return JSON responses with appropriate HTTP status codes

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
| `limit` | integer | No | Maximum number of repositories to return. Set to 0 for no limit. When limit is 0 and posted is not specified, returns all records regardless of posted status. |
| `posted` | boolean | No | Filter repositories by posted status. If not specified and limit is 0, returns all records regardless of posted status. |
| `sort_by` | string | No | Field to sort results by. Valid values: `id`, `date_added`, `date_posted`. Default: `date_added` for unposted repositories, `date_posted` for posted repositories. When sorting by `date_posted`, repositories without a publication date (null) will be displayed according to the sorting order. |
| `sort_order` | string | No | Order of sorting. Valid values: `ASC` (ascending), `DESC` (descending). Default: `DESC`. |

**Request Example:**

```json
{
  "limit": 0,
  "posted": null,
  "sort_by": "date_added",
  "sort_order": "DESC"
}
```

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
    ]
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
2. Install MariaDB 10.5 or later
3. Clone the repository
4. Install dependencies: `go mod download`

### Running Locally

1. Set up your .env file
2. Start MariaDB
3. Run the server:
   ```bash
   go run ./cmd/server/main.go
   ```

### Building

```bash
go build -o content-alchemist ./cmd/server/main.go
```
