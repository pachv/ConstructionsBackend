package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pachv/constructions/constructions/internal/domain/entity"
)

type CertificatesAdminService struct {
	db *sqlx.DB

	storageDir   string // ./uploads/certificates
	fileURLPrefx string // /admin/certificates/file/
}

func NewCertificatesAdminService(db *sqlx.DB) *CertificatesAdminService {
	return &CertificatesAdminService{
		db:           db,
		storageDir:   "./uploads/certificates",
		fileURLPrefx: "/admin/certificates/file/",
	}
}

func (s *CertificatesAdminService) GetAll(ctx context.Context) ([]entity.Certificate, error) {
	const q = `
		SELECT id, title, file_path, created_at
		FROM certificates
		ORDER BY created_at DESC;
	`
	var out []entity.Certificate
	if err := s.db.SelectContext(ctx, &out, q); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *CertificatesAdminService) GetByID(ctx context.Context, id string) (entity.Certificate, error) {
	const q = `
		SELECT id, title, file_path, created_at
		FROM certificates
		WHERE id = $1
		LIMIT 1;
	`
	var out entity.Certificate
	if err := s.db.GetContext(ctx, &out, q, id); err != nil {
		return entity.Certificate{}, err
	}
	return out, nil
}

func (s *CertificatesAdminService) Create(ctx context.Context, title string, file *multipart.FileHeader) (entity.Certificate, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		title = "Сертификат"
	}
	if file == nil {
		return entity.Certificate{}, fmt.Errorf("file is required")
	}

	if err := os.MkdirAll(s.storageDir, 0o755); err != nil {
		return entity.Certificate{}, err
	}

	filename, err := s.saveUploadedFile(file)
	if err != nil {
		return entity.Certificate{}, err
	}

	id := "cert-" + randomHex(8)
	filePath := s.fileURLPrefx + filename

	const q = `
		INSERT INTO certificates (id, title, file_path, created_at)
		VALUES ($1, $2, $3, $4);
	`

	createdAt := time.Now()
	if _, err := s.db.ExecContext(ctx, q, id, title, filePath, createdAt); err != nil {
		_ = os.Remove(filepath.Join(s.storageDir, filename))
		return entity.Certificate{}, err
	}

	return entity.Certificate{
		ID:        id,
		Title:     title,
		FilePath:  filePath,
		CreatedAt: &createdAt,
	}, nil
}

func (s *CertificatesAdminService) Update(ctx context.Context, id string, newTitle *string, newFile *multipart.FileHeader) (entity.Certificate, error) {
	existing, err := s.GetByID(ctx, id)
	if err != nil {
		return entity.Certificate{}, err
	}

	title := existing.Title
	if newTitle != nil {
		t := strings.TrimSpace(*newTitle)
		if t != "" {
			title = t
		}
	}

	filePath := existing.FilePath
	oldFilename := s.filenameFromFilePath(existing.FilePath)

	// если меняем файл — сначала сохраняем новый, потом обновляем БД, потом удаляем старый
	var newFilename string
	if newFile != nil {
		if err := os.MkdirAll(s.storageDir, 0o755); err != nil {
			return entity.Certificate{}, err
		}

		newFilename, err = s.saveUploadedFile(newFile)
		if err != nil {
			return entity.Certificate{}, err
		}
		filePath = s.fileURLPrefx + newFilename
	}

	const q = `
		UPDATE certificates
		SET title = $2,
		    file_path = $3
		WHERE id = $1;
	`

	if _, err := s.db.ExecContext(ctx, q, id, title, filePath); err != nil {
		// откат по файлу: если мы успели записать новый файл — удалим его
		if newFilename != "" {
			_ = os.Remove(filepath.Join(s.storageDir, newFilename))
		}
		return entity.Certificate{}, err
	}

	// удаляем старый файл best-effort
	if newFilename != "" && oldFilename != "" {
		_ = os.Remove(filepath.Join(s.storageDir, oldFilename))
	}

	updated, err := s.GetByID(ctx, id)
	if err != nil {
		// БД уже обновили, но вернуть что-то надо
		return entity.Certificate{ID: id, Title: title, FilePath: filePath, CreatedAt: existing.CreatedAt}, nil
	}
	return updated, nil
}

func (s *CertificatesAdminService) Delete(ctx context.Context, id string) error {
	existing, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}

	const q = `DELETE FROM certificates WHERE id = $1;`
	if _, err := s.db.ExecContext(ctx, q, id); err != nil {
		return err
	}

	// best-effort удалить файл
	oldFilename := s.filenameFromFilePath(existing.FilePath)
	if oldFilename != "" {
		_ = os.Remove(filepath.Join(s.storageDir, oldFilename))
	}
	return nil
}

func (s *CertificatesAdminService) saveUploadedFile(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext == "" || len(ext) > 10 {
		ext = ".bin"
	}

	filename := "cert-" + randomHex(12) + ext
	dstPath := filepath.Join(s.storageDir, filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		_ = os.Remove(dstPath)
		return "", err
	}

	return filename, nil
}

func (s *CertificatesAdminService) filenameFromFilePath(p string) string {
	// ожидаем "/admin/certificates/file/<name>"
	if p == "" {
		return ""
	}
	p = strings.TrimSpace(strings.ReplaceAll(p, "\\", "/"))
	i := strings.LastIndex(p, "/")
	if i < 0 || i == len(p)-1 {
		return ""
	}
	return p[i+1:]
}

func randomHex(nBytes int) string {
	b := make([]byte, nBytes)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
