package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pachv/constructions/constructions/internal/domain/entity"
)

type SiteSectionsRepository struct {
	db *sqlx.DB
}

func NewSiteSectionsRepository(db *sqlx.DB) *SiteSectionsRepository {
	return &SiteSectionsRepository{db: db}
}

// GET /api/v1/sections
func (r *SiteSectionsRepository) GetAll(ctx context.Context) ([]entity.SiteSectionSummary, error) {
	const q = `
		SELECT id, title, label, slug, image_url, has_gallery, has_catalog
		FROM site_sections
		ORDER BY title
	`

	var out []entity.SiteSectionSummary
	if err := r.db.SelectContext(ctx, &out, q); err != nil {
		return nil, fmt.Errorf("site_sections get all: %w", err)
	}

	// defaults
	for i := range out {
		if out[i].Label == "" {
			out[i].Label = out[i].Title
		}
	}
	return out, nil
}

// GET /api/v1/sections/:slug
func (r *SiteSectionsRepository) GetBySlugFull(ctx context.Context, slug string) (*entity.SiteSection, error) {
	const sectionQ = `
		SELECT id, title, label, slug, image_url, advanteges_text, has_gallery, has_catalog
		FROM site_sections
		WHERE slug = $1
		LIMIT 1
	`

	var s struct {
		ID       string `db:"id"`
		Title    string `db:"title"`
		Label    string `db:"label"`
		Slug     string `db:"slug"`
		ImageURL string `db:"image_url"`
		AdvText  string `db:"advanteges_text"`
		HasGal   bool   `db:"has_gallery"`
		HasCat   bool   `db:"has_catalog"`
	}

	if err := r.db.GetContext(ctx, &s, sectionQ, slug); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("site_sections get by slug: %w", err)
	}

	out := &entity.SiteSection{
		ID:              s.ID,
		Title:           s.Title,
		Label:           s.Label,
		Slug:            s.Slug,
		Image:           s.ImageURL,
		AdvantegesText:  s.AdvText,
		AdvantegesArray: []string{},
		HasGallery:      s.HasGal,
		HasCatalog:      s.HasCat,
		Gallery:         []entity.SiteSectionGallery{},
	}

	if out.Label == "" {
		out.Label = out.Title
	}
	if out.AdvantegesText == "" {
		out.AdvantegesText = ""
	}

	// advanteges array
	const advQ = `
		SELECT text
		FROM site_section_advanteges
		WHERE section_id = $1
		ORDER BY sort_order, text
	`
	if err := r.db.SelectContext(ctx, &out.AdvantegesArray, advQ, out.ID); err != nil {
		return nil, fmt.Errorf("site_section_advanteges select: %w", err)
	}

	// gallery
	if out.HasGallery {
		const galleryQ = `
			SELECT id, name, url, sort_order
			FROM site_section_gallery
			WHERE section_id = $1
			ORDER BY sort_order
		`
		if err := r.db.SelectContext(ctx, &out.Gallery, galleryQ, out.ID); err != nil {
			return nil, fmt.Errorf("site_section_gallery select: %w", err)
		}
	}

	// catalog
	if out.HasCatalog {
		cat := &entity.SiteSectionCatalog{
			Categories: []entity.SiteSectionCatalogCategory{},
			Items:      []entity.SiteSectionCatalogItem{},
		}

		// categories (ожидает, что catalog_categories уже есть)
		const categoriesQ = `
			SELECT c.id, c.title, c.slug, scc.sort_order
			FROM site_section_catalog_categories scc
			JOIN catalog_categories c ON c.id = scc.category_id
			WHERE scc.section_id = $1
			ORDER BY scc.sort_order, c.title
		`
		if err := r.db.SelectContext(ctx, &cat.Categories, categoriesQ, out.ID); err != nil {
			return nil, fmt.Errorf("catalog categories select: %w", err)
		}

		const itemsQ = `
			SELECT id, category_id, title, price_rub, image_url, sort_order
			FROM site_section_catalog_items
			WHERE section_id = $1
			ORDER BY sort_order, title
		`
		if err := r.db.SelectContext(ctx, &cat.Items, itemsQ, out.ID); err != nil {
			return nil, fmt.Errorf("catalog items select: %w", err)
		}

		// badges + specs на каждый item
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
