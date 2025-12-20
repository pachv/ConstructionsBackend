package services

import (
	"context"
	"errors"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pachv/constructions/constructions/internal/domain/entity"
)

var ErrReviewNotFound = errors.New("review not found")

type AdminReviewService struct {
	db *sqlx.DB
}

func NewAdminReviewService(db *sqlx.DB) *AdminReviewService {
	return &AdminReviewService{db: db}
}

func (s *AdminReviewService) GetAll(ctx context.Context) ([]entity.Review, error) {
	const q = `
		SELECT
			id, name, position, text, rating, image_path, consent, can_publish, created_at
		FROM reviews
		ORDER BY created_at DESC, id ASC;
	`

	out := make([]entity.Review, 0)
	if err := s.db.SelectContext(ctx, &out, q); err != nil {
		return nil, err
	}
	return out, nil
}

type UpdateReviewRequest struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Position   string `json:"position"`
	Text       string `json:"text"`
	Rating     int    `json:"rating"`
	ImagePath  string `json:"imagePath"`
	Consent    bool   `json:"consent"`
	CanPublish bool   `json:"canPublish"`
}

// UpdateOne обновляет запись по id.
// Важно: если хочешь частичный update (patch) — скажи, сделаю через COALESCE и *string/*int.
func (s *AdminReviewService) UpdateOne(ctx context.Context, req UpdateReviewRequest) error {
	req.ID = strings.TrimSpace(req.ID)
	if req.ID == "" {
		return errors.New("id is required")
	}

	const q = `
		UPDATE reviews
		SET
			name = $2,
			position = $3,
			text = $4,
			rating = $5,
			image_path = $6,
			consent = $7,
			can_publish = $8
		WHERE id = $1;
	`

	res, err := s.db.ExecContext(ctx, q,
		req.ID,
		req.Name,
		req.Position,
		req.Text,
		req.Rating,
		req.ImagePath,
		req.Consent,
		req.CanPublish,
	)
	if err != nil {
		return err
	}

	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if aff == 0 {
		return ErrReviewNotFound
	}
	return nil
}

func (s *AdminReviewService) DeleteOne(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return errors.New("id is required")
	}

	const q = `DELETE FROM reviews WHERE id = $1;`

	res, err := s.db.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}

	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if aff == 0 {
		return ErrReviewNotFound
	}
	return nil
}
