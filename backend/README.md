# Backend

Minimal Go REST API for the Business Task Platform MVP.

## Run

```sh
go run ./cmd/api
```

The server listens on port `8080` by default. Set `PORT` to override it.

## Endpoints

- `GET /api` — API availability check
- `GET /health` — service health check
- `POST /api/tasks` — create a task draft from an initial description
- `GET /api/tasks/{id}` — get a task
- `GET /api/tasks` — list published tasks; supports `topic`, `readiness_level`, and readiness sorting
- `PUT /api/tasks/{id}` — replace editable fields and recalculate readiness
- `POST /api/tasks/{id}/confirm` — confirm a draft and finalize its readiness
- `POST /api/tasks/{id}/publish` — publish a confirmed task
- `GET /api/teams` — list demo teams

## Test

```sh
go test ./...
```

## Database

From the repository root, start PostgreSQL with:

```sh
docker compose up -d postgres
```

The default connection string is:

```text
postgres://shmaloogles:shmaloogles@localhost:5432/shmaloogles?sslmode=disable
```

Copy `.env.example` to `.env` to override the local defaults. The initial
migration and demo seed data in `backend/db` are applied automatically when the
database volume is created for the first time.

To stop the database:

```sh
docker compose down
```

To also delete local database data and reapply the schema on the next start:

```sh
docker compose down -v
```
