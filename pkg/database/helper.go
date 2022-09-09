package database

import (
	"database/sql"
	"embed"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/logrusadapter"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/sirupsen/logrus"
	"net/url"
)


// Directory stores a directory of objects.
type Directory struct {
	logger  *logrus.Logger
	db      *sql.DB
	sb      squirrel.StatementBuilderType
	querier Querier
}

// NewDirectory creates a new Directory, connecting it to the postgres server on
// the URL provided.
func NewDirectory(logger *logrus.Logger, pgURL *url.URL) (*Directory, error) {
	c, err := pgx.ParseConfig(pgURL.String())
	if err != nil {
		return nil, fmt.Errorf("parsing postgres URI: %w", err)
	}

	c.Logger = logrusadapter.NewLogger(logger)
	db := stdlib.OpenDB(*c)

	err = validateSchema(db)
	if err != nil {
		return nil, fmt.Errorf("validating schema: %w", err)
	}

	return &Directory{
		logger:  logger,
		db:      db,
		sb:      squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).RunWith(db),
		querier: New(db),
	}, nil
}

// Close releases any resources.
func (d Directory) Close() error {
	return d.db.Close()
}


//go:embed migrations/*.sql
var fs embed.FS

// version defines the current migration version. This ensures the app
// is always compatible with the version of the database.
const version = 1

// Migrate migrates the Postgres schema to the current version.
func validateSchema(db *sql.DB) error {
	sourceInstance, err := iofs.New(fs, "migrations")
	if err != nil {
		return err
	}
	targetInstance, err := postgres.WithInstance(db, new(postgres.Config))
	if err != nil {
		return err
	}
	m, err := migrate.NewWithInstance("iofs", sourceInstance, "postgres", targetInstance)
	if err != nil {
		return err
	}
	err = m.Migrate(version) // current version
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return sourceInstance.Close()
}
