# Product Specification

## Goal

Build a working MVP that connects businesses with student teams.

Core flow:

Business draft
→ AI clarification
→ Editable task card
→ Readiness score
→ Publication
→ Team proposal
→ Business accepts/rejects proposals

The end-to-end flow is more important than visual polish or extra features.

---

## Actors

### Business

Can:

* enter a task/problem description;
* answer AI clarification questions;
* edit and confirm the generated task card;
* see readiness score and missing information;
* publish the task;
* review proposals;
* accept or reject proposals manually.

### Student Team

Can:

* browse published tasks;
* filter and sort the catalog;
* open task details;
* submit proposals.

Any team may respond to any published task.

---

## Task Creation Flow

1. Business enters a free-form description.
2. System analyzes missing information.
3. AI generates at least 3 relevant clarification questions.
4. Business answers them.
5. System generates a structured task card.
6. Business edits and confirms the card.
7. Readiness score is calculated.
8. Business publishes the task.

AI-generated content must always remain editable.

---

## Task Card

Required fields:

* title
* context
* need
* users
* data
* constraints
* expected result
* success criteria
* contact
* interaction format
* topic

Implementation metadata may include:

* id
* status
* readiness score
* readiness level
* timestamps

---

## Readiness Score

Score range:

`0–100`

The score represents how ready the task is for student work.

It does not represent company popularity or team interest.

### Weights

| Category               | Points |
| ---------------------- | -----: |
| Context and need       |     20 |
| Data and materials     |     20 |
| Expected result        |     15 |
| Success criteria       |     15 |
| Constraints            |     10 |
| Users                  |     10 |
| Business communication |     10 |

Total: `100`

Scoring must be deterministic.

The LLM must not directly choose the score.

For the MVP, simple field/category completeness checks are sufficient.

---

## Readiness Levels

| Score  | Level    |
| ------ | -------- |
| 0–39   | Draft    |
| 40–69  | Workable |
| 70–89  | Ready    |
| 90–100 | Priority |

The UI must show:

* total score;
* readiness level;
* score breakdown;
* missing information;
* suggestions for increasing the score.

Score is recalculated after confirmed edits.

Low readiness must not prevent publication or proposals.

---

## Catalog

The catalog contains all published tasks.

Required functionality:

* show all published tasks;
* open task details;
* sort by readiness score;
* filter by topic;
* filter by readiness level.

Higher readiness may affect catalog position.

Low-scoring published tasks remain visible.

---

## Proposals

A proposal contains:

* team;
* solution idea;
* plan;
* expected timeline;
* prototype link.

There is no limit on the number of proposals per task.

Business can manually:

* accept a proposal;
* reject a proposal.

Business may accept:

* one team;
* multiple teams;
* no teams.

Automatic team assignment is forbidden.

---

## AI

AI is used for:

1. analyzing the initial description;
2. detecting missing information;
3. generating at least 3 clarification questions;
4. converting the description and answers into a structured task card.

Prefer structured JSON output.

AI must not invent facts not provided by the business.

Missing information should remain missing instead of being guessed.

LLM output must be validated before use.

If the external AI API is unavailable, a deterministic mock/fallback is allowed.

The project must still be able to demonstrate:

* prompt;
* input format;
* output format;
* invalid-response handling.

---

## Demo Data

If organizers do not provide data, create at least:

* 5 task drafts;
* 5 task cards;
* 5 team profiles;
* 5 proposals.

Tasks should have different completeness/readiness levels.

---

## Out of Scope

Do not implement unless explicitly requested:

* production authentication;
* OAuth;
* password recovery;
* complex roles;
* real-time chat;
* notifications;
* calendars;
* file storage;
* custom ML training;
* RAG;
* vector databases;
* recommendation engines;
* production infrastructure;
* mobile optimization;
* full project tracking;
* automatic team assignment.

Business and Team identities may be simulated.

---

## Demo Scenario

The final demo must show:

1. weak business description;
2. AI clarification questions;
3. business answers;
4. generated editable task card;
5. readiness score;
6. missing information;
7. task improvement;
8. increased score;
9. confirmation;
10. publication;
11. task visible in catalog;
12. team submits proposal;
13. business accepts or rejects proposal.

All transitions must work in the actual MVP.

---

## Acceptance Criteria

The MVP is complete when:

* free-form task input works;
* AI generates at least 3 relevant questions;
* answers become a structured editable task card;
* human confirmation is required;
* readiness score is calculated from 0 to 100;
* score breakdown and missing information are shown;
* score recalculates after edits;
* confirmed task can be published;
* published tasks appear in the catalog;
* catalog sorting and filters work;
* team can submit a proposal;
* business can manually accept/reject proposals;
* low readiness does not block proposals;
* AI does not invent business facts;
* complete demo flow works end-to-end.
