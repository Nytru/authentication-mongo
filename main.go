package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"net/http"
	"simpleserver/handleFuncs"
)

var logger *zap.Logger

func init() {
	var err error
	err = godotenv.Load("envs/.env")
	if err != nil {
		panic(err)
	}
	var conf = zap.NewDevelopmentConfig()
	conf.OutputPaths = []string{"stdout", ".log"}
	logger, err = conf.Build()
	if err != nil {
		panic(err)
	}
}

func initHandlers() {
	http.HandleFunc("/", handleFuncs.Greet)
	http.HandleFunc("/login", handleFuncs.Auth)
	//http.HandleFunc("/reg", handleFuncs.Reg)
	//http.HandleFunc("/license", handleFuncs.License)
	http.HandleFunc("/info", handleFuncs.Info)
	http.HandleFunc("/check", handleFuncs.Check)
}

func main() {
	logger.Info("Got envs")
	defer exitHandler()
	initHandlers()
	logger.Info("Server started")
	_ = http.ListenAndServe(":8080", nil)
	//var prov = entities.NewAccessProvider(entities.User{})
	//prov.CheckTokens("")
}

func exitHandler() {
	if r := recover(); r != nil {
		logger.Info(fmt.Sprint("Exited with error: ", r))
	} else {
		logger.Info("Exit with no errors")
	}
}
