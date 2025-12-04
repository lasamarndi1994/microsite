package main

import (
	"micro-site/api/routes"
	"micro-site/config"
	"micro-site/database"
)

func main() {
	router := routes.SetupRouter()
	cfg := config.LoadConfig() //Load configuration from .env
	database.InitDB(cfg)       // init database

	defer database.CloseDB()

	// Listen and Server in 0.0.0.0:8080
	router.Run(":8080")
}
