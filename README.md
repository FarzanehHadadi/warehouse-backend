# Warehouse API

Backend service for a warehouse management system. Built with Go and Gin, it exposes a REST API for managing inventory, stores, orders, departments, reports, and dashboards.

## Features

- **Inventory management** — categories, units, products, stores
- **Operations** — orders, departments, managers
- **Reporting** — threshold proximity and store product quantity reports with Excel export
- **Dashboard** — aggregated stats and recent activity feed
- **Authentication** — JWT access/refresh tokens with API key protection on all routes
- **Rate limiting** — Redis-backed per-route limits with `X-RateLimit-*` headers
- **Caching** — Redis cache for dashboard stats, activities, and read-heavy entities (units, categories)
- **Activity logging** — domain events for create/update/delete actions
- **API docs** — Swagger UI at `/swagger/index.html`
- **CORS** — configurable origins for frontend clients

## Tech Stack

| Layer        | Technology                          |
| ------------ | ----------------------------------- |
| Language     | Go 1.25                             |
| HTTP         | [Gin](https://github.com/gin-gonic/gin) |
| Database     | PostgreSQL + [GORM](https://gorm.io) |
| Cache / Rate limit | Redis                         |
| Auth         | JWT (`golang-jwt/jwt/v5`)           |
| Docs         | Swaggo                              |
| Logging      | Zap                                 |

## Prerequisites

- Go 1.25+
- PostgreSQL
- Redis
- (Optional) [Air](https://github.com/air-verse/air) for hot reload
- (Optional) [Swag](https://github.com/swaggo/swag) for regenerating API docs

## Getting Started

### 1. Clone and install dependencies

```bash
git clone <repository-url>
cd backend
go mod download
```

### 2. Configure environment

Create a `.env` file in the project root:

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=123456
DB_NAME=warehouse
DB_SSLMODE=disable

# Redis
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=

# Auth (minimum 32 characters for API and admin keys)
API_SECRET_KEY=your-api-secret-key-at-least-32-chars
ADMIN_REGISTRATION_KEY=your-admin-registration-key-32-chars
JWT_SECRET_KEY=your-jwt-secret-key
REFRESH_SECRET_KEY=your-refresh-secret-key

# CORS (comma-separated)
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
```

Generate random keys:

```bash
go run ./pkg/scripts/generate_key.go
```

### 3. Run locally

```bash
go run .
```

The server starts on **http://localhost:8080**.

### 4. Run with hot reload (Docker)

```bash
docker build -f Dockerfile.dev -t warehouse-api-dev .
docker run --env-file .env -p 8080:8080 warehouse-api-dev
```

The dev image uses Air, which rebuilds on file changes and regenerates Swagger docs before each build.

## Environment Variables

| Variable                 | Required | Default              | Description                              |
| ------------------------ | -------- | -------------------- | ---------------------------------------- |
| `DB_HOST`                | No       | `localhost`          | PostgreSQL host                          |
| `DB_PORT`                | No       | `5432`               | PostgreSQL port                          |
| `DB_USER`                | No       | `postgres`           | Database user                            |
| `DB_PASSWORD`            | No       | `123456`             | Database password                        |
| `DB_NAME`                | No       | `warehouse`          | Database name                            |
| `DB_SSLMODE`             | No       | `disable`            | PostgreSQL SSL mode                      |
| `REDIS_ADDR`             | No       | `localhost:6379`     | Redis address                            |
| `REDIS_PASSWORD`         | No       | _(empty)_            | Redis password                           |
| `API_SECRET_KEY`         | Yes      | —                    | API key for `X-API-Key` header (≥ 32 chars) |
| `ADMIN_REGISTRATION_KEY` | Yes      | —                    | Key for `X-Admin-Key` on `/auth/register` (≥ 32 chars) |
| `JWT_SECRET_KEY`         | Yes      | —                    | Secret for signing access tokens         |
| `REFRESH_SECRET_KEY`     | Yes      | —                    | Secret for signing refresh tokens        |
| `CORS_ALLOWED_ORIGINS`   | No       | localhost dev URLs   | Comma-separated allowed frontend origins |

## API Documentation

With the server running, open:

**http://localhost:8080/swagger/index.html**

Regenerate docs manually:

```bash
swag init -g ./main.go
swag fmt
```

## Authentication

All `/v1` routes require the API key header:

```
X-API-Key: <API_SECRET_KEY>
```

Protected routes additionally require a JWT:

```
Authorization: Bearer <access_token>
```

Obtain tokens via `POST /v1/auth/login`. Refresh with `POST /v1/auth/refresh`.

Admin registration (`POST /v1/auth/register`) also requires:

```
X-Admin-Key: <ADMIN_REGISTRATION_KEY>
```

## API Overview

| Group          | Base path            | Notes                                      |
| -------------- | -------------------- | ------------------------------------------ |
| Health         | `GET /v1/ping`       | Health check                               |
| Auth           | `/v1/auth`           | Login, refresh, register                   |
| Categories     | `/v1/categories`     | Public list; CRUD requires JWT             |
| Units          | `/v1/units`          | Public list; CRUD requires JWT             |
| Departments    | `/v1/departments`    | Public list; CRUD requires JWT             |
| Managers       | `/v1/managers`       | Public list; CRUD requires JWT             |
| Products       | `/v1/products`       | Public list/search; CRUD requires JWT      |
| Stores         | `/v1/stores`         | Public list; CRUD requires JWT             |
| Orders         | `/v1/orders`         | Public list; CRUD + export require JWT     |
| Reports        | `/v1/reports`        | JWT required                               |
| Dashboard      | `/v1/dashboard`      | JWT required                               |

List endpoints support cursor pagination, filtering, search, and sorting via query parameters.

## Project Structure

```
.
├── main.go                  # Application entry point
├── grace-shutdown.go        # Graceful HTTP shutdown
├── docs/                    # Generated Swagger docs
├── pkg/
│   ├── api/
│   │   ├── handlers/        # HTTP handlers
│   │   ├── middleware/      # Auth, CORS, rate limiting
│   │   ├── filter/          # Query filtering and cursor pagination
│   │   ├── dto/             # Request/response types
│   │   └── router.go        # Route definitions
│   ├── cache/               # Redis cache layer
│   ├── database/            # PostgreSQL and Redis setup
│   ├── events/              # In-process event bus
│   ├── listeners/           # Event subscribers (activity logger)
│   ├── models/              # GORM models
│   ├── repository/          # Data access layer
│   └── utils/               # Shared helpers
└── Dockerfile.dev           # Dev container with Air + Swag
```

## Development

### Hot reload with Air

```bash
air -c .air.toml
```

### Run tests

```bash
go test ./...
```

### Build

```bash
go build -o bin/warehouse .
```

## Architecture Notes

- **Repository pattern** — handlers call repositories; a generic `BaseRepository` provides CRUD, filtering, and optional Redis caching.
- **Event bus** — create/update/delete operations publish events; an activity logger persists audit records.
- **Rate limiting** — global limit on `/v1`, stricter limits on auth endpoints, higher limits on JWT-protected routes.
- **Graceful shutdown** — the server handles `SIGINT`/`SIGTERM` with a 30-second drain period.

## License

Private — all rights reserved.
