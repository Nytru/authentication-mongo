package interfaces

import "simpleserver/entities"

type UserRepository interface {
	AddMany(users ...entities.User) error
	GetMany(ids ...int) []entities.User
	GetManyLogins(logins ...string) []entities.User
	GetAll() []entities.User
	UpdateMany(users ...entities.User) error
	DeleteMany(ids ...int) error
	DeleteAll() error
}
