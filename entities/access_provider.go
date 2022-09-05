package entities

import (
	"errors"
	"os"
)

type accessProvider struct {
	mongoCrud *mongoUserRepository
	tokenProv *tokenProvider
}

func NewAccessProvider() *accessProvider {
	var provider *accessProvider
	provider.mongoCrud = NewMongoUserRepository(getConnectionOptions())
	return provider
}

func getConnectionOptions() (options MongoConnectionOptions) {
	options.Connection_string = tryGetEnv("MONGO_FULL_PASS")
	options.Database_Name = tryGetEnv("DB_NAME")
	options.Collection_Name = tryGetEnv("COLLECTION_NAME")
	return options
}

func tryGetEnv(key string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	} else {
		panic(errors.New("env require: " + key))
	}
}
