package entities

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	accessLevel "simpleserver/enums"
	"strconv"
	"time"
)

type RefreshToken struct {
	Token     string
	ExpiresAt int64
}

type tokenProvider struct {
	user        User
	key         string
	accessToken string
	//RefreshToken RefreshToken
}

type LiteClaims struct {
	jwt.StandardClaims
	Name  string                  `json:"name,omitempty"`
	Id    int                     `json:"_id"`
	Level accessLevel.AccessLevel `json:"level"`
}

func NewTokenProvider(user User, key string) *tokenProvider {
	return &tokenProvider{user: user, key: key}
}

func (tc *tokenProvider) checkTokens(accToken, refToken string) accessLevel.AccessLevel {
	var token, _ = jwt.ParseWithClaims(accToken, &LiteClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tryGetEnv("KEY")), nil
	})
	if claims, ok := token.Claims.(*LiteClaims); ok && token.Valid {
		return claims.Level
	} else {
		return accessLevel.None
	}
}

func (tc *tokenProvider) GetNewTokens(user User) TokenPair {
	var err = tc.newAccess()
	if err != nil {
		return TokenPair{}
	}
	err = tc.newRefresh()
	if err != nil {
		return TokenPair{}
	}
	return TokenPair{
		AccessToken:  tc.accessToken,
		RefreshToken: &tc.user.RefreshToken,
	}
}

func (tc *tokenProvider) newAccess() error {
	var accessDuration, err = strconv.Atoi(tryGetEnv("ACCESS_DURATION"))

	var token = jwt.NewWithClaims(jwt.SigningMethodHS512, LiteClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(accessDuration)).Unix(),
		},
		Name: tc.user.Name,
		Id:   tc.user.Id,
	})
	tc.accessToken, err = token.SignedString([]byte(tc.key))
	return err
}

func (tc *tokenProvider) newRefresh() error {
	var token, er = bcrypt.GenerateFromPassword([]byte(tc.accessToken), bcrypt.DefaultCost)
	if er != nil {
		return er
	}
	tc.user.RefreshToken.Token = string(token) // refreshToken.Token = string(token)
	var refreshDuration, _ = strconv.Atoi(tryGetEnv("REFRESH_DURATION"))
	tc.user.RefreshToken.ExpiresAt = time.Now().Add(time.Duration(refreshDuration)).Unix()
	return nil
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
