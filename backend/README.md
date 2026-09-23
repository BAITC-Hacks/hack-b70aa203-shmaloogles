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

## Test

```sh
go test ./...
```
