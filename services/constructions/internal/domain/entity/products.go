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

type CatalogProduct struct {
	ID           string     `json:"id"`
	Title        string     `json:"title"`
	Slug         string     `json:"slug"`
	CategorySlug string     `json:"categorySlug"`
	SectionSlug  string     `json:"sectionSlug"`
	Brand        string     `json:"brand"`
	Type         string     `json:"type"`
	Price        int        `json:"price"`
	OldPrice     *int       `json:"oldPrice,omitempty"`
	InStock      bool       `json:"inStock"`
	Badges       []string   `json:"badges,omitempty"`
	SalePercent  *int       `json:"salePercent,omitempty"`
	ImagePath    *string    `json:"image,omitempty"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
}
