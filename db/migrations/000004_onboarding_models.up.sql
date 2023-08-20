CREATE TABLE IF NOT EXISTS users_interests (
    id BYTEA NOT NULL,
    user_id BYTEA NOT NULL,
    category_id BYTEA NOT NULL,
    deleted_at TIMESTAMPTZ,
    PRIMARY KEY(id),
    CONSTRAINT users_interests_unique UNIQUE NULLS NOT DISTINCT (user_id, category_id, deleted_at)
);

ALTER TABLE IF EXISTS users
ADD COLUMN IF NOT EXISTS role VARCHAR(50),
ADD COLUMN IF NOT EXISTS education_level VARCHAR(50);