package main

import (
	"time"
	"strconv"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/motionwerkGmbH/cpo-backend-api/tools"
	_ "github.com/mattn/go-sqlite3"
	"github.com/motionwerkGmbH/cpo-backend-api/configs"
	"log"
)

var router *gin.Engine

func main() {

	// Configs
	Config := configs.Load()

	// Gin Configuration
	if (Config.GetString("environment")) == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router = gin.New()
	router.Use(gin.Recovery())

	InitializeRoutes()

	// Establish database connection
	tools.Connect("_theDb.db")

	log.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
	log.Println("Running on http://localhost:9090/api/v1/account/info")
	log.Println("Running on http://18.195.223.26:9090/api/v1/account/info")
	log.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")


	// Serve 'em...
	server := &http.Server{
		Addr:           ":" + strconv.Itoa(Config.GetInt("port")),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	server.SetKeepAlivesEnabled(false)
	server.ListenAndServe()


}
