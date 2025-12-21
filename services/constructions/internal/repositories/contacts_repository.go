package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/pachv/constructions/constructions/internal/domain/entity"
)

type ContactsRepository struct {
	db *sqlx.DB
}

func NewContactsRepository(db *sqlx.DB) *ContactsRepository {
	return &ContactsRepository{db: db}
}

func (r *ContactsRepository) GetEmail(ctx context.Context) (entity.ContactsEmailSetting, error) {
	const q = `
		SELECT id, email, created_at, updated_at
		FROM contacts_email_settings
		WHERE id = 'singleton'
		LIMIT 1;
	`

	var out entity.ContactsEmailSetting
	err := r.db.GetContext(ctx, &out, q)
	if err == nil {
		return out, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		// на всякий случай: если миграции не прогнались
		return entity.ContactsEmailSetting{ID: "singleton", Email: ""}, nil
	}
	return entity.ContactsEmailSetting{}, err
}

func (r *ContactsRepository) GetNumbers(ctx context.Context) ([]entity.ContactNumber, error) {
	const q = `
		SELECT id, phone, label, sort_order, created_at, updated_at
		FROM contacts_numbers
		ORDER BY sort_order ASC, created_at ASC;
	`

	var out []entity.ContactNumber
	if err := r.db.SelectContext(ctx, &out, q); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *ContactsRepository) GetAddresses(ctx context.Context) ([]entity.ContactAddress, error) {
	const q = `
		SELECT id, title, address, lat, lon, sort_order, created_at, updated_at
		FROM contacts_addresses
		ORDER BY sort_order ASC, created_at ASC;
	`

	var out []entity.ContactAddress
	if err := r.db.SelectContext(ctx, &out, q); err != nil {
		return nil, err
	}
	return out, nil
}
