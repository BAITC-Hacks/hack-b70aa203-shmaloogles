# AGENTS.md

## Project

Hackathon MVP connecting businesses with student teams.

Core flow:

Business problem
→ AI clarification
→ Task card
→ Readiness score
→ Publication
→ Team proposal
→ Business accepts/rejects proposal

See `docs/MVP.md` for product requirements.

## Priorities

1. Working end-to-end flow is the highest priority.
2. Keep implementation minimal.
3. Do not add features outside `docs/MVP.md`.
4. Prefer simple solutions over abstractions.
5. Do not prematurely optimize.
6. Do not introduce infrastructure unless necessary.

## Stack

Frontend:
- Next.js
- TypeScript

Backend:
- Go
- REST API

Database:
- PostgreSQL

## Backend

Keep the Go backend simple.

Prefer:
- standard library where practical;
- explicit code;
- small packages;
- simple SQL/database access.

Avoid unnecessary:
- repository abstractions;
- dependency injection frameworks;
- event buses;
- microservices;
- generic abstractions.

## Frontend

Keep UI simple and demo-oriented.

Prioritize:
- task creation;
- AI clarification;
- task editing;
- readiness score;
- catalog;
- proposal submission;
- proposal acceptance/rejection.

Do not spend significant time on visual polish until the complete flow works.

## AI

The LLM is used to:
1. analyze the initial business description;
2. identify missing information;
3. generate at least 3 clarification questions;
4. convert the answers into a structured task card.

Prefer structured JSON output.

Never invent information that the business did not provide.

## Scope

Do NOT implement unless explicitly requested:
- full authentication;
- OAuth;
- password recovery;
- real-time chat;
- notifications;
- file storage;
- RAG/vector databases;
- recommendation engines;
- microservices;
- Redis;
- message queues;
- Kubernetes.

For the MVP, Business/Team roles may be simulated without real authentication.

## Development

Before implementing a large change:
1. inspect existing code;
2. choose the smallest implementation;
3. preserve the existing architecture;
4. run relevant tests/build/lint after changes.

Do not refactor unrelated code.