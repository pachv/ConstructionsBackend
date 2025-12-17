package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/is_backend/services/admin/internal/domain/entity"
	"github.com/jmoiron/sqlx"
)

const (
	monthInHours = time.Hour * 24 * 30
)

type SessionRepository struct {
	db *sqlx.DB
}

func NewSessionRepository(db *sqlx.DB) *SessionRepository {
	return &SessionRepository{
		db: db,
	}
}

func (r *SessionRepository) GetSessionBySessionId(sessionId string) (*entity.UserSession, error) {

	var userSession entity.UserSession

	selectSessionQuery := `
		SELECT id,user_name,user_id,created_at,expires_at
		FROM user_sessions
		WHERE id = $1
		LIMIT 1;
	`

	err := r.db.Get(&userSession, selectSessionQuery, sessionId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session doesnt exist")
		}
	}

	return &userSession, nil
}

func (r *SessionRepository) CreateSession(username string, userId string) (sessionId string, err error) {

	sessionId = uuid.NewString()
	createdAt := time.Now()

	expiresAt := createdAt.Add(monthInHours)

	query := `
		INSERT INTO user_sessions(id,user_id,user_name,created_at,expires_at)
		VALUES($1,$2,$3,$4,$5)
	`

	_, err = r.db.Exec(query, sessionId, userId, username, createdAt, expiresAt)
	if err != nil {
		return "", fmt.Errorf("cant create session : %s", err.Error())
	}

	return
}

func (r *SessionRepository) DeleteSession(id string) error {

	query := `
		DELETE FROM user_sessions
		WHERE id = $1
	`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("cant delete session : " + err.Error())
	}

	return nil
}
