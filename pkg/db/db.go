package db

import (
	"database/sql"
	"fmt"
	"log"
	cfg "onboarding/pkg/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func InitDB(d cfg.Database) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", d.Name)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Migrations
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatal(err)
	}

	migration := fmt.Sprintf("file://%s", d.MigrationFile)
	m, err := migrate.NewWithDatabaseInstance(migration, "sqlite3", driver)
	if err != nil {
		return nil, fmt.Errorf("Couldn't create migration: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("Couldn't run migration up: %w", err)
	}

	return db, nil
}
