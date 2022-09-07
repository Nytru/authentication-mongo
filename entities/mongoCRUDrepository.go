package entities

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	opts "go.mongodb.org/mongo-driver/mongo/options"
)

type mongoUserRepository struct {
	collection *mongo.Collection
}

func NewMongoUserRepository(options MongoConnectionOptions) *mongoUserRepository {
	var repos mongoUserRepository
	var client, _ = mongo.NewClient(opts.Client().ApplyURI(options.Connection_string))
	var err = client.Connect(context.TODO())
	if err != nil {
		return nil
	}
	repos.collection = client.Database(options.Database_Name).Collection(options.Collection_Name)
	return &repos
}

func (mr mongoUserRepository) AddMany(users ...User) error {
	var interfaceSlice = make([]interface{}, len(users))
	for i, v := range users {
		interfaceSlice[i] = v
	}
	var _, err = mr.collection.InsertMany(context.TODO(), interfaceSlice)
	return err
}

func (mr mongoUserRepository) GetMany(ids ...int) []User {
	var arr []User
	for _, v := range ids {
		var user User
		var cur = mr.collection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: v}})
		var err = cur.Decode(&user)
		if err != nil {
			return nil
		}
		arr = append(arr, user)
	}
	return arr
}

func (mr mongoUserRepository) GetManyLogins(logins ...string) []User {
	var arr []User
	for _, v := range logins {
		var user User
		var cur = mr.collection.FindOne(context.TODO(), bson.D{{Key: "name", Value: v}})
		var err = cur.Decode(&user)
		if err != nil {
			return nil
		}
		arr = append(arr, user)
	}
	return arr
}

func (mr mongoUserRepository) GetAll() []User {
	var arr []User
	var cur, _ = mr.collection.Find(context.TODO(), bson.D{})
	_ = cur.All(context.TODO(), &arr)
	return arr
}

func (mr mongoUserRepository) UpdateMany(users ...User) error {
	for _, v := range users {
		var _, err = mr.collection.ReplaceOne(context.TODO(), bson.D{{Key: "_id", Value: v.Id}}, v) // bson.D{{Key: "name", Value: v.Name}}
		if err != nil {
			return err
		}
	}
	return nil
}

func (mr mongoUserRepository) DeleteMany(ids ...int) error {
	for _, v := range ids {
		var _, err = mr.collection.DeleteOne(context.TODO(), bson.D{{Key: "_id", Value: v}})
		if err != nil {
			return err
		}
	}
	return nil
}

func (mr mongoUserRepository) DeleteAll() error {
	return mr.collection.Drop(context.TODO())
}
