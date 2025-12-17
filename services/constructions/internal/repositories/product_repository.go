package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/pachv/constructions/constructions/internal/domain/entity"
)

type ProductRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

type categoryRow struct {
	ID        string       `db:"id"`
	Title     string       `db:"title"`
	Slug      string       `db:"slug"`
	ImagePath *string      `db:"image_path"`
	CreatedAt sql.NullTime `db:"created_at"`
}

func (r *ProductRepository) GetAllCategories(ctx context.Context) ([]entity.CatalogCategory, error) {
	const q = `
		SELECT id, title, slug, image_path, created_at
		FROM catalog_categories
		ORDER BY title
	`

	var rows []categoryRow
	if err := r.db.SelectContext(ctx, &rows, q); err != nil {
		return nil, fmt.Errorf("get all categories: %w", err)
	}

	res := make([]entity.CatalogCategory, 0, len(rows))
	for _, row := range rows {
		var createdAtPtr *time.Time
		if row.CreatedAt.Valid {
			t := row.CreatedAt.Time
			createdAtPtr = &t
		}

		res = append(res, entity.CatalogCategory{
			ID:        row.ID,
			Title:     row.Title,
			Slug:      row.Slug,
			ImagePath: row.ImagePath,
			CreatedAt: createdAtPtr,
		})
	}

	return res, nil
}

// --- sections ---

type sectionRow struct {
	ID                 string       `db:"id"`
	Title              string       `db:"title"`
	Slug               string       `db:"slug"`
	ParentCategorySlug string       `db:"parent_category_slug"`
	ImagePath          *string      `db:"image_path"`
	CreatedAt          sql.NullTime `db:"created_at"`
}

func (r *ProductRepository) GetAllSections(ctx context.Context) ([]entity.CatalogSection, error) {
	const q = `
		SELECT
			s.id,
			s.title,
			s.slug,
			COALESCE(MIN(c.slug), '') AS parent_category_slug,
			s.image_path,
			s.created_at
		FROM catalog_sections s
		LEFT JOIN catalog_category_sections cs ON cs.section_id = s.id
		LEFT JOIN catalog_categories c ON c.id = cs.category_id
		GROUP BY s.id, s.title, s.slug, s.image_path, s.created_at
		ORDER BY s.title
	`

	var rows []sectionRow
	if err := r.db.SelectContext(ctx, &rows, q); err != nil {
		return nil, fmt.Errorf("get all sections: %w", err)
	}

	res := make([]entity.CatalogSection, 0, len(rows))
	for _, row := range rows {
		var createdAtPtr *time.Time
		if row.CreatedAt.Valid {
			t := row.CreatedAt.Time
			createdAtPtr = &t
		}

		res = append(res, entity.CatalogSection{
			ID:                 row.ID,
			Title:              row.Title,
			Slug:               row.Slug,
			ParentCategorySlug: row.ParentCategorySlug,
			ImagePath:          row.ImagePath,
			CreatedAt:          createdAtPtr,
		})
	}

	return res, nil
}
