CREATE TABLE IF NOT EXISTS friends (
    id BYTEA NOT NULL,
    requester_id BYTEA NOT NULL,
    requestee_id BYTEA NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    deleted_at TIMESTAMPTZ,

    PRIMARY KEY (id)
);

CREATE UNIQUE INDEX IF NOT EXISTS friends_user_ids_unique ON friends(GREATEST(requester_id, requestee_id), LEAST(requester_id, requestee_id), deleted_at) NULLS NOT DISTINCT
WHERE deleted_at IS NULL;