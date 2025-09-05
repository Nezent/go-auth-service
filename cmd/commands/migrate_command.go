package commands

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/Nezent/auth-service/internal/infrastructure/config"
	"github.com/Nezent/auth-service/internal/infrastructure/persistence"
	"github.com/spf13/cobra"
)

func init() {
	migrateCommand.AddCommand(migrateUpToCommand)
	migrateCommand.AddCommand(migrateDownCommand)
	migrateCommand.AddCommand(migrateDownToCommand)
	migrateCommand.AddCommand(migrateRedoCommand)
	migrateCommand.AddCommand(migrateResetCommand)
	migrateCommand.AddCommand(migrateStatusCommand)
	migrateCommand.AddCommand(migrateVersionCommand)
	migrateCommand.AddCommand(migrateCreateCommand)
	migrateCommand.AddCommand(migrateFixCommand)
	migrateCommand.AddCommand(migrateUpCommand)
}

var migrateCommand = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate the database",
	Long:  "Manage database migrations",
}

// up-to VERSION        Migrate the DB to a specific VERSION
var migrateUpToCommand = &cobra.Command{
	Use:   "up-to",
	Short: "Migrate the DB to a specific version",
	Long:  "Apply migrations until the database reaches the specified version",
	RunE: func(cmd *cobra.Command, args []string) error {
		var version int64

		if len(args) != 1 {
			return fmt.Errorf("expected 1 argument, got %d", len(args))
		}

		version, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid version: %w", err)
		}

		migrator := initMigrator()

		return migrator.UpTo(version)
	},
}

// down                 Roll back the version by 1
var migrateDownCommand = &cobra.Command{
	Use:   "down",
	Short: "Roll back the DB version by 1",
	Long:  "Undo the last applied migration",
	RunE: func(cmd *cobra.Command, args []string) error {
		migrator := initMigrator()
		return migrator.Down()
	},
}

// down-to VERSION      Roll back to a specific VERSION
var migrateDownToCommand = &cobra.Command{
	Use:   "down-to",
	Short: "Roll back to a specific version",
	Long:  "Revert database migrations until it reaches the specified version",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("expected 1 argument, got %d", len(args))
		}
		version, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid version: %w", err)
		}

		migrator := initMigrator()
		return migrator.DownTo(version)
	},
}

// redo                 Re-run the latest migration
var migrateRedoCommand = &cobra.Command{
	Use:   "redo",
	Short: "Re-run the latest migration",
	Long:  "Undo the last applied migration and reapply it",
	RunE: func(cmd *cobra.Command, args []string) error {
		migrator := initMigrator()
		return migrator.Redo()
	},
}

// reset                Roll back all migrations
var migrateResetCommand = &cobra.Command{
	Use:   "reset",
	Short: "Roll back all migrations",
	Long:  "Undo all applied migrations, resetting the database",
	RunE: func(cmd *cobra.Command, args []string) error {
		migrator := initMigrator()
		return migrator.Reset()
	},
}

// status               Dump the migration status for the current DB
var migrateStatusCommand = &cobra.Command{
	Use:   "status",
	Short: "Show migration status",
	Long:  "Display the current state of applied and pending migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		migrator := initMigrator()
		return migrator.Status()
	},
}

// version              Print the current version of the database
var migrateVersionCommand = &cobra.Command{
	Use:   "version",
	Short: "Show current database version",
	Long:  "Prints the current version of the applied database migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		migrator := initMigrator()
		return migrator.Version()
	},
}

func initMigrator() *persistence.Migrator {
	cfg := config.NewConfig()
	db, err := persistence.NewDatabase(cfg)
	if err != nil {
		log.Fatal(err)
	}
	dir := cfg.Postgres.MigrationDir
	sqlDb := db.RawSQLDB()
	migrator, err := persistence.NewMigrator(sqlDb, dir, cfg.Postgres.Driver)

	if err != nil {
		log.Fatal(err)
	}

	return migrator
}

// create NAME [sql|go] Creates new migration file with the current timestamp
var migrateCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Create a new migration file NAME [sql|go]",
	Long:  "Generate a new migration file with the specified name and format (SQL or Go)",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			name   string
			format string
		)

		if len(args) == 0 {
			return errors.New("migration name is required")
		}

		name = args[0]
		format = "sql"

		if len(args) > 1 {
			format = args[1]
		}

		if format != "sql" && format != "go" {
			return errors.New("invalid format, must be sql or go")
		}

		migrator := initMigrator()
		migrator.Create(name, format)
		return nil
	},
}

// fix                  Apply sequential ordering to migrations
var migrateFixCommand = &cobra.Command{
	Use:   "fix",
	Short: "Fix migration sequence",
	Long:  "Reorders migration files to maintain sequential execution order",
	RunE: func(cmd *cobra.Command, args []string) error {
		migrator := initMigrator()
		return migrator.Fix()
	},
}

// up                   Migrate the DB to the most recent version available
var migrateUpCommand = &cobra.Command{
	Use:   "up",
	Short: "Migrate the DB to the latest version",
	Long:  "Applies all pending migrations to bring the database up to date",
	RunE: func(cmd *cobra.Command, args []string) error {
		migrator := initMigrator()
		return migrator.Up()
	},
}
