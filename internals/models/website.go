package models

import "time"

type Website struct {
	Url     string `bson:"url"`
	Content string `bson:"content"`

	Title     []string  `bson:"title"`
	Headings  string    `bson:"headings"`
	Outlinks  []string  `bson:"outlinks"`
	Images    []string  `bson:"images"`
	CrawledAt time.Time `bson:"crawledAt"`
}
