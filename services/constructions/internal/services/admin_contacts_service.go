package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pachv/constructions/constructions/internal/domain/entity"
)

const contactsEmailSingletonID = "singleton"

var ErrAdminContactsDBNil = errors.New("admin contacts service: db is nil")

type AdminContactsService struct {
	db     *sqlx.DB
	logger *slog.Logger
}

func NewAdminContactsService(db *sqlx.DB, logger *slog.Logger) *AdminContactsService {
	if logger == nil {
		logger = slog.Default()
	}
	return &AdminContactsService{
		db:     db,
		logger: logger.With("component", "admin_contacts_service"),
	}
}

func (s *AdminContactsService) checkDB(method string) error {
	if s == nil || s.db == nil {
		if s != nil && s.logger != nil {
			s.logger.Error("db is nil", "method", method)
		}
		return ErrAdminContactsDBNil
	}
	return nil
}

// -------- Email (singleton) --------

func (s *AdminContactsService) GetContactsEmail(ctx context.Context) (string, error) {
	if err := s.checkDB("GetContactsEmail"); err != nil {
		return "", err
	}

	s.logger.Info("GetContactsEmail: start", "id", contactsEmailSingletonID)

	var email string
	err := s.db.GetContext(ctx, &email, `
		SELECT COALESCE(email, '')
		FROM contacts_email_settings
		WHERE id = $1
		LIMIT 1
	`, contactsEmailSingletonID)

	if err != nil {
		low := strings.ToLower(err.Error())
		if strings.Contains(low, "no rows") {
			s.logger.Warn("GetContactsEmail: no rows (return empty)", "id", contactsEmailSingletonID)
			return "", nil
		}
		s.logger.Error("GetContactsEmail: query failed", "err", err)
		return "", err
	}

	s.logger.Info("GetContactsEmail: ok", "email", email)
	return email, nil
}

func (s *AdminContactsService) SetContactsEmail(ctx context.Context, email string) error {
	if err := s.checkDB("SetContactsEmail"); err != nil {
		return err
	}

	email = strings.TrimSpace(email)
	s.logger.Info("SetContactsEmail: start", "id", contactsEmailSingletonID, "email", email)

	res, err := s.db.ExecContext(ctx, `
		UPDATE contacts_email_settings
		SET email = $2, updated_at = now()
		WHERE id = $1
	`, contactsEmailSingletonID, email)
	if err != nil {
		s.logger.Error("SetContactsEmail: update failed", "err", err)
		return err
	}

	aff, _ := res.RowsAffected()
	s.logger.Info("SetContactsEmail: update rows affected", "rows", aff)

	if aff > 0 {
		return nil
	}

	_, err = s.db.ExecContext(ctx, `
		INSERT INTO contacts_email_settings (id, email, created_at, updated_at)
		VALUES ($1, $2, now(), now())
	`, contactsEmailSingletonID, email)
	if err != nil {
		s.logger.Error("SetContactsEmail: insert failed", "err", err)
		return err
	}

	s.logger.Info("SetContactsEmail: inserted singleton", "id", contactsEmailSingletonID)
	return nil
}

// Если тебе нужно именно entity (редко, но вдруг)
func (s *AdminContactsService) GetContactsEmailSetting(ctx context.Context) (*entity.ContactsEmailSetting, error) {
	if err := s.checkDB("GetContactsEmailSetting"); err != nil {
		return nil, err
	}

	s.logger.Info("GetContactsEmailSetting: start", "id", contactsEmailSingletonID)

	var row entity.ContactsEmailSetting
	err := s.db.GetContext(ctx, &row, `
		SELECT id, email, created_at, updated_at
		FROM contacts_email_settings
		WHERE id = $1
		LIMIT 1
	`, contactsEmailSingletonID)
	if err != nil {
		low := strings.ToLower(err.Error())
		if strings.Contains(low, "no rows") {
			s.logger.Warn("GetContactsEmailSetting: no rows (return empty entity)", "id", contactsEmailSingletonID)
			return &entity.ContactsEmailSetting{ID: contactsEmailSingletonID, Email: ""}, nil
		}
		s.logger.Error("GetContactsEmailSetting: query failed", "err", err)
		return nil, err
	}

	s.logger.Info("GetContactsEmailSetting: ok", "id", row.ID, "email", row.Email)
	return &row, nil
}

// -------- Phones --------

func (s *AdminContactsService) ListContactNumbers(ctx context.Context) ([]entity.ContactNumber, error) {
	if err := s.checkDB("ListContactNumbers"); err != nil {
		return nil, err
	}

	s.logger.Info("ListContactNumbers: start")

	items := make([]entity.ContactNumber, 0)
	err := s.db.SelectContext(ctx, &items, `
		SELECT id, phone, label, COALESCE(sort_order, 0) AS sort_order, created_at, updated_at
		FROM contacts_numbers
		ORDER BY sort_order ASC NULLS LAST, created_at ASC
	`)
	if err != nil {
		s.logger.Error("ListContactNumbers: query failed", "err", err)
		return nil, err
	}

	s.logger.Info("ListContactNumbers: ok", "count", len(items))
	return items, nil
}

func (s *AdminContactsService) UpsertContactNumber(ctx context.Context, item entity.ContactNumber) error {
	if err := s.checkDB("UpsertContactNumber"); err != nil {
		return err
	}

	item.ID = strings.TrimSpace(item.ID)
	if item.ID == "" {
		return errors.New("id is required")
	}

	s.logger.Info("UpsertContactNumber: start",
		"id", item.ID,
		"phone", strings.TrimSpace(item.Phone),
		"label", strings.TrimSpace(item.Label),
		"sortOrder", item.SortOrder,
	)

	res, err := s.db.ExecContext(ctx, `
		UPDATE contacts_numbers
		SET phone = $2,
		    label = $3,
		    sort_order = $4,
		    updated_at = now()
		WHERE id = $1
	`, item.ID, strings.TrimSpace(item.Phone), strings.TrimSpace(item.Label), item.SortOrder)
	if err != nil {
		s.logger.Error("UpsertContactNumber: update failed", "err", err, "id", item.ID)
		return err
	}

	aff, _ := res.RowsAffected()
	s.logger.Info("UpsertContactNumber: update rows affected", "id", item.ID, "rows", aff)

	if aff > 0 {
		return nil
	}

	_, err = s.db.ExecContext(ctx, `
		INSERT INTO contacts_numbers (id, phone, label, sort_order, created_at, updated_at)
		VALUES ($1, $2, $3, $4, now(), now())
	`, item.ID, strings.TrimSpace(item.Phone), strings.TrimSpace(item.Label), item.SortOrder)
	if err != nil {
		s.logger.Error("UpsertContactNumber: insert failed", "err", err, "id", item.ID)
		return err
	}

	s.logger.Info("UpsertContactNumber: inserted", "id", item.ID)
	return nil
}

func (s *AdminContactsService) DeleteContactNumber(ctx context.Context, id string) error {
	if err := s.checkDB("DeleteContactNumber"); err != nil {
		return err
	}

	id = strings.TrimSpace(id)
	if id == "" {
		return errors.New("id is required")
	}

	s.logger.Info("DeleteContactNumber: start", "id", id)

	res, err := s.db.ExecContext(ctx, `DELETE FROM contacts_numbers WHERE id = $1`, id)
	if err != nil {
		s.logger.Error("DeleteContactNumber: delete failed", "err", err, "id", id)
		return err
	}
	aff, _ := res.RowsAffected()
	s.logger.Info("DeleteContactNumber: ok", "id", id, "rows", aff)
	return nil
}

func (s *AdminContactsService) SetContactNumbersOrder(ctx context.Context, ids []string) error {
	if err := s.checkDB("SetContactNumbersOrder"); err != nil {
		return err
	}

	s.logger.Info("SetContactNumbersOrder: start", "count", len(ids))

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		s.logger.Error("SetContactNumbersOrder: begin tx failed", "err", err)
		return err
	}
	defer func() { _ = tx.Rollback() }()

	for i, id := range ids {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}

		if _, err := tx.ExecContext(ctx, `
			UPDATE contacts_numbers
			SET sort_order = $2, updated_at = now()
			WHERE id = $1
		`, id, i); err != nil {
			s.logger.Error("SetContactNumbersOrder: update failed", "err", err, "id", id, "pos", i)
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		s.logger.Error("SetContactNumbersOrder: commit failed", "err", err)
		return err
	}

	s.logger.Info("SetContactNumbersOrder: ok")
	return nil
}

// -------- Addresses --------

func (s *AdminContactsService) ListContactAddresses(ctx context.Context) ([]entity.ContactAddress, error) {
	if err := s.checkDB("ListContactAddresses"); err != nil {
		return nil, err
	}

	s.logger.Info("ListContactAddresses: start")

	items := make([]entity.ContactAddress, 0)
	err := s.db.SelectContext(ctx, &items, `
		SELECT id, title, address, lat, lon, COALESCE(sort_order, 0) AS sort_order, created_at, updated_at
		FROM contacts_addresses
		ORDER BY sort_order ASC NULLS LAST, created_at ASC
	`)
	if err != nil {
		s.logger.Error("ListContactAddresses: query failed", "err", err)
		return nil, err
	}

	s.logger.Info("ListContactAddresses: ok", "count", len(items))
	return items, nil
}

func (s *AdminContactsService) UpsertContactAddress(ctx context.Context, item entity.ContactAddress) error {
	if err := s.checkDB("UpsertContactAddress"); err != nil {
		return err
	}

	item.ID = strings.TrimSpace(item.ID)
	if item.ID == "" {
		return errors.New("id is required")
	}

	s.logger.Info("UpsertContactAddress: start",
		"id", item.ID,
		"title", strings.TrimSpace(item.Title),
		"address", strings.TrimSpace(item.Address),
		"lat", item.Lat,
		"lon", item.Lon,
		"sortOrder", item.SortOrder,
	)

	res, err := s.db.ExecContext(ctx, `
		UPDATE contacts_addresses
		SET title = $2,
		    address = $3,
		    lat = $4,
		    lon = $5,
		    sort_order = $6,
		    updated_at = now()
		WHERE id = $1
	`, item.ID, strings.TrimSpace(item.Title), strings.TrimSpace(item.Address), item.Lat, item.Lon, item.SortOrder)
	if err != nil {
		s.logger.Error("UpsertContactAddress: update failed", "err", err, "id", item.ID)
		return err
	}

	aff, _ := res.RowsAffected()
	s.logger.Info("UpsertContactAddress: update rows affected", "id", item.ID, "rows", aff)

	if aff > 0 {
		return nil
	}

	_, err = s.db.ExecContext(ctx, `
		INSERT INTO contacts_addresses (id, title, address, lat, lon, sort_order, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, now(), now())
	`, item.ID, strings.TrimSpace(item.Title), strings.TrimSpace(item.Address), item.Lat, item.Lon, item.SortOrder)
	if err != nil {
		s.logger.Error("UpsertContactAddress: insert failed", "err", err, "id", item.ID)
		return err
	}

	s.logger.Info("UpsertContactAddress: inserted", "id", item.ID)
	return nil
}

func (s *AdminContactsService) DeleteContactAddress(ctx context.Context, id string) error {
	if err := s.checkDB("DeleteContactAddress"); err != nil {
		return err
	}

	id = strings.TrimSpace(id)
	if id == "" {
		return errors.New("id is required")
	}

	s.logger.Info("DeleteContactAddress: start", "id", id)

	res, err := s.db.ExecContext(ctx, `DELETE FROM contacts_addresses WHERE id = $1`, id)
	if err != nil {
		s.logger.Error("DeleteContactAddress: delete failed", "err", err, "id", id)
		return err
	}
	aff, _ := res.RowsAffected()
	s.logger.Info("DeleteContactAddress: ok", "id", id, "rows", aff)
	return nil
}

func (s *AdminContactsService) SetContactAddressesOrder(ctx context.Context, ids []string) error {
	if err := s.checkDB("SetContactAddressesOrder"); err != nil {
		return err
	}

	s.logger.Info("SetContactAddressesOrder: start", "count", len(ids))

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		s.logger.Error("SetContactAddressesOrder: begin tx failed", "err", err)
		return err
	}
	defer func() { _ = tx.Rollback() }()

	for i, id := range ids {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		if _, err := tx.ExecContext(ctx, `
			UPDATE contacts_addresses
			SET sort_order = $2, updated_at = now()
			WHERE id = $1
		`, id, i); err != nil {
			s.logger.Error("SetContactAddressesOrder: update failed", "err", err, "id", id, "pos", i)
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		s.logger.Error("SetContactAddressesOrder: commit failed", "err", err)
		return err
	}

	s.logger.Info("SetContactAddressesOrder: ok")
	return nil
}

// helper: иногда удобно логировать кратко slice id
func (s *AdminContactsService) debugIDs(ids []string) string {
	if len(ids) == 0 {
		return "[]"
	}
	lim := 10
	if len(ids) < lim {
		lim = len(ids)
	}
	return fmt.Sprintf("%v (len=%d)", ids[:lim], len(ids))
}
