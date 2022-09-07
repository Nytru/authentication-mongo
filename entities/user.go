package entities

import accessLevel "simpleserver/enums"

type User struct {
	// aka login
	Name         string                  `bson:"name"`
	Id           int                     `bson:"_id"`
	Password     string                  `bson:"password"`
	RefreshToken RefreshToken            `bson:"RefreshToken"`
	Level        accessLevel.AccessLevel `bson:"level"`
}
