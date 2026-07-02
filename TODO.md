# TODO

## Bugs / correctness

- [ ] Validate input DTOs in `handler.go` (`create`, `update`) before hitting the DB — reject empty `Name`, enforce length limits matching the DB schema (`name VARCHAR(255)`), return `400` instead of letting Postgres constraint errors bubble up as `500`.
- [ ] Stop leaking raw `err.Error()` to clients in `handler.go` — log the full error server-side, return a generic message to the client.
- [ ] Add `ORDER BY` to `GetAll` in `postgres.go` (e.g. `ORDER BY created_at DESC`) — row order is currently undefined without it.

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
