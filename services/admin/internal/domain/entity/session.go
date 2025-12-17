package entity

import "time"

type UserSession struct {
	Id        string    `db:"id"`
	UserId    string    `db:"user_id"`
	UserName  string    `db:"user_name"`
	CreatedAt time.Time `db:"created_at"`
	ExpiresAt time.Time `db:"expires_at"`
}
