package entity

type AskQuestion struct {
	Id      string  `db:"id" json:"id"`
	Message string  `db:"message" json:"message"`
	Name    string  `db:"name" json:"name"`
	Phone   string  `db:"phone" json:"phone"`
	Email   *string `db:"email" json:"email,omitempty"`
	Product *string `db:"product" json:"product,omitempty"`
	Consent bool    `db:"consent" json:"consent"`
}
