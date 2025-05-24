package main

import (
	"database/sql"
	"log"
)

func InitSchema(db *sql.DB) error {
	query := `

CREATE TABLE IF NOT EXISTS hit_counts (
    key VARCHAR(255) PRIMARY KEY,
    count INTEGER NOT NULL DEFAULT 0
);`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	log.Println("Database schema initialized (or already exists).")
	return nil
}

func IncrementHit(db *sql.DB, key string) (int, error) {
	var currentCount int

	query := `
INSERT INTO hit_counts (key, count) VALUES ($1, 1)
ON CONFLICT (key) DO UPDATE
SET count = hit_counts.count + 1
RETURNING count;`

	err := db.QueryRow(query, key).Scan(&currentCount)
	if err != nil {
		return 0, err
	}
	return currentCount, nil
}

func GetCount(db *sql.DB, key string) (int, error) {
	var count int
	query := `
SELECT count FROM hit_counts WHERE key = $1;`
	err := db.QueryRow(query, key).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return count, nil
}
