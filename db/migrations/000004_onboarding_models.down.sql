ALTER TABLE IF EXISTS users
DROP COLUMN IF EXISTS role,
DROP COLUMN IF EXISTS education_level;

DROP TABLE IF EXISTS users_interests;