package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pachv/constructions/constructions/internal/domain/entity"
)

var ErrAdminFontsDBNil = errors.New("admin fonts service: db is nil")

type AdminFontsService struct {
	db     *sqlx.DB
	logger *slog.Logger

	fontsDir string // "fonts"
	buildDir string // "build"
}

func NewAdminFontsService(db *sqlx.DB, logger *slog.Logger) *AdminFontsService {
	if logger == nil {
		logger = slog.Default()
	}
	return &AdminFontsService{
		db:       db,
		logger:   logger.With("component", "admin_fonts_service"),
		fontsDir: "fonts",
		buildDir: "build",
	}
}

func (s *AdminFontsService) checkDB(method string) error {
	if s == nil || s.db == nil {
		if s != nil && s.logger != nil {
			s.logger.Error("db is nil", "method", method)
		}
		return ErrAdminFontsDBNil
	}
	return nil
}

// -------- Public API --------

func (s *AdminFontsService) ListFonts(ctx context.Context) ([]entity.AdminFont, error) {
	if err := s.checkDB("ListFonts"); err != nil {
		return nil, err
	}

	items := make([]entity.AdminFont, 0)
	if err := s.db.SelectContext(ctx, &items, `
		SELECT id, name, file_path, selected, created_at, updated_at
		FROM admin_fonts
		ORDER BY created_at ASC
	`); err != nil {
		s.logger.Error("ListFonts: query failed", "err", err)
		return nil, err
	}

	return items, nil
}

func (s *AdminFontsService) DeleteFont(ctx context.Context, id string) error {
	if err := s.checkDB("DeleteFont"); err != nil {
		return err
	}
	id = strings.TrimSpace(id)
	if id == "" {
		return errors.New("id is required")
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	var filePath string
	var wasSelected bool
	err = tx.GetContext(ctx, &struct {
		FilePath    *string `db:"file_path"`
		WasSelected bool    `db:"selected"`
	}{}, `SELECT file_path, selected FROM admin_fonts WHERE id = $1`, id)
	if err != nil {
		s.logger.Error("DeleteFont: select failed", "err", err, "id", id)
		return err
	}

	// реальный get (не через анонимку) — чтобы не плясать с nil
	if err := tx.GetContext(ctx, &filePath, `SELECT file_path FROM admin_fonts WHERE id = $1`, id); err != nil {
		return err
	}
	if err := tx.GetContext(ctx, &wasSelected, `SELECT selected FROM admin_fonts WHERE id = $1`, id); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM admin_fonts WHERE id = $1`, id); err != nil {
		s.logger.Error("DeleteFont: delete failed", "err", err, "id", id)
		return err
	}

	// Если удалили selected — выбираем любой оставшийся и делаем selected
	if wasSelected {
		var nextID string
		var nextPath string
		err := tx.GetContext(ctx, &struct {
			ID   *string `db:"id"`
			Path *string `db:"file_path"`
		}{}, `SELECT id, file_path FROM admin_fonts ORDER BY created_at ASC LIMIT 1`)
		if err == nil {
			// снова заберём нормальными переменными
			_ = tx.GetContext(ctx, &nextID, `SELECT id FROM admin_fonts ORDER BY created_at ASC LIMIT 1`)
			_ = tx.GetContext(ctx, &nextPath, `SELECT file_path FROM admin_fonts ORDER BY created_at ASC LIMIT 1`)
			if _, err := tx.ExecContext(ctx, `UPDATE admin_fonts SET selected = FALSE, updated_at = now()`); err != nil {
				return err
			}
			if _, err := tx.ExecContext(ctx, `
				UPDATE admin_fonts SET selected = TRUE, updated_at = now()
				WHERE id = $1
			`, nextID); err != nil {
				return err
			}
			// commit -> sync after
			if err := tx.Commit(); err != nil {
				return err
			}

			_ = os.Remove(filePath)
			return s.syncMainFont(nextPath)
		}
		// нет оставшихся — commit и просто удалим файл
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	_ = os.Remove(filePath)
	return nil
}

// CreateFontFromForm: name + file (multipart)
func (s *AdminFontsService) CreateFontFromForm(ctx context.Context, name string, fileHeader *multipart.FileHeader) (*entity.AdminFont, error) {
	if err := s.checkDB("CreateFontFromForm"); err != nil {
		return nil, err
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("name is required")
	}
	if fileHeader == nil {
		return nil, errors.New("file is required")
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !isAllowedFontExt(ext) {
		return nil, fmt.Errorf("unsupported font extension: %s", ext)
	}

	id := newID()
	safeBase := sanitizeFileBase(strings.TrimSuffix(fileHeader.Filename, filepath.Ext(fileHeader.Filename)))
	if safeBase == "" {
		safeBase = "font"
	}

	if err := os.MkdirAll(s.fontsDir, 0o755); err != nil {
		return nil, err
	}

	dstRel := filepath.ToSlash(filepath.Join(s.fontsDir, fmt.Sprintf("%s_%s%s", id, safeBase, ext)))
	dstAbs := filepath.Clean(dstRel)

	// save file first (чтобы в БД не было битых путей)
	src, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	out, err := os.Create(dstAbs)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	if _, err := io.Copy(out, src); err != nil {
		_ = os.Remove(dstAbs)
		return nil, err
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		_ = os.Remove(dstAbs)
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	// если это первый шрифт — делаем selected=true
	var count int
	if err := tx.GetContext(ctx, &count, `SELECT COUNT(1) FROM admin_fonts`); err != nil {
		_ = os.Remove(dstAbs)
		return nil, err
	}
	selected := (count == 0)

	if selected {
		if _, err := tx.ExecContext(ctx, `UPDATE admin_fonts SET selected = FALSE, updated_at = now()`); err != nil {
			_ = os.Remove(dstAbs)
			return nil, err
		}
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO admin_fonts (id, name, file_path, selected, created_at, updated_at)
		VALUES ($1, $2, $3, $4, now(), now())
	`, id, name, dstRel, selected); err != nil {
		_ = os.Remove(dstAbs)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		_ = os.Remove(dstAbs)
		return nil, err
	}

	row := &entity.AdminFont{
		ID:       id,
		Name:     name,
		FilePath: dstRel,
		Selected: selected,
	}

	if selected {
		_ = s.syncMainFont(dstRel)
	}

	return row, nil
}

func (s *AdminFontsService) SelectFont(ctx context.Context, id string) error {
	if err := s.checkDB("SelectFont"); err != nil {
		return err
	}
	id = strings.TrimSpace(id)
	if id == "" {
		return errors.New("id is required")
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	var path string
	if err := tx.GetContext(ctx, &path, `SELECT file_path FROM admin_fonts WHERE id = $1`, id); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `UPDATE admin_fonts SET selected = FALSE, updated_at = now()`); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `
		UPDATE admin_fonts SET selected = TRUE, updated_at = now()
		WHERE id = $1
	`, id); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return s.syncMainFont(path)
}

// -------- FS helpers --------

func (s *AdminFontsService) syncMainFont(selectedFontPath string) error {
	if selectedFontPath == "" {
		return nil
	}

	if err := os.MkdirAll(s.buildDir, 0o755); err != nil {
		return err
	}

	ext := strings.ToLower(filepath.Ext(selectedFontPath))
	if ext == "" {
		ext = ".ttf"
	}

	// Удаляем build/main_font*
	glob := filepath.Join(s.buildDir, "main_font*")
	matches, _ := filepath.Glob(glob)
	for _, m := range matches {
		_ = os.Remove(m)
	}

	src := filepath.Clean(selectedFontPath)
	dst := filepath.Join(s.buildDir, "main_font"+ext)

	return copyFile(src, dst)
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}

func isAllowedFontExt(ext string) bool {
	switch strings.ToLower(ext) {
	case ".ttf", ".otf", ".woff", ".woff2":
		return true
	default:
		return false
	}
}

func sanitizeFileBase(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	s = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			return r
		}
		if r == ' ' {
			return '-'
		}
		return -1
	}, s)
	s = strings.Trim(s, "-_")
	if len(s) > 60 {
		s = s[:60]
	}
	return s
}

func newID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}
