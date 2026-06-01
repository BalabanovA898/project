# Account Service

Service template for user management and authentication/authorization.

## Purpose

This service is intended to combine user account management and auth responsibilities:
- user lifecycle operations (register, login, profile read/update)
- authentication flows (login, refresh token)
- authorization-ready JWT issuance

## Structure

- `go.mod` - module declaration and dependencies
- `src/config` - configuration loader
- `src/domain` - domain models and shared errors
- `src/repository/postgres` - Postgres repository implementations
- `src/usecase` - business logic and application flow
- `src/delivery` - HTTP delivery layer
- `src/main.go` - entrypoint wiring config, database and server
- `schema.sql` - database schema for users and refresh tokens
- `Dockerfile`, `docker-compose.yml` - container setup for service and Postgres

## Run with Docker

```bash
cd server/account
docker compose up --build
```

Service will be available at `http://localhost:8080`.
