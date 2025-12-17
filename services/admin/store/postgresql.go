package store

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/is_backend/services/admin/internal/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgreSQLStore struct {
	db  *sqlx.DB
	dsn string
}

func NewPostgreSQLStore(c *config.Config) *PostgreSQLStore {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.Database.User, c.Database.Password,
		c.Database.Host, c.Database.Port, c.Database.Name)

	return &PostgreSQLStore{
		dsn: dsn,
	}

}

func (s *PostgreSQLStore) Connect() {
	const maxAttempts = 5
	var db *sqlx.DB
	var err error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		time.Sleep(time.Duration(attempt*2) * time.Second)

		db, err = sqlx.Connect("postgres", s.dsn)
		if err != nil {
			log.Printf("Attempt %d: sqlx.Connect error: %v", attempt, err)
			continue
		}

		if err = db.Ping(); err != nil {
			log.Printf("Attempt %d: Ping failed: %v", attempt, err)
			continue
		}

		log.Println("Connected to postgresql")
		s.db = db
		return
	}

	log.Fatalf("All %d attempts to connect to the database failed", maxAttempts)
}

func (s *PostgreSQLStore) GetDB() *sqlx.DB {
	return s.db
}

func (s *PostgreSQLStore) MakeMigrations() {
	// Готовим миграции
	driver, err := postgres.WithInstance(s.db.DB, &postgres.Config{})
	if err != nil {
		log.Fatalf("migration fail: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", // Путь к папке с миграциями
		"postgres", driver,
	)
	if err != nil {
		log.Fatalf("migration faild: %v", err)
	}

	// Применяем миграции
	if err := m.Up(); err != nil && err.Error() != "no change" {
		log.Fatalf("migration faild: %v", err)
	}

	log.Println("migrations complete")
}
