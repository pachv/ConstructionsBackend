package entity

type User struct {
	Id             string `json:"id" db:"id"`
	Username       string `json:"username" db:"username"`
	HashedPassword string `json:"password" db:"hashed_password"`
}
