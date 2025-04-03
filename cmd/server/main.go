package main

import (
	_ "hihand/docs"
	config "hihand/internal/configs/dev"
	"hihand/internal/router"
	"log"
)

var (
	logger = log.New(log.Writer(), "[server/main.go] ", log.LstdFlags|log.Lshortfile)
)

// @title Hihand API
// @version 1.0
// @description API For Order Service
// @host localhost:8080
// @BasePath /
func main() {
	_, cfgErr := config.Instance()
	if cfgErr != nil {
		logger.Println("Can not get config:", cfgErr)
	}

	app := router.NewRouter()

	app.Run(":8080")

	logger.Println("App Running!")
}
