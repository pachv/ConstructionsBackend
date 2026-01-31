package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pachv/constructions/constructions/internal/domain/entity"
)

type SiteSectionsRepository struct {
	db     *sqlx.DB
	domain string
}

func NewSiteSectionsRepository(db *sqlx.DB, domain string) *SiteSectionsRepository {
	return &SiteSectionsRepository{
		db:     db,
		domain: strings.TrimRight(strings.TrimSpace(domain), "/"),
	}
}

const MAIN_PICTURE_URL = "/sections/picture/"
const GALLERY_PICTURE_URL = "/sections/gallery/picture/"
const CATALOG_PICTURE_URL = "/catalog/picture/"

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

	for i := range out {
		if strings.TrimSpace(out[i].Label) == "" {
			out[i].Label = out[i].Title
		}

		fn := strings.TrimSpace(out[i].Image)
		if fn != "" {
			out[i].Image = r.withDomainPicture(MAIN_PICTURE_URL, fn)
		} else {
			out[i].Image = ""
		}
	}
	return out, nil
}

// GET /api/v1/sections/:slug
func (r *SiteSectionsRepository) GetBySlugFull(ctx context.Context, slug string) (*entity.SiteSection, error) {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return nil, nil
	}

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
		Image:           "",
		AdvantegesText:  s.AdvText,
		AdvantegesArray: []string{},
		HasGallery:      s.HasGal,
		HasCatalog:      s.HasCat,
		Gallery:         []entity.SiteSectionGallery{},
		Catalog: &entity.SiteSectionCatalog{
			Categories: []entity.SiteSectionCatalogCategory{},
			Items:      []entity.SiteSectionCatalogItem{},
		},
	}

	if strings.TrimSpace(out.Label) == "" {
		out.Label = out.Title
	}
	if out.AdvantegesText == "" {
		out.AdvantegesText = ""
	}
	if strings.TrimSpace(s.ImageURL) != "" {
		out.Image = r.withDomainPicture(MAIN_PICTURE_URL, s.ImageURL)
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
	if out.AdvantegesArray == nil {
		out.AdvantegesArray = []string{}
	}

	// gallery
	if out.HasGallery {
		const galleryQ = `
			SELECT id, section_id, name, url, sort_order
			FROM site_section_gallery
			WHERE section_id = $1
			ORDER BY sort_order
		`
		if err := r.db.SelectContext(ctx, &out.Gallery, galleryQ, out.ID); err != nil {
			return nil, fmt.Errorf("site_section_gallery select: %w", err)
		}

		for i := range out.Gallery {
			u := strings.TrimSpace(out.Gallery[i].URL)
			if u != "" {
				out.Gallery[i].URL = r.withDomainPicture(GALLERY_PICTURE_URL, u)
			}
		}
	} else {
		out.Gallery = []entity.SiteSectionGallery{}
	}

	// catalog
	if out.HasCatalog {
		cat := &entity.SiteSectionCatalog{
			Categories: []entity.SiteSectionCatalogCategory{},
			Items:      []entity.SiteSectionCatalogItem{},
		}

		// ✅ categories: теперь без catalog_categories, берём title/slug прямо из site_section_catalog_categories
		// Важно: делаем "category_id AS id", чтобы entity.SiteSectionCatalogCategory.ID заполнялся как раньше.
		const categoriesQ = `
			SELECT
				category_id AS id,
				title,
				slug,
				sort_order
			FROM site_section_catalog_categories
			WHERE section_id = $1
			ORDER BY sort_order, title
		`
		if err := r.db.SelectContext(ctx, &cat.Categories, categoriesQ, out.ID); err != nil {
			return nil, fmt.Errorf("catalog categories select: %w", err)
		}
		if cat.Categories == nil {
			cat.Categories = []entity.SiteSectionCatalogCategory{}
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
		if cat.Items == nil {
			cat.Items = []entity.SiteSectionCatalogItem{}
		}

		for i := range cat.Items {
			if strings.TrimSpace(cat.Items[i].ImageURL) != "" {
				cat.Items[i].ImageURL = r.withDomainPicture(CATALOG_PICTURE_URL, cat.Items[i].ImageURL)
			} else {
				cat.Items[i].ImageURL = ""
			}
		}

		// specs на каждый item (как у тебя было)
		for i := range cat.Items {
			itemID := cat.Items[i].ID

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
			if specs == nil {
				specs = []entity.SiteSectionItemSpec{}
			}
			cat.Items[i].Specs = specs
		}

		out.Catalog = cat
	} else {
		out.Catalog = &entity.SiteSectionCatalog{
			Categories: []entity.SiteSectionCatalogCategory{},
			Items:      []entity.SiteSectionCatalogItem{},
		}
	}

	return out, nil
}

// helpers

func (r *SiteSectionsRepository) withDomainPicture(prefix, v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}
	if strings.HasPrefix(v, "http://") || strings.HasPrefix(v, "https://") {
		return v
	}
	if strings.HasPrefix(v, "/") {
		return r.domain + v
	}
	return r.domain + prefix + v
}
