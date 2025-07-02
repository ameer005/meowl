package storage

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func CreateSqlIndex(db *sql.DB) error {
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
