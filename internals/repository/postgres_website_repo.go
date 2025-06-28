package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ameer005/meowl/internals/models"
	"github.com/lib/pq"
)

type PostgresWebsiteRepo struct {
	db *sql.DB
}

func NewPostgresWebsiteRepo(db *sql.DB) *PostgresWebsiteRepo {
	err := initTables(db)

	if err != nil {
		os.Exit(1)
	}

	return &PostgresWebsiteRepo{db}
}

func (p *PostgresWebsiteRepo) InsertWebsite(ctx context.Context, data *models.Website) error {
	query, err := LoadSQLQuery("insert_website.sql")
	if err != nil {
		return err
	}

	_, err = p.db.ExecContext(ctx, query,
		data.Url,
		data.Content,
		data.Title,
		data.Headings,
		pq.Array(data.InternalLinks),
		pq.Array(data.ExternalLinks),
		pq.Array(data.Images),
	)

	return err
}

func (p *PostgresWebsiteRepo) GetByurl(ctx context.Context, url string) (*models.Website, error) {
	query, err := LoadSQLQuery("get_website_by_url.sql")
	if err != nil {
		return nil, err
	}

	var website models.Website
	row := p.db.QueryRowContext(ctx, query, url)
	err = row.Scan(
		&website.Id,
		&website.Url,
		&website.Content,
		&website.Title,
		&website.Headings,
		pq.Array(&website.InternalLinks),
		pq.Array(&website.ExternalLinks),
		pq.Array(&website.Images),
		&website.CrawledAt,
	)

	if err != nil {
		return nil, err
	}

	return &website, nil
}

func initTables(db *sql.DB) error {
	query, err := LoadSQLQuery("create_table_websites.sql")
	if err != nil {
		fmt.Println(err)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := db.ExecContext(ctx, string(query))

	fmt.Println("tables created ", res)

	return nil
}

func LoadSQLQuery(filename string) (string, error) {
	basePath := "internals/storage/sql"
	fullPath := filepath.Join(basePath, filename)

	query, err := os.ReadFile(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to read SQL file %s: %w", fullPath, err)
	}

	return string(query), nil
}
