package models

type Product struct {
	Name string
	Model string
	OwnerID int64
}

type ListProducts struct {
	Id int64
	Name string
	Model string
	OwnerID int64
}