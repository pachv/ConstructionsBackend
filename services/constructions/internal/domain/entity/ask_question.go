package entity

type AskQuestion struct {
	Id      string `db:"id"`
	Message string `db:"message"`
	Name    string `db:"phone_number"`
	Email   string `db:"email"`
	Product string `db:"product"`
}
