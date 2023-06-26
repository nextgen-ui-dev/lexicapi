CREATE TABLE IF NOT EXISTS article_categories (
  id BYTEA NOT NULL,
  name VARCHAR(100) NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  updated_at TIMESTAMPTZ,
  deleted_at TIMESTAMPTZ,

  CONSTRAINT article_categories_name_unique UNIQUE NULLS NOT DISTINCT (name, deleted_at),
  PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS articles (
  id BYTEA NOT NULL,
  category_id BYTEA NOT NULL,
  title VARCHAR(255) NOT NULL,
  thumbnail_url TEXT,
  original_url TEXT NOT NULL,
  source VARCHAR(255) NOT NULL,
  author VARCHAR(255),
  is_published BOOLEAN DEFAULT FALSE NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  updated_at TIMESTAMPTZ,
  deleted_at TIMESTAMPTZ,

  PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS article_texts (
  id BYTEA NOT NULL,
  article_id BYTEA NOT NULL,
  content TEXT NOT NULL,
  difficulty VARCHAR(25) NOT NULL,
  is_adapted BOOLEAN NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  updated_at TIMESTAMPTZ,
  deleted_at TIMESTAMPTZ,

  PRIMARY KEY(id)
);
