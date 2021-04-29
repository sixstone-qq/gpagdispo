package pg

import (
	"context"
	"fmt"

	migrate "github.com/golang-migrate/migrate/v4"
	migratepg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"

	"github.com/sixstone-qq/gpagdispo/recorder/pkg/domain"
)

// Store holds the DB connection.
type Store struct {
	DB *sqlx.DB
}

// NewStore connects to DB for to use store
func NewStore(dsn string) (*Store, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("can't connect to DB: %w", err)
	}

	return &Store{DB: db}, nil
}

// CreateSchema creates the Database Schema in the desired DB
func (s *Store) CreateSchema(sourcePath string) error {
	m, err := s.migrate(sourcePath)
	if err != nil {
		return err
	}

	// Migrate all the way up
	if err := m.Up(); err != nil {
		return fmt.Errorf("can't perform up migrations: %w", err)
	}

	return nil
}

// DropSchema drops the Database schema.
func (s *Store) DropSchema(sourcePath string) error {
	m, err := s.migrate(sourcePath)
	if err != nil {
		return err
	}

	// Migrate all the way down
	if err := m.Down(); err != nil {
		return fmt.Errorf("can't perform down migrations: %w", err)
	}

	return nil
}

func (s *Store) migrate(sourcePath string) (*migrate.Migrate, error) {
	instance, err := migratepg.WithInstance(s.DB.DB, new(migratepg.Config))
	if err != nil {
		return nil, fmt.Errorf("can't create migrate instance: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+sourcePath, "postgres", instance)
	if err != nil {
		return nil, fmt.Errorf("can't create datbase instance: %w", err)
	}

	return m, nil
}

// InsertWebsiteResult inserts the website and website_results in the respective tables.
func (s *Store) InsertWebsiteResult(ctx context.Context, wp domain.WebsiteParams, wr domain.WebsiteResult) error {
	tx, err := s.DB.Beginx()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	res, err := tx.NamedExecContext(ctx, `
                    INSERT INTO websites(id, url, method, match_regexp) VALUES
                    (:id, :url, :method, :match_regexp)
                    ON CONFLICT DO NOTHING;
                    `, wp)
	if err != nil {
		return fmt.Errorf("can't insert website: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("can't get number of affected rows: %w", err)
	}
	if n == 1 {
		log.Info().Msgf("Added website %s", wp.URL)
	}

	res, err = tx.NamedExecContext(ctx, `
                   INSERT INTO websites_results(website_id, elapsed_time, status, matched, unreachable, at) VALUES
                   (:id, :elapsed_time, :status, :matched, :unreachable, :at)
                   ON CONFLICT DO NOTHING`,
		map[string]interface{}{
			"id":           wp.ID,
			"elapsed_time": wr.Elapsed.Seconds(),
			"status":       wr.Status,
			"matched":      wr.Matched,
			"unreachable":  wr.Unreachable,
			"at":           wr.At,
		})
	if err != nil {
		return fmt.Errorf("can't insert website result: %w", err)
	}
	n, err = res.RowsAffected()
	if err != nil {
		return fmt.Errorf("can't get number of affected rows from website results: %w", err)
	}
	if n == 1 {
		log.Info().Msgf("Added website result from %s", wp.URL)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("can't commit tx: %w", err)
	}

	return nil
}

// Close closes the connection
func (s *Store) Close() error {
	return s.DB.Close()
}
