package handleFuncs

import (
	"fmt"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"simpleserver/consts"
	"simpleserver/entities"
	accessLevel "simpleserver/enums"
	"time"
)

var logger *zap.Logger

func init() {
	var err error
	err = godotenv.Load("envs/.env")
	if err != nil {
		panic(err)
	}

	var conf = zap.NewDevelopmentConfig()
	conf.EncoderConfig.TimeKey = "timestamp"
	conf.Encoding = "console"
	conf.OutputPaths = []string{"stdout", ".log"}

	logger, err = conf.Build()
	if err != nil {
		panic(err)
	}
}

func Greet(w http.ResponseWriter, r *http.Request) {
	var err = loadPage(w, "pages/test_page.html")
	if err != nil {
		pageErrorHandler(w, err)
		return
	}
}

func Auth(w http.ResponseWriter, r *http.Request) {
	var accValue, accOk = r.Header[consts.AccessTokenName]
	var refValue, refOk = r.Header[consts.RefreshTokenName]
	logger.Info(fmt.Sprintf("Trying login with tokens:\nacc: %s\nref: %s", accValue, refValue))
	if accOk && refOk {
		var accessProvider = entities.NewAccessProvider(entities.User{
			Name:     "",
			Id:       0,
			Password: "",
			RefreshToken: entities.RefreshToken{
				Token:     refValue[0],
				ExpiresAt: 0,
			},
			Level: 0,
		})
		var level = accessProvider.CheckTokens(accValue[0], refValue[0])
		logger.Info(fmt.Sprintf("Auth checked tokens. Level: %s", level))
		switch level {
		case accessLevel.Admin:
			http.Redirect(w, r, "/info", http.StatusContinue)
			return
		case accessLevel.User:
			http.Redirect(w, r, "/info", http.StatusContinue)
			return
		}
	}
	var err = loadPage(w, "pages/test_auth_page.html")
	if err != nil {
		pageErrorHandler(w, err)
		return
	}
}

//func Reg(w http.ResponseWriter, r *http.Request) {
//	var err = loadPage(w, "pages/test_reg_page.html")
//	if err != nil {
//		pageErrorHandler(w, err)
//		return
//	}
//}
//
//func License(w http.ResponseWriter, r *http.Request) {
//	var err = loadPage(w, "pages/license.html")
//	if err != nil {
//		pageErrorHandler(w, err)
//		return
//	}
//}

func Info(w http.ResponseWriter, r *http.Request) {
	logger.Info("I got:")

	for s, strings := range r.Header {
		logger.Info(fmt.Sprintf("name: %s value: %s", s, strings))
	}

	var accValue, accOk = r.Header[consts.AccessTokenName]
	var refValue, refOk = r.Header[consts.RefreshTokenName]
	logger.Info(fmt.Sprintf("Got info request with\n%s\naccess token and\n%s\nrefresh token\n",
		accValue, refValue))

	if accOk && refOk {
		var accessProvider = entities.NewAccessProvider(entities.User{
			RefreshToken: entities.RefreshToken{Token: refValue[0]},
		})
		var level = accessProvider.CheckTokens(accValue[0], refValue[0])
		logger.Info(fmt.Sprintf("info checked tokens with result: %s", level))
		switch level {
		case accessLevel.Admin:
			var err = loadPage(w, "pages/test_greet_page.html")
			if err != nil {
				pageErrorHandler(w, err)
			}
			return
		case accessLevel.User:
			var err = loadPage(w, "pages/test_greet_page.html")
			if err != nil {
				pageErrorHandler(w, err)
			}
			return
		}
	}
	http.Redirect(w, r, "/Auth", http.StatusUnauthorized)
}

func Check(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		logger.Info(fmt.Sprintf("Trying access with %q login and %q password",
			r.Header[consts.LoginName], r.Header[consts.PasswordName]))

		if len(r.Header[consts.PasswordName]) == 0 || len(r.Header[consts.LoginName]) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var login = r.Header[consts.LoginName][0]
		var password = r.Header[consts.LoginName][0]
		var accessProvider = entities.NewAccessProvider(entities.User{Name: login, Password: password})
		var pair = accessProvider.TryLog(login, password)
		logger.Info(fmt.Sprintf("got new tokens:\nacc: %s\nref:%s + %s\n", pair.AccessToken, pair.RefreshToken.Token, time.Duration(pair.RefreshToken.ExpiresAt)))
		if pair.IsValid() {
			logger.Info("got tokens:\n" + fmt.Sprintf("access: %s\nrefrsh: %s\n",
				pair.AccessToken, pair.RefreshToken.Token))

			logger.Info(fmt.Sprintf("exp at: %s\n", time.Duration(pair.RefreshToken.ExpiresAt)))
			var cookie = http.Cookie{
				Name:    consts.AccessTokenName,
				Value:   pair.AccessToken,
				Expires: time.Now().Add(time.Duration(time.Minute * 15)),
			}
			http.SetCookie(w, &cookie)
			cookie = http.Cookie{
				Name:    consts.RefreshTokenName,
				Value:   pair.RefreshToken.Token,
				Expires: time.Now().Add(time.Duration(pair.RefreshToken.ExpiresAt)),
			}
			http.SetCookie(w, &cookie)
			io.WriteString(w, "success")
			//http.Redirect(w, r, "/info", http.StatusContinue)
		} else {
			logger.Info(fmt.Sprintf("no user with login: %s and password: %s", login, password))
			w.WriteHeader(http.StatusUnauthorized)
			io.WriteString(w, "bad auth")
		}
		return
	}
	logger.Info(fmt.Sprintf("get not post check requst. Method: %s", r.Header))
	w.WriteHeader(http.StatusNotFound)
}

func loadPage(w http.ResponseWriter, adr string) error {
	var file, err = os.Open(adr)
	defer file.Close()
	if err != nil {
		return err
	}
	if _, err = io.Copy(w, file); err != nil {
		return err
	}
	return nil
}

func pageErrorHandler(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	logger.Info("Page cannot be laded with error: " + err.Error())
	io.WriteString(w, "Something went wrong try again later")
}
