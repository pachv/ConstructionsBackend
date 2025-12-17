package repositories

import (
	"database/sql"
	"errors"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/pachv/constructions/constructions/internal/domain/entity"
)

type UserRepository struct {
	logger *slog.Logger
	db     *sqlx.DB
}

func NewUserRepository(logger *slog.Logger, db *sqlx.DB) *UserRepository {
	return &UserRepository{
		logger: logger.With("component", "UserRepository"),
		db:     db,
	}
}

func (r *UserRepository) DoesUserExist(login string) (bool, error) {

	userExistQuery := `
		SELECT EXISTS(
			SELECT 1
			FROM users
			WHERE login_name = $1
		)
	`

	var userExist bool

	err := r.db.Get(&userExist, userExistQuery, login)
	if err != nil {
		return false, err
	}

	return userExist, nil
}

func (r *UserRepository) RegisterUser(id, surname, name, login, fathername, email, phoneNumber, hashedPassword string) error {
	createUserQuery := `
		INSERT INTO users(id,surname,username,login_name,father_name,email,phone_number,hashed_password)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8);
	`

	_, err := r.db.Exec(createUserQuery, id, surname, name, login, fathername, email, phoneNumber, hashedPassword)
	if err != nil {
		return err
	}

	return nil
}

var ErrUserNotFound = errors.New("user not found")

func (r *UserRepository) GetUserByLogin(login string) (entity.User, error) {
	const q = `
		SELECT
			id,
			surname,
			login_name,
			father_name,
			email,
			phone_number,
			hashed_password
		FROM users
		WHERE login_name = $1
		LIMIT 1;
	`

	var user entity.User
	err := r.db.Get(&user, q, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, ErrUserNotFound
		}
		return entity.User{}, err
	}

	return user, nil
}

func (r *UserRepository) GetUserByID(id string) (entity.User, error) {
	const q = `
		SELECT id, surname, login_name, father_name, email, phone_number, hashed_password
		FROM users
		WHERE id = $1
		LIMIT 1;
	`

	var u entity.User
	if err := r.db.Get(&u, q, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, ErrUserNotFound
		}
		return entity.User{}, err
	}

	return u, nil
}

func (r *UserRepository) UpdateUserPassword(userID string, hashedPassword string) error {
	const q = `
		UPDATE users
		SET hashed_password = $2
		WHERE id = $1;
	`

	res, err := r.db.Exec(q, userID, hashedPassword)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrUserNotFound
	}

	return nil
}
