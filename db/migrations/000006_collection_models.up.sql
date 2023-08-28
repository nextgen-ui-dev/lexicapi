CREATE TABLE IF NOT EXISTS collections (
    id BYTEA NOT NULL,
    creator_id BYTEA NOT NULL,
    name VARCHAR(100) NOT NULL,
    visibility VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS collection_articles (
    id BYTEA NOT NULL,
    collection_id BYTEA NOT NULL,
    article_id BYTEA NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    deleted_at TIMESTAMPTZ,

    CONSTRAINT collection_articles_unique_ids UNIQUE NULLS NOT DISTINCT (collection_id, article_id, deleted_at),
    PRIMARY KEY(id)
);