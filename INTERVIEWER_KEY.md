# Interviewer Key

## main.go

| # | Problem | Line(s) |
|---|---------|---------|
| 1 | **Shouldn't exist** — off-the-shelf solutions (Auth0, Keycloak, Cognito) handle user auth; rolling your own is a security anti-pattern | — |
| 2 | **Inconsistent logging** — `fmt.Println` for startup, `log.*` everywhere else | 42, 95 |
| 3 | **Non-cryptographic RNG** — `math/rand` used for salt; `rand.Seed(time.Now().UnixNano())` is also deprecated since Go 1.20 | 11, 47 |
| 4 | **Weak password hashing** — HMAC-SHA512 is not a password hashing function; use bcrypt or argon2 | 56–63 |
| 5 | **Connection not reused** — `sql.Open` inside the handler on every request shadows the package-level `db` initialised in `main`; defeats the connection pool | 30, 66, 72 |
| 6 | **Wrong SQL placeholders** — `?` is MySQL syntax; `lib/pq` (Postgres) requires `$1, $2, ...`; runtime error on every insert | 87 |
| 7 | **Collision-prone ID** — `rand.Int63()` seeded with wall-clock time; concurrent requests can produce identical IDs | 83 |
| 8 | **Divergent timestamps** — two `time.Now()` calls for `created_at` / `updated_at`; they will differ | 88 |
| 9 | **No input validation** — empty or malformed `email`/`password` accepted silently | 74–76 |
| 10 | **No body size limit** — missing `http.MaxBytesReader`; trivial DoS vector | 74 |
| 11 | **No HTTP method guard** — handler runs for GET, DELETE, etc. | 65 |
| 12 | **PII in logs** — `%+v` on `User` writes plaintext password to stdout | 95 |
| 13 | **PII in message queue** — `json.Marshal(user)` publishes plaintext password to RabbitMQ | 127 |
| 14 | **AMQP connection not reused** — `amqp.Dial` called on every publish; should be initialised once | 115 |
| 15 | **Marshal error ignored** — `body, _ :=` silently drops a JSON error | 127 |
| 16 | **`WriteHeader` after `Write`** — response headers already sent by `json.NewEncoder`; `w.WriteHeader` is a no-op; response always returns 200 | 110–111 |
| 17 | **No `Content-Type` header** — response never sets `application/json` | 110 |
| 18 | **No context propagation** — `db.Exec` instead of `db.ExecContext(r.Context(), ...)` | 85 |
| 19 | **Divergent log timestamp** — `time.Now()` in the success log will differ from the `created_at`/`updated_at` timestamps inserted into the DB; all three should use the same pre-captured value | 88, 95 |
| 20 | **Wrong success status** — changed `http.StatusCreated` (201) to `http.StatusOK` (200) for a resource-creation endpoint; 201 is the correct code | 111 |
| 21 | **Hardcoded credentials** — DB and AMQP connection strings with plaintext username/password baked into source code; should use env vars | 35, 66, 115 |
| 22 | **Race condition on ID** — package-level `counter++` with no mutex; concurrent requests will corrupt the counter and produce duplicate IDs | 31, 82–83 |
| 23 | **Credentials read from query params** — `email` and `password` taken from `r.URL.Query()`; query params appear in server logs, proxies, and browser history in plaintext | 74–75 |
| 24 | **Nil pointer dereference** — `var rows *sql.Rows; defer rows.Close()` defers a method call on a nil pointer; panics at runtime | 78–79 |
| 25 | **Unnecessary `string([]byte(...))` roundtrips** — `[]byte(string([]byte(x)))` is a no-op that wastes allocations | 57–58 |
| 26 | **DB error swallowed** — `db.Exec(...)` return value ignored entirely; insert failures are silently discarded | 85 |
| 27 | **Raw password stored in DB** — `password` column inserted alongside the hash; plaintext credential persisted permanently | 86–88 |
| 28 | **Error response returns 200** — on DB connection failure, `fmt.Fprintf` writes a JSON error body with no `WriteHeader`; client receives HTTP 200 with an error payload | 68–69 |
| 29 | **Sensitive data in response** — response body includes `password`, `password_hash`, and `salt` in plaintext | 97–110 |
| 30 | **Unused import** — `"os"` is imported but never referenced; code will not compile | 13 |

## tasks.go

| # | Problem | Line(s) |
|---|---------|---------|
| 19 | **`Reverse` not implemented** — returns `nil` | 5–7 |
| 20 | **`Intersect` is O(n²)** — nested loops; use a map for O(n) | 10–18 |

## Dockerfile

| # | Problem | Line(s) |
|---|---------|---------|
| 21 | **Floating base tag** — `golang:latest` is unpinned | 1 |
| 22 | **No multi-stage build** — ships the full Go toolchain (~1 GB) in the final image | — |
| 23 | **Runs as root** — no `USER` directive | — |
