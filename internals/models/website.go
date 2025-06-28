package models

import (
	"time"

	"github.com/lib/pq"
)

type Website struct {
	Id            int            `db:"id"`
	Url           string         `db:"url"`
	Content       string         `db:"content"`
	Title         string         `db:"title"`
	Headings      string         `db:"headings"`
	InternalLinks pq.StringArray `db:"internal_links"`
	ExternalLinks pq.StringArray `db:"external_links"`
	Images        pq.StringArray `db:"images"`
	CrawledAt     time.Time      `db:"crawled_at"`
}
