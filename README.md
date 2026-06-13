# URL Shortener

A REST API that turns long URLs into short, shareable links. Send a URL, get back a short code, and anyone who visits that code gets redirected to the original destination.

## What it's used for

URL shorteners are common whenever you need a compact link. Think of social posts, messages, QR codes, or anywhere character count matters. This project also tracks click counts so you can see how often a link is used.

## How it works

1. **Create a link** - `POST /shorten` accepts a URL, validates it, generates a random 6-character code, and stores the mapping in PostgreSQL.
2. **Redirect** - `GET /:code` looks up the short code, increments the click count, and responds with a `302` redirect to the original URL.
3. **View stats** - `GET /stats/:code` returns metadata for a link (original URL, click count, creation time).

Short codes are alphanumeric (`a-z`, `A-Z`, `0-9`). If a generated code already exists, a new one is created automatically.

## Tech stack


| Layer    | Choice                        |
| -------- | ----------------------------- |
| Language | Go                            |
| HTTP     | `net/http` (standard library) |
| Database | PostgreSQL                    |
| Driver   | `github.com/jackc/pgx/v5`     |
| Config   | `github.com/joho/godotenv`    |


## Project structure

```
├── main.go       Server setup, routing, database connection
├── handlers.go   HTTP handler functions
├── db.go         Database queries
├── models.go     Request/response types
└── schema.sql    PostgreSQL table definition
```

## Getting started

### Prerequisites

- [Go](https://go.dev/dl/) 1.26+
- [PostgreSQL](https://www.postgresql.org/download/) running locally

### 1. Clone the repository

```bash
git clone https://github.com/zainbharde/url-shortener.git
cd url-shortener
```

### 2. Set up the database

Create the database and table:

```bash
psql postgres -f schema.sql
```

Or run the statements in `schema.sql` manually against your PostgreSQL instance.

### 3. Configure environment variables

Create a `.env` file in the project root:

```env
DATABASE_URL=postgres://user:password@localhost:5432/url_shortener?sslmode=disable
```

Replace `user`, `password`, and other values with your local PostgreSQL credentials.

### 4. Install dependencies and run

```bash
go mod download
go run .
```

The server starts on `http://localhost:8080`.

## API usage

**Shorten a URL:**

```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.google.com"}'
```

**Follow a redirect:**

```bash
curl -L http://localhost:8080/abc123
```

**Get link stats:**

```bash
curl http://localhost:8080/stats/abc123
```

