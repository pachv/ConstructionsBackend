package entity

import "time"

type Certificate struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	FilePath  string     `json:"file_path"` // ссылка
	CreatedAt *time.Time `json:"created_at,omitempty"`
}
