package repositories

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pachv/constructions/constructions/internal/domain/entity"
)

type SiteSectionRepository struct {
	db *sqlx.DB
}

func NewSiteSectionRepository(db *sqlx.DB) *SiteSectionRepository {
	return &SiteSectionRepository{db: db}
}

func (r *SiteSectionRepository) GetAll(ctx context.Context) ([]entity.SiteSectionSummary, error) {
	const q = `
		SELECT id, title, slug, has_gallery, has_catalog
		FROM site_sections
		ORDER BY title
	`

	var rows []struct {
		ID         string `db:"id"`
		Title      string `db:"title"`
		Slug       string `db:"slug"`
		HasGallery bool   `db:"has_gallery"`
		HasCatalog bool   `db:"has_catalog"`
	}

	if err := r.db.SelectContext(ctx, &rows, q); err != nil {
		return nil, fmt.Errorf("site_sections get all: %w", err)
	}

	res := make([]entity.SiteSectionSummary, 0, len(rows))
	for _, row := range rows {
		res = append(res, entity.SiteSectionSummary{
			ID:         row.ID,
			Title:      row.Title,
			Slug:       row.Slug,
			HasGallery: row.HasGallery,
			HasCatalog: row.HasCatalog,
		})
	}

	return res, nil
}

func (r *SiteSectionRepository) GetBySlug(ctx context.Context, slug string) (*entity.SiteSection, error) {
	const sectionQ = `
		SELECT id, title, slug, has_gallery, has_catalog
		FROM site_sections
		WHERE slug = $1
	`

	var sec struct {
		ID         string `db:"id"`
		Title      string `db:"title"`
		Slug       string `db:"slug"`
		HasGallery bool   `db:"has_gallery"`
		HasCatalog bool   `db:"has_catalog"`
	}

	if err := r.db.GetContext(ctx, &sec, sectionQ, slug); err != nil {
		return nil, fmt.Errorf("site_sections get by slug: %w", err)
	}

	out := &entity.SiteSection{
		ID:         sec.ID,
		Title:      sec.Title,
		Slug:       sec.Slug,
		HasGallery: sec.HasGallery,
		HasCatalog: sec.HasCatalog,
		Gallery:    []entity.SiteSectionGallery{},
	}

	// gallery
	if sec.HasGallery {
		const galleryQ = `
			SELECT id, name, url, sort_order
			FROM site_section_gallery
			WHERE section_id = $1
			ORDER BY sort_order
		`
		var g []entity.SiteSectionGallery
		if err := r.db.SelectContext(ctx, &g, galleryQ, sec.ID); err != nil {
			return nil, fmt.Errorf("gallery select: %w", err)
		}
		out.Gallery = g
	}

	// catalog
	if sec.HasCatalog {
		cat := &entity.SiteSectionCatalog{
			Categories: []entity.SiteSectionCatalogCategory{},
			Items:      []entity.SiteSectionCatalogItem{},
		}

		const categoriesQ = `
			SELECT c.id, c.title, c.slug, scc.sort_order
			FROM site_section_catalog_categories scc
			JOIN catalog_categories c ON c.id = scc.category_id
			WHERE scc.section_id = $1
			ORDER BY scc.sort_order, c.title
		`
		if err := r.db.SelectContext(ctx, &cat.Categories, categoriesQ, sec.ID); err != nil {
			return nil, fmt.Errorf("catalog categories select: %w", err)
		}

		const itemsQ = `
			SELECT id, category_id, title, price_rub, image_url, sort_order
			FROM site_section_catalog_items
			WHERE section_id = $1
			ORDER BY sort_order, title
		`
		if err := r.db.SelectContext(ctx, &cat.Items, itemsQ, sec.ID); err != nil {
			return nil, fmt.Errorf("catalog items select: %w", err)
		}

		// badges + specs на каждый item (простая реализация)
		for i := range cat.Items {
			itemID := cat.Items[i].ID

			const badgesQ = `
				SELECT badge
				FROM site_section_catalog_item_badges
				WHERE item_id = $1
				ORDER BY sort_order, badge
			`
			var badges []string
			if err := r.db.SelectContext(ctx, &badges, badgesQ, itemID); err != nil {
				return nil, fmt.Errorf("item badges select: %w", err)
			}
			cat.Items[i].Badges = badges

			const specsQ = `
				SELECT key, value
				FROM site_section_catalog_item_specs
				WHERE item_id = $1
				ORDER BY sort_order, key
			`
			var specs []entity.SiteSectionItemSpec
			if err := r.db.SelectContext(ctx, &specs, specsQ, itemID); err != nil {
				return nil, fmt.Errorf("item specs select: %w", err)
			}
			cat.Items[i].Specs = specs
		}

		out.Catalog = cat
	}

	return out, nil
}
