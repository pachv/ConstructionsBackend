package sender

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// =========================
// URLs (как в твоём примере)
// =========================

const (
	// PublicAPIBaseURL = "http://constructions_service:8080"

	AdminFontsBaseURL   = PublicAPIBaseURL + "/admin/fonts"
	AdminFontsDeleteURL = PublicAPIBaseURL + "/admin/fonts/%s"
	AdminFontsSelectURL = PublicAPIBaseURL + "/admin/fonts/%s/select"
)

// =========================
// DTO
// =========================

type AdminFont struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	FilePath string `json:"filePath"`
	Selected bool   `json:"selected"`
	// createdAt/updatedAt можно добавить если надо, но обычно не нужно в UI
}

type okResp struct {
	Ok bool `json:"ok"`
}

type createFontResp struct {
	ID string `json:"id"`
}

// =========================
// GET list
// GET /admin/fonts
// =========================

func GetAdminFonts(ctx context.Context) ([]AdminFont, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, AdminFontsBaseURL, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("bad status: %s: %s", resp.Status, string(b))
	}

	var out []AdminFont
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	if out == nil {
		out = []AdminFont{}
	}
	return out, nil
}

// =========================
// CREATE (multipart)
// POST /admin/fonts
// form-data: name + file
// filePath — путь на диске (локально, где запускается sender)
// =========================

func CreateAdminFont(ctx context.Context, name string, filePath string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", fmt.Errorf("name is required")
	}
	filePath = strings.TrimSpace(filePath)
	if filePath == "" {
		return "", fmt.Errorf("filePath is required")
	}

	// Собираем multipart
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	_ = w.WriteField("name", name)

	f, err := os.Open(filePath)
	if err != nil {
		_ = w.Close()
		return "", err
	}
	defer f.Close()

	part, err := w.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		_ = w.Close()
		return "", err
	}
	if _, err := io.Copy(part, f); err != nil {
		_ = w.Close()
		return "", err
	}
	_ = w.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, AdminFontsBaseURL, &buf)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("bad status: %s: %s", resp.Status, string(b))
	}

	var out createFontResp
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	out.ID = strings.TrimSpace(out.ID)
	if out.ID == "" {
		return "", fmt.Errorf("empty id in response")
	}
	return out.ID, nil
}

// =========================
// SELECT
// POST /admin/fonts/:id/select
// =========================

func SelectAdminFont(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("id is required")
	}

	u := fmt.Sprintf(AdminFontsSelectURL, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, nil)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("bad status: %s: %s", resp.Status, string(b))
	}

	return nil
}

// =========================
// DELETE
// DELETE /admin/fonts/:id
// =========================

func DeleteAdminFont(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("id is required")
	}

	u := fmt.Sprintf(AdminFontsDeleteURL, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("bad status: %s: %s", resp.Status, string(b))
	}

	return nil
}
