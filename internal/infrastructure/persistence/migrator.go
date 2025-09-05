// Package persistence provides database migration helpers using Goose.
package persistence

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/pressly/goose/v3"
)

// Migrator handles database schema migrations using Goose.
type Migrator struct {
	db            *sql.DB
	migrationsDir string
}

// NewMigrator initializes a new Migrator instance using an existing *sql.DB connection.
// driver should be "postgres" for PostgreSQL.
func NewMigrator(db *sql.DB, migrationsDir, driver string) (*Migrator, error) {
	if err := goose.SetDialect(driver); err != nil {
		return nil, fmt.Errorf("failed to set goose dialect: %w", err)
	}
	return &Migrator{
		db:            db,
		migrationsDir: migrationsDir,
	}, nil
}

// Close closes the database connection.
func (m *Migrator) Close() error {
	if err := m.db.Close(); err != nil {
		log.Printf("Failed to close database connection: %v", err)
		return err
	}
	return nil
}

// Up migrates the database to the latest version.
func (m *Migrator) Up() error {
	return goose.Up(m.db, m.migrationsDir)
}

// Down rolls back the database by one migration.
func (m *Migrator) Down() error {
	return goose.Down(m.db, m.migrationsDir)
}

// UpTo migrates the database up to a specific version.
func (m *Migrator) UpTo(version int64) error {
	return goose.UpTo(m.db, m.migrationsDir, version)
}

// DownTo migrates the database down to a specific version.
func (m *Migrator) DownTo(version int64) error {
	return goose.DownTo(m.db, m.migrationsDir, version)
}

// Reset rolls back all migrations and then migrates to the latest version.
func (m *Migrator) Reset() error {
	return goose.Reset(m.db, m.migrationsDir)
}

// Redo rolls back the last migration and re-applies it.
func (m *Migrator) Redo() error {
	return goose.Redo(m.db, m.migrationsDir)
}

// Status prints the status of all migrations.
func (m *Migrator) Status() error {
	return goose.Status(m.db, m.migrationsDir)
}

// Version prints the current migration version.
func (m *Migrator) Version() error {
	return goose.Version(m.db, m.migrationsDir)
}

// Create generates a new migration file.
func (m *Migrator) Create(name, migrationType string) error {
	return goose.Create(m.db, m.migrationsDir, name, migrationType)
}

// Fix renames migrations to resolve ordering issues.
func (m *Migrator) Fix() error {
	return goose.Fix(m.migrationsDir)
}
