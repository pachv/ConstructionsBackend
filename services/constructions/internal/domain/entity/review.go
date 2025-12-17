package entity

import "time"

type Review struct {
	Id         string    `db:"id" json:"id"`
	Name       string    `db:"name" json:"name"`
	Position   string    `db:"position" json:"position"`
	Text       string    `db:"text" json:"text"`
	Rating     int       `db:"rating" json:"rating"`
	ImagePath  string    `db:"image_path" json:"imagePath"`
	Consent    bool      `db:"consent" json:"consent"`
	CanPublish bool      `db:"can_publish" json:"canPublish"`
	CreatedAt  time.Time `db:"created_at" json:"createdAt"`
}
