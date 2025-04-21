-- 0001_create_urls_table.up.sql
CREATE TABLE IF NOT EXISTS urls (
  code          VARCHAR(32)      PRIMARY KEY,
  canonical_url VARCHAR(2048)    NOT NULL UNIQUE,
  created_at    TIMESTAMP        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  expires_at    TIMESTAMP        NULL,
  -- Automatically purge expired rows:
  TTL = 'expires_at + INTERVAL 0 SECOND'
);
