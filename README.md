## Table of Contents
- [Table of Contents](#table-of-contents)
- [Description](#description)
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
  - [run](#run)
  - [build](#build)

## Description

An API server app that generates texts using AI, stores them in a database, and provides functionality to manage the stored records.

I once had the idea to create an app that would manage a Telegram channel by searching for and posting interesting repositories using AI. That idea eventually grew into this project. You can read more details about how the idea was born here: [Part 1](https://drukarnia.com.ua/articles/yak-chatgpt-vede-za-mene-kanal-v-telegram-i-u-nogo-ce-maizhe-vikhodit-chastina-1-VywRW) and [Part 2](https://drukarnia.com.ua/articles/yak-chatgpt-vede-za-mene-kanal-v-telegram-i-u-nogo-ce-maizhe-vikhodit-chastina-2-X9Yjz). (The articles are in Ukrainian, but I think you’ll manage to translate them into your preferred language.)

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

```properties
MISTRAL_TOKEN=<mistral api key>
MISTRAL_AGENT=<get agent api id https://console.mistral.ai/build/agents>
DB_CONNECTION=<db connection string e.g. user:password@tcp(localhost:3306)>
BEARER_TOKEN=<your server token>
```

### Deploy
- create network `docker network create chappie_network` (to allow the app access to the database run mariadb in this network)
- run [mariadb](https://hub.docker.com/_/mariadb) `docker run -d --name mariadb --network chappie_network -e MYSQL_ROOT_PASSWORD=your_password -p 3306:3306 mariadb:latest`
- build `docker build -t chappie_server:latest -f Dockerfile .`
- run `docker run --name chappie_server --restart always --env-file .env -e TZ=Europe/Kiev --network chappie_network chappie_server:latest`
- or via docker compose `docker compose up -d`


## API

### /api/manual-generate/

**Endpoint:** `/think-root/api/manual-generate/`

**Method:** `POST`

**Description:** This endpoint is used to manually generate content for a provided repository URL.

**Request Example:**
```json
{
  "url": "https://github.com/example/repo"
}
```

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

**Description:** This endpoint is used to automatically parse trending repositories and generate description based on certain parameters.

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

**Request Example:**
```json
{
  "limit": 10,
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

**Method:** `POST`

**Description:** This endpoint updates the posted status of a repository identified by its URL.

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

- install [go](https://go.dev/dl/)
- install [mariadb](https://mariadb.org/download/)

### run
```shell
 go run ./cmd/server/main.go  
```

### build
```shell
go build -o chappie_server ./cmd/server/main.go
```

