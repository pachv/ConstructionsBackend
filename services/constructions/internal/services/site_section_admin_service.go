package services

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/jmoiron/sqlx"

	"github.com/pachv/constructions/constructions/internal/domain/entity"
)

// NOTE:
// Этот сервис предполагает, что в таблице site_section_catalog_categories есть поля:
//   - title VARCHAR(255)
//   - slug  VARCHAR(255)
// чтобы можно было хранить категории прямо в site_section_catalog_categories,
// без использования catalog_categories.
//
// Миграция (пример):
//   ALTER TABLE site_section_catalog_categories ADD COLUMN IF NOT EXISTS title VARCHAR(255);
//   ALTER TABLE site_section_catalog_categories ADD COLUMN IF NOT EXISTS slug  VARCHAR(255);

type SiteSectionsAdminService struct {
	db     *sqlx.DB
	domain string
}

func NewSiteSectionsAdminService(db *sqlx.DB, domain string) *SiteSectionsAdminService {
	return &SiteSectionsAdminService{
		db:     db,
		domain: strings.TrimRight(strings.TrimSpace(domain), "/"),
	}
}

// ==== inputs ====

type AdminToggleInput struct {
	Enabled bool `json:"enabled"`
}

type AdminAddGalleryPhotoInput struct {
	Name      string `json:"name"`
	URL       string `json:"url"`       // "/sections/gallery/picture/x.jpg" | "x.jpg" | "https://..."
	SortOrder int    `json:"sortOrder"` // optional
}

type AdminAddCatalogItemSpecInput struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	SortOrder int    `json:"sortOrder"`
}

type AdminAddCatalogItemBadgeInput struct {
	Badge     string `json:"badge"`
	SortOrder int    `json:"sortOrder"`
}

type AdminAddCatalogCategoryInput struct {
	CategoryID string `json:"categoryId"`
	Title      string `json:"title"`
	Slug       string `json:"slug"`
	SortOrder  int    `json:"sortOrder"`
}

type AdminAddCatalogItemInput struct {
	CategoryID string `json:"categoryId"`
	Title      string `json:"title"`
	PriceRub   int    `json:"priceRub"`
	Image      string `json:"image"` // "/catalog/picture/x.jpg" | "x.jpg" | "https://..."
	SortOrder  int    `json:"sortOrder"`

	Specs  []AdminAddCatalogItemSpecInput  `json:"specs"`
	Badges []AdminAddCatalogItemBadgeInput `json:"badges"`
}

// ==== constants ====

const (
	MAIN_PICTURE_URL    = "/sections/picture/"
	GALLERY_PICTURE_URL = "/sections/gallery/picture/"
	CATALOG_PICTURE_URL = "/catalog/picture/"
)

// ==== public methods ====

func (s *SiteSectionsAdminService) ToggleGallery(ctx context.Context, sectionID string, enabled bool) error {
	sectionID = strings.TrimSpace(sectionID)
	if sectionID == "" {
		return fmt.Errorf("sectionID is required")
	}

	res, err := s.db.ExecContext(ctx, `
		UPDATE site_sections
		SET has_gallery = $2
		WHERE id = $1
	`, sectionID, enabled)
	if err != nil {
		return fmt.Errorf("toggle gallery: %w", err)
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return fmt.Errorf("section not found")
	}
	return nil
}

func (s *SiteSectionsAdminService) ToggleCatalog(ctx context.Context, sectionID string, enabled bool) error {
	sectionID = strings.TrimSpace(sectionID)
	if sectionID == "" {
		return fmt.Errorf("sectionID is required")
	}

	res, err := s.db.ExecContext(ctx, `
		UPDATE site_sections
		SET has_catalog = $2
		WHERE id = $1
	`, sectionID, enabled)
	if err != nil {
		return fmt.Errorf("toggle catalog: %w", err)
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return fmt.Errorf("section not found")
	}
	return nil
}

func (s *SiteSectionsAdminService) AddGalleryPhoto(ctx context.Context, sectionID string, in AdminAddGalleryPhotoInput) (string, error) {
	sectionID = strings.TrimSpace(sectionID)
	in.Name = strings.TrimSpace(in.Name)
	in.URL = strings.TrimSpace(in.URL)

	if sectionID == "" {
		return "", fmt.Errorf("sectionID is required")
	}
	if in.URL == "" {
		return "", fmt.Errorf("url is required")
	}
	if in.SortOrder == 0 {
		in.SortOrder = 1
	}

	photoID := "gal-" + randID()
	urlToStore := stripPicturePrefix(in.URL)

	if err := s.ensureSectionExists(ctx, sectionID); err != nil {
		return "", err
	}

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO site_section_gallery (id, section_id, name, url, sort_order)
		VALUES ($1,$2,$3,$4,$5)
	`, photoID, sectionID, in.Name, urlToStore, in.SortOrder)
	if err != nil {
		return "", fmt.Errorf("add gallery photo: %w", err)
	}

	return photoID, nil
}

func (s *SiteSectionsAdminService) DeleteGalleryPhoto(ctx context.Context, sectionID string, photoID string) error {
	sectionID = strings.TrimSpace(sectionID)
	photoID = strings.TrimSpace(photoID)
	if sectionID == "" || photoID == "" {
		return fmt.Errorf("sectionID and photoID are required")
	}

	res, err := s.db.ExecContext(ctx, `
		DELETE FROM site_section_gallery
		WHERE id = $1 AND section_id = $2
	`, photoID, sectionID)
	if err != nil {
		return fmt.Errorf("delete gallery photo: %w", err)
	}

	aff, _ := res.RowsAffected()
	if aff == 0 {
		return fmt.Errorf("photo not found")
	}
	return nil
}

func (s *SiteSectionsAdminService) AddCatalogCategory(ctx context.Context, sectionID string, in AdminAddCatalogCategoryInput) (string, error) {
	sectionID = strings.TrimSpace(sectionID)
	in.CategoryID = strings.TrimSpace(in.CategoryID)
	in.Title = strings.TrimSpace(in.Title)
	in.Slug = strings.TrimSpace(in.Slug)

	if sectionID == "" {
		return "", fmt.Errorf("sectionID is required")
	}
	if in.CategoryID == "" {
		return "", fmt.Errorf("categoryId is required")
	}
	if in.Title == "" {
		return "", fmt.Errorf("title is required")
	}
	if in.SortOrder == 0 {
		in.SortOrder = 1
	}
	if in.Slug == "" {
		in.Slug = slugify(in.Title)
	}

	if err := s.ensureSectionExists(ctx, sectionID); err != nil {
		return "", err
	}

	id := "ssc-cat-" + randID()

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO site_section_catalog_categories (id, section_id, category_id, title, slug, sort_order)
		VALUES ($1,$2,$3,$4,$5,$6)
	`, id, sectionID, in.CategoryID, in.Title, in.Slug, in.SortOrder)
	if err != nil {
		return "", fmt.Errorf("add catalog category: %w", err)
	}

	return id, nil
}

func (s *SiteSectionsAdminService) DeleteCatalogCategory(ctx context.Context, sectionID, categoryID string) error {
	sectionID = strings.TrimSpace(sectionID)
	categoryID = strings.TrimSpace(categoryID)
	if sectionID == "" || categoryID == "" {
		return fmt.Errorf("sectionID and categoryID are required")
	}

	res, err := s.db.ExecContext(ctx, `
		DELETE FROM site_section_catalog_categories
		WHERE section_id = $1 AND category_id = $2
	`, sectionID, categoryID)
	if err != nil {
		return fmt.Errorf("delete catalog category: %w", err)
	}

	aff, _ := res.RowsAffected()
	if aff == 0 {
		return fmt.Errorf("category not found")
	}
	return nil
}

func (s *SiteSectionsAdminService) AddCatalogItem(ctx context.Context, sectionID string, in AdminAddCatalogItemInput) (string, error) {
	sectionID = strings.TrimSpace(sectionID)
	in.CategoryID = strings.TrimSpace(in.CategoryID)
	in.Title = strings.TrimSpace(in.Title)
	in.Image = strings.TrimSpace(in.Image)

	if sectionID == "" {
		return "", fmt.Errorf("sectionID is required")
	}
	if in.CategoryID == "" || in.Title == "" {
		return "", fmt.Errorf("categoryId and title are required")
	}
	if in.SortOrder == 0 {
		in.SortOrder = 1
	}

	if err := s.ensureSectionExists(ctx, sectionID); err != nil {
		return "", err
	}

	itemID := "prd-" + randID()
	imageToStore := stripPicturePrefix(in.Image)

	tx, err := s.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return "", fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO site_section_catalog_items
			(id, section_id, category_id, title, price_rub, image_url, sort_order)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
	`, itemID, sectionID, in.CategoryID, in.Title, in.PriceRub, imageToStore, in.SortOrder)
	if err != nil {
		return "", fmt.Errorf("insert item: %w", err)
	}

	// badges
	for i, b := range in.Badges {
		badge := strings.TrimSpace(b.Badge)
		if badge == "" {
			continue
		}
		sortOrder := b.SortOrder
		if sortOrder == 0 {
			sortOrder = i + 1
		}
		bID := "bad-" + randID()
		_, err = tx.ExecContext(ctx, `
			INSERT INTO site_section_catalog_item_badges (id, item_id, badge, sort_order)
			VALUES ($1,$2,$3,$4)
		`, bID, itemID, badge, sortOrder)
		if err != nil {
			return "", fmt.Errorf("insert badge: %w", err)
		}
	}

	// specs
	for i, sp := range in.Specs {
		key := strings.TrimSpace(sp.Key)
		value := strings.TrimSpace(sp.Value)
		if key == "" {
			continue
		}
		sortOrder := sp.SortOrder
		if sortOrder == 0 {
			sortOrder = i + 1
		}
		sID := "spec-" + randID()
		_, err = tx.ExecContext(ctx, `
			INSERT INTO site_section_catalog_item_specs (id, item_id, key, value, sort_order)
			VALUES ($1,$2,$3,$4,$5)
		`, sID, itemID, key, value, sortOrder)
		if err != nil {
			return "", fmt.Errorf("insert spec: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("commit: %w", err)
	}

	return itemID, nil
}

func (s *SiteSectionsAdminService) DeleteCatalogItem(ctx context.Context, sectionID, itemID string) error {
	sectionID = strings.TrimSpace(sectionID)
	itemID = strings.TrimSpace(itemID)
	if sectionID == "" || itemID == "" {
		return fmt.Errorf("sectionID and itemID are required")
	}

	tx, err := s.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	// проверяем, что item принадлежит секции
	var exists bool
	if err := tx.GetContext(ctx, &exists, `
		SELECT EXISTS(
			SELECT 1 FROM site_section_catalog_items
			WHERE id = $1 AND section_id = $2
		)
	`, itemID, sectionID); err != nil {
		return fmt.Errorf("check item: %w", err)
	}
	if !exists {
		return fmt.Errorf("item not found")
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM site_section_catalog_item_specs WHERE item_id = $1`, itemID)
	if err != nil {
		return fmt.Errorf("delete specs: %w", err)
	}
	_, err = tx.ExecContext(ctx, `DELETE FROM site_section_catalog_item_badges WHERE item_id = $1`, itemID)
	if err != nil {
		return fmt.Errorf("delete badges: %w", err)
	}
	_, err = tx.ExecContext(ctx, `
		DELETE FROM site_section_catalog_items
		WHERE id = $1 AND section_id = $2
	`, itemID, sectionID)
	if err != nil {
		return fmt.Errorf("delete item: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}
	return nil
}

func (s *SiteSectionsAdminService) GetAllSectionsSummary() ([]*entity.SiteSectionSummary, error) {
	q := `
		SELECT
			s.id,
			s.title,
			s.label,
			s.slug,
			s.image_url,
			s.has_gallery,
			s.has_catalog
		FROM site_sections s
		ORDER BY s.created_at DESC
	`

	var items []*entity.SiteSectionSummary
	if err := s.db.Select(&items, q); err != nil {
		return nil, fmt.Errorf("get all sections: %w", err)
	}

	for i := range items {
		it := items[i]
		if it == nil {
			continue
		}
		if strings.TrimSpace(it.Label) == "" {
			it.Label = it.Title
		}
		fn := strings.TrimSpace(it.Image)
		if fn != "" {
			it.Image = s.withDomainPicture(MAIN_PICTURE_URL, fn)
		}
	}

	if items == nil {
		items = []*entity.SiteSectionSummary{}
	}
	return items, nil
}

// =========================
// GET FULL BY SLUG
// =========================

func (s *SiteSectionsAdminService) GetSectionFullBySlug(slugStr string) (*entity.SiteSection, error) {
	slugStr = strings.TrimSpace(slugStr)
	if slugStr == "" {
		return nil, fmt.Errorf("slug is required")
	}

	var base struct {
		ID         string `db:"id"`
		Title      string `db:"title"`
		Label      string `db:"label"`
		Slug       string `db:"slug"`
		ImageURL   string `db:"image_url"`
		AdvText    string `db:"advanteges_text"`
		HasGallery bool   `db:"has_gallery"`
		HasCatalog bool   `db:"has_catalog"`
	}

	if err := s.db.Get(&base, `
		SELECT id, title, label, slug, image_url, advanteges_text, has_gallery, has_catalog
		FROM site_sections
		WHERE slug = $1
		LIMIT 1
	`, slugStr); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("section not found")
		}
		return nil, fmt.Errorf("get section: %w", err)
	}

	out := &entity.SiteSection{
		ID:              base.ID,
		Title:           base.Title,
		Label:           base.Label,
		Slug:            base.Slug,
		Image:           strings.TrimSpace(base.ImageURL),
		AdvantegesText:  base.AdvText,
		HasGallery:      base.HasGallery,
		HasCatalog:      base.HasCatalog,
		AdvantegesArray: []string{},
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
	if out.Image != "" {
		out.Image = s.withDomainPicture(MAIN_PICTURE_URL, out.Image)
	}

	// advantages array
	if err := s.db.Select(&out.AdvantegesArray, `
		SELECT a.text
		FROM site_section_advanteges a
		WHERE a.section_id = $1
		ORDER BY a.sort_order
	`, out.ID); err != nil {
		return nil, fmt.Errorf("get advantages: %w", err)
	}
	if out.AdvantegesArray == nil {
		out.AdvantegesArray = []string{}
	}

	// gallery
	if out.HasGallery {
		var gallery []entity.SiteSectionGallery
		if err := s.db.Select(&gallery, `
			SELECT id, section_id, name, url, sort_order
			FROM site_section_gallery
			WHERE section_id = $1
			ORDER BY sort_order
		`, out.ID); err != nil {
			return nil, fmt.Errorf("get gallery: %w", err)
		}

		for i := range gallery {
			u := strings.TrimSpace(gallery[i].URL)
			if u != "" {
				gallery[i].URL = s.withDomainPicture(GALLERY_PICTURE_URL, u)
			}
		}

		out.Gallery = gallery
	} else {
		out.Gallery = []entity.SiteSectionGallery{}
	}

	// catalog
	if out.HasCatalog {
		// categories: берём напрямую из site_section_catalog_categories (без catalog_categories)
		var cats []entity.SiteSectionCatalogCategory
		if err := s.db.Select(&cats, `
			SELECT
				category_id AS id,
				title,
				slug,
				sort_order
			FROM site_section_catalog_categories
			WHERE section_id = $1
			ORDER BY sort_order
		`, out.ID); err != nil {
			return nil, fmt.Errorf("get catalog categories: %w", err)
		}
		if cats == nil {
			cats = []entity.SiteSectionCatalogCategory{}
		}
		out.Catalog.Categories = cats

		// items
		var items []entity.SiteSectionCatalogItem
		if err := s.db.Select(&items, `
			SELECT
				i.id,
				i.category_id,
				i.title,
				i.price_rub,
				i.image_url,
				i.sort_order
			FROM site_section_catalog_items i
			WHERE i.section_id = $1
			ORDER BY i.sort_order
		`, out.ID); err != nil {
			return nil, fmt.Errorf("get catalog items: %w", err)
		}
		if items == nil {
			items = []entity.SiteSectionCatalogItem{}
		}

		// badges пачкой
		type badgeRow struct {
			ItemID string `db:"item_id"`
			Badge  string `db:"badge"`
		}
		var badgeRows []badgeRow
		if err := s.db.Select(&badgeRows, `
			SELECT item_id, badge
			FROM site_section_catalog_item_badges
			WHERE item_id IN (SELECT id FROM site_section_catalog_items WHERE section_id = $1)
			ORDER BY sort_order
		`, out.ID); err != nil {
			return nil, fmt.Errorf("get item badges: %w", err)
		}
		badgeMap := map[string][]string{}
		for _, r := range badgeRows {
			badgeMap[r.ItemID] = append(badgeMap[r.ItemID], r.Badge)
		}

		// specs пачкой
		type specRow struct {
			ItemID string `db:"item_id"`
			Key    string `db:"key"`
			Value  string `db:"value"`
		}
		var specRows []specRow
		if err := s.db.Select(&specRows, `
			SELECT item_id, key, value
			FROM site_section_catalog_item_specs
			WHERE item_id IN (SELECT id FROM site_section_catalog_items WHERE section_id = $1)
			ORDER BY sort_order
		`, out.ID); err != nil {
			return nil, fmt.Errorf("get item specs: %w", err)
		}
		specMap := map[string][]entity.SiteSectionItemSpec{}
		for _, r := range specRows {
			specMap[r.ItemID] = append(specMap[r.ItemID], entity.SiteSectionItemSpec{
				Key:   r.Key,
				Value: r.Value,
			})
		}

		for i := range items {
			// image full url
			fn := strings.TrimSpace(items[i].ImageURL)
			if fn != "" {
				items[i].ImageURL = s.withDomainPicture(CATALOG_PICTURE_URL, fn)
			}

			// badges
			items[i].Badges = badgeMap[items[i].ID]
			if items[i].Badges == nil {
				items[i].Badges = []string{}
			}

			// specs
			items[i].Specs = specMap[items[i].ID]
			if items[i].Specs == nil {
				items[i].Specs = []entity.SiteSectionItemSpec{}
			}
		}

		out.Catalog.Items = items
	} else {
		out.Catalog = &entity.SiteSectionCatalog{
			Categories: []entity.SiteSectionCatalogCategory{},
			Items:      []entity.SiteSectionCatalogItem{},
		}
	}

	return out, nil
}

// =========================
// DELETE SECTION
// =========================

func (s *SiteSectionsAdminService) DeleteSection(ctx context.Context, sectionID string) error {
	sectionID = strings.TrimSpace(sectionID)
	if sectionID == "" {
		return fmt.Errorf("sectionID is required")
	}

	tx, err := s.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	var exists bool
	if err := tx.GetContext(ctx, &exists, `SELECT EXISTS(SELECT 1 FROM site_sections WHERE id=$1)`, sectionID); err != nil {
		return fmt.Errorf("check section: %w", err)
	}
	if !exists {
		return fmt.Errorf("section not found")
	}

	_, _ = tx.ExecContext(ctx, `DELETE FROM site_section_advanteges WHERE section_id = $1`, sectionID)
	_, _ = tx.ExecContext(ctx, `DELETE FROM site_section_gallery WHERE section_id = $1`, sectionID)

	_, _ = tx.ExecContext(ctx, `
		DELETE FROM site_section_catalog_item_specs 
		WHERE item_id IN (SELECT id FROM site_section_catalog_items WHERE section_id = $1)
	`, sectionID)
	_, _ = tx.ExecContext(ctx, `
		DELETE FROM site_section_catalog_item_badges 
		WHERE item_id IN (SELECT id FROM site_section_catalog_items WHERE section_id = $1)
	`, sectionID)

	_, _ = tx.ExecContext(ctx, `DELETE FROM site_section_catalog_items WHERE section_id = $1`, sectionID)
	_, _ = tx.ExecContext(ctx, `DELETE FROM site_section_catalog_categories WHERE section_id = $1`, sectionID)

	_, err = tx.ExecContext(ctx, `DELETE FROM site_sections WHERE id = $1`, sectionID)
	if err != nil {
		return fmt.Errorf("delete section: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}

// =========================
// CREATE / UPDATE FROM FORM
// =========================

type AdminCreateSectionFormInput struct {
	Title          string
	ImageFilename  string   // уже сохранённый файл
	AdvantegesText string   // textarea
	Advanteges     []string // список преимуществ
}

func (s *SiteSectionsAdminService) CreateSectionFromForm(ctx context.Context, in AdminCreateSectionFormInput) (string, error) {
	in.Title = strings.TrimSpace(in.Title)
	in.ImageFilename = strings.TrimSpace(in.ImageFilename)
	in.AdvantegesText = strings.TrimSpace(in.AdvantegesText)

	if in.Title == "" {
		return "", fmt.Errorf("title is required")
	}
	if in.ImageFilename == "" {
		return "", fmt.Errorf("image is required")
	}

	id := "sec-" + randID()
	slugStr := slugify(in.Title)

	tx, err := s.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return "", fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO site_sections (id, title, label, slug, image_url, advanteges_text, has_gallery, has_catalog)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`, id, in.Title, in.Title, slugStr, in.ImageFilename, in.AdvantegesText, false, false)
	if err != nil {
		return "", fmt.Errorf("insert section: %w", err)
	}

	order := 1
	for _, t := range in.Advanteges {
		t = strings.TrimSpace(t)
		if t == "" {
			continue
		}
		advID := "adv-" + randID()
		_, err = tx.ExecContext(ctx, `
			INSERT INTO site_section_advanteges (id, section_id, text, sort_order)
			VALUES ($1,$2,$3,$4)
		`, advID, id, t, order)
		if err != nil {
			return "", fmt.Errorf("insert advantage: %w", err)
		}
		order++
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("commit: %w", err)
	}

	return id, nil
}

type AdminUpdateSectionFormInput struct {
	ID             string
	Title          string
	Slug           string
	AdvantegesText string
	Advanteges     []string
	ImageFilename  string // "" => не менять картинку
}

func (s *SiteSectionsAdminService) UpdateSectionFromForm(ctx context.Context, in AdminUpdateSectionFormInput) error {
	in.ID = strings.TrimSpace(in.ID)
	in.Title = strings.TrimSpace(in.Title)
	in.Slug = strings.TrimSpace(in.Slug)
	in.AdvantegesText = strings.TrimSpace(in.AdvantegesText)
	in.ImageFilename = strings.TrimSpace(in.ImageFilename)

	if in.ID == "" {
		return fmt.Errorf("id is required")
	}
	if in.Title == "" {
		return fmt.Errorf("title is required")
	}
	if in.Slug == "" {
		in.Slug = slugify(in.Title)
	}

	tx, err := s.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	var exists bool
	if err := tx.GetContext(ctx, &exists, `SELECT EXISTS(SELECT 1 FROM site_sections WHERE id=$1)`, in.ID); err != nil {
		return fmt.Errorf("check section: %w", err)
	}
	if !exists {
		return fmt.Errorf("section not found")
	}

	if in.ImageFilename != "" {
		res, err := tx.ExecContext(ctx, `
			UPDATE site_sections
			SET title = $2,
			    label = $2,
			    slug = $3,
			    advanteges_text = $4,
			    image_url = $5
			WHERE id = $1
		`, in.ID, in.Title, in.Slug, in.AdvantegesText, in.ImageFilename)
		if err != nil {
			return fmt.Errorf("update section: %w", err)
		}
		aff, _ := res.RowsAffected()
		if aff == 0 {
			return fmt.Errorf("section not found")
		}
	} else {
		res, err := tx.ExecContext(ctx, `
			UPDATE site_sections
			SET title = $2,
			    label = $2,
			    slug = $3,
			    advanteges_text = $4
			WHERE id = $1
		`, in.ID, in.Title, in.Slug, in.AdvantegesText)
		if err != nil {
			return fmt.Errorf("update section: %w", err)
		}
		aff, _ := res.RowsAffected()
		if aff == 0 {
			return fmt.Errorf("section not found")
		}
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM site_section_advanteges WHERE section_id = $1`, in.ID)
	if err != nil {
		return fmt.Errorf("delete advantages: %w", err)
	}

	order := 1
	for _, t := range in.Advanteges {
		t = strings.TrimSpace(t)
		if t == "" {
			continue
		}
		advID := "adv-" + randID()
		_, err = tx.ExecContext(ctx, `
			INSERT INTO site_section_advanteges (id, section_id, text, sort_order)
			VALUES ($1,$2,$3,$4)
		`, advID, in.ID, t, order)
		if err != nil {
			return fmt.Errorf("insert advantage: %w", err)
		}
		order++
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}

// ==== helpers ====

func (s *SiteSectionsAdminService) ensureSectionExists(ctx context.Context, sectionID string) error {
	var ok bool
	if err := s.db.GetContext(ctx, &ok, `SELECT EXISTS(SELECT 1 FROM site_sections WHERE id=$1)`, sectionID); err != nil {
		return fmt.Errorf("check section exists: %w", err)
	}
	if !ok {
		return fmt.Errorf("section not found")
	}
	return nil
}

// Если пришло "/sections/picture/x.jpg" или "http://domain/sections/picture/x.jpg" — сохраняем "x.jpg"
func stripPicturePrefix(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}

	if strings.HasPrefix(v, "http://") || strings.HasPrefix(v, "https://") {
		if u, err := url.Parse(v); err == nil && u.Path != "" {
			v = u.Path
		}
	}

	v = strings.TrimPrefix(v, MAIN_PICTURE_URL)
	v = strings.TrimPrefix(v, GALLERY_PICTURE_URL)
	v = strings.TrimPrefix(v, CATALOG_PICTURE_URL)

	if strings.Contains(v, "/") {
		if i := strings.LastIndex(v, "/"); i >= 0 && i+1 < len(v) {
			return v[i+1:]
		}
	}
	return strings.TrimPrefix(v, "/")
}

func (s *SiteSectionsAdminService) withDomainPicture(prefix, v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}
	if strings.HasPrefix(v, "http://") || strings.HasPrefix(v, "https://") {
		return v
	}
	if strings.HasPrefix(v, "/") {
		return s.domain + v
	}
	return s.domain + prefix + v
}

func randID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b) + "-" + fmt.Sprint(time.Now().UnixNano())
}

func slugify(s string) string {
	out := slug.Make(s)
	if out == "" {
		return "section"
	}
	return out
}
