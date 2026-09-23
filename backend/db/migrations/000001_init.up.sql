BEGIN;

CREATE TABLE tasks (
    id                  BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    initial_description TEXT,
    clarification       JSONB NOT NULL DEFAULT '{"questions": [], "answers": []}'::JSONB,
    title               TEXT,
    topic               TEXT,
    context             TEXT,
    need                TEXT,
    users               TEXT,
    data                TEXT,
    constraints         TEXT,
    expected_result     TEXT,
    success_criteria    TEXT,
    contact             TEXT,
    interaction_format  TEXT,
    status              TEXT NOT NULL DEFAULT 'draft'
                        CHECK (status IN ('draft', 'confirmed', 'published')),
    readiness_score     SMALLINT NOT NULL DEFAULT 0
                        CHECK (readiness_score BETWEEN 0 AND 100),
    readiness_level     TEXT NOT NULL DEFAULT 'draft'
                        CHECK (readiness_level IN ('draft', 'workable', 'ready', 'priority')),
    readiness_breakdown JSONB NOT NULL DEFAULT '{}'::JSONB,
    missing_information TEXT[] NOT NULL DEFAULT '{}',
    suggestions         TEXT[] NOT NULL DEFAULT '{}',
    confirmed_at        TIMESTAMPTZ,
    published_at        TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CHECK (status <> 'confirmed' OR confirmed_at IS NOT NULL),
    CHECK (status <> 'published' OR (confirmed_at IS NOT NULL AND published_at IS NOT NULL))
);

CREATE TABLE teams (
    id           BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name         TEXT NOT NULL,
    interests    TEXT[] NOT NULL DEFAULT '{}',
    skills       TEXT[] NOT NULL DEFAULT '{}',
    technologies TEXT[] NOT NULL DEFAULT '{}',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE proposals (
    id            BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    task_id       BIGINT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    team_id       BIGINT NOT NULL REFERENCES teams(id) ON DELETE RESTRICT,
    solution_idea TEXT NOT NULL,
    plan          TEXT NOT NULL,
    timeline      TEXT NOT NULL,
    prototype_url TEXT,
    status        TEXT NOT NULL DEFAULT 'pending'
                  CHECK (status IN ('pending', 'accepted', 'rejected')),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX tasks_catalog_idx
    ON tasks (readiness_score DESC, published_at DESC)
    WHERE status = 'published';

CREATE INDEX tasks_catalog_filters_idx
    ON tasks (topic, readiness_level)
    WHERE status = 'published';

CREATE INDEX proposals_task_id_idx ON proposals (task_id);
CREATE INDEX proposals_team_id_idx ON proposals (team_id);

COMMIT;
