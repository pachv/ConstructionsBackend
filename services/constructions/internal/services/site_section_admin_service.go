package services

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pachv/constructions/constructions/internal/domain/entity"
)

const (
	MAIN_PICTURE_URL    = "/sections/picture/"
	GALLERY_PICTURE_URL = "/sections/gallery/picture/"
	CATALOG_PICTURE_URL = "/catalog/picture/"
)

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

type AdminSectionsPage struct {
	Items      []*entity.SiteSectionSummary `json:"items"`
	Page       int                          `json:"page"`
	PageAmount int                          `json:"pageAmount"`
	Total      int                          `json:"total"`
}

type AdminUpsertSectionInput struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	Label          string `json:"label"`
	Slug           string `json:"slug"`
	Image          string `json:"image"` // может прийти как "/sections/picture/x.jpg" или "x.jpg" или "https://..."
	AdvantegesText string `json:"advantegesText"`
	HasGallery     bool   `json:"hasGallery"`
	HasCatalog     bool   `json:"hasCatalog"`
}

func (s *SiteSectionsAdminService) GetSectionsSummary(page int, search, orderBy string) ([]*entity.SiteSectionSummary, int, int, error) {
	const pageSize = 10
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	allowedOrderBy := map[string]bool{
		"id":         true,
		"title":      true,
		"label":      true,
		"slug":       true,
		"created_at": true,
	}
	if !allowedOrderBy[orderBy] {
		orderBy = "created_at"
	}

	search = strings.TrimSpace(search)

	var args []any
	where := ""
	if search != "" {
		where = `
			WHERE
				s.id ILIKE $1 OR
				s.title ILIKE $1 OR
				s.label ILIKE $1 OR
				s.slug ILIKE $1
		`
		args = append(args, "%"+search+"%")
	}

	// total
	countQuery := `SELECT COUNT(*) FROM site_sections s ` + where
	var total int
	if err := s.db.Get(&total, countQuery, args...); err != nil {
		return nil, 0, 0, fmt.Errorf("count site sections: %w", err)
	}

	if total == 0 {
		return []*entity.SiteSectionSummary{}, 0, 0, nil
	}

	pageAmount := (total + pageSize - 1) / pageSize

	// IMPORTANT: возвращаем image_url (без AS image), чтобы sqlx маппился в db:"image_url"
	query := `
		SELECT
			s.id,
			s.title,
			s.label,
			s.slug,
			s.image_url,
			s.has_gallery,
			s.has_catalog
		FROM site_sections s
	` + where + fmt.Sprintf(" ORDER BY s.%s LIMIT %d OFFSET %d", orderBy, pageSize, offset)

	var items []*entity.SiteSectionSummary
	if err := s.db.Select(&items, query, args...); err != nil {
		return nil, 0, 0, fmt.Errorf("get site sections: %w", err)
	}

	// paths like repository example
	for i := range items {
		it := items[i]
		if it == nil {
			continue
		}

		if strings.TrimSpace(it.Label) == "" {
			it.Label = it.Title
		}

		// ВАЖНО: это поле должно быть db:"image_url" в entity.SiteSectionSummary
		// Обычно в entity оно называется Image, но db tag = image_url
		filename := strings.TrimSpace(it.Image)
		if filename != "" {
			it.Image = s.withDomainPicture(MAIN_PICTURE_URL, filename)
		}
	}

	return items, pageAmount, total, nil
}

func (s *SiteSectionsAdminService) GetSectionBySlug(slug string) (*entity.SiteSection, error) {
	slug = strings.TrimSpace(slug)
	if slug == "" {
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

	q := `
		SELECT
			s.id,
			s.title,
			s.label,
			s.slug,
			s.image_url,
			s.advanteges_text,
			s.has_gallery,
			s.has_catalog
		FROM site_sections s
		WHERE s.slug = $1
		LIMIT 1
	`
	if err := s.db.Get(&base, q, slug); err != nil {
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
		// Catalog: nil (пока)
	}

	if strings.TrimSpace(out.Label) == "" {
		out.Label = out.Title
	}
	if out.AdvantegesText == "" {
		out.AdvantegesText = ""
	}

	// main image full url
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

	// gallery
	if out.HasGallery {
		var gallery []entity.SiteSectionGallery
		if err := s.db.Select(&gallery, `
			SELECT
				g.id,
				g.section_id,
				g.name,
				g.url,
				g.sort_order
			FROM site_section_gallery g
			WHERE g.section_id = $1
			ORDER BY g.sort_order
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
	}

	return out, nil
}

func (s *SiteSectionsAdminService) CreateSection(in AdminUpsertSectionInput) error {
	in.ID = strings.TrimSpace(in.ID)
	in.Title = strings.TrimSpace(in.Title)
	in.Slug = strings.TrimSpace(in.Slug)
	in.Label = strings.TrimSpace(in.Label)
	in.Image = strings.TrimSpace(in.Image)
	in.AdvantegesText = strings.TrimSpace(in.AdvantegesText)

	if in.ID == "" || in.Title == "" || in.Slug == "" {
		return fmt.Errorf("id, title, slug are required")
	}

	imageToStore := stripPicturePrefix(in.Image)

	q := `
		INSERT INTO site_sections (
			id,
			title,
			label,
			slug,
			image_url,
			advanteges_text,
			has_gallery,
			has_catalog
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`
	if _, err := s.db.Exec(q,
		in.ID,
		in.Title,
		in.Label,
		in.Slug,
		imageToStore,
		in.AdvantegesText,
		in.HasGallery,
		in.HasCatalog,
	); err != nil {
		return fmt.Errorf("create section: %w", err)
	}
	return nil
}

func (s *SiteSectionsAdminService) UpdateSection(id string, in AdminUpsertSectionInput) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("id is required")
	}

	in.Title = strings.TrimSpace(in.Title)
	in.Slug = strings.TrimSpace(in.Slug)
	in.Label = strings.TrimSpace(in.Label)
	in.Image = strings.TrimSpace(in.Image)
	in.AdvantegesText = strings.TrimSpace(in.AdvantegesText)

	if in.Title == "" || in.Slug == "" {
		return fmt.Errorf("title and slug are required")
	}

	imageToStore := stripPicturePrefix(in.Image)

	q := `
		UPDATE site_sections
		SET
			title = $2,
			label = $3,
			slug = $4,
			image_url = $5,
			advanteges_text = $6,
			has_gallery = $7,
			has_catalog = $8
		WHERE id = $1
	`

	res, err := s.db.Exec(q,
		id,
		in.Title,
		in.Label,
		in.Slug,
		imageToStore,
		in.AdvantegesText,
		in.HasGallery,
		in.HasCatalog,
	)
	if err != nil {
		return fmt.Errorf("update section: %w", err)
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return fmt.Errorf("section not found")
	}
	return nil
}

/*
=========================
 helpers
=========================
*/

// Делает "domain + prefix + filename", но если уже абсолютный/с prefix — не портит
func (s *SiteSectionsAdminService) withDomainPicture(prefix, v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}

	if strings.HasPrefix(v, "http://") || strings.HasPrefix(v, "https://") {
		return v
	}

	// "/sections/picture/x.jpg" -> домен + этот путь
	if strings.HasPrefix(v, "/") {
		return s.domain + v
	}

	// "x.jpg" -> домен + prefix + x.jpg
	return s.domain + prefix + v
}

// Если пришло "/sections/picture/x.jpg" или "http://domain/sections/picture/x.jpg" — сохраняем "x.jpg"
func stripPicturePrefix(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}

	// абсолютный url -> распарсим и возьмем path
	if strings.HasPrefix(v, "http://") || strings.HasPrefix(v, "https://") {
		if u, err := url.Parse(v); err == nil && u.Path != "" {
			v = u.Path
		}
	}

	// убираем известные префиксы
	v = strings.TrimPrefix(v, MAIN_PICTURE_URL)
	v = strings.TrimPrefix(v, GALLERY_PICTURE_URL)
	v = strings.TrimPrefix(v, CATALOG_PICTURE_URL)

	// если все равно путь вида "/something/x.jpg" -> basename
	if strings.Contains(v, "/") {
		if i := strings.LastIndex(v, "/"); i >= 0 && i+1 < len(v) {
			return v[i+1:]
		}
	}

	return strings.TrimPrefix(v, "/")
}
