-- 0002_create_url_previews_table.up.sql
CREATE TABLE IF NOT EXISTS url_previews (
  url_code    VARCHAR(32)   PRIMARY KEY,
  title       VARCHAR(512)  NULL,
  descr       TEXT          NULL,
  fetched_at  TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_url
    FOREIGN KEY (url_code)
    REFERENCES urls(code)
    ON DELETE CASCADE
);
