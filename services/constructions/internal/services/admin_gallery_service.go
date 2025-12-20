package services

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pachv/constructions/constructions/internal/domain/entity"
)

var (
	ErrGalleryNotFound     = errors.New("not found")
	ErrGalleryLimitReached = errors.New("categories limit reached (max 5)")
	ErrGalleryBadImage     = errors.New("photo must be an image")
)

type AdminGalleryService struct {
	db        *sqlx.DB
	uploadDir string // "./uploads/gallery"
}

func NewAdminGalleryService(db *sqlx.DB, uploadDir string) *AdminGalleryService {
	return &AdminGalleryService{db: db, uploadDir: uploadDir}
}

// ================= PUBLIC =================

// <= 5
func (s *AdminGalleryService) ListCategories(ctx context.Context) ([]entity.GalleryCategory, error) {
	const q = `
		SELECT id, title, slug, created_at
		FROM gallery_categories
		ORDER BY created_at DESC;
	`
	var out []entity.GalleryCategory
	if err := s.db.SelectContext(ctx, &out, q); err != nil {
		return nil, err
	}
	if len(out) > 5 {
		out = out[:5]
	}
	return out, nil
}

func (s *AdminGalleryService) ListPhotosByCategorySlug(ctx context.Context, slug string) ([]entity.GalleryPhoto, error) {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return nil, errors.New("slug is required")
	}

	const q = `
		SELECT
			p.id,
			c.slug AS category_slug,
			p.alt,
			p.image_path AS image,
			COALESCE(p.sort_order, 0) AS sort_order,
			p.created_at
		FROM gallery_photos p
		JOIN gallery_categories c ON c.id = p.category_id
		WHERE c.slug = $1
		ORDER BY p.sort_order ASC, p.created_at DESC;
	`
	var out []entity.GalleryPhoto
	if err := s.db.SelectContext(ctx, &out, q, slug); err != nil {
		return nil, err
	}
	return out, nil
}

// image = только имя файла (baths-1.jpg)
// возвращает полный путь для gin c.File(...)
func (s *AdminGalleryService) OpenPicturePath(image string) (string, error) {
	image = strings.TrimSpace(image)
	// запрещаем путь/директории
	if image == "" || strings.Contains(image, "/") || strings.Contains(image, `\`) || strings.Contains(image, "..") {
		return "", errors.New("invalid image name")
	}

	full := filepath.Join(s.uploadDir, image)
	if _, err := os.Stat(full); err != nil {
		if os.IsNotExist(err) {
			return "", ErrGalleryNotFound
		}
		return "", err
	}
	return full, nil
}

// ================= ADMIN: categories =================

// id категории = "gal-" + slug(title)
func (s *AdminGalleryService) CreateCategory(ctx context.Context, title string) (entity.GalleryCategory, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return entity.GalleryCategory{}, errors.New("title is required")
	}

	// limit <= 5
	{
		const qCount = `SELECT COUNT(*) FROM gallery_categories;`
		var n int
		if err := s.db.GetContext(ctx, &n, qCount); err != nil {
			return entity.GalleryCategory{}, err
		}
		if n >= 5 {
			return entity.GalleryCategory{}, ErrGalleryLimitReached
		}
	}

	slug := MakeSlug(title)
	id := "gal-" + slug

	const q = `INSERT INTO gallery_categories (id, title, slug) VALUES ($1, $2, $3);`
	if _, err := s.db.ExecContext(ctx, q, id, title, slug); err != nil {
		return entity.GalleryCategory{}, err
	}

	return s.getCategoryByID(ctx, id)
}

// title меняем, slug пересоздаём, id НЕ меняем
func (s *AdminGalleryService) UpdateCategory(ctx context.Context, id string, newTitle string) (entity.GalleryCategory, error) {
	id = strings.TrimSpace(id)
	newTitle = strings.TrimSpace(newTitle)

	if id == "" {
		return entity.GalleryCategory{}, errors.New("id is required")
	}
	if newTitle == "" {
		return entity.GalleryCategory{}, errors.New("title is required")
	}

	newSlug := MakeSlug(newTitle)

	const q = `
		UPDATE gallery_categories
		SET title = $2, slug = $3
		WHERE id = $1;
	`
	res, err := s.db.ExecContext(ctx, q, id, newTitle, newSlug)
	if err != nil {
		return entity.GalleryCategory{}, err
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return entity.GalleryCategory{}, ErrGalleryNotFound
	}

	return s.getCategoryByID(ctx, id)
}

// удаляет категорию + все фото (и файлы на диске best-effort)
func (s *AdminGalleryService) DeleteCategory(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return errors.New("id is required")
	}

	// image_path в БД = только имя файла
	const qPaths = `SELECT image_path FROM gallery_photos WHERE category_id = $1;`
	var paths []string
	if err := s.db.SelectContext(ctx, &paths, qPaths, id); err != nil {
		return err
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(ctx, `DELETE FROM gallery_photos WHERE category_id = $1;`, id); err != nil {
		return err
	}

	res, err := tx.ExecContext(ctx, `DELETE FROM gallery_categories WHERE id = $1;`, id)
	if err != nil {
		return err
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return ErrGalleryNotFound
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	// best-effort: удалить файлы
	for _, name := range paths {
		_ = os.Remove(filepath.Join(s.uploadDir, name))
	}
	return nil
}

// ================= ADMIN: photos =================

// Добавить фото в категорию (categoryID)
// fileHeader обязателен. image_path в БД кладём только имя файла.
func (s *AdminGalleryService) AddPhoto(ctx context.Context, categoryID string, alt string, sortOrderStr string, fileHeader *multipart.FileHeader) (entity.GalleryPhoto, error) {
	categoryID = strings.TrimSpace(categoryID)
	if categoryID == "" {
		return entity.GalleryPhoto{}, errors.New("category id is required")
	}
	if fileHeader == nil {
		return entity.GalleryPhoto{}, errors.New("photo is required")
	}

	// проверим категорию и получим slug
	const qCat = `SELECT slug FROM gallery_categories WHERE id = $1 LIMIT 1;`
	var catSlug string
	if err := s.db.GetContext(ctx, &catSlug, qCat, categoryID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.GalleryPhoto{}, ErrGalleryNotFound
		}
		return entity.GalleryPhoto{}, err
	}

	sortOrder := 10
	if strings.TrimSpace(sortOrderStr) != "" {
		if n, err := strconv.Atoi(strings.TrimSpace(sortOrderStr)); err == nil {
			sortOrder = n
		}
	}

	alt = strings.TrimSpace(alt)
	if alt == "" {
		alt = "Фото"
	}

	filename, err := s.saveUploadedImage(fileHeader)
	if err != nil {
		return entity.GalleryPhoto{}, err
	}

	photoID := makeID("ph")

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		_ = os.Remove(filepath.Join(s.uploadDir, filename))
		return entity.GalleryPhoto{}, err
	}
	defer func() { _ = tx.Rollback() }()

	const qIns = `
		INSERT INTO gallery_photos (id, category_id, alt, image_path, sort_order)
		VALUES ($1, $2, $3, $4, $5);
	`
	if _, err := tx.ExecContext(ctx, qIns, photoID, categoryID, alt, filename, sortOrder); err != nil {
		_ = os.Remove(filepath.Join(s.uploadDir, filename))
		return entity.GalleryPhoto{}, err
	}

	if err := tx.Commit(); err != nil {
		_ = os.Remove(filepath.Join(s.uploadDir, filename))
		return entity.GalleryPhoto{}, err
	}

	// вернуть сущность в формате твоего API (category_slug + image)
	_ = catSlug
	return s.getPhotoByID(ctx, photoID)
}

func (s *AdminGalleryService) DeletePhoto(ctx context.Context, photoID string) error {
	photoID = strings.TrimSpace(photoID)
	if photoID == "" {
		return errors.New("photo id is required")
	}

	const qGet = `SELECT image_path FROM gallery_photos WHERE id = $1 LIMIT 1;`
	var imageName string
	if err := s.db.GetContext(ctx, &imageName, qGet, photoID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrGalleryNotFound
		}
		return err
	}

	res, err := s.db.ExecContext(ctx, `DELETE FROM gallery_photos WHERE id = $1;`, photoID)
	if err != nil {
		return err
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return ErrGalleryNotFound
	}

	_ = os.Remove(filepath.Join(s.uploadDir, imageName))
	return nil
}

// ================= internal =================

func (s *AdminGalleryService) getCategoryByID(ctx context.Context, id string) (entity.GalleryCategory, error) {
	const q = `
		SELECT id, title, slug, created_at
		FROM gallery_categories
		WHERE id = $1
		LIMIT 1;
	`
	var out entity.GalleryCategory
	if err := s.db.GetContext(ctx, &out, q, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.GalleryCategory{}, ErrGalleryNotFound
		}
		return entity.GalleryCategory{}, err
	}
	return out, nil
}

func (s *AdminGalleryService) getPhotoByID(ctx context.Context, id string) (entity.GalleryPhoto, error) {
	const q = `
		SELECT
			p.id,
			c.slug AS category_slug,
			p.alt,
			p.image_path AS image,
			COALESCE(p.sort_order, 0) AS sort_order,
			p.created_at
		FROM gallery_photos p
		JOIN gallery_categories c ON c.id = p.category_id
		WHERE p.id = $1
		LIMIT 1;
	`
	var out entity.GalleryPhoto
	if err := s.db.GetContext(ctx, &out, q, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.GalleryPhoto{}, ErrGalleryNotFound
		}
		return entity.GalleryPhoto{}, err
	}
	return out, nil
}

func (s *AdminGalleryService) saveUploadedImage(fh *multipart.FileHeader) (string, error) {
	src, err := fh.Open()
	if err != nil {
		return "", errors.New("failed to open uploaded file")
	}
	defer src.Close()

	// sniff
	head := make([]byte, 512)
	n, _ := io.ReadFull(src, head)
	head = head[:n]
	ct := http.DetectContentType(head)
	if !strings.HasPrefix(ct, "image/") {
		return "", ErrGalleryBadImage
	}

	// reopen (без seek)
	_ = src.Close()
	src, err = fh.Open()
	if err != nil {
		return "", errors.New("failed to reopen uploaded file")
	}
	defer src.Close()

	if err := os.MkdirAll(s.uploadDir, 0o755); err != nil {
		return "", errors.New("failed to create upload dir")
	}

	ext := strings.ToLower(filepath.Ext(fh.Filename))
	if ext == "" {
		switch ct {
		case "image/jpeg":
			ext = ".jpg"
		case "image/png":
			ext = ".png"
		case "image/webp":
			ext = ".webp"
		default:
			ext = ".img"
		}
	}

	// ВАЖНО: filename = только имя файла (в БД тоже только имя)
	filename := makeID("img") + ext
	dstPath := filepath.Join(s.uploadDir, filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		return "", errors.New("failed to save file")
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		_ = os.Remove(dstPath)
		return "", errors.New("failed to write file")
	}

	return filename, nil
}

func makeID(prefix string) string {
	b := make([]byte, 10)
	_, _ = rand.Read(b)
	return prefix + "-" + hex.EncodeToString(b)
}
