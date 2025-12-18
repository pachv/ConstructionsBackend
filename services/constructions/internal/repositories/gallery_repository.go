package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pachv/constructions/constructions/internal/domain/entity"
)

type GalleryRepository struct {
	db     *sqlx.DB
	logger *slog.Logger
}

func NewGalleryRepository(db *sqlx.DB, logger *slog.Logger) *GalleryRepository {
	return &GalleryRepository{db: db, logger: logger}
}

type galleryCategoryRow struct {
	ID        string       `db:"id"`
	Title     string       `db:"title"`
	Slug      string       `db:"slug"`
	CreatedAt sql.NullTime `db:"created_at"`
}

func (r *GalleryRepository) GetAllCategories(ctx context.Context) ([]entity.GalleryCategory, error) {
	const q = `
		SELECT id, title, slug, created_at
		FROM gallery_categories
		ORDER BY title
	`

	r.logger.Debug("Gallery.GetAllCategories: start")

	var rows []galleryCategoryRow
	if err := r.db.SelectContext(ctx, &rows, q); err != nil {
		r.logger.Error("Gallery.GetAllCategories: select failed", "err", err, "query", q)
		return nil, fmt.Errorf("get gallery categories: %w", err)
	}

	res := make([]entity.GalleryCategory, 0, len(rows))
	for _, row := range rows {
		var createdAt *time.Time
		if row.CreatedAt.Valid {
			t := row.CreatedAt.Time
			createdAt = &t
		}
		res = append(res, entity.GalleryCategory{
			ID:        row.ID,
			Title:     row.Title,
			Slug:      row.Slug,
			CreatedAt: createdAt,
		})
	}

	r.logger.Debug("Gallery.GetAllCategories: done", "count", len(res))
	return res, nil
}

type galleryPhotoRow struct {
	ID           string       `db:"id"`
	CategorySlug string       `db:"category_slug"`
	Alt          string       `db:"alt"`
	ImagePath    string       `db:"image_path"`
	SortOrder    int          `db:"sort_order"`
	CreatedAt    sql.NullTime `db:"created_at"`
}

func (r *GalleryRepository) GetPhotosByCategorySlug(ctx context.Context, categorySlug string) ([]entity.GalleryPhoto, error) {
	const q = `
		SELECT
			p.id,
			c.slug AS category_slug,
			p.alt,
			p.image_path,
			p.sort_order,
			p.created_at
		FROM gallery_photos p
		JOIN gallery_categories c ON c.id = p.category_id
		WHERE c.slug = $1
		ORDER BY p.sort_order, p.created_at
	`

	r.logger.Debug("Gallery.GetPhotosByCategorySlug: start", "category_slug", categorySlug)

	var rows []galleryPhotoRow
	if err := r.db.SelectContext(ctx, &rows, q, categorySlug); err != nil {
		r.logger.Error("Gallery.GetPhotosByCategorySlug: select failed",
			"err", err,
			"category_slug", categorySlug,
			"query", q,
		)
		return nil, fmt.Errorf("get gallery photos: %w", err)
	}

	res := make([]entity.GalleryPhoto, 0, len(rows))
	for _, row := range rows {
		var createdAt *time.Time
		if row.CreatedAt.Valid {
			t := row.CreatedAt.Time
			createdAt = &t
		}

		res = append(res, entity.GalleryPhoto{
			ID:           row.ID,
			CategorySlug: row.CategorySlug,
			Alt:          row.Alt,
			ImagePath:    row.ImagePath,
			SortOrder:    row.SortOrder,
			CreatedAt:    createdAt,
		})
	}

	r.logger.Debug("Gallery.GetPhotosByCategorySlug: done", "count", len(res), "category_slug", categorySlug)
	return res, nil
}
