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

## 🔌 API

API documentation is available here: [API.md](API.md).

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
