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

type ProductRepository struct {
	db     *sqlx.DB
	logger *slog.Logger
}

func NewProductRepository(db *sqlx.DB, logger *slog.Logger) *ProductRepository {
	return &ProductRepository{db: db, logger: logger}
}

// --- Categories ---

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
		r.logger.Error("GetAllCategories: select failed", "err", err)
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

	r.logger.Debug("GetAllCategories: done", "count", len(res))
	return res, nil
}

// --- Sections ---

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
		r.logger.Error("GetAllSections: select failed", "err", err)
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

	r.logger.Debug("GetAllSections: done", "count", len(res))
	return res, nil
}

// --- Products ---

type productRow struct {
	ID           string        `db:"id"`
	Title        string        `db:"title"`
	Slug         string        `db:"slug"`
	CategorySlug string        `db:"category_slug"`
	SectionSlug  string        `db:"section_slug"`
	Brand        string        `db:"brand"`
	Type         string        `db:"type"`
	Price        int           `db:"price"`
	OldPrice     sql.NullInt64 `db:"old_price"`
	InStock      sql.NullBool  `db:"in_stock"`
	SalePercent  sql.NullInt64 `db:"sale_percent"`
	ImagePath    *string       `db:"image_path"`
	CreatedAt    sql.NullTime  `db:"created_at"`
}

type badgeLinkRow struct {
	ProductID string `db:"product_id"`
	Code      string `db:"code"`
}

func (r *ProductRepository) GetAllProducts(ctx context.Context) ([]entity.CatalogProduct, error) {
	const productsQ = `
		SELECT
			id, title, slug,
			category_slug, section_slug,
			brand, type,
			price, old_price,
			in_stock,
			sale_percent,
			image_path,
			created_at
		FROM catalog_products
		ORDER BY title
	`

	r.logger.Debug("GetAllProducts: start")

	var prodRows []productRow
	if err := r.db.SelectContext(ctx, &prodRows, productsQ); err != nil {
		r.logger.Error("GetAllProducts: products select failed", "err", err)
		return nil, fmt.Errorf("get all products: %w", err)
	}

	r.logger.Debug("GetAllProducts: products loaded", "count", len(prodRows))

	const badgesQ = `
		SELECT
			l.product_id,
			b.code
		FROM product_badge_links l
		JOIN product_badges b ON b.id = l.badge_id
	`

	var linkRows []badgeLinkRow
	if err := r.db.SelectContext(ctx, &linkRows, badgesQ); err != nil {
		r.logger.Error("GetAllProducts: badges select failed", "err", err)
		return nil, fmt.Errorf("get product badges: %w", err)
	}

	r.logger.Debug("GetAllProducts: badge links loaded", "count", len(linkRows))

	badgesByProduct := map[string][]string{}
	for _, lr := range linkRows {
		badgesByProduct[lr.ProductID] = append(badgesByProduct[lr.ProductID], lr.Code)
	}

	res := make([]entity.CatalogProduct, 0, len(prodRows))

	for i, rrow := range prodRows {
		if rrow.ID == "" || rrow.Slug == "" {
			r.logger.Warn("GetAllProducts: suspicious row (empty id/slug)",
				"row_index", i,
				"id", rrow.ID,
				"slug", rrow.Slug,
				"title", rrow.Title,
			)
		}

		var oldPrice *int
		if rrow.OldPrice.Valid {
			v := int(rrow.OldPrice.Int64)
			oldPrice = &v
		}

		inStock := false
		if rrow.InStock.Valid {
			inStock = rrow.InStock.Bool
		} else {
			r.logger.Warn("GetAllProducts: in_stock is NULL, defaulting to false",
				"product_id", rrow.ID,
				"slug", rrow.Slug,
			)
		}

		var salePercent *int
		if rrow.SalePercent.Valid {
			v := int(rrow.SalePercent.Int64)
			salePercent = &v
		}

		var createdAt *time.Time
		if rrow.CreatedAt.Valid {
			t := rrow.CreatedAt.Time
			createdAt = &t
		}

		res = append(res, entity.CatalogProduct{
			ID:           rrow.ID,
			Title:        rrow.Title,
			Slug:         rrow.Slug,
			CategorySlug: rrow.CategorySlug,
			SectionSlug:  rrow.SectionSlug,
			Brand:        rrow.Brand,
			Type:         rrow.Type,
			Price:        rrow.Price,
			OldPrice:     oldPrice,
			InStock:      inStock,
			Badges:       badgesByProduct[rrow.ID],
			SalePercent:  salePercent,
			ImagePath:    rrow.ImagePath,
			CreatedAt:    createdAt,
		})
	}

	r.logger.Debug("GetAllProducts: done", "count", len(res))
	return res, nil
}
