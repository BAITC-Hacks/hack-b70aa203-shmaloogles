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

Frontend integration: [API contract and complete flow](API.md).

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
- `POST /api/tasks/score` — preview readiness after editing; body is the card
  itself (for example `{"context":"Current process","need":"Desired change"}`);
  no AI call, persistence or confirmation

AI defaults to an explicitly labelled deterministic mock. OpenAI configuration,
request/response contracts and scoring integration are in [AI_SCORING.md](AI_SCORING.md).
AI routes do not persist or confirm tasks. Pass the returned `card` to the task PUT.
Export AI variables into the API process environment; Go does not load `.env` itself.

## Test

```sh
go test ./...
```

PostgreSQL integration tests are opt-in. Use a migrated disposable test database:

```sh
AI_TEST_DATABASE_URL='postgres://shmaloogles:shmaloogles@localhost:5432/shmaloogles?sslmode=disable' go test ./... -count=1
```

These tests create and remove their own fixtures. The HTTP flow test covers mock AI,
card persistence/editing, confirmation, publication, catalog, proposals and independent
accept/reject decisions, including publication and proposals at zero readiness.
It uses mock AI explicitly: no API key or paid model requests are needed.

## Database

### Recalculate existing readiness

Seed changes only affect newly created database volumes. To update existing
derived scores, export `DATABASE_URL` and run from `backend`:

```sh
go run ./cmd/recalculate-readiness
go run ./cmd/recalculate-readiness -apply
```

The first command is a dry-run (transaction rolled back); the second commits.
Both lock task rows briefly, so run during a quiet period. The command times out
after 30 seconds and uses the same deterministic scorer as the API, without AI.
It updates only score, level, breakdown, missing fields and suggestions for all
tasks. Cards, lifecycle states, timestamps and proposals are preserved. Repeating
the command is safe. It does not load `.env` automatically. Back up important data
before applying maintenance commands.

### CI

The Backend GitHub Actions workflow runs on pushes and pull requests: formatting,
tests (including PostgreSQL integration), vet and build. Its PostgreSQL service is
initialized from the migration and seed; it uses mock AI and needs no API secrets.

### Local database

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
