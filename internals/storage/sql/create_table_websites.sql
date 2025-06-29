CREATE TABLE IF NOT EXISTS websites (
  id SERIAL PRIMARY KEY,
  url TEXT UNIQUE,
  content TEXT,
  title TEXT,
  headings TEXT[],
  internal_links TEXT[],
  external_links TEXT[],
  images TEXT[],
  crawled_at TIMESTAMP DEFAULT NOW(),
  description TEXT
);

