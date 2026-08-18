# Interviewer Key

## main.go

| # | Problem | Line(s) |
|---|---------|---------|
| 1 | **Shouldn't exist** — off-the-shelf solutions (Auth0, Keycloak, Cognito) handle user auth; rolling your own is a security anti-pattern | — |
| 2 | **Inconsistent logging** — `fmt.Println` for startup, `log.*` everywhere else | 41, 97 |
| 3 | **Non-cryptographic RNG** — `math/rand` used for salt; `rand.Seed(time.Now().UnixNano())` is also deprecated since Go 1.20 | 11, 46 |
| 4 | **Weak password hashing** — HMAC-SHA512 is not a password hashing function; use bcrypt or argon2 | 56–61 |
| 5 | **Connection not reused** — `sql.Open` inside the handler on every request shadows the package-level `db` initialised in `main`; defeats the connection pool | 30, 65, 71 |
| 6 | **Wrong SQL placeholders** — `?` is MySQL syntax; `lib/pq` (Postgres) requires `$1, $2, ...`; runtime error on every insert | 84 |
| 7 | **Collision-prone ID** — `rand.Int63()` seeded with wall-clock time; concurrent requests can produce identical IDs | 80 |
| 8 | **Divergent timestamps** — two `time.Now()` calls for `created_at` / `updated_at`; they will differ | 85 |
| 9 | **No input validation** — empty or malformed `email`/`password` accepted silently | 73–77 |
| 10 | **No body size limit** — missing `http.MaxBytesReader`; trivial DoS vector | 74 |
| 11 | **No HTTP method guard** — handler runs for GET, DELETE, etc. | 64 |
| 12 | **PII in logs** — `%+v` on `User` writes plaintext password to stdout | 97 |
| 13 | **PII in message queue** — `json.Marshal(user)` publishes plaintext password to RabbitMQ | 116 |
| 14 | **AMQP connection not reused** — `amqp.Dial` called on every publish; should be initialised once | 104 |
| 15 | **Marshal error ignored** — `body, _ :=` silently drops a JSON error | 116 |
| 16 | **`WriteHeader` after `Write`** — response headers already sent by `fmt.Fprintf`; `w.WriteHeader(201)` is a no-op; response always returns 200 | 99–100 |
| 17 | **No `Content-Type` header** — response never sets `application/json` | 99 |
| 18 | **No context propagation** — `db.Exec` instead of `db.ExecContext(r.Context(), ...)` | 82 |

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
