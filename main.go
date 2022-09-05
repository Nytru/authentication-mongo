package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"simpleserver/entities"
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

func greet(w http.ResponseWriter, r *http.Request) {
	var page = getPage("pages/test_page.html")
	io.Copy(w, page)
}

func auth(w http.ResponseWriter, r *http.Request) {
	var page = getPage("pages/test_auth_page.html")
	io.Copy(w, page)
}

func reg(w http.ResponseWriter, r *http.Request) {
	var page = getPage("pages/test_reg_page.html")
	io.Copy(w, page)
}

func license(w http.ResponseWriter, r *http.Request) {
	var page = getPage("pages/license.html")
	io.Copy(w, page)
}

func info(w http.ResponseWriter, r *http.Request) {
	var page = getPage("pages/test_greet_page.html")
	io.Copy(w, page)
}

func selphone(w http.ResponseWriter, r *http.Request) {
	var page = getPage("pages/test_selphone_page.html")
	io.Copy(w, page)
}

func server(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var ac = entities.NewAccessProvider()
		w.Write([]byte(fmt.Sprint(ac)))
	}
	w.WriteHeader(http.StatusNotFound)
}

func initHandlers() {
	http.HandleFunc("/", greet)
	http.HandleFunc("/login", auth)
	http.HandleFunc("/reg", reg)
	http.HandleFunc("/license", license)
	http.HandleFunc("/info", info)
	http.HandleFunc("/check", server)
}

func main() {
	logger.Info("Got envs")
	defer errHandler()
	initHandlers()

	logger.Info("Server started")
	_ = http.ListenAndServe(":8080", nil)
}

func errHandler() {
	if r := recover(); r != nil {
		logger.Info(fmt.Sprint("Exited with error: ", r))
	} else {
		logger.Info("Exit with no errors")
	}
}

func getPage(adr string) io.Reader {
	var file, err = os.Open(adr)
	defer file.Close()
	if err != nil {
		return nil
	}
	return file
}
