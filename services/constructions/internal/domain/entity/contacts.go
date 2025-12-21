package entity

import "time"

type ContactsEmailSetting struct {
	ID        string     `db:"id" json:"id"`
	Email     string     `db:"email" json:"email"`
	CreatedAt *time.Time `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at,omitempty"`
}

type ContactNumber struct {
	ID        string     `db:"id" json:"id"`
	Phone     string     `db:"phone" json:"phone"`
	Label     string     `db:"label" json:"label"`
	SortOrder int        `db:"sort_order" json:"sortOrder"`
	CreatedAt *time.Time `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at,omitempty"`
}

type ContactAddress struct {
	ID        string     `db:"id" json:"id"`
	Title     string     `db:"title" json:"title"`
	Address   string     `db:"address" json:"address"`
	Lat       float64    `db:"lat" json:"lat"`
	Lon       float64    `db:"lon" json:"lon"`
	SortOrder int        `db:"sort_order" json:"sortOrder"`
	CreatedAt *time.Time `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at,omitempty"`
}
