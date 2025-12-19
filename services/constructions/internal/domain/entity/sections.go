package entity

type SiteSection struct {
	ID         string               `json:"id"`
	Title      string               `json:"title"`
	Slug       string               `json:"slug"`
	HasGallery bool                 `json:"hasGallery"`
	HasCatalog bool                 `json:"hasCatalog"`
	Gallery    []SiteSectionGallery `json:"gallery"`
	Catalog    *SiteSectionCatalog  `json:"catalog,omitempty"`
}

type SiteSectionSummary struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Slug       string `json:"slug"`
	HasGallery bool   `json:"hasGallery"`
	HasCatalog bool   `json:"hasCatalog"`
}

type SiteSectionGallery struct {
	ID        string `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	URL       string `json:"url" db:"url"`
	SortOrder int    `json:"sortOrder" db:"sort_order"`
}

type SiteSectionCatalog struct {
	Categories []SiteSectionCatalogCategory `json:"categories"`
	Items      []SiteSectionCatalogItem     `json:"items"`
}

type SiteSectionCatalogCategory struct {
	ID        string `json:"id" db:"id"`
	Title     string `json:"title" db:"title"`
	Slug      string `json:"slug" db:"slug"`
	SortOrder int    `json:"sortOrder" db:"sort_order"`
}

type SiteSectionCatalogItem struct {
	ID         string                `json:"id" db:"id"`
	CategoryID string                `json:"categoryId" db:"category_id"`
	Title      string                `json:"title" db:"title"`
	PriceRub   int                   `json:"priceRub" db:"price_rub"`
	ImageURL   *string               `json:"imageUrl" db:"image_url"`
	Badges     []string              `json:"badges"`
	Specs      []SiteSectionItemSpec `json:"specs"`
	SortOrder  int                   `json:"sortOrder" db:"sort_order"`
}

type SiteSectionItemSpec struct {
	Key   string `json:"key" db:"key"`
	Value string `json:"value" db:"value"`
}
