package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pachv/constructions/constructions/internal/domain/entity"
)

type AdminEmailRepository struct {
	db *sqlx.DB
}

func NewAdminEmailRepository(db *sqlx.DB) *AdminEmailRepository {
	return &AdminEmailRepository{db: db}
}

func (r *AdminEmailRepository) GetEmail() string {
	const q = `
		SELECT email
		FROM admin_email_settings
		WHERE id = 'singleton'
		LIMIT 1;
	`

	var email string
	_ = r.db.Get(&email, q)
	return email
}

func (r *AdminEmailRepository) Get(ctx context.Context) (entity.AdminEmailSetting, error) {
	const q = `
		SELECT id, email, created_at, updated_at
		FROM admin_email_settings
		WHERE id = 'singleton'
		LIMIT 1;
	`

	var out entity.AdminEmailSetting
	err := r.db.GetContext(ctx, &out, q)
	if err != nil {
		// если строки нет — вернём дефолт, а строку создаст Set() или миграция
		return entity.AdminEmailSetting{ID: "singleton", Email: ""}, nil
	}
	return out, nil
}

func (r *AdminEmailRepository) Set(ctx context.Context, email string) error {
	// обновим, а если строки нет — вставим (без ON CONFLICT)
	const updateQ = `
		UPDATE admin_email_settings
		SET email = $1,
		    updated_at = now()
		WHERE id = 'singleton';
	`

	res, err := r.db.ExecContext(ctx, updateQ, email)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected > 0 {
		return nil
	}

	const insertQ = `
		INSERT INTO admin_email_settings (id, email, created_at, updated_at)
		VALUES ('singleton', $1, now(), now());
	`
	_, err = r.db.ExecContext(ctx, insertQ, email)
	return err
}
