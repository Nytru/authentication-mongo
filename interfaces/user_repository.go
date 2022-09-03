package interfaces

import "simpleserver/entities"

type UserRepository interface {
	AddMany(users ...entities.User)
	GetMany(ids ...int) []entities.User
	UpdateMany(users ...entities.User)
	DeleteMany(ids ...int)
	DeleteAll()
}
