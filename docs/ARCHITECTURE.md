# Architecture

## Goal

Keep the architecture as small as possible while supporting the complete MVP flow.

No production-scale infrastructure is required.

---

## System Overview

```text
Web Client
    |
 HTTP / JSON
    |
    v
Application Backend
    |
    +---- Database
    |
    +---- LLM Provider
```

The backend remains a single application.

---

## Components

### Web Client

Responsible for:

* business task creation;
* clarification UI;
* task-card editing;
* score display;
* catalog;
* task details;
* proposal submission;
* business proposal review.

Business/Team roles may be simulated.

### Application Backend

Responsible for:

* task lifecycle;
* AI orchestration;
* validating AI output;
* readiness scoring;
* publication;
* catalog queries;
* proposals;
* proposal decisions;
* persistence.

### Database

Stores:

* tasks;
* teams;
* proposals.

Clarification data may also be stored if convenient.

### LLM Provider

Used only for clarification and structured task generation.

Treat LLM output as untrusted input.

---

## Core Entities

### Task

```text
Task
- id
- title
- topic
- context
- need
- users
- data
- constraints
- expected_result
- success_criteria
- contact
- interaction_format
- status
- readiness_score
- readiness_level
```

Minimal statuses:

```text
draft
confirmed
published
```

---

### Team

```text
Team
- id
- name
- interests
- skills
- technologies
```

For MVP, teams may be seeded demo records.

---

### Proposal

```text
Proposal
- id
- task_id
- team_id
- solution_idea
- plan
- timeline
- prototype_url
- status
```

Statuses:

```text
pending
accepted
rejected
```

Proposal statuses are independent so multiple teams may be accepted.

---

## Task Flow

```text
Description
    ↓
AI clarification
    ↓
Questions
    ↓
Answers
    ↓
Structured task card
    ↓
Validation
    ↓
Human editing
    ↓
Confirmation
    ↓
Readiness scoring
    ↓
Publication
```

Human confirmation is required before publication.

---

## AI Boundary

Suggested clarification input:

```json
{
  "description": "..."
}
```

Suggested output:

```json
{
  "questions": [
    {
      "field": "expected_result",
      "question": "..."
    }
  ]
}
```

Task generation should return the fields defined in `SPEC.md`.

Missing values should be `null` or another explicit missing value.

Backend must validate:

* valid JSON;
* expected fields;
* expected types.

A simple fallback implementation may replace the real LLM while preserving the same contract.

---

## Readiness Scoring

Scoring is backend application logic.

```text
Task Card
    ↓
Score Calculator
    ↓
score + level + breakdown + missing fields
```

Weights and product rules are defined in `SPEC.md`.

Do not duplicate them here.

The LLM must not calculate the final score.

---

## Catalog

Catalog queries operate on published tasks.

Required operations:

```text
list published tasks
filter by topic
filter by readiness level
sort by readiness score
get task details
```

---

## Proposal Flow

```text
Published Task
    ↓
Team submits proposal
    ↓
pending
    ↓
Business decision
   / \
accept reject
```

No automatic assignment exists.

---

## API

Exact routes may change during implementation.

Suggested minimal surface:

```text
POST /api/tasks/clarify
POST /api/tasks/generate

POST /api/tasks
GET  /api/tasks
GET  /api/tasks/:id
PUT  /api/tasks/:id

POST /api/tasks/:id/confirm
POST /api/tasks/:id/publish

POST /api/tasks/:id/proposals
GET  /api/tasks/:id/proposals

PATCH /api/proposals/:id
```

Avoid adding endpoints unless needed by the UI.

---

## Validation

Backend must reject or handle:

* empty task descriptions;
* malformed AI responses;
* invalid task data;
* nonexistent tasks;
* invalid proposal statuses.

Readiness score must always remain in:

```text
0..100
```

---

## Seed Data

If no organizer dataset exists, seed:

```text
5 task drafts
5 task cards
5 teams
5 proposals
```

Use different readiness levels to make catalog behavior visible during the demo.

---

## Constraints

Do not introduce:

* microservices;
* Redis;
* message queues;
* event buses;
* Kubernetes;
* vector databases;
* RAG;
* real-time infrastructure;
* file-storage infrastructure.

See `SPEC.md` for product scope.

---

## Stack

TBD.
