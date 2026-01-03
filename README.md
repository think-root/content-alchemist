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

This is a ready-made solution in the form of an API server that generates social media posts containing descriptions of GitHub repositories using AI and stores them for later use. It is a standalone solution that you can manage using your own utility or by leveraging existing tools available in this organization's repositories, for example [content-maestro](https://github.com/think-root/content-maestro)

### Key Features

- RESTful API for content generation, management, and editing
- Automatic generation from GitHub trends or manual via repository URL
- Multilingual support for content creation and retrieval
- PostgreSQL database storage
- Multiple AI providers (Mistral AI, OpenAI, OpenRouter, Chutes.ai)

### Technology Stack

- Go 1.23
- PostgreSQL 17
- Multiple AI Providers:
  - Mistral AI [API](https://docs.mistral.ai/api/)
  - OpenAI [API](https://platform.openai.com/docs/api-reference)
  - OpenRouter [API](https://openrouter.ai/docs)
  - Chutes.ai [API](https://chutes.ai)

## Installation

### Prerequisites

- Go 1.23 or later
- [PostgreSQL 17](https://www.postgresql.org/download/) or later
- API key for at least one of the supported AI providers

### 1. Clone the repository

```shell
git clone https://github.com/think-root/content-alchemist.git
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Set up PostgreSQL

Install PostgreSQL from the [official website](https://www.postgresql.org/download/) and create a database for the application.

### 4. Configure environment variables

Create a **.env** file in the project root directory:

```properties
# Required for database and API protection
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
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

### 5. Run the application

Start the server directly:

```bash
go run ./cmd/server/main.go
```

Or build and run the binary:

```bash
go build -o content-alchemist ./cmd/server/main.go
./content-alchemist
```

## API

API documentation is available here: [API.md](API.md).

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
