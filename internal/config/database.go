package config

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func InitializeDatabase(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	log.Println("Database connection established successfully")

	// Run migrations automatically
	if err := runMigrations(db); err != nil {
		log.Printf("Warning: Migration error: %v\n", err)
	}

	return db, nil
}

// runMigrations executes all pending migrations
func runMigrations(db *sqlx.DB) error {
	// Create tracking table
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`); err != nil {
		return err
	}

	// Get all migration files
	files, err := os.ReadDir("db/migrations")
	if err != nil {
		return err
	}

	// Get executed migrations
	executed := make(map[string]bool)
	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var version string
			rows.Scan(&version)
			executed[version] = true
		}
	}

	// Sort and execute pending migrations
	var sqlFiles []string
	for _, f := range files {
		if !f.IsDir() && filepath.Ext(f.Name()) == ".sql" {
			sqlFiles = append(sqlFiles, f.Name())
		}
	}
	sort.Strings(sqlFiles)

	for _, filename := range sqlFiles {
		version := filename[:len(filename)-4] // Remove .sql
		if executed[version] {
			continue
		}

		sql, err := os.ReadFile(filepath.Join("db/migrations", filename))
		if err != nil {
			return err
		}

		if _, err := db.Exec(string(sql)); err != nil {
			return err
		}

		if _, err := db.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", version); err != nil {
			return err
		}

		log.Printf("Migration: %s\n", filename)
	}

	return nil
}
