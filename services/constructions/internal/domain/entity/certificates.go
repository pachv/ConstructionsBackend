package entity

import "time"

type Certificate struct {
	ID        string     `db:"id" json:"id"`
	Title     string     `db:"title" json:"title"`
	FilePath  string     `db:"file_path" json:"file_path"`
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
}
