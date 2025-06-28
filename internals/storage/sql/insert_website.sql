INSERT INTO websites (
  url,
  content,
  title,
  headings,
  internal_links,
  external_links,
  images
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
);
