package persistence

import (
	"database/sql"
	"runtime"

	"github.com/Nezent/auth-service/internal/infrastructure/config"
	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

type Database struct {
	DB  *bun.DB
	cfg *config.Config
}

// NewDatabase initializes and returns a Bun DB instance for PostgreSQL.
// Pass the DSN string directly for now; later, you can use config.
func NewDatabase(cfg *config.Config) (*Database, error) {
	sqldb, err := sql.Open(cfg.Postgres.Driver, cfg.Postgres.BuildDsn())
	if err != nil {
		return nil, err
	}

	maxOpenConns := 4 * runtime.GOMAXPROCS(0)
	sqldb.SetMaxOpenConns(maxOpenConns)
	sqldb.SetMaxIdleConns(maxOpenConns)

	db := bun.NewDB(sqldb, pgdialect.New())

	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	return &Database{
		DB:  db,
		cfg: cfg,
	}, nil
}

// Close closes the database connection.
func (db *Database) Close() error {
	return db.DB.DB.Close()
}

// RawSQLDB returns the underlying *sql.DB instance.
func (db *Database) RawSQLDB() *sql.DB {
	return db.DB.DB
}
