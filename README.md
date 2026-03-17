# echo-realworld-api

A RESTful API implementing the [RealWorld](https://github.com/gothinkster/realworld) spec, built with **Go**, **Echo**, **GORM**, and **PostgreSQL**, following Clean Architecture principles.

---

## Tech Stack

| Layer      | Technology                          |
|------------|-------------------------------------|
| Framework  | [Echo v4](https://echo.labstack.com) |
| ORM        | [GORM v2](https://gorm.io)          |
| Database   | PostgreSQL                          |
| Auth       | JWT (golang-jwt/jwt v5)             |
| Migrations | [golang-migrate](https://github.com/golang-migrate/migrate) |
| Docs       | Swagger UI (swaggo/echo-swagger)    |

---

## Project Structure

```
.
├── main.go                    # Entry point
├── migrations/                # Versioned SQL migration files
├── internal/
│   ├── config/                # Environment config loader
│   ├── database/              # DB connection & migration runner
│   ├── domain/                # Domain models (User, Article, Comment, Tag)
│   ├── repository/            # Data access layer (interfaces + GORM impl)
│   ├── service/               # Business logic
│   ├── handler/               # HTTP handlers
│   ├── middleware/            # JWT auth middleware
│   └── router/                # Route registration
├── pkg/
│   ├── apperrors/             # Shared error types
│   ├── utils/                 # Slug generator
│   └── validator/             # Request validator
└── docs/                      # Swagger generated files (gitignored)
```

---

## Prerequisites

- Go 1.24+
- PostgreSQL running locally (or via Docker)
- `swag` CLI for regenerating Swagger docs:
  ```bash
  go install github.com/swaggo/swag/cmd/swag@latest
  ```
- `air` for hot-reload (optional):
  ```bash
  go install github.com/air-verse/air@latest
  ```

---

## Setup

**1. Clone & install dependencies**
```bash
git clone <repo-url>
cd echo-realworld-api
go mod tidy
```

**2. Configure environment**
```bash
cp .env.example .env
# Edit .env with your database credentials
```

**3. Create the database**
```sql
CREATE DATABASE realworld_api_echo;
```

---

## Commands

### Build

```bash
# Compile binary to ./server
make build

# Or directly with Go
go build -o server .
```

### Run

```bash
# Build then run binary
make run

# Or run directly (no build step)
go run main.go

# Hot-reload with air
make dev
```

### Database Migrations

Migrations run **automatically at startup**. To run them manually, first install the CLI:

```bash
go install -tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

```bash
make migrate-up          # apply all pending migrations
make migrate-down        # roll back the last migration
make migrate-force V=3   # force-set version (use after a failed migration)
```

> DB connection is read from `.env` automatically — no need to pass a long connection string.

### Swagger Docs

```bash
# Regenerate swagger docs after changing handler annotations
make swagger

# Then start the server and open:
# http://localhost:8080/swagger/index.html
```

### Other

```bash
# Run tests
make test

# Tidy Go modules
make tidy

# Start PostgreSQL via Docker Compose
make docker-up

# Stop containers
make docker-down

# Remove build binary
make clean
```

---

## API Endpoints

| Method | Path                              | Auth | Description          |
|--------|-----------------------------------|------|----------------------|
| POST   | `/api/users`                      |      | Register             |
| POST   | `/api/users/login`                |      | Login                |
| GET    | `/api/user`                       | ✓    | Get current user     |
| GET    | `/api/profiles/:username`         |      | Get profile          |
| GET    | `/api/articles`                   |      | List articles        |
| GET    | `/api/articles/:slug`             |      | Get article          |
| POST   | `/api/articles/:slug/comments`    | ✓    | Add comment          |
| GET    | `/api/articles/:slug/comments`    |      | Get comments         |
| GET    | `/api/tags`                       |      | List tags            |

Auth header: `Authorization: Token <jwt>`

---

## Environment Variables

| Variable          | Default                    | Description            |
|-------------------|----------------------------|------------------------|
| `APP_PORT`        | `8080`                     | HTTP server port       |
| `APP_ENV`         | `development`              | Environment name       |
| `DB_HOST`         | `localhost`                | PostgreSQL host        |
| `DB_PORT`         | `5432`                     | PostgreSQL port        |
| `DB_USER`         | `postgres`                 | Database user          |
| `DB_PASSWORD`     | *(empty)*                  | Database password      |
| `DB_NAME`         | `realworld`                | Database name          |
| `DB_SSLMODE`      | `disable`                  | SSL mode               |
| `JWT_SECRET`      | `change-me-in-production`  | JWT signing secret     |
| `JWT_EXPIRE_HOURS`| `72`                       | JWT expiry in hours    |
