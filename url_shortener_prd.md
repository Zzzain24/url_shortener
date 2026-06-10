# PRD — URL Shortener API

**Project:** url-shortener  
**Language:** Go  
**Database:** PostgreSQL  
**Status:** Not Started

---

## Goal

Build a fully functional URL shortener REST API from scratch in Go — no AI assistance, no tutorials, no frameworks. The objective is not the product itself but what you prove in building it: you can design a schema, wire up an HTTP server, talk to a database, and handle errors — all with just documentation and your own reasoning.

---

## Success Criteria

You're done when:

- [ ] All three core endpoints work correctly
- [ ] The database schema is written by hand and intentionally designed
- [ ] All errors are handled — no unhandled panics, no silent failures
- [ ] The server can be started fresh, used via curl, and behaves predictably
- [ ] You can explain every line of code you wrote and why it's there

---

## Constraints

These are non-negotiable for the learning to work:

- **No AI code generation.** No Copilot, no Claude, no Cursor for writing code. Use AI only after the project is complete to review and critique your work.
- **No web frameworks.** Only Go's standard `net/http` package for routing and handling. No Gin, Echo, Chi, Fiber, or similar.
- **No ORMs.** Only `database/sql` with raw SQL queries. No GORM, sqlx (for now), or query builders.
- **No tutorial following.** If you're stuck, reach for the Go docs or PostgreSQL docs — not a walkthrough.
- **Allowed third-party packages:** One postgres driver only — `github.com/lib/pq`

---

## Tech Stack


| Layer     | Choice                        |
| --------- | ----------------------------- |
| Language  | Go                            |
| HTTP      | `net/http` (standard library) |
| Database  | PostgreSQL (local)            |
| DB Driver | `github.com/lib/pq`           |
| Testing   | `curl` from terminal          |


---

## File Structure

```
url-shortener/
├── go.mod
├── go.sum
├── main.go          ← server setup, route registration, DB connection
├── handlers.go      ← HTTP handler functions
├── db.go            ← all database query functions
└── schema.sql       ← table definition (write this first)
```

---

## Database Schema

**Write this before any Go code.** One table is all you need.

Ask yourself before writing it:

- What data do I need to store to make all three endpoints work?
- Which column will I query most often?
- Should that column have an index?
- What should my primary key be — the short code, or a separate ID?
- What are the right column types for each field?

Your `schema.sql` file should be runnable — someone should be able to run it against a fresh PostgreSQL database and have everything set up.

---

## API Specification

### POST /shorten

Creates a new short link.

**Request body:**

```json
{
  "url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ"
}
```

**Success response — 201 Created:**

```json
{
  "short_code": "abc123",
  "short_url": "http://localhost:8080/abc123",
  "original_url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
  "created_at": "2026-06-10T14:32:00Z"
}
```

**Error responses:**

- `400 Bad Request` — missing or empty URL field
- `400 Bad Request` — URL is not a valid URL (no scheme, no host)
- `500 Internal Server Error` — database failure

**Behavior:**

- Validate the incoming URL before doing anything else
- Generate a random 6-character alphanumeric short code
- If the generated code already exists in the database, generate a new one and retry
- Insert the record and return the response

---

### GET /:code

Redirects to the original URL.

**Example request:** `GET /abc123`

**Success response — 302 Found:**

- No body
- `Location` header set to the original URL
- Browser follows the redirect automatically

**Error responses:**

- `404 Not Found` — short code does not exist in the database

```json
{ "error": "short code not found" }
```

**Behavior:**

- Look up the short code in the database
- If found, increment the click count by 1, then redirect
- The increment and redirect should both happen — if the increment fails, log the error but still redirect (don't punish the user for a non-critical failure)

---

### GET /stats/:code

Returns metadata about a short link.

**Example request:** `GET /stats/abc123`

**Success response — 200 OK:**

```json
{
  "short_code": "abc123",
  "original_url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
  "short_url": "http://localhost:8080/abc123",
  "clicks": 47,
  "created_at": "2026-06-10T14:32:00Z"
}
```

**Error responses:**

- `404 Not Found` — short code does not exist

```json
{ "error": "short code not found" }
```

---

## Routing Gotcha

Go's default `http.HandleFunc` router matches by prefix, not exactly. This means if you register `/stats/` it will also catch `/shorten` in some cases, and `/:code` is a catch-all that needs to be handled carefully.

The pattern that works cleanly with standard `net/http`:

- Register `/shorten` explicitly
- Register `/stats/` explicitly  
- Register `/` as the catch-all — inside that handler, strip the leading `/` to get the code and handle it

This is a real problem you'll need to think through. Don't look up a solution first — try to reason through it.

---

## Validation Rules

**URL validation (POST /shorten):**

- The `url` field must be present and non-empty
- Must have a valid scheme (`http` or `https`)
- Must have a non-empty host
- Use Go's `net/url` package (`url.Parse`) to check — it's in the standard library

**Short code generation:**

- 6 characters
- Alphanumeric only: `a-z`, `A-Z`, `0-9`
- Randomly generated — use `math/rand` with a proper seed, or `crypto/rand` if you want to go further

---

## Error Handling Rules

- Every error must be handled — no `_` discarding errors unless you have a deliberate reason and comment explaining why
- HTTP handlers should never panic — all errors get converted to appropriate HTTP responses
- Database errors that aren't "not found" should return `500` with a generic message — don't leak internal error details to the caller
- Log all errors server-side with `log.Println` or `log.Printf` so you can see what's happening

---

## How to Test Each Endpoint

**Create a short link:**

```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.google.com"}'
```

**Follow a redirect:**

```bash
curl -L http://localhost:8080/abc123
# -L flag follows the redirect
```

**Check just the redirect response without following it:**

```bash
curl -v http://localhost:8080/abc123
# Look for "Location:" header in the response
```

**Get stats:**

```bash
curl http://localhost:8080/stats/abc123
```

**Test a missing code:**

```bash
curl http://localhost:8080/stats/doesnotexist
```

**Test bad input:**

```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": ""}'

curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "not-a-url"}'
```

---

## Stretch Goals

Only attempt these after all three core endpoints work correctly and you're satisfied with the code quality.

**In order of priority:**

### 1. Link Expiration

Add an optional `expires_at` field to `POST /shorten`:

```json
{
  "url": "https://www.google.com",
  "expires_at": "2026-12-31T00:00:00Z"
}
```

- Store it in the database (nullable timestamp column)
- On redirect, check if the link is expired — if so, return `410 Gone` instead of redirecting
- Expired links should still show stats, but note they're expired

**What this teaches:** nullable fields in SQL and Go, time comparison logic, new HTTP status codes

---

### 2. Custom Short Codes

Add an optional `custom_code` field to `POST /shorten`:

```json
{
  "url": "https://www.google.com",
  "custom_code": "myvideo"
}
```

- If provided, use it instead of generating a random code
- If that code is already taken, return `409 Conflict`
- Apply the same alphanumeric validation to custom codes

**What this teaches:** conditional logic branching, conflict handling, input validation

---

### 3. Paginated Link List

New endpoint: `GET /links?page=1&limit=10`

Returns all links, newest first, with pagination.

```json
{
  "links": [...],
  "page": 1,
  "limit": 10,
  "total": 42
}
```

**What this teaches:** SQL `LIMIT`, `OFFSET`, `COUNT(*)`, query parameter parsing in Go

---

## Order of Operations

This is the recommended sequence. Resist the urge to jump around.

1. **Write `schema.sql`** — table definition, think through indexes
2. **Set up PostgreSQL locally** and run your schema against it
3. **Write `db.go`** — functions for insert, lookup by code, increment clicks
4. **Write `main.go`** — DB connection, route registration, server start
5. **Write `handlers.go`** — one handler at a time, starting with `POST /shorten`
6. **Test with curl** after each handler is complete before moving to the next
7. **Handle edge cases and errors** — go back through each handler and harden it
8. **Stretch goals** — only after core is solid

---

## Definition of Done

Before calling this project complete, you should be able to answer these questions out loud without looking at the code:

- Why did you choose the primary key you chose?
- Why does the redirect endpoint return a 302 and not a 200?
- What happens in your code if two requests come in at the same time with the same randomly generated short code?
- What does `defer db.Close()` do and why is it there?
- What's the difference between `log.Fatal` and `log.Println` and when did you use each?
- What would break first if this API got 10,000 requests per second?

