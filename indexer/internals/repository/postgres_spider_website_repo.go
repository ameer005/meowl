package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ameer005/meowl/internals/models"
	"github.com/ameer005/meowl/internals/storage"
)

type PostgresSpiderWebsiteRepo struct {
	db *sql.DB
}

func NewPostgresSpiderWebsiteRepo(db *sql.DB) *PostgresSpiderWebsiteRepo {
	return &PostgresSpiderWebsiteRepo{db: db}
}

func (p *PostgresSpiderWebsiteRepo) GetWebsitesBatch(id, batchSize int) ([]models.SpiderWebsite, error) {
	query, err := storage.LoadSQLQuery("fetch_rows_paginated.sql")
	var websites []models.SpiderWebsite

	if err != nil {
		return websites, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	row, err := p.db.QueryContext(ctx, query, id, batchSize)
	if err != nil {
		return websites, err
	}

	for row.Next() {
		var site models.SpiderWebsite
		err = row.Scan(
			&site.Id,
			&site.Url,
			&site.Content,
			&site.Title,
			&site.Headings,
			&site.InternalLinks,
			&site.ExternalLinks,
			&site.Images,
			&site.CrawledAt,
			&site.Description,
		)

		if err != nil {
			return websites, fmt.Errorf("Failed to scan row %w", err)
		}

		websites = append(websites, site)
	}

	return websites, nil
}
