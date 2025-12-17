package entity

type Sertificate struct {
	Name    string `db:""`
	Slug    string `db:"sert_slug"`
	FileUrl string `db:"file_url"`
}
