package repositories

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pachv/constructions/constructions/internal/domain/entity"
)

type CertificateRepository struct {
	db     *sqlx.DB
	logger *slog.Logger
}

func NewCertificateRepository(db *sqlx.DB, logger *slog.Logger) *CertificateRepository {
	return &CertificateRepository{db: db, logger: logger}
}

type certificateRow struct {
	ID        string    `db:"id"`
	Title     string    `db:"title"`
	FilePath  string    `db:"file_path"`
	CreatedAt time.Time `db:"created_at"`
}

func (r *CertificateRepository) GetAll(ctx context.Context) ([]entity.Certificate, error) {
	const q = `
		SELECT id, title, file_path, created_at
		FROM certificates
		ORDER BY created_at DESC
	`

	r.logger.Debug("CertificateRepository.GetAll: start")

	var rows []certificateRow
	if err := r.db.SelectContext(ctx, &rows, q); err != nil {
		r.logger.Error("CertificateRepository.GetAll: select failed", "err", err, "query", q)
		return nil, fmt.Errorf("get all certificates: %w", err)
	}

	res := make([]entity.Certificate, 0, len(rows))
	for _, row := range rows {
		t := row.CreatedAt
		res = append(res, entity.Certificate{
			ID:        row.ID,
			Title:     row.Title,
			FilePath:  row.FilePath,
			CreatedAt: &t,
		})
	}

	r.logger.Debug("CertificateRepository.GetAll: done", "count", len(res))
	return res, nil
}
