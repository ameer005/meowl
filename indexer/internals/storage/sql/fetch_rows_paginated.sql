SELECT
  id,
  url,
  content,
  title,
  headings,
  internal_links,
  external_links,
  images,
  crawled_at,
  description
FROM websites
WHERE id > $1 LIMIT $2;
