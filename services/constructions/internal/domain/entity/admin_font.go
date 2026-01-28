package entity

import "time"

type AdminFont struct {
	ID        string    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	FilePath  string    `db:"file_path" json:"filePath"`
	Selected  bool      `db:"selected" json:"selected"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}
