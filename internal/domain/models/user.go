package models

type User struct {
	ID             int64
	Email          string
	Nickname       string
	PasswordHashed []byte
}
