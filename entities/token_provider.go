package entities

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type tokenProvider struct {
	user        User
	key         string
	token       jwt.Token
	accessToken string
}

func NewTokenProvider(user User, key string) *tokenProvider {
	return &tokenProvider{user: user, key: key}
}

func Sign(key any) string {
	var token = jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.StandardClaims{
		Audience:  "j",
		ExpiresAt: time.Now().Add(time.Minute).Unix(),
		Id:        "78",
	})
	var signedToken, err = token.SignedString([]byte(fmt.Sprint(key)))
	if err != nil {
		return err.Error()
	}
	return signedToken
}
