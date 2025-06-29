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
		fmt.Println("Failed to create table", err)
		os.Exit(1)
	}
	err = createIndex(db)

	if err != nil {
		fmt.Println("Failed to create index", err)
		os.Exit(1)
	}

	return &PostgresWebsiteRepo{db}
}

func (p *PostgresWebsiteRepo) InsertWebsite(data *models.Website) error {
	query, err := LoadSQLQuery("insert_website.sql")
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_, err = p.db.ExecContext(ctx, query,
		data.Url,
		data.Content,
		data.Title,
		pq.Array(data.Headings),
		pq.Array(data.InternalLinks),
		pq.Array(data.ExternalLinks),
		pq.Array(data.Images),
		data.Description,
	)

	return err
}

func (p *PostgresWebsiteRepo) GetByurl(url string) (*models.Website, error) {
	query, err := LoadSQLQuery("get_website_by_url.sql")
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var website models.Website
	row := p.db.QueryRowContext(ctx, query, url)
	err = row.Scan(
		&website.Id,
		&website.Url,
		&website.Content,
		&website.Title,
		&website.Headings,
		&website.InternalLinks,
		&website.ExternalLinks,
		&website.Images,
		&website.CrawledAt,
		&website.Description,
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

func createIndex(db *sql.DB) error {
	query, err := LoadSQLQuery("create_index_websites_url.sql")
	if err != nil {
		fmt.Println(err)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = db.ExecContext(ctx, query)

	return err
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
