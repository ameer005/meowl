package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // registers the "postgres" driver
)

func NewPostgresConnection(connectionStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionStr)
	if err != nil {
		fmt.Printf("Postgres failed to connect %v\n", err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		fmt.Printf("Postgres ping failed %v\n", err)
		return nil, err
	}

	fmt.Println("✅ Connected to PostgreSQL")
	return db, nil
}
