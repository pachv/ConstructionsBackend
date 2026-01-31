package services

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
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

	// Нельзя удалять, если шрифт в базе только один
	var total int
	if err := tx.GetContext(ctx, &total, `SELECT COUNT(1) FROM admin_fonts`); err != nil {
		return err
	}
	if total <= 1 {
		return errors.New("cannot delete the last remaining font")
	}

	// 1) Забираем file_path и selected одним запросом
	var cur struct {
		FilePath string `db:"file_path"`
		Selected bool   `db:"selected"`
	}
	if err := tx.GetContext(ctx, &cur, `
		SELECT file_path, selected
		FROM admin_fonts
		WHERE id = $1
	`, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("font not found: %s", id)
		}
		s.logger.Error("DeleteFont: select failed", "err", err, "id", id)
		return err
	}

	// 2) Удаляем запись
	if _, err := tx.ExecContext(ctx, `DELETE FROM admin_fonts WHERE id = $1`, id); err != nil {
		s.logger.Error("DeleteFont: delete failed", "err", err, "id", id)
		return err
	}

	// 3) Если удалили selected — выбираем следующую (если есть) и помечаем selected=true
	var nextPath string
	if cur.Selected {
		var next struct {
			ID       string `db:"id"`
			FilePath string `db:"file_path"`
		}

		err := tx.GetContext(ctx, &next, `
			SELECT id, file_path
			FROM admin_fonts
			ORDER BY created_at ASC
			LIMIT 1
		`)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return err
			}
		} else {
			if _, err := tx.ExecContext(ctx, `UPDATE admin_fonts SET selected = FALSE, updated_at = now()`); err != nil {
				return err
			}
			if _, err := tx.ExecContext(ctx, `
				UPDATE admin_fonts
				SET selected = TRUE, updated_at = now()
				WHERE id = $1
			`, next.ID); err != nil {
				return err
			}
			nextPath = next.FilePath
		}
	}

	// 4) Коммитим транзакцию
	if err := tx.Commit(); err != nil {
		return err
	}

	// 5) Удаляем файл (после commit)
	_ = s.removeFontFileSafe(cur.FilePath)

	// 6) Если был selected и нашли следующий — обновляем main.css по шаблону
	if nextPath != "" {
		return s.syncMainFontFromTemplate(nextPath)
	}

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

	// Файл всегда кладём в /app/build/static/media/{file}
	filename := fmt.Sprintf("%s_%s%s", id, safeBase, ext)

	// То, что храним в БД (веб-относительный путь)
	dstRel := filepath.ToSlash(filepath.Join("static", "media", filename))

	// Абсолютный путь в контейнере (куда реально пишем файл)
	dstAbs := filepath.Clean(filepath.Join("/app/build", "static", "media", filename))

	if err := os.MkdirAll(filepath.Dir(dstAbs), 0o755); err != nil {
		return nil, err
	}

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
		// оставил как у тебя: syncMainFontFromTemplate получает "путь как в БД" (static/media/...)
		_ = s.syncMainFontFromTemplate(dstRel)
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

	return s.syncMainFontFromTemplate(path)
}

// -------- FS helpers --------
func (s *AdminFontsService) syncMainFontFromTemplate(fontRelPath string) error {
	fontRelPath = strings.TrimSpace(fontRelPath)
	if fontRelPath == "" {
		return errors.New("syncMainFontFromTemplate: fontRelPath is empty")
	}

	// "static/media/xxx.ttf" -> "/static/media/xxx.ttf"
	webPath := "/" + filepath.ToSlash(strings.TrimPrefix(fontRelPath, "/"))
	newSrcDecl := "\tsrc: url(" + webPath + ");"

	// 1) читаем шаблон
	tplPath := "/app/css_templates/main.css" // поменяй, если путь другой
	tplBytes, err := os.ReadFile(tplPath)
	if err != nil {
		return fmt.Errorf("syncMainFontFromTemplate: read template %s: %w", tplPath, err)
	}
	css := string(tplBytes)

	// 2) берём первый @font-face блок
	reFontFace := regexp.MustCompile(`(?s)@font-face\s*\{.*?\}`)
	blockLoc := reFontFace.FindStringIndex(css)
	if blockLoc == nil {
		return fmt.Errorf("syncMainFontFromTemplate: @font-face block not found in template %s", tplPath)
	}
	block := css[blockLoc[0]:blockLoc[1]]

	// 3) находим "src:" внутри блока (без требования ";")
	reSrcStart := regexp.MustCompile(`(?i)\bsrc\s*:`)
	srcLoc := reSrcStart.FindStringIndex(block)
	if srcLoc == nil {
		return fmt.Errorf("syncMainFontFromTemplate: src: not found inside @font-face in template %s", tplPath)
	}

	// srcLoc[0] — начало "src", srcLoc[1] — позиция сразу после ":"
	start := srcLoc[0]
	afterColon := srcLoc[1]

	// Найдём конец декларации src:
	// приоритет: ';' (если есть), иначе конец строки, иначе перед '}'
	rest := block[afterColon:]

	semi := strings.Index(rest, ";")
	nl := strings.IndexAny(rest, "\r\n")
	brace := strings.Index(rest, "}")

	// выбираем ближайший "разделитель" (но brace обязан существовать в блоке)
	endRel := -1
	endKeep := "" // что оставить после заменённого участка: ";" или "\n" или ""
	if semi >= 0 {
		endRel = semi
		endKeep = ";" // сохраним ';'
	} else if nl >= 0 {
		endRel = nl
		endKeep = "" // перенос строки уже в rest[ nl: ] — мы его не съедаем
	} else if brace >= 0 {
		endRel = brace
		endKeep = "" // не съедаем '}'
	} else {
		// на всякий случай: если сломанный шаблон
		return fmt.Errorf("syncMainFontFromTemplate: malformed @font-face block in template %s", tplPath)
	}

	// absolute end внутри block
	end := afterColon + endRel
	// Если нашли ';' — включим его в заменяемую часть
	if endKeep == ";" {
		end++ // съедаем ';'
	}

	// заменяем ровно декларацию src: ....
	newBlock := block[:start] + newSrcDecl + block[end:]

	// подменяем блок в css
	css = css[:blockLoc[0]] + newBlock + css[blockLoc[1]:]

	// 4) перезаписываем итоговый main.css
	dstPath := "/app/build/static/css/main.0293452a.css"
	if err := os.MkdirAll(filepath.Dir(dstPath), 0o755); err != nil {
		return fmt.Errorf("syncMainFontFromTemplate: mkdir %s: %w", filepath.Dir(dstPath), err)
	}
	if err := os.WriteFile(dstPath, []byte(css), 0o644); err != nil {
		return fmt.Errorf("syncMainFontFromTemplate: write %s: %w", dstPath, err)
	}

	return nil
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

func (s *AdminFontsService) removeFontFileSafe(fontRelPath string) error {
	fontRelPath = strings.TrimSpace(fontRelPath)
	if fontRelPath == "" {
		return nil
	}

	// ожидаем "static/media/..."
	abs := filepath.Clean(filepath.Join("/app/build", filepath.FromSlash(strings.TrimPrefix(fontRelPath, "/"))))

	allowedRoot := filepath.Clean("/app/build/static/media") + string(os.PathSeparator)
	if !strings.HasPrefix(abs+string(os.PathSeparator), allowedRoot) {
		return fmt.Errorf("removeFontFileSafe: refused to delete outside %s: %s", allowedRoot, abs)
	}

	if err := os.Remove(abs); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}
