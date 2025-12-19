package entity

type SiteSectionSummary struct {
	ID         string `json:"id" db:"id"`
	Title      string `json:"title" db:"title"`
	Label      string `json:"label" db:"label"`
	Slug       string `json:"slug" db:"slug"`
	Image      string `json:"image" db:"image_url"`
	HasGallery bool   `json:"hasGallery" db:"has_gallery"`
	HasCatalog bool   `json:"hasCatalog" db:"has_catalog"`
}

type SiteSection struct {
	ID              string               `json:"id"`
	Title           string               `json:"title"`
	Label           string               `json:"label"`
	Slug            string               `json:"slug"`
	Image           string               `json:"image"`
	AdvantegesText  string               `json:"advantegesText"`
	AdvantegesArray []string             `json:"advantegesArray"`
	HasGallery      bool                 `json:"hasGallery"`
	HasCatalog      bool                 `json:"hasCatalog"`
	Gallery         []SiteSectionGallery `json:"gallery,omitempty"`
	Catalog         *SiteSectionCatalog  `json:"catalog,omitempty"`
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
	ImageURL   string                `json:"imageUrl" db:"image_url"`
	SortOrder  int                   `json:"sortOrder" db:"sort_order"`
	Badges     []string              `json:"badges,omitempty"`
	Specs      []SiteSectionItemSpec `json:"specs,omitempty"`
}

type SiteSectionItemSpec struct {
	Key   string `json:"key" db:"key"`
	Value string `json:"value" db:"value"`
}
