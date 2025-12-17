package repositories

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/pachv/constructions/constructions/internal/domain/entity"
)

type CallbackRepository struct {
	logger *slog.Logger
	db     *sqlx.DB
}

func NewCallbackRepository(logger *slog.Logger, db *sqlx.DB) *CallbackRepository {
	return &CallbackRepository{
		logger: logger.With("component", "CallbackRepository"),
		db:     db,
	}
}

func (r *CallbackRepository) Save(cb entity.CallbackRequest) error {
	const q = `
		INSERT INTO callback_requests (id, name, phone, consent,created_at)
		VALUES ($1, $2, $3, $4,NOW());
	`
	_, err := r.db.Exec(q, cb.Id, cb.Name, cb.Phone, cb.Consent)
	return err
}
