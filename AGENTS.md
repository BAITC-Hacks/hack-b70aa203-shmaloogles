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

## Source of Truth

Product requirements are defined in `docs/MVP.md`.

If implementation details conflict with `docs/MVP.md`, follow `docs/MVP.md`.

Do not change product behavior without explicit instruction.

## Priorities

1. Working end-to-end flow is the highest priority.
2. Keep implementation minimal.
3. Do not add features outside `docs/MVP.md`.
4. Prefer simple solutions over abstractions.
5. Do not prematurely optimize.
6. Do not introduce infrastructure unless necessary.

## Stack

Frontend:

* Next.js
* TypeScript

Backend:

* Go
* REST API

Database:

* PostgreSQL

## Backend

Keep the Go backend simple.

Prefer:

* standard library where practical;
* explicit code;
* small packages;
* simple SQL/database access.

Avoid unnecessary:

* repository abstractions;
* dependency injection frameworks;
* event buses;
* microservices;
* generic abstractions.

## Frontend

Keep UI simple and demo-oriented.

Prioritize:

* task creation;
* AI clarification;
* task editing;
* readiness score;
* catalog;
* proposal submission;
* proposal acceptance/rejection.

Do not spend significant time on visual polish until the complete flow works.

## AI

The LLM is used to:

1. analyze the initial business description;
2. identify missing information;
3. generate at least 3 clarification questions;
4. convert the answers into a structured task card.

Use structured JSON output with a defined schema.

Never invent information that the business did not provide.

If information is missing, represent it explicitly as missing or `null` rather than guessing.

LLM output must be validated before storing it or returning it to the frontend.

## Scope

Do NOT implement unless explicitly requested:

* full authentication;
* OAuth;
* password recovery;
* real-time chat;
* notifications;
* file storage;
* RAG/vector databases;
* recommendation engines;
* microservices;
* Redis;
* message queues;
* Kubernetes.

For the MVP, Business/Team roles may be simulated without real authentication.

## Ambiguity

If a requirement is unclear:

1. inspect `docs/MVP.md`;
2. inspect the existing implementation;
3. choose the simplest interpretation consistent with the MVP;
4. do not invent new product requirements.

Ask for clarification only if multiple reasonable interpretations would materially change product behavior.

## Development

Before implementing any non-trivial change:

1. inspect existing code;
2. choose the smallest implementation;
3. preserve the existing architecture;
4. run relevant tests, build, and lint after changes.

Do not refactor unrelated code.

Do not introduce new dependencies unless they provide clear value and the existing stack cannot reasonably solve the problem.

## Definition of Done

A task is complete only when:

* the requested behavior is implemented;
* relevant tests pass;
* the frontend builds successfully when frontend code was changed;
* the backend builds successfully when backend code was changed;
* relevant linting passes;
* no unrelated code was changed.

Do not stop after generating code if the affected part of the project does not build.

If validation fails because of your changes, investigate and fix the failure before declaring the task complete.
