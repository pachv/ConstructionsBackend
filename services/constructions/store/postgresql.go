package store

import (
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/pachv/constructions/constructions/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgreSQLStore struct {
	db     *sqlx.DB
	logger *slog.Logger
	dsn    string
}

func NewPostgreSQLStore(c *config.Config, logger *slog.Logger) *PostgreSQLStore {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.Database.User, c.Database.Password,
		c.Database.Host, c.Database.Port, c.Database.Name)

	return &PostgreSQLStore{
		dsn:    dsn,
		logger: logger,
	}

}

func (s *PostgreSQLStore) MustConnect() {
	const maxAttempts = 5
	var db *sqlx.DB
	var err error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		time.Sleep(time.Duration(attempt*2) * time.Second)

		db, err = sqlx.Connect("postgres", s.dsn)
		if err != nil {
			s.logger.Error(fmt.Sprintf("Attempt %d: sqlx.Connect error: %v", attempt, err))
			continue
		}

		if err = db.Ping(); err != nil {
			s.logger.Error(fmt.Sprintf("Attempt %d: Ping failed: %v", attempt, err))
			continue
		}

		s.logger.Info("Connected to postgresql")
		s.db = db
		return
	}

	panic(fmt.Errorf("all %d attempts to connect to the database failed", maxAttempts))
}

func (s *PostgreSQLStore) GetDB() *sqlx.DB {
	return s.db
}

func (s *PostgreSQLStore) MakeMigrations() {

	driver, err := postgres.WithInstance(s.db.DB, &postgres.Config{})
	if err != nil {
		log.Fatalf("migration fail: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver,
	)
	if err != nil {
		s.logger.Error(fmt.Sprintf("migration faild: %v", err))
	}

	if err := m.Up(); err != nil && err.Error() != "no change" {
		s.logger.Error(fmt.Sprintf("migration faild: %v", err))
	}

	s.logger.Info("Migrations complete")
}
