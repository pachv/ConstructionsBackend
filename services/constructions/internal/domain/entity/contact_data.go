package entity

type ContactData struct {
	Id    string `db:"id"`
	Email string `db:"contact_email"`
}
