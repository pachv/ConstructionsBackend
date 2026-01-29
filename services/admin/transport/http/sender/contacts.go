package sender

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// =========================
// URLs
// =========================

const (
	// PublicAPIBaseURL = "http://constructions_service:8080"

	AdminContactsBaseURL             = PublicAPIBaseURL + "/admin/contacts"
	AdminContactsEmailURL            = PublicAPIBaseURL + "/admin/contacts/email"
	AdminContactsPhonesURL           = PublicAPIBaseURL + "/admin/contacts/phones"
	AdminContactsPhoneDeleteURL      = PublicAPIBaseURL + "/admin/contacts/phones/%s"
	AdminContactsPhonesReorderURL    = PublicAPIBaseURL + "/admin/contacts/phones/reorder"
	AdminContactsAddressesURL        = PublicAPIBaseURL + "/admin/contacts/addresses"
	AdminContactsAddressDeleteURL    = PublicAPIBaseURL + "/admin/contacts/addresses/%s"
	AdminContactsAddressesReorderURL = PublicAPIBaseURL + "/admin/contacts/addresses/reorder"
)

// =========================
// DTO (как в handler)
// =========================

type AdminContactsSnapshotResp struct {
	Email     string                `json:"email"`
	Phones    []AdminContactNumber  `json:"phones"`
	Addresses []AdminContactAddress `json:"addresses"`
}

type AdminContactNumber struct {
	ID        string `json:"id"`
	Phone     string `json:"phone"`
	Label     string `json:"label"`
	SortOrder int    `json:"sortOrder"`
}

type AdminContactAddress struct {
	ID        string  `json:"id"`
	Title     string  `json:"title"`
	Address   string  `json:"address"`
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	SortOrder int     `json:"sortOrder"`
}

type adminSetEmailReq struct {
	Email string `json:"email"`
}

type adminUpsertPhoneReq struct {
	ID        string `json:"id"`
	Phone     string `json:"phone"`
	Label     string `json:"label"`
	SortOrder int    `json:"sortOrder"`
}

type adminUpsertAddressReq struct {
	ID        string  `json:"id"`
	Title     string  `json:"title"`
	Address   string  `json:"address"`
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	SortOrder int     `json:"sortOrder"`
}

type adminReorderReq struct {
	IDs []string `json:"ids"`
}

type okRespp struct {
	Ok bool `json:"ok"`
}

// =========================
// Helpers
// =========================

func doJSONNN(ctx context.Context, method, url string, in any, out any, timeout time.Duration) error {
	var body io.Reader
	if in != nil {
		b, err := json.Marshal(in)
		if err != nil {
			return err
		}
		body = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return err
	}
	if in != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("bad status: %s: %s", resp.Status, string(b))
	}

	if out != nil {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return err
		}
	}
	return nil
}

// =========================
// GET /admin/contacts (snapshot)
// =========================

func GetAdminContacts(ctx context.Context) (AdminContactsSnapshotResp, error) {
	var out AdminContactsSnapshotResp
	err := doJSONNN(ctx, http.MethodGet, AdminContactsBaseURL, nil, &out, 10*time.Second)
	if err != nil {
		return AdminContactsSnapshotResp{}, err
	}
	if out.Phones == nil {
		out.Phones = []AdminContactNumber{}
	}
	if out.Addresses == nil {
		out.Addresses = []AdminContactAddress{}
	}
	return out, nil
}

// =========================
// PUT /admin/contacts/email
// =========================

func SetAdminContactsEmail(ctx context.Context, email string) error {
	email = strings.TrimSpace(email)
	// handler сейчас не валидирует email — оставим мягко, но пустое не шлём
	if email == "" {
		return fmt.Errorf("email is required")
	}
	return doJSONNN(ctx, http.MethodPut, AdminContactsEmailURL, adminSetEmailReq{Email: email}, &okRespp{}, 10*time.Second)
}

// =========================
// GET /admin/contacts/phones
// =========================

func ListAdminContactPhones(ctx context.Context) ([]AdminContactNumber, error) {
	var out []AdminContactNumber
	err := doJSONNN(ctx, http.MethodGet, AdminContactsPhonesURL, nil, &out, 10*time.Second)
	if err != nil {
		return nil, err
	}
	if out == nil {
		out = []AdminContactNumber{}
	}
	return out, nil
}

// =========================
// PUT /admin/contacts/phones (upsert)
// =========================

func UpsertAdminContactPhone(ctx context.Context, item AdminContactNumber) error {
	item.ID = strings.TrimSpace(item.ID)
	if item.ID == "" {
		return fmt.Errorf("id is required")
	}

	req := adminUpsertPhoneReq{
		ID:        item.ID,
		Phone:     item.Phone,
		Label:     item.Label,
		SortOrder: item.SortOrder,
	}
	return doJSONNN(ctx, http.MethodPut, AdminContactsPhonesURL, req, &okRespp{}, 10*time.Second)
}

// =========================
// DELETE /admin/contacts/phones/:id
// =========================

func DeleteAdminContactPhone(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("id is required")
	}
	u := fmt.Sprintf(AdminContactsPhoneDeleteURL, id)
	return doJSONNN(ctx, http.MethodDelete, u, nil, &okRespp{}, 10*time.Second)
}

// =========================
// PUT /admin/contacts/phones/reorder
// =========================

func ReorderAdminContactPhones(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return fmt.Errorf("ids is required")
	}
	// чуть подчистим
	clean := make([]string, 0, len(ids))
	for _, id := range ids {
		id = strings.TrimSpace(id)
		if id != "" {
			clean = append(clean, id)
		}
	}
	if len(clean) == 0 {
		return fmt.Errorf("ids is required")
	}

	return doJSONNN(ctx, http.MethodPut, AdminContactsPhonesReorderURL, adminReorderReq{IDs: clean}, &okRespp{}, 10*time.Second)
}

// =========================
// GET /admin/contacts/addresses
// =========================

func ListAdminContactAddresses(ctx context.Context) ([]AdminContactAddress, error) {
	var out []AdminContactAddress
	err := doJSONNN(ctx, http.MethodGet, AdminContactsAddressesURL, nil, &out, 10*time.Second)
	if err != nil {
		return nil, err
	}
	if out == nil {
		out = []AdminContactAddress{}
	}
	return out, nil
}

// =========================
// PUT /admin/contacts/addresses (upsert)
// =========================

func UpsertAdminContactAddress(ctx context.Context, item AdminContactAddress) error {
	item.ID = strings.TrimSpace(item.ID)
	if item.ID == "" {
		return fmt.Errorf("id is required")
	}

	req := adminUpsertAddressReq{
		ID:        item.ID,
		Title:     item.Title,
		Address:   item.Address,
		Lat:       item.Lat,
		Lon:       item.Lon,
		SortOrder: item.SortOrder,
	}
	return doJSONNN(ctx, http.MethodPut, AdminContactsAddressesURL, req, &okRespp{}, 10*time.Second)
}

// =========================
// DELETE /admin/contacts/addresses/:id
// =========================

func DeleteAdminContactAddress(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("id is required")
	}
	u := fmt.Sprintf(AdminContactsAddressDeleteURL, id)
	return doJSONNN(ctx, http.MethodDelete, u, nil, &okRespp{}, 10*time.Second)
}

// =========================
// PUT /admin/contacts/addresses/reorder
// =========================

func ReorderAdminContactAddresses(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return fmt.Errorf("ids is required")
	}
	clean := make([]string, 0, len(ids))
	for _, id := range ids {
		id = strings.TrimSpace(id)
		if id != "" {
			clean = append(clean, id)
		}
	}
	if len(clean) == 0 {
		return fmt.Errorf("ids is required")
	}

	return doJSONNN(ctx, http.MethodPut, AdminContactsAddressesReorderURL, adminReorderReq{IDs: clean}, &okRespp{}, 10*time.Second)
}
