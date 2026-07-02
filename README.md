# go-samurai

Learning project: porting [nest-samurai](https://github.com/sugarflocky/nest-samurai) (a NestJS blog platform) to idiomatic Go, feature by feature.

## Stack

- [chi](https://github.com/go-chi/chi) — lightweight HTTP router
- [pgx](https://github.com/jackc/pgx) — PostgreSQL driver (raw SQL, no ORM)
- [golang-migrate](https://github.com/golang-migrate/migrate) — versioned SQL migrations
- [godotenv](https://github.com/joho/godotenv) — `.env` config loading
- PostgreSQL 17 (Docker)

## Architecture

Each feature lives in its own package under `internal/`, layered internally as:

```
handler.go   → HTTP layer (chi routes, JSON decode/encode, status codes)
service.go   → business logic, depends only on the Repository interface
postgres.go  → concrete Repository implementation (raw SQL via pgx)
blog.go      → domain model + Repository interface + sentinel errors
```

No DI container — dependencies are wired explicitly in `cmd/api/main.go`.

## Getting started

1. Copy `.env.example` to `.env` and adjust if needed:
   ```
   cp .env.example .env
   ```
2. Start PostgreSQL:
   ```
   docker compose up -d
   ```
3. Apply migrations:
   ```
   migrate -path migrations -database "$DATABASE_URL" up
   ```
4. Run the server:
   ```
   go run ./cmd/api
   ```
   Server listens on `http://localhost:8080`.

## API

| Method | Path          | Description       |
|--------|---------------|--------------------|
| POST   | /blogs        | Create a blog      |
| GET    | /blogs        | List all blogs     |
| GET    | /blogs/{id}   | Get a blog by id   |
| PUT    | /blogs/{id}   | Update a blog      |
| DELETE | /blogs/{id}   | Delete a blog      |

## Status

Only the `blogs` feature is implemented so far. See [TODO.md](./TODO.md) for
known gaps (validation, timeouts, tests, etc.) and planned next features
(`posts`, `comments`, `likes`, `user-accounts` with JWT auth).
