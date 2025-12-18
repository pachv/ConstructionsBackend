package entity

import "time"

type GalleryCategory struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Slug      string     `json:"slug"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

type GalleryPhoto struct {
	ID           string     `json:"id"`
	CategorySlug string     `json:"category_slug"`
	Alt          string     `json:"alt"`
	ImagePath    string     `json:"image"`
	SortOrder    int        `json:"sort_order"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
}
