# go-httpserver

HTTP server with JWT authentication, chirp management, and admin features.

![Demo](demo.gif)

## Features

- JWT Authentication (Access/Refresh tokens)
- Chirp CRUD operations with profanity filtering
- User management system
- Admin metrics dashboard
- Polka webhook integration
- Makefile-driven development workflow
- PostgreSQL backend

## Install

### Prerequisites

- Go 1.21+
- PostgreSQL
- [golangci-lint](https://golangci-lint.run/)
- [goreleaser](https://goreleaser.com/)

```bash
git clone https://github.com/alnah/go-httpserver
cd go-httpserver
```

## Configure

1. Create `.env` file:

```bash
echo "DB_URL=postgres://user:pass@localhost:5432/chirpy?sslmode=disable
JWT_SECRET=your_jwt_secret
POLKA_KEY=your_polka_key
PLATFORM=dev" > .env
```

2. Initialize database:

```bash
createdb chirpy
psql chirpy -c "CREATE EXTENSION IF NOT EXISTS pgcrypto;"
```

## Build & Run

```bash
# Install dependencies and build
make install build

# Start server (listens on :8080)
make run
```

## Complete API Endpoints

### Authentication

| Method | Path         | Description          | Headers                    | Body                | Status Codes  |
| ------ | ------------ | -------------------- | -------------------------- | ------------------- | ------------- |
| POST   | /api/login   | User login           | None                       | `{email, password}` | 200, 401, 500 |
| POST   | /api/refresh | Refresh access token | `Authorization: Bearer...` | None                | 200, 401, 500 |
| POST   | /api/revoke  | Revoke refresh token | `Authorization: Bearer...` | None                | 204, 400, 500 |

### Users

| Method | Path       | Description             | Headers                    | Body                | Status Codes  |
| ------ | ---------- | ----------------------- | -------------------------- | ------------------- | ------------- |
| POST   | /api/users | Create new user         | None                       | `{email, password}` | 201, 400, 500 |
| PUT    | /api/users | Update user credentials | `Authorization: Bearer...` | `{email, password}` | 200, 401, 500 |

### Chirps

| Method | Path                  | Description        | Headers                    | Body     | Parameters                       | Status Codes       |
| ------ | --------------------- | ------------------ | -------------------------- | -------- | -------------------------------- | ------------------ |
| POST   | /api/chirps           | Create new chirp   | `Authorization: Bearer...` | `{body}` | None                             | 201, 400, 401, 500 |
| GET    | /api/chirps           | List chirps        | None                       | None     | `?author_id=UUID&sort=asc\|desc` | 200, 500           |
| GET    | /api/chirps/{chirpID} | Get specific chirp | None                       | None     | None                             | 200, 400, 404      |
| DELETE | /api/chirps/{chirpID} | Delete chirp       | `Authorization: Bearer...` | None     | None                             | 204, 400, 401, 404 |

### Admin

| Method | Path           | Description              | Headers | Body | Status Codes |
| ------ | -------------- | ------------------------ | ------- | ---- | ------------ |
| GET    | /admin/metrics | Get server metrics       | None    | None | 200          |
| POST   | /admin/reset   | Reset metrics & database | None    | None | 200, 403     |

### Webhooks

| Method | Path                | Description             | Headers                    | Body                                  | Status Codes       |
| ------ | ------------------- | ----------------------- | -------------------------- | ------------------------------------- | ------------------ |
| POST   | /api/polka/webhooks | Upgrade user membership | `Authorization: ApiKey...` | `{event:"user.upgraded", data:{...}}` | 204, 401, 404, 500 |

---

## Detailed Examples

### User Registration

```http
POST /api/users
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securePass123!"
}

Response (201 Created):
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "user@example.com",
  "is_chirpy_red": false,
  "created_at": "2024-03-20T15:04:05Z",
  "updated_at": "2024-03-20T15:04:05Z"
}
```

### Chirp Creation

```http
POST /api/chirps
Authorization: Bearer eyJhbGci...
Content-Type: application/json

{
  "body": "Hello Chirpy world!"
}

Response (201 Created):
{
  "id": "3d3a002e-4807-4a7e-bd2a-7c93239d48a7",
  "created_at": "2024-03-20T15:04:05Z",
  "updated_at": "2024-03-20T15:04:05Z",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "body": "Hello Chirpy world!"
}
```

### Error Response Format

```json
{
  "error": "Detailed error message"
}
```

## Development

### Makefile Commands

```bash
make dev    # Run full checks (fmt+lint+test+coverage+benchmark)
make test   # Run unit tests with verbose output
make lint   # Run static analysis
make build  # Build production binary to ./bin
make clean  # Remove build artifacts
```

### Testing

```bash
# Run tests with coverage
make coverage

# Performance benchmarks
make benchmark

# Generate release builds
make release  # Outputs to ./dist
```

## Database Schema

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    hashed_password TEXT NOT NULL,
    is_chirpy_red BOOLEAN DEFAULT false
);

CREATE TABLE chirps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    body TEXT NOT NULL CHECK (LENGTH(body) <= 140),
    user_id UUID REFERENCES users(id)
);
```

## Dependencies

- [godotenv](https://github.com/joho/godotenv) - Environment loading
- [pq](https://github.com/lib/pq) - PostgreSQL driver
- [jwt-go](https://github.com/golang-jwt/jwt) - JWT authentication (v5)
- [google/uuid](https://github.com/google/uuid) - UUID generation
- [x/crypto](https://pkg.go.dev/golang.org/x/crypto) - Bcrypt password hashing

## Licence

[MIT Licence](LICENSE)
