package storage

import (
	"fmt"

	cfg "onboarding/pkg/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const DSN = "postgres://%s:%s@%s:%s/%s?sslmode=%s"

func InitDB(d cfg.Database) (*gorm.DB, error) {
	source := fmt.Sprintf(
		DSN,
		d.User,
		d.Password,
		d.Host,
		d.Port,
		d.Name,
		d.SslMode,
	)

	db, err := gorm.Open(
		postgres.Open(source),
		&gorm.Config{
			SkipDefaultTransaction: true,
			PrepareStmt:            true,
			Logger:                 logger.Default.LogMode(logger.Info),
		},
	)
	if err != nil {
		return nil, err
	}

	migration, err := migrate.New(d.MigrationURL, source)
	if err != nil {
		return nil, fmt.Errorf("Couldn't create migration: %w", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("Couldn't run migration up: %w", err)
	}

	return db, nil
}
