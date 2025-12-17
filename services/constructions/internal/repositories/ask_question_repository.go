package repositories

import (
	"errors"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/pachv/constructions/constructions/internal/domain/entity"
)

var ErrAskQuestionNotFound = errors.New("ask question not found")

type AskQuestionRepository struct {
	logger *slog.Logger
	db     *sqlx.DB
}

func NewAskQuestionRepository(logger *slog.Logger, db *sqlx.DB) *AskQuestionRepository {
	return &AskQuestionRepository{
		logger: logger.With("component", "AskQuestionRepository"),
		db:     db,
	}
}

func (r *AskQuestionRepository) Save(q entity.AskQuestion) error {
	const query = `
		INSERT INTO ask_questions (id, message, name, phone, email, product, consent,created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7,NOW());
	`

	_, err := r.db.Exec(query, q.Id, q.Message, q.Name, q.Phone, q.Email, q.Product, q.Consent)
	return err
}
