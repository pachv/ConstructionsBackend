package entity

type Section struct {
	Id                string `db:"id"`
	Name              string `db:"section_name"`
	SectionPictureURL string `db:"section_picture"`
	SectionText       string `db:"section_text"`
}
