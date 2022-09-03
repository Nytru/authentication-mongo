package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"simpleserver/entities"
	mongorepository "simpleserver/mongo_repository"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
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
	var file, err = os.Open("pages/test_page.html")
	defer file.Close()

	if err != nil {
		panic(err)
	}
	io.Copy(w, file)
}

func auth(w http.ResponseWriter, r *http.Request) {
	var file, err = os.Open("pages/test_auth_page.html")
	defer file.Close()
	if err != nil {
		panic(err)
	}
	io.Copy(w, file)
}

func reg(w http.ResponseWriter, r *http.Request) {
	var file, err = os.Open("pages/test_reg_page.html")
	defer file.Close()
	if err != nil {
		panic(err)
	}
	io.Copy(w, file)
}

func license(w http.ResponseWriter, r *http.Request) {
	var file, _ = os.Open("pages/license.html")
	defer file.Close()
	io.Copy(w, file)
}

func server(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		w.Write([]byte("hello"))
	}
	w.WriteHeader(http.StatusNotFound)
}

func getEnv(options *mongorepository.MongoConnectionOptions) {
	if v, ok := os.LookupEnv("MONGO_FULL_PASS"); ok {
		options.Connection_string = v
	}

	if v, ok := os.LookupEnv("DB_NAME"); ok {
		options.Database_Name = v
	}

	if v, ok := os.LookupEnv("COLLECTION_NAME"); ok {
		options.Collection_Name = v
	}
}

func main() {
	var options = mongorepository.MongoConnectionOptions{}
	getEnv(&options)

	var rep = mongorepository.NewMongoUserRepository(options)

	rep.AddMany(entities.User{Id: 12, Name: "igor'"})
	// rep.DeleteAll()

	return
	defer errHandler()
	http.HandleFunc("/", greet)
	http.HandleFunc("/login", auth)	
	http.HandleFunc("/reg", reg)
	http.HandleFunc("/license", license)

	http.HandleFunc("/chek", server)

	logger.Info("Server started")
	http.ListenAndServe(":8080", nil)
}

func errHandler() {
	if r := recover(); r != nil {
		logger.Info(fmt.Sprint("Exited with error: ", r))
	} else {
		logger.Info("Exit with no errors")
	}
}
