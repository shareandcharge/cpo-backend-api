package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/motionwerkGmbH/cpo-backend-api/configs"
	"github.com/motionwerkGmbH/cpo-backend-api/tools"
	"net/http"
	"strconv"
	"time"
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

	// allow all origins
	router.Use(cors.Default())

	InitializeRoutes()

	// Establish database connection
	tools.Connect("_theDb.db")

	log.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
	log.Println("Running on http://localhost:9090/api/v1/account/info")
	log.Println("Running on http://18.195.223.26:9090/api/v1/account/info")
	log.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")

	// Serve 'em...
	server := &http.Server{
		Addr:           Config.GetString("hostname") + ":" + strconv.Itoa(Config.GetInt("port")),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	server.SetKeepAlivesEnabled(false)
	server.ListenAndServe()

}
