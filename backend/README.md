# Backend

Minimal Go REST API for the Business Task Platform MVP.

## Run

From the repository root, copy the example environment and start the complete
backend stack:

```sh
cp .env.example .env
docker compose up --build
```

This starts PostgreSQL on port `5432` and the API on port `8080`. Compose reads
the root `.env` and passes the database and AI settings to the API container.

To run the API directly instead, export the variables from the root `.env` and
run:

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
- `GET /api/tasks?scope=all` — list every task for the simulated business workspace
- `PUT /api/tasks/{id}` — replace editable fields and recalculate readiness
- `POST /api/tasks/{id}/confirm` — confirm a draft and finalize its readiness
- `POST /api/tasks/{id}/publish` — publish a confirmed task
- `GET /api/teams` — list demo teams
- `POST /api/tasks/{id}/proposals` — submit a proposal to a published task
- `GET /api/tasks/{id}/proposals` — list proposals for a task
- `GET /api/proposals` — list all proposals for the simulated business workspace
- `PATCH /api/proposals/{id}` — accept or reject a proposal
- `POST /api/tasks/clarify` — AI clarification; body `{"description":"..."}`
- `POST /api/tasks/generate` — generate a card and a preliminary readiness score;
  body `{"description":"...","answers":[{"field":"need","answer":"..."}]}`

AI defaults to an explicitly labelled deterministic mock. OpenAI configuration,
request/response contracts and scoring integration are in [AI_SCORING.md](AI_SCORING.md).
AI routes do not persist or confirm tasks. Pass the returned `card` to the task PUT.
Export AI variables into the API process environment; Go does not load `.env` itself.

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
