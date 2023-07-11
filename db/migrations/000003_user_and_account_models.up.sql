CREATE TABLE IF NOT EXISTS users (
  id BYTEA NOT NULL,
  name VARCHAR(255),
  email VARCHAR(255),
  image_url TEXT,
  status INTEGER NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  updated_at TIMESTAMPTZ,
  deleted_at TIMESTAMPTZ,

  PRIMARY KEY(id),
  CONSTRAINT users_email_unique UNIQUE NULLS NOT DISTINCT (email, deleted_at)
);

CREATE TABLE IF NOT EXISTS accounts (
  id BYTEA NOT NULL,
  user_id BYTEA NOT NULL,
  type VARCHAR(20) NOT NULL,
  password_hash VARCHAR(255),
  provider VARCHAR(30) NOT NULL,
  provider_account_id VARCHAR(255) NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  updated_at TIMESTAMPTZ,
  deleted_at TIMESTAMPTZ,

  PRIMARY KEY(id),
  CONSTRAINT accounts_provider_provider_account_id_unique UNIQUE NULLS NOT DISTINCT (provider, provider_account_id, deleted_at)
);
