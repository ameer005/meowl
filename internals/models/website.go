package models

import "time"

type Website struct {
	Url     string `bson:"url"`
	Content string `bson:"content"`

	Title         []string  `bson:"title"`
	Headings      string    `bson:"headings"`
	Outlinks      []string  `bson:"outlinks"`
	InternalLinks []string  `bson:"internalLinks"`
	ExternalLinks []string  `bson:"externalLinks"`
	Images        []string  `bson:"images"`
	CrawledAt     time.Time `bson:"crawledAt"`
}
