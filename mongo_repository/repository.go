package mongorepository

import (
	"context"
	"simpleserver/entities"

	"go.mongodb.org/mongo-driver/mongo"
	opts "go.mongodb.org/mongo-driver/mongo/options"
)

type MongoUserRepository struct {
	collection *mongo.Collection
}

func NewMongoUserRepository(options MongoConnectionOptions) (*MongoUserRepository) {
	var repos MongoUserRepository
	var client, _ = mongo.NewClient(opts.Client().ApplyURI(options.Connection_string))
	var err = client.Connect(context.TODO())
	if err != nil {
		return nil
	}
	repos.collection = client.Database(options.Database_Name).Collection(options.Collection_Name)
	return &repos
}

func (mr MongoUserRepository)AddMany(users ...entities.User) error {
	var interfaceSlice = make([]interface{}, len(users))
	for i, v := range users {
		interfaceSlice[i] = v
	}
	var _, err = mr.collection.InsertMany(context.TODO(), interfaceSlice)
	return err
}

func (mr MongoUserRepository)GetMany(ids ...int) []entities.User {
	return []entities.User{}
}

func (mr MongoUserRepository)UpdateMany(users ...entities.User) {
	for _, v := range users {
		
	}
	mongo.
	mr.collection.
}

func (mr MongoUserRepository)DeleteMany(ids ...int) {

}

func (mr MongoUserRepository)DeleteAll() error {
	return mr.collection.Drop(context.TODO())
}