package models

type User struct {
	LastName string
	Name string
}

type RegisterUserModel struct {
	Name string
	LastName string
	Email string
	Userame string
	Password string
	Code string
}
