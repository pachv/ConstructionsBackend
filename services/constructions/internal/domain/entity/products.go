package entity

import "time"

type CatalogCategory struct {
	ID        string     `db:"id" json:"id"`
	Title     string     `db:"title" json:"title"`
	Slug      string     `db:"slug" json:"slug"`
	ImagePath *string    `db:"image_path" json:"imagePath,omitempty"` // filename
	CreatedAt *time.Time `db:"created_at" json:"createdAt,omitempty"`
}

type CatalogSection struct {
	ID    string `db:"id" json:"id"`
	Title string `db:"title" json:"title"`
	Slug  string `db:"slug" json:"slug"`

	// У тебя в таблице catalog_sections НЕТ parent_category_slug.
	// Если реально нужно — это должно приходить из join таблицы,
	// либо убери это поле из sqlx выборок.
	ParentCategorySlug string `db:"parent_category_slug" json:"parentCategorySlug,omitempty"`

	ImagePath *string    `db:"image_path" json:"imagePath,omitempty"` // filename
	CreatedAt *time.Time `db:"created_at" json:"createdAt,omitempty"`
}

type CatalogProduct struct {
	ID    string `db:"id" json:"id"`
	Title string `db:"title" json:"title"`
	Slug  string `db:"slug" json:"slug"`

	CategorySlug string `db:"category_slug" json:"categorySlug"`
	SectionSlug  string `db:"section_slug" json:"sectionSlug"`

	Brand string `db:"brand" json:"brand"`
	Type  string `db:"type" json:"type"`

	Price    int  `db:"price" json:"price"`
	OldPrice *int `db:"old_price" json:"oldPrice,omitempty"`
	InStock  bool `db:"in_stock" json:"inStock"`
	// Badges НЕ лежат в catalog_products, поэтому sqlx сюда не заполнит.
	// Заполняй отдельно (join/агрегация) или оставляй как есть для JSON.
	Badges []string `json:"badges,omitempty"`

	SalePercent *int    `db:"sale_percent" json:"salePercent,omitempty"`
	ImagePath   *string `db:"image_path" json:"imagePath,omitempty"` // filename

	CreatedAt *time.Time `db:"created_at" json:"createdAt,omitempty"`
}

type ProductBadge struct {
	ID        string    `db:"id" json:"id"`
	Code      string    `db:"code" json:"code"`
	Title     string    `db:"title" json:"title"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}
