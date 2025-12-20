package services

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pachv/constructions/constructions/internal/domain/entity"
)

type AdminReviewService struct {
	db     *sqlx.DB
	domain string
}

func NewAdminReviewService(db *sqlx.DB, domain string) *AdminReviewService {
	return &AdminReviewService{db: db, domain: domain}
}

const REVIEW_PUBLISHED_URL = "/reviews/picture/"

func (s *AdminReviewService) GetPaged(ctx context.Context, page int, search, orderBy string) ([]entity.Review, int, error) {
	const pageSize = 10
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	allowedOrderBy := map[string]bool{
		"created_at":  true,
		"rating":      true,
		"name":        true,
		"position":    true,
		"can_publish": true,
		"consent":     true,
	}
	if !allowedOrderBy[orderBy] {
		orderBy = "created_at"
	}

	search = strings.TrimSpace(search)

	var args []any
	where := ""
	if search != "" {
		where = "WHERE (r.name ILIKE $1 OR r.position ILIKE $1)"
		args = append(args, "%"+search+"%")
	}

	countQ := `SELECT COALESCE(COUNT(*), 0) FROM reviews r ` + where

	var total int
	if err := s.db.GetContext(ctx, &total, countQ, args...); err != nil {
		return nil, 0, fmt.Errorf("reviews count: %w", err)
	}
	if total == 0 {
		return []entity.Review{}, 0, nil
	}

	pageAmount := (total + pageSize - 1) / pageSize

	// LIMIT/OFFSET добавим как отдельные параметры, чтобы не плясать с $2/$3
	baseQ := `
		SELECT
			r.id, r.name, r.position, r.text, r.rating, r.image_path, r.consent, r.can_publish, r.created_at
		FROM reviews r
	` + " " + where + fmt.Sprintf(" ORDER BY r.%s DESC LIMIT %d OFFSET %d", orderBy, pageSize, offset)

	var out []entity.Review
	if err := s.db.SelectContext(ctx, &out, baseQ, args...); err != nil {
		return nil, 0, fmt.Errorf("reviews select: %w", err)
	}

	filename := ""
	for i := range out {
		filename = out[i].ImagePath
		out[i].ImagePath = s.domain + REVIEW_PUBLISHED_URL + filename
	}

	return out, pageAmount, nil
}

func (s *AdminReviewService) GetByID(ctx context.Context, id string) (*entity.Review, error) {
	const q = `
		SELECT id, name, position, text, rating, image_path, consent, can_publish, created_at
		FROM reviews
		WHERE id = $1
		LIMIT 1
	`
	var out entity.Review
	if err := s.db.GetContext(ctx, &out, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("review get by id: %w", err)
	}
	return &out, nil
}

func (s *AdminReviewService) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM reviews WHERE id = $1`
	if _, err := s.db.ExecContext(ctx, q, id); err != nil {
		return fmt.Errorf("review delete: %w", err)
	}
	return nil
}

type BulkUpdateReview struct {
	ID         string  `json:"id"`
	Name       *string `json:"name"`
	Position   *string `json:"position"`
	Text       *string `json:"text"`
	Rating     *int    `json:"rating"`
	Consent    *bool   `json:"consent"`
	CanPublish *bool   `json:"canPublish"`
	// imagePath тут не трогаем (без файла). Если надо — скажешь, добавим.
}

func (s *AdminReviewService) BulkUpdate(ctx context.Context, items []BulkUpdateReview) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("bulk update begin: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	const q = `
		UPDATE reviews
		SET
			name        = COALESCE($2, name),
			position    = COALESCE($3, position),
			text        = COALESCE($4, text),
			rating      = COALESCE($5, rating),
			consent     = COALESCE($6, consent),
			can_publish = COALESCE($7, can_publish)
		WHERE id = $1
	`

	for _, it := range items {
		if strings.TrimSpace(it.ID) == "" {
			return fmt.Errorf("bulk update: empty id")
		}
		if _, err := tx.ExecContext(ctx, q, it.ID, it.Name, it.Position, it.Text, it.Rating, it.Consent, it.CanPublish); err != nil {
			return fmt.Errorf("bulk update %s: %w", it.ID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("bulk update commit: %w", err)
	}
	return nil
}

func (s *AdminReviewService) Create(ctx context.Context, name, position, text string, rating int, imagePath string, consent bool) (string, error) {
	id := "rev-" + uuid.NewString()
	now := time.Now()

	const q = `
		INSERT INTO reviews (id, name, position, text, rating, image_path, consent, can_publish, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, FALSE, $8)
	`
	if _, err := s.db.ExecContext(ctx, q, id, name, position, text, rating, imagePath, consent, now); err != nil {
		return "", fmt.Errorf("review create: %w", err)
	}
	return id, nil
}
