package repositories

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/pachv/constructions/constructions/internal/domain/entity"
)

type ReviewRepository struct {
	logger *slog.Logger
	db     *sqlx.DB
}

func NewReviewRepository(logger *slog.Logger, db *sqlx.DB) *ReviewRepository {
	return &ReviewRepository{
		logger: logger.With("component", "ReviewRepository"),
		db:     db,
	}
}

func (r *ReviewRepository) Create(rv entity.Review) error {
	const q = `
		INSERT INTO reviews
		(id, name, position, text, rating, image_path, consent, can_publish, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9);
	`

	_, err := r.db.Exec(
		q,
		rv.Id,
		rv.Name,
		rv.Position,
		rv.Text,
		rv.Rating,
		rv.ImagePath,
		rv.Consent,
		rv.CanPublish,
		rv.CreatedAt,
	)

	return err
}

func (r *ReviewRepository) GetAllPublished() ([]entity.Review, error) {
	const q = `
		SELECT id, name, position, text, rating, image_path, consent, can_publish, created_at
		FROM reviews
		WHERE can_publish = true
		ORDER BY created_at DESC;
	`

	var items []entity.Review
	if err := r.db.Select(&items, q); err != nil {
		return nil, err
	}

	return items, nil
}
