package entity

type User struct {
	Id             string `db:"id"`
	Surname        string `db:"surname"`
	LoginName      string `db:"login_name"`
	Fathername     string `db:"father_name"`
	Email          string `db:"email"`
	PhoneNumber    string `db:"phone_number"`
	HashedPassword string `db:"hashed_password"`
}
