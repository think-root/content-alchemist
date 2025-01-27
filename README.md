# chappie

An API server app that generates texts using AI, stores them in a database, and provides functionality to manage the stored records.

## Table of Contents
- [chappie](#chappie)
  - [Table of Contents](#table-of-contents)
  - [Details](#details)
  - [API](#api)
    - [/api/manual-generate/](#apimanual-generate)
    - [/api/auto-generate/](#apiauto-generate)
    - [/api/get-repository/](#apiget-repository)
    - [/api/update-posted/](#apiupdate-posted)

## Details

I once had the idea to create an app that would manage a Telegram channel by searching for and posting interesting repositories using AI. That idea eventually grew into this project. You can read more details about how the idea was born here: [Part 1](https://drukarnia.com.ua/articles/yak-chatgpt-vede-za-mene-kanal-v-telegram-i-u-nogo-ce-maizhe-vikhodit-chastina-1-VywRW) and [Part 2](https://drukarnia.com.ua/articles/yak-chatgpt-vede-za-mene-kanal-v-telegram-i-u-nogo-ce-maizhe-vikhodit-chastina-2-X9Yjz). (The articles are in Ukrainian, but I think you’ll manage to translate them into your preferred language.)

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

**Description:** This endpoint is used to automatically generate content for trending repositories based on certain parameters.

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

