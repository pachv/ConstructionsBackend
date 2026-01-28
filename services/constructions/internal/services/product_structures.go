package services

// ===== Categories =====

type AdminCreateCategoryReq struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Slug      string `json:"slug"`
	ImagePath string `json:"imagePath"`
}

type AdminUpdateCategoryReq struct {
	Title     string `json:"title"`
	Slug      string `json:"slug"`
	ImagePath string `json:"imagePath"`
}

// ===== Sections =====

type AdminCreateSectionReq struct {
	ID                 string `json:"id"`
	Title              string `json:"title"`
	Slug               string `json:"slug"`
	ParentCategorySlug string `json:"parentCategorySlug"`
	ImagePath          string `json:"imagePath"`
}

type AdminUpdateSectionReq struct {
	Title              string `json:"title"`
	Slug               string `json:"slug"`
	ParentCategorySlug string `json:"parentCategorySlug"`
	ImagePath          string `json:"imagePath"` // пусто => не менять
}

// ===== Products =====

type AdminCreateProductReq struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Slug         string   `json:"slug"`
	CategorySlug string   `json:"categorySlug"`
	SectionSlug  string   `json:"sectionSlug"`
	Brand        string   `json:"brand"`
	Type         string   `json:"type"`
	Price        int      `json:"price"`
	InStock      bool     `json:"inStock"`
	ImagePath    string   `json:"imagePath"`
	Badges       []string `json:"badges"`
	Discount     string   `json:"discount"` // sale_20 или ""
}

type AdminUpdateProductReq struct {
	Title        string   `json:"title"`
	Slug         string   `json:"slug"`
	CategorySlug string   `json:"categorySlug"`
	SectionSlug  string   `json:"sectionSlug"`
	Brand        string   `json:"brand"`
	Type         string   `json:"type"`
	Price        int      `json:"price"`
	InStock      bool     `json:"inStock"`
	ImagePath    string   `json:"imagePath"`
	Badges       []string `json:"badges"`
	Discount     string   `json:"discount"` // sale_20 => применить, "" => убрать
}
