package entity

type ReturnCall struct {
	Id          string `db:"id"`
	Name        string `db:"name"`
	PhoneNumber string `db:"phone_number"`
}
