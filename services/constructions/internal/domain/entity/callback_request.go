package entity

type CallbackRequest struct {
	Id      string `db:"id" json:"id"`
	Name    string `db:"name" json:"name"`
	Phone   string `db:"phone" json:"phone"`
	Consent bool   `db:"consent" json:"consent"`
}
