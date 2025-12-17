package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pachv/constructions/constructions/internal/domain/entity"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetAllCategories(ctx context.Context) ([]entity.CatalogCategory, error) {
	const q = `
		SELECT id, title, slug, image_path, created_at
		FROM catalog_categories
		ORDER BY title
	`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("get all categories: %w", err)
	}
	defer rows.Close()

	var res []entity.CatalogCategory

	for rows.Next() {
		var c entity.CatalogCategory
		var nt sql.NullTime

		if err := rows.Scan(
			&c.ID,
			&c.Title,
			&c.Slug,
			&c.ImagePath,
			&nt,
		); err != nil {
			return nil, fmt.Errorf("scan category: %w", err)
		}

		if nt.Valid {
			t := nt.Time
			c.CreatedAt = &t
		}

		res = append(res, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return res, nil
}

func (r *ProductRepository) GetAllSections(ctx context.Context) ([]entity.CatalogSection, error) {
	// parentCategorySlug берём как MIN(slug) из связанных категорий (детерминированно)
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

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("get all sections: %w", err)
	}
	defer rows.Close()
	var res []entity.CatalogSection
	for rows.Next() {
		var s entity.CatalogSection
		var nt sql.NullTime

		if err := rows.Scan(
			&s.ID,
			&s.Title,
			&s.Slug,
			&s.ParentCategorySlug,
			&s.ImagePath,
			&nt,
		); err != nil {
			return nil, fmt.Errorf("scan section: %w", err)
		}

		if nt.Valid {
			t := nt.Time
			s.CreatedAt = &t
		}

		res = append(res, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows sections: %w", err)
	}

	return res, nil
}
