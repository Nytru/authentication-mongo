package entities

import (
	"errors"
	"os"
	accessLevel "simpleserver/enums"
)

type accessProvider struct {
	mongoCrud *mongoUserRepository
	tokenProv *tokenProvider
}

func NewAccessProvider(user User) *accessProvider {
	var provider = new(accessProvider)
	provider.mongoCrud = NewMongoUserRepository(getConnectionOptions())
	provider.tokenProv = NewTokenProvider(user, tryGetEnv("KEY"))
	return provider
}

func (ac accessProvider) TryLog(login, password string) TokenPair {
	var user = ac.mongoCrud.GetManyLogins(login)
	if len(user) != 1 {
		return TokenPair{}
	}
	var tokenPair = ac.tokenProv.GetNewTokens(user[0])
	return tokenPair
}

func (ac accessProvider) CheckTokens(accToken, refToken string) accessLevel.AccessLevel {
	return ac.tokenProv.checkTokens(accToken, refToken)
}

func (ac accessProvider) TryRefresh(token any) {

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
