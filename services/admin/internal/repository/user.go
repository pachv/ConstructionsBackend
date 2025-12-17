package repository

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/is_backend/services/admin/internal/domain/entity"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetUserByUsername(username string) (*entity.User, error) {

	var user entity.User

	selectUserQuery := `
		SELECT username, hashed_password
		FROM users
		WHERE username = $1
		LIMIT 1;
	`

	err := r.db.Get(&user, selectUserQuery, username)
	if err != nil {
		return nil, fmt.Errorf("selecting user error : %s", err.Error())
	}

	return &user, nil

}

func (r *UserRepository) GetUsersAmount() (usersAmount int, err error) {

	query := `
		SELECT COUNT(id)
		FROM users
	`

	err = r.db.Get(&usersAmount, query)
	if err != nil {
		return 0, err
	}

	return
}

func (r *UserRepository) GetUsers(page int, search, orderBy string) ([]*entity.User, int, error) {
	const pageSize = 10
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	if orderBy == "" {
		orderBy = "id"
	}

	var (
		users       []*entity.User
		totalUsers  int
		args        []interface{}
		whereClause string
	)

	baseQuery := `
		SELECT id, username, hashed_password
		FROM users
	`
	countQuery := `
		SELECT COUNT(*)
		FROM users
	`

	if search != "" {
		whereClause = "WHERE id ILIKE $1 OR username ILIKE $1"
		args = append(args, "%"+search+"%")
	}

	query := baseQuery
	if whereClause != "" {
		query += " " + whereClause
		countQuery += " " + whereClause
	}

	query += fmt.Sprintf(" ORDER BY %s LIMIT %d OFFSET %d", orderBy, pageSize, offset)

	fmt.Println("query : " + query)

	// Получаем общее количество пользователей
	if err := r.db.Get(&totalUsers, countQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Если пользователей нет
	if totalUsers == 0 {
		return []*entity.User{}, 0, nil
	}

	pageAmount := (totalUsers + pageSize - 1) / pageSize

	// Если страница выходит за диапазон — возвращаем пусто
	if page > pageAmount {
		return []*entity.User{}, pageAmount, nil
	}

	// Получаем пользователей
	if err := r.db.Select(&users, query, args...); err != nil {
		return nil, 0, fmt.Errorf("failed to get users: %w", err)
	}

	return users, pageAmount, nil
}

func (r *UserRepository) CreateUser(username, hashedPassword string) error {

	var userAlreadyExist bool

	alreadyExistQuery := `
		SELECT EXISTS (
			SELECT id
			FROM users
			WHERE username = $1
		)
	`

	err := r.db.Get(&userAlreadyExist, alreadyExistQuery, username)
	if err != nil {
		return fmt.Errorf("user with such username already exist: %s", err.Error())
	}

	createQuery := `
		INSERT INTO users(id,username,hashed_password)
		VALUES($1,$2,$3)
	`

	_, err = r.db.Exec(createQuery, uuid.NewString(), username, hashedPassword)
	if err != nil {
		return fmt.Errorf("cant create user : %s", err.Error())
	}

	return nil
}

func (r *UserRepository) GetUserById(id string) (user *entity.User, err error) {

	fmt.Println("GetUserById")
	fmt.Println("id is " + id)

	var userData entity.User

	query := `
		SELECT id,username, hashed_password
		FROM users
		WHERE id = $1
	`
	err = r.db.Get(&userData, query, id)
	if err != nil {
		return nil, fmt.Errorf("cant get user : %s", err.Error())
	}

	return &userData, nil

}

func (r *UserRepository) DeleteUser(id string) error {

	var userAmount int

	usersMoreThanOne := `
		SELECT COUNT(id)
		FROM users
	`

	err := r.db.Get(&userAmount, usersMoreThanOne)
	if err != nil {
		return err
	}

	if userAmount == 1 {
		return fmt.Errorf("only one user cant delete")
	}

	var exists bool

	suchUserExist := `
		SELECT EXISTS (
			SELECT 1
			FROM users
			WHERE id = $1
		)
	`

	err = r.db.Get(&exists, suchUserExist, id)
	if err != nil {
		return fmt.Errorf("such user dosnt exist : %s", err.Error())
	}

	deleteUser := `
		DELETE FROM users
		WHERE id = $1
	`

	_, err = r.db.Exec(deleteUser, id)
	if err != nil {
		return fmt.Errorf("cant delete user : %s", err.Error())
	}

	return nil
}

func (r *UserRepository) UpdatePassword(id, hashedPassword string) error {

	query := `
		UPDATE users
		SET hashed_password = $1
		WHERE id = $2
	`

	_, err := r.db.Exec(query, hashedPassword, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) UpdateUsername(id, username string) error {
	query := `
		UPDATE users
		SET username = $1
		WHERE id = $2
	`

	_, err := r.db.Exec(query, username, id)
	if err != nil {
		return err
	}

	return nil
}
