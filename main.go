package main

import (
	"fmt"
	"time"
	"strconv"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/motionwerkGmbH/cpo-backend-api/tools"
	_ "github.com/mattn/go-sqlite3"
)

var router *gin.Engine

func main() {

	// Configs
	config, err := tools.ReadConfig("api_config", map[string]interface{}{
		"port":     9090,
		"hostname": "localhost",
		"auth": map[string]string{
			"username": "user",
			"password": "pass",
		},
	})
	if err != nil {
		panic(fmt.Errorf("Error when reading config: %v\n", err))
	}

	// Gin Configuration
	if (config.GetString("environment")) == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router = gin.New()
	router.Use(gin.Recovery())

	InitializeRoutes()

	// Establish database connection
	tools.Connect("_theDb.db")



	// Serve 'em...
	server := &http.Server{
		Addr:           ":" + strconv.Itoa(config.GetInt("port")),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	server.SetKeepAlivesEnabled(false)
	server.ListenAndServe()

}
