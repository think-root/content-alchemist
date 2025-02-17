# Chappie Server

[![Go Version](https://img.shields.io/github/go-mod/go-version/Think-Root/chappie_server)](https://github.com/Think-Root/chappie_server)
[![License](https://img.shields.io/github/license/Think-Root/chappie_server)](LICENSE)
[![Version](https://img.shields.io/github/v/release/Think-Root/chappie_server)](https://github.com/Think-Root/chappie_server/releases)
[![Changelog](https://img.shields.io/badge/changelog-view-blue)](CHANGELOG.md)
[![Build Status](https://github.com/Think-Root/chappie_server/workflows/Build/badge.svg)](https://github.com/Think-Root/chappie_server/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/Think-Root/chappie_server)](https://goreportcard.com/report/github.com/Think-Root/chappie_server)

## Table of Contents

- [Chappie Server](#chappie-server)
  - [Table of Contents](#table-of-contents)
  - [Description](#description)
    - [Key Features](#key-features)
    - [Technology Stack](#technology-stack)
  - [Blog Articles](#blog-articles)
  - [Installation](#installation)
    - [Prerequisites](#prerequisites)
  - [How to run](#how-to-run)
    - [Requirements](#requirements)
    - [Preparation](#preparation)
    - [Clone repo](#clone-repo)
    - [Config](#config)
    - [Deploy](#deploy)
  - [API](#api)
    - [/api/manual-generate/](#apimanual-generate)
    - [/api/auto-generate/](#apiauto-generate)
    - [/api/get-repository/](#apiget-repository)
    - [/api/update-posted/](#apiupdate-posted)
  - [Contribution](#contribution)
    - [Development Setup](#development-setup)
    - [Running Locally](#running-locally)
    - [Building](#building)

## Description

Chappie Server is an AI-powered API service that automatically generates and manages repository descriptions for Telegram channels.

### Key Features
- AI-powered text generation using Mistral AI
- Automatic GitHub repository parsing and description generation
- RESTful API for content management
- Database storage for generated descriptions

### Technology Stack
- Go 1.23
- MariaDB/MySQL
- Docker & Docker Compose
- Mistral AI API

## Blog Articles
Read about the project's journey and development:
- [How ChatGPT Manages My Telegram Channel - Part 1](https://drukarnia.com.ua/articles/yak-chatgpt-vede-za-mene-kanal-v-telegram-i-u-nogo-ce-maizhe-vikhodit-chastina-1-VywRW)
- [How ChatGPT Manages My Telegram Channel - Part 2](https://drukarnia.com.ua/articles/yak-chatgpt-vede-za-mene-kanal-v-telegram-i-u-nogo-ce-maizhe-vikhodit-chastina-2-X9Yjz)

(Articles are in Ukrainian)

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
Ти слухняний і корисний помічник, який суворо дотримується усіх нищезазначених вимог. Твоя основна задача — генерувати короткі підсумки українською мовою для текстів, які будуть надані. Ці тексти є описами (README) GitHub репозиторіїв. При генеруванні тексту обов'язково дотримуйся таких вимог:

1. Починай текст словами: "Цей репозиторій".  
2. В описі має бути не більше трьох ключових функцій репозиторію.  
3. Використовуй простий і зрозумілий стиль тексту без перерахувань. Інформацію про функції репозиторію вплітай у зв'язний текст.  
4. Не згадуй інформацію про сумісність, платформи, авторів або назву репозиторію.  
5. Не використовуй розмітку тексту, таку як HTML теги, Markdown розмітку тощо.  
6. Опис має бути лаконічним і точним, обсягом від 300 до 600 символів (з урахуванням пробілів та інших символів).  
7. Якщо зустрічаються технічні терміни, такі як назви мов програмування, бібліотек, команд або інструментів, видів програмування, залишай їх англійською мовою без перекладу.  
8. Перед генерацією тексту переконайся, що він повністю відповідає усім вищезазначеним вимогам.  

Далі тобі буде надано назву GitHub репозиторію та його README. Твоє завдання — створити чіткий і зрозумілий підсумок, який відповідає всім вищезазначеним вимогам.
```

### Clone repo
```shell
git clone https://github.com/Think-Root/chappie_server.git
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
   docker network create chappie_network
   ```

2. Deploy MariaDB:
   ```bash
   docker run -d --name mariadb --network chappie_network -e MYSQL_ROOT_PASSWORD=your_password -p 3306:3306 mariadb:latest
   ```

3. Deploy Chappie Server:
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
  "added": [
    "https://github.com/example/repo"
  ],
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
  "added": [
    "https://github.com/example/repo1",
    "https://github.com/example/repo2"
  ],
  "dont_added": [
    "https://github.com/example/repo3"
  ]
}
```

---

### /api/get-repository/

**Endpoint:** `/think-root/api/get-repository/`

**Method:** `POST`

**Description:** This endpoint retrieves a list of repositories based on the provided limit and posted status.

**Curl Example:**
```bash
curl -X POST \
  'http://localhost:9111/think-root/api/get-repository/' \
  -H 'Authorization: Bearer <BEARER_TOKEN>' \
  -H 'Content-Type: application/json' \
  -d '{
    "limit": 1,
    "posted": false
  }'
```

**Request Example:**
```json
{
  "limit": 1,
  "posted": false
}
```

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
        "id": "1",
        "posted": false,
        "url": "https://github.com/example/repo",
        "text": "Repository description here."
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
go build -o chappie_server ./cmd/server/main.go
```
