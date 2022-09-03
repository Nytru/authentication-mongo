package entities

// import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Name string `bson:"name"`
	Id   int `bson:"_id"`
}