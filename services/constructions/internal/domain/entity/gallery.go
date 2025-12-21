package entity

import "time"

type GalleryCategory struct {
	ID        string     `json:"id" db:"id"`
	Title     string     `json:"title" db:"title"`
	Slug      string     `json:"slug" db:"slug"`
	CreatedAt *time.Time `json:"created_at,omitempty" db:"created_at"`
}

type GalleryPhoto struct {
	ID           string     `json:"id" db:"id"`
	CategorySlug string     `json:"category_slug" db:"category_slug"`
	Alt          string     `json:"alt" db:"alt"`
	ImagePath    string     `json:"image" db:"image"`
	SortOrder    int        `json:"sort_order" db:"sort_order"`
	CreatedAt    *time.Time `json:"created_at,omitempty" db:"created_at"`
}
