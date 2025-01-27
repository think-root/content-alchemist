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

**Description:** This endpoint is used to manually generate something. The exact details of what is generated are not provided in the `manual_generate.go` file.

### /api/auto-generate/

**Endpoint:** `/think-root/api/auto-generate/`

**Method:** `POST`

**Description:** This endpoint is used to automatically generate something. The exact details of what is generated are not provided in the `auto_generate.go` file.

### /api/get-repository/

**Endpoint:** `/think-root/api/get-repository/`

**Method:** `POST`

**Request Body:**
```json
{
  "limit": int,
  "posted": bool
}
```

**Response:**
```json
{
  "status": "ok",
  "message": "Repositories fetched successfully",
  "data": {
    "all": int,
    "posted": int,
    "unposted": int,
    "items": [
      {
        "id": string,
        "posted": bool,
        "url": string,
        "text": string
      }
    ]
  }
}
```

**Description:** This endpoint retrieves a list of repositories based on the provided limit and posted status. It returns the total count of all, posted, and unposted repositories, along with the repository details.

### /api/update-posted/

**Endpoint:** `/think-root/api/update-posted/`

**Method:** `POST`

**Request Body:**
```json
{
  "url": string,
  "posted": bool
}
```

**Response:**
```json
{
  "status": "ok",
  "message": "Posted status updated successfully"
}
```

**Description:** This endpoint updates the posted status of a repository identified by its URL.