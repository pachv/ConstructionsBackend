package entity

import "time"

type CatalogCategory struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Slug      string     `json:"slug"`
	ImagePath *string    `json:"image_path,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

type CatalogSection struct {
	ID                 string     `json:"id"`
	Title              string     `json:"title"`
	Slug               string     `json:"slug"`
	ParentCategorySlug string     `json:"parentCategorySlug"`
	ImagePath          *string    `json:"image_path,omitempty"`
	CreatedAt          *time.Time `json:"created_at,omitempty"`
}
