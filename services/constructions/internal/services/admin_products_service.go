package services

import (
	"database/sql"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const PicturePrefix = "/api/v1/products/picture/"

type AdminProductService struct {
	db     *sqlx.DB
	domain string // "http://localhost:80"
}

func NewAdminProductService(db *sqlx.DB, domain string) *AdminProductService {
	return &AdminProductService{db: db, domain: strings.TrimRight(domain, "/")}
}

// =========================
// DTOs (API shapes)
// =========================

type CategoryDTO struct {
	ID        string    `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Slug      string    `json:"slug" db:"slug"`
	ImagePath string    `json:"imagePath" db:"image_path"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

type SectionDTO struct {
	ID                 string    `json:"id" db:"id"`
	Title              string    `json:"title" db:"title"`
	Slug               string    `json:"slug" db:"slug"`
	ParentCategorySlug string    `json:"parentCategorySlug" db:"parent_category_slug"`
	ImagePath          string    `json:"imagePath" db:"image_path"`
	CreatedAt          time.Time `json:"createdAt" db:"created_at"`
}

type ProductDTO struct {
	ID           string    `json:"id" db:"id"`
	Title        string    `json:"title" db:"title"`
	Slug         string    `json:"slug" db:"slug"`
	CategorySlug string    `json:"categorySlug" db:"category_slug"`
	SectionSlug  string    `json:"sectionSlug" db:"section_slug"`
	Brand        string    `json:"brand" db:"brand"`
	Type         string    `json:"type" db:"type"`
	Price        int       `json:"price" db:"price"`
	OldPrice     *int      `json:"oldPrice" db:"old_price"`
	InStock      bool      `json:"inStock" db:"in_stock"`
	SalePercent  int       `json:"salePercent" db:"sale_percent"`
	ImagePath    string    `json:"imagePath" db:"image_path"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	Badges       []string  `json:"badges"`
}

// =========================
// Requests
// =========================

type CreateCategoryReq struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Slug      string `json:"slug"`
	ImagePath string `json:"imagePath"` // filename or url
}

type UpdateCategoryReq struct {
	Title     string `json:"title"`
	Slug      string `json:"slug"`
	ImagePath string `json:"imagePath"`
}

type CreateSectionReq struct {
	ID                 string `json:"id"`
	Title              string `json:"title"`
	Slug               string `json:"slug"`
	ParentCategorySlug string `json:"parentCategorySlug"`
	ImagePath          string `json:"imagePath"`
}

type UpdateSectionReq struct {
	Title              string `json:"title"`
	Slug               string `json:"slug"`
	ParentCategorySlug string `json:"parentCategorySlug"`
	ImagePath          string `json:"imagePath"` // если пусто - не меняем
}

type CreateProductReq struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Slug         string   `json:"slug"`
	CategorySlug string   `json:"categorySlug"`
	SectionSlug  string   `json:"sectionSlug"`
	Brand        string   `json:"brand"`
	Type         string   `json:"type"`
	Price        int      `json:"price"`   // base price
	InStock      bool     `json:"inStock"` // true/false
	ImagePath    string   `json:"imagePath"`
	Badges       []string `json:"badges"`
	Discount     int      `json:"discount"` // число процентов (0-99)
}

// ✅ UpdateProductReq теперь с int для discount
type UpdateProductReq struct {
	Title        string   `json:"title"`
	Slug         string   `json:"slug"`
	CategorySlug string   `json:"categorySlug"`
	SectionSlug  string   `json:"sectionSlug"`
	Brand        string   `json:"brand"`
	Type         string   `json:"type"`
	Price        int      `json:"price"`   // base price if discount set
	InStock      bool     `json:"inStock"` // true/false
	ImagePath    string   `json:"imagePath"`
	Badges       []string `json:"badges"`
	Discount     int      `json:"discount"` // число процентов (0-99) или 0 чтобы убрать скидку
}

// =========================
// Helpers
// =========================

func (s *AdminProductService) fullImageURL(filename string) string {
	filename = strings.TrimSpace(filename)
	if filename == "" {
		return ""
	}
	if strings.HasPrefix(filename, "http://") || strings.HasPrefix(filename, "https://") {
		return filename
	}
	return filename
}

func filenameOnly(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}
	if i := strings.LastIndex(v, "/"); i >= 0 {
		return v[i+1:]
	}
	return v
}

var reDiscount = regexp.MustCompile(`^sale_(\d{1,2})$`)

func parseDiscountPercent(discount string) (int, bool) {
	discount = strings.TrimSpace(discount)
	if discount == "" {
		return 0, false
	}
	m := reDiscount.FindStringSubmatch(discount)
	if len(m) != 2 {
		return 0, false
	}
	p, err := strconv.Atoi(m[1])
	if err != nil || p <= 0 || p >= 100 {
		return 0, false
	}
	return p, true
}

func applyDiscount(basePrice int, percent int) int {
	x := float64(basePrice) * float64(100-percent) / 100.0
	return int(math.Round(x))
}

// =========================
// Categories
// =========================

func (s *AdminProductService) CreateCategory(title, slug, imageFilename, id string) error {
	title = strings.TrimSpace(title)
	slug = strings.TrimSpace(slug)
	if title == "" || slug == "" {
		return fmt.Errorf("title/slug required")
	}
	imageFilename = filenameOnly(imageFilename)

	// ✅ Если id пустой — генерим
	id = strings.TrimSpace(id)
	if id == "" {
		id = uuid.NewString()
	}

	_, err := s.db.Exec(`
		INSERT INTO catalog_categories (id, title, slug, image_path, created_at)
		VALUES ($1, $2, $3, $4, now())
	`, id, title, slug, imageFilename)
	return err
}

func (s *AdminProductService) UpdateCategory(id, title, slug, imageFilename string) error {
	id = strings.TrimSpace(id)
	title = strings.TrimSpace(title)
	slug = strings.TrimSpace(slug)
	if id == "" || title == "" || slug == "" {
		return fmt.Errorf("id/title/slug required")
	}

	imageFilename = strings.TrimSpace(imageFilename)

	var res sql.Result
	var err error

	// ✅ Если imagePath передан — обновляем, иначе не трогаем
	if imageFilename != "" {
		imageFilename = filenameOnly(imageFilename)
		res, err = s.db.Exec(`
			UPDATE catalog_categories
			SET title=$1, slug=$2, image_path=$3
			WHERE id=$4
		`, title, slug, imageFilename, id)
	} else {
		res, err = s.db.Exec(`
			UPDATE catalog_categories
			SET title=$1, slug=$2
			WHERE id=$3
		`, title, slug, id)
	}

	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *AdminProductService) DeleteCategory(id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("id required")
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	// 1. Получаем все секции этой категории
	var sectionIDs []string
	if err := tx.Select(&sectionIDs, `
		SELECT section_id 
		FROM catalog_category_sections 
		WHERE category_id = $1
	`, id); err != nil {
		return err
	}

	// 2. Удаляем badges всех товаров этих секций
	for _, secID := range sectionIDs {
		_, _ = tx.Exec(`
			DELETE FROM product_badge_links 
			WHERE product_id IN (
				SELECT id FROM catalog_products 
				WHERE section_slug IN (
					SELECT slug FROM catalog_sections WHERE id = $1
				)
			)
		`, secID)
	}

	// 3. Удаляем товары всех секций
	_, _ = tx.Exec(`
		DELETE FROM catalog_products 
		WHERE category_slug IN (
			SELECT slug FROM catalog_categories WHERE id = $1
		)
	`, id)

	// 4. Удаляем связь категория-секция
	_, _ = tx.Exec(`DELETE FROM catalog_category_sections WHERE category_id = $1`, id)

	// 5. Удаляем секции
	for _, secID := range sectionIDs {
		_, _ = tx.Exec(`DELETE FROM catalog_sections WHERE id = $1`, secID)
	}

	// 6. Удаляем саму категорию
	if _, err := tx.Exec(`DELETE FROM catalog_categories WHERE id = $1`, id); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *AdminProductService) GetAllCategories() ([]CategoryDTO, error) {
	var rows []CategoryDTO
	if err := s.db.Select(&rows, `
		SELECT id, title, slug, image_path, created_at
		FROM catalog_categories
		ORDER BY created_at DESC
	`); err != nil {
		return nil, err
	}
	for i := range rows {
		rows[i].ImagePath = s.fullImageURL(rows[i].ImagePath)
	}
	return rows, nil
}

func (s *AdminProductService) GetCategoryByTitle(title string) (*CategoryDTO, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, fmt.Errorf("title required")
	}
	var row CategoryDTO
	if err := s.db.Get(&row, `
		SELECT id, title, slug, image_path, created_at
		FROM catalog_categories
		WHERE title ILIKE $1
		LIMIT 1
	`, title); err != nil {
		return nil, err
	}
	row.ImagePath = s.fullImageURL(row.ImagePath)
	return &row, nil
}

// ✅ Сохранение изображения категории
func (s *AdminProductService) SaveCategoryImage(r io.Reader, originalName string) (string, error) {
	originalName = strings.TrimSpace(originalName)
	if originalName == "" {
		return "", fmt.Errorf("empty filename")
	}

	ext := filepath.Ext(originalName)
	if ext == "" {
		ext = ".bin"
	}

	filename := uuid.NewString() + ext
	baseDir := "/app/uploads/categories"
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return "", err
	}

	fullPath := filepath.Join(baseDir, filename)

	f, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := io.Copy(f, r); err != nil {
		return "", err
	}

	return "uploads/categories/" + filename, nil
}

// =========================
// Sections
// =========================

func (s *AdminProductService) CreateSection(id, title, slug, parentCategorySlug, imageFilename string) error {
	title = strings.TrimSpace(title)
	slug = strings.TrimSpace(slug)
	parentCategorySlug = strings.TrimSpace(parentCategorySlug)
	if title == "" || slug == "" || parentCategorySlug == "" {
		return fmt.Errorf("title/slug/parentCategorySlug required")
	}

	// ✅ Если id пустой — генерим
	id = strings.TrimSpace(id)
	if id == "" {
		id = uuid.NewString()
	}

	imageFilename = filenameOnly(imageFilename)

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	_, err = tx.Exec(`
		INSERT INTO catalog_sections (id, title, slug, image_path, created_at)
		VALUES ($1, $2, $3, $4, now())
	`, id, title, slug, imageFilename)
	if err != nil {
		return err
	}

	var catID string
	if err := tx.Get(&catID, `SELECT id FROM catalog_categories WHERE slug=$1 LIMIT 1`, parentCategorySlug); err != nil {
		return fmt.Errorf("category not found by slug=%s: %w", parentCategorySlug, err)
	}

	_, err = tx.Exec(`
		INSERT INTO catalog_category_sections (category_id, section_id, created_at)
		VALUES ($1, $2, now())
	`, catID, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *AdminProductService) UpdateSection(id, title, slug, parentCategorySlug, imageFilename string) error {
	id = strings.TrimSpace(id)
	title = strings.TrimSpace(title)
	slug = strings.TrimSpace(slug)
	parentCategorySlug = strings.TrimSpace(parentCategorySlug)
	if id == "" || title == "" || slug == "" || parentCategorySlug == "" {
		return fmt.Errorf("id/title/slug/parentCategorySlug required")
	}

	imageFilename = strings.TrimSpace(imageFilename)

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	// ✅ Если imagePath передан — обновляем, иначе не трогаем
	if imageFilename != "" {
		imageFilename = filenameOnly(imageFilename)
		if _, err := tx.Exec(`UPDATE catalog_sections SET title=$1, slug=$2, image_path=$3 WHERE id=$4`,
			title, slug, imageFilename, id); err != nil {
			return err
		}
	} else {
		if _, err := tx.Exec(`UPDATE catalog_sections SET title=$1, slug=$2 WHERE id=$3`,
			title, slug, id); err != nil {
			return err
		}
	}

	var catID string
	if err := tx.Get(&catID, `SELECT id FROM catalog_categories WHERE slug=$1 LIMIT 1`, parentCategorySlug); err != nil {
		return fmt.Errorf("category not found by slug=%s: %w", parentCategorySlug, err)
	}

	if _, err := tx.Exec(`DELETE FROM catalog_category_sections WHERE section_id=$1`, id); err != nil {
		return err
	}
	if _, err := tx.Exec(`
		INSERT INTO catalog_category_sections (category_id, section_id, created_at)
		VALUES ($1, $2, now())
	`, catID, id); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *AdminProductService) DeleteSection(id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("id required")
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	// 1. Получаем slug секции
	var sectionSlug string
	if err := tx.Get(&sectionSlug, `SELECT slug FROM catalog_sections WHERE id = $1`, id); err != nil {
		return err
	}

	// 2. Удаляем badges всех товаров этой секции
	_, _ = tx.Exec(`
		DELETE FROM product_badge_links 
		WHERE product_id IN (
			SELECT id FROM catalog_products WHERE section_slug = $1
		)
	`, sectionSlug)

	// 3. Удаляем товары этой секции
	_, _ = tx.Exec(`DELETE FROM catalog_products WHERE section_slug = $1`, sectionSlug)

	// 4. Удаляем связь категория-секция
	_, _ = tx.Exec(`DELETE FROM catalog_category_sections WHERE section_id = $1`, id)

	// 5. Удаляем саму секцию
	if _, err := tx.Exec(`DELETE FROM catalog_sections WHERE id = $1`, id); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *AdminProductService) GetAllSections() ([]SectionDTO, error) {
	var rows []SectionDTO
	err := s.db.Select(&rows, `
		SELECT
			s.id, s.title, s.slug, s.image_path, s.created_at,
			c.slug AS parent_category_slug
		FROM catalog_sections s
		LEFT JOIN catalog_category_sections cs ON cs.section_id = s.id
		LEFT JOIN catalog_categories c ON c.id = cs.category_id
		ORDER BY s.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	for i := range rows {
		rows[i].ImagePath = s.fullImageURL(rows[i].ImagePath)
	}
	return rows, nil
}

func (s *AdminProductService) GetSectionByTitle(title string) (*SectionDTO, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, fmt.Errorf("title required")
	}
	var row SectionDTO
	err := s.db.Get(&row, `
		SELECT
			s.id, s.title, s.slug, s.image_path, s.created_at,
			c.slug AS parent_category_slug
		FROM catalog_sections s
		LEFT JOIN catalog_category_sections cs ON cs.section_id = s.id
		LEFT JOIN catalog_categories c ON c.id = cs.category_id
		WHERE s.title ILIKE $1
		LIMIT 1
	`, title)
	if err != nil {
		return nil, err
	}
	row.ImagePath = s.fullImageURL(row.ImagePath)
	return &row, nil
}

func (s *AdminProductService) GetCategoryOfSection(sectionID string) (*CategoryDTO, error) {
	sectionID = strings.TrimSpace(sectionID)
	if sectionID == "" {
		return nil, fmt.Errorf("section id required")
	}
	var row CategoryDTO
	err := s.db.Get(&row, `
		SELECT c.id, c.title, c.slug, c.image_path, c.created_at
		FROM catalog_categories c
		JOIN catalog_category_sections cs ON cs.category_id = c.id
		WHERE cs.section_id = $1
		LIMIT 1
	`, sectionID)
	if err != nil {
		return nil, err
	}
	row.ImagePath = s.fullImageURL(row.ImagePath)
	return &row, nil
}

// ✅ Сохранение изображения секции
func (s *AdminProductService) SaveSectionImage(r io.Reader, originalName string) (string, error) {
	originalName = strings.TrimSpace(originalName)
	if originalName == "" {
		return "", fmt.Errorf("empty filename")
	}

	ext := filepath.Ext(originalName)
	if ext == "" {
		ext = ".bin"
	}

	filename := uuid.NewString() + ext
	baseDir := "/app/uploads/sections"
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return "", err
	}

	fullPath := filepath.Join(baseDir, filename)

	f, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := io.Copy(f, r); err != nil {
		return "", err
	}

	return "uploads/sections/" + filename, nil
}

// =========================
// Products
// =========================

func (s *AdminProductService) GetAllProducts() ([]ProductDTO, error) {
	var items []ProductDTO
	if err := s.db.Select(&items, `
		SELECT
			id, title, slug, category_slug, section_slug,
			brand, type, price, old_price, in_stock, sale_percent,
			image_path, created_at
		FROM catalog_products
		ORDER BY created_at DESC
	`); err != nil {
		return nil, err
	}

	type badgeRow struct {
		ProductID string `db:"product_id"`
		Code      string `db:"code"`
	}
	var br []badgeRow
	_ = s.db.Select(&br, `
		SELECT l.product_id, b.code
		FROM product_badge_links l
		JOIN product_badges b ON b.id = l.badge_id
	`)

	bmap := map[string][]string{}
	for _, r := range br {
		bmap[r.ProductID] = append(bmap[r.ProductID], r.Code)
	}

	for i := range items {
		items[i].Badges = bmap[items[i].ID]
		items[i].ImagePath = s.fullImageURL(items[i].ImagePath)
	}
	return items, nil
}

func (s *AdminProductService) GetProduct(id string) (*ProductDTO, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("id required")
	}
	var p ProductDTO
	if err := s.db.Get(&p, `
		SELECT
			id, title, slug, category_slug, section_slug,
			brand, type, price, old_price, in_stock, sale_percent,
			image_path, created_at
		FROM catalog_products
		WHERE id=$1
	`, id); err != nil {
		return nil, err
	}

	var badges []string
	_ = s.db.Select(&badges, `
		SELECT b.code
		FROM product_badges b
		JOIN product_badge_links l ON l.badge_id = b.id
		WHERE l.product_id = $1
		ORDER BY b.code ASC
	`, id)

	p.Badges = badges
	p.ImagePath = s.fullImageURL(p.ImagePath)
	return &p, nil
}

func (s *AdminProductService) SaveProductImage(r io.Reader, originalName string) (string, error) {
	originalName = strings.TrimSpace(originalName)
	if originalName == "" {
		return "", fmt.Errorf("empty filename")
	}

	// гарантируем безопасное имя файла
	ext := filepath.Ext(originalName)
	if ext == "" {
		ext = ".bin"
	}

	// уникальное имя
	filename := uuid.NewString() + ext

	// абсолютный путь на диске
	baseDir := "/app/uploads/products"
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return "", err
	}

	fullPath := filepath.Join(baseDir, filename)

	f, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := io.Copy(f, r); err != nil {
		return "", err
	}

	// ❗ В БД лучше хранить ОТНОСИТЕЛЬНЫЙ путь
	// чтобы не хардкодить /app
	return "uploads/products/" + filename, nil
}

func (s *AdminProductService) CreateProduct(req CreateProductReq) error {
	// ✅ ID: если пустой — генерим
	req.ID = strings.TrimSpace(req.ID)
	if req.ID == "" {
		req.ID = uuid.NewString()
	}

	req.Title = strings.TrimSpace(req.Title)
	req.CategorySlug = strings.TrimSpace(req.CategorySlug)
	req.SectionSlug = strings.TrimSpace(req.SectionSlug)
	req.Brand = strings.TrimSpace(req.Brand)
	req.Type = strings.TrimSpace(req.Type)

	if req.Title == "" || req.CategorySlug == "" || req.SectionSlug == "" || req.Type == "" {
		return fmt.Errorf("title/categorySlug/sectionSlug/type required")
	}
	if req.Price < 0 {
		return fmt.Errorf("price must be >= 0")
	}

	// ✅ slug ВСЕГДА из title (игнорируем req.Slug)
	req.Slug = slugify(req.Title)
	if req.Slug == "" {
		return fmt.Errorf("cannot build slug from title")
	}

	// discount — просто число процентов
	if req.Discount < 0 || req.Discount > 99 {
		return fmt.Errorf("discount must be in range 0..99")
	}

	imageFilename := filenameOnly(req.ImagePath)

	price := req.Price
	oldPricePtr := req.Price
	salePercent := 0

	if req.Discount > 0 {
		price = applyDiscountPercent(req.Price, req.Discount)

		// оставляю как у тебя было: отрицательное
		salePercent = -req.Discount
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	_, err = tx.Exec(`
		INSERT INTO catalog_products (
			id, title, slug, category_slug, section_slug,
			brand, type, price, old_price, in_stock, sale_percent,
			image_path, created_at
		) VALUES (
			$1,$2,$3,$4,$5,
			$6,$7,$8,$9,$10,$11,
			$12, now()
		)
	`, req.ID, req.Title, req.Slug, req.CategorySlug, req.SectionSlug,
		req.Brand, req.Type, price, oldPricePtr, req.InStock, salePercent,
		imageFilename)
	if err != nil {
		return err
	}

	if err := s.syncBadgesTx(tx, req.ID, req.Badges); err != nil {
		return err
	}

	return tx.Commit()
}

func applyDiscountPercent(price int, discount int) int {
	// округление до ближайшего целого
	return (price*(100-discount) + 50) / 100
}

// ✅ UpdateProduct с int discount
func (s *AdminProductService) UpdateProduct(id string, req UpdateProductReq) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("id required")
	}

	req.Title = strings.TrimSpace(req.Title)
	req.Slug = strings.TrimSpace(req.Slug)
	req.CategorySlug = strings.TrimSpace(req.CategorySlug)
	req.SectionSlug = strings.TrimSpace(req.SectionSlug)
	req.Brand = strings.TrimSpace(req.Brand)
	req.Type = strings.TrimSpace(req.Type)

	if req.Title == "" || req.Slug == "" || req.CategorySlug == "" || req.SectionSlug == "" || req.Type == "" {
		return fmt.Errorf("title/slug/categorySlug/sectionSlug/type required")
	}
	if req.Price < 0 {
		return fmt.Errorf("price must be >= 0")
	}
	if req.Discount < 0 || req.Discount > 99 {
		return fmt.Errorf("discount must be in range 0..99")
	}

	imageFilename := strings.TrimSpace(req.ImagePath)

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	var cur struct {
		Price    int `db:"price"`
		OldPrice int `db:"old_price"`
	}
	if err := tx.Get(&cur, `SELECT price, old_price FROM catalog_products WHERE id=$1`, id); err != nil {
		return err
	}

	var (
		price       int
		oldPricePtr int
		salePercent int
	)

	// ✅ Работаем с числовым discount
	if req.Discount > 0 {

		price = applyDiscountPercent(cur.OldPrice, req.Discount)
		salePercent = -req.Discount
	} else {
		price = cur.OldPrice
	}

	// ✅ Если imagePath передан — обновляем, иначе не трогаем
	if imageFilename != "" {
		imageFilename = filenameOnly(imageFilename)
		_, err = tx.Exec(`
			UPDATE catalog_products
			SET title=$1, slug=$2, category_slug=$3, section_slug=$4,
				brand=$5, type=$6, price=$7, old_price=$8,
				in_stock=$9, sale_percent=$10, image_path=$11
			WHERE id=$12
		`, req.Title, req.Slug, req.CategorySlug, req.SectionSlug,
			req.Brand, req.Type, price, oldPricePtr,
			req.InStock, salePercent, imageFilename, id)
	} else {
		_, err = tx.Exec(`
			UPDATE catalog_products
			SET title=$1, slug=$2, category_slug=$3, section_slug=$4,
				brand=$5, type=$6, price=$7, old_price=$8,
				in_stock=$9, sale_percent=$10
			WHERE id=$11
		`, req.Title, req.Slug, req.CategorySlug, req.SectionSlug,
			req.Brand, req.Type, price, oldPricePtr,
			req.InStock, salePercent, id)
	}
	if err != nil {
		return err
	}

	if err := s.syncBadgesTx(tx, id, req.Badges); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *AdminProductService) DeleteProduct(id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("id required")
	}
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	_, _ = tx.Exec(`DELETE FROM product_badge_links WHERE product_id=$1`, id)
	if _, err := tx.Exec(`DELETE FROM catalog_products WHERE id=$1`, id); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *AdminProductService) syncBadgesTx(tx *sqlx.Tx, productID string, badgeCodes []string) error {
	if _, err := tx.Exec(`DELETE FROM product_badge_links WHERE product_id=$1`, productID); err != nil {
		return err
	}
	uniq := map[string]struct{}{}
	for _, c := range badgeCodes {
		c = strings.TrimSpace(c)
		if c == "" {
			continue
		}
		uniq[c] = struct{}{}
	}
	for code := range uniq {
		if _, err := tx.Exec(`
			INSERT INTO product_badge_links (product_id, badge_id, created_at)
			SELECT $1, b.id, now()
			FROM product_badges b
			WHERE b.code = $2
		`, productID, code); err != nil {
			return err
		}
	}
	return nil
}
