# Chirpy HTTP Server (Go)

A lightweight REST API server for a Chirpy-style social app, built with Go and PostgreSQL.

## Features

- User registration and login
- JWT-based access tokens
- Refresh token issuance, refresh, and revocation
- Create, list, fetch, and delete chirps
- Chirp content filtering (blocked words are masked)
- Webhook endpoint for Chirpy Red upgrades
- Admin metrics and reset endpoints

## Tech Stack

- Go 1.25+
- `net/http` (`http.ServeMux`)
- PostgreSQL
- SQLC-generated query layer (`internal/database`)
- JWT (`github.com/golang-jwt/jwt/v5`)
- Argon2 password hashing (`github.com/alexedwards/argon2id`)

## Project Structure

```text
.
├── main.go
├── *.go                      # HTTP handlers
├── internal/
│   ├── auth/                 # password + token helpers
│   └── database/             # sqlc generated db code
├── sql/
│   ├── schema/               # DB migrations (goose format)
│   └── queries/              # SQLC query definitions
├── assets/
├── .env.example
└── go.mod
```

## Requirements

- Go installed
- PostgreSQL running locally or remotely
- A database created (for example: `chirpy`)

## Environment Variables

Create a `.env` file in the project root:

```env
DB_URL="postgres://username:password@localhost:5432/chirpy?sslmode=disable"
JWT_SECRET="replace-with-a-long-random-secret"
POLKA_KEY="optional-webhook-api-key"
```

- `DB_URL` (required): PostgreSQL connection string.
- `JWT_SECRET` (required): Secret used to sign/validate JWT access tokens.
- `POLKA_KEY` (optional but needed for webhook auth): API key for `/api/polka/webhooks`.

## Database Setup

Apply schema files in order from `sql/schema`:

1. `001_users.sql`
2. `002_chirps.sql`
3. `003_users_hashed_password.sql`
4. `004_refresh_tokens.sql`
5. `005_users_is_chirpy_red.sql`

You can run them manually with `psql` (or with `goose` if you already use it).

## Run

```bash
go mod download
go run .
```

Server starts on:

- `http://localhost:8080`

## API Endpoints

### Health

- `GET /api/healthz`
  - Returns plain text `OK`.

### Admin

- `GET /admin/metrics`
  - Returns HTML with static file hit count.
- `POST /admin/reset`
  - Resets hit counter.

### Users & Auth

- `POST /api/users`
  - Create user.
  - Body: `{ "email": "...", "password": "..." }`

- `POST /api/login`
  - Login user.
  - Body: `{ "email": "...", "password": "..." }`
  - Returns access token + refresh token.

- `PUT /api/users`
  - Update current user's email/password.
  - Auth: `Authorization: Bearer <access_token>`

- `POST /api/refresh`
  - Exchange refresh token for new access token.
  - Auth: `Authorization: Bearer <refresh_token>`

- `POST /api/revoke`
  - Revoke refresh token.
  - Auth: `Authorization: Bearer <refresh_token>`

### Chirps

- `POST /api/chirps`
  - Create chirp (max 140 chars).
  - Auth: `Authorization: Bearer <access_token>`
  - Body: `{ "body": "your chirp" }`

- `GET /api/chirps`
  - List chirps (ascending by `created_at`).

- `GET /api/chirps/{chirpID}`
  - Fetch one chirp by ID.

- `DELETE /api/chirps/{chirpID}`
  - Delete chirp if caller owns it.
  - Auth: `Authorization: Bearer <access_token>`

### Webhooks

- `POST /api/polka/webhooks`
  - Auth header: `Authorization: ApiKey <POLKA_KEY>`
  - Expected event payload shape:

```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": "<uuid>"
  }
}
```

If event is `user.upgraded`, user is upgraded to Chirpy Red.

## Example cURL Flow

Create user:

```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"secret123"}'
```

Login:

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"secret123"}'
```

Create chirp (replace `ACCESS_TOKEN`):

```bash
curl -X POST http://localhost:8080/api/chirps \
  -H "Authorization: Bearer ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"body":"hello chirpy"}'
```

## Notes

- Bad words currently filtered in chirps: `kerfuffle`, `sharbert`, `fornax`.
- Access token lifetime is currently set to 1 hour.
- Refresh tokens are stored in DB and expire after 60 days.
