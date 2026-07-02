# TODO

## Bugs / correctness

- [x] Validate input DTOs in `handler.go` (`create`, `update`) — done via `go-playground/validator` struct tags + `validate.Struct(dto)`.
- [x] Validate `{id}` path param format (UUID) before hitting the DB — `validateIDMiddleware`, returns `400` on malformed id.
- [x] Stop leaking raw `err.Error()` to clients in `handler.go` — unexpected errors are logged server-side (`log.Printf`), clients get a generic `"internal server error"`.
- [ ] Add `ORDER BY` to `GetAll` in `postgres.go` — folded into the pagination/sorting/search work below, not done as a standalone fix.

## In progress — pagination, sorting, search on `GetAll`

Design agreed:
- Query params on `GET /blogs`: `pageNumber` (default 1), `pageSize` (default 10), `sortBy` (whitelist: `name`, `createdAt`; default `createdAt`), `sortDirection` (`asc`/`desc`; default `desc`), `searchNameTerm` (optional substring filter via `ILIKE`).
- `sortBy` must stay whitelisted in Go code, never interpolated from the request directly into SQL — column names can't be parameterized with `$n` (only values can), so an unvalidated `sortBy` would be a SQL-injection vector.
- Pagination style: `LIMIT`/`OFFSET` via `pageNumber`+`pageSize` (not keyset — table is small, simplicity preferred, matches nest-samurai's approach).
- Response envelope: `{items, totalCount, page, pageSize, pagesCount}` (mirrors nest-samurai's `PaginatedViewDto`).
- Next concrete step: write `getBlogsQueryParams` struct + `parseGetBlogsQueryParams(r)` + `queryParamInt(r, key, default)` helper in `handler.go`. Then thread the params through `Repository.GetAll` (interface signature changes in `blog.go`) → `postgres.go` (SQL with `WHERE name ILIKE $1`, whitelisted `ORDER BY`, `LIMIT`/`OFFSET`, plus a count query for `totalCount`) → `service.go` → `handler.go` response envelope.

## Reliability / production readiness

- [ ] Give DB calls a bounded `context` timeout instead of relying only on `r.Context()` (which only cancels on client disconnect, never on a time budget).
- [ ] Replace `http.ListenAndServe` with an `http.Server` that sets `ReadTimeout`/`WriteTimeout`/`IdleTimeout`.
- [ ] Add graceful shutdown (listen for `SIGINT`/`SIGTERM`, call `srv.Shutdown(ctx)`, close the DB pool cleanly).
- [ ] Add `chi` middleware: `middleware.Logger` (request logging) and `middleware.Recoverer` (catch panics instead of dropping the connection raw).
- [ ] Add structured logging for errors (currently nothing is logged server-side on failure).

## Config / secrets

- [ ] Move the Postgres connection string out of `main.go` into environment variables / `.env` (e.g. via `github.com/joho/godotenv`), matching the `.env.*` pattern already used in nest-samurai.
- [ ] Set explicit `pgxpool.Config` connection limits (`MaxConns` etc.) instead of relying on driver defaults.

## Scaling

- [ ] Add pagination to `GetAll` (`LIMIT`/`OFFSET` or keyset pagination on `created_at`/`id`) — currently returns the entire table.
- [ ] Add indexes once queries beyond `GetByID`/`GetAll` appear (e.g. filtering/sorting by `name` or `created_at`).

## Testing

- [ ] Write unit tests for `service.go` using an in-memory fake `Repository` (interface already supports this — currently unused).
- [ ] Write integration tests for `postgres.go` against the real `go-samurai-postgres` container.
- [ ] Write HTTP-level tests for `handler.go` (`httptest.NewRecorder`).

## Next features

- [ ] Port `posts`, `comments`, `likes` from nest-samurai following the same `internal/<feature>` pattern as `blogs`.
- [ ] Port `user-accounts`: JWT auth (access + refresh tokens), email confirmation, password recovery, session/device management.
