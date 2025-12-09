package main

import (
	"micro-site/api/routes"
	"micro-site/config"
	"micro-site/database"
)

/*
* Main entry point of the application
* @return void
 */
func main() {
	router := routes.SetupRouter()
	cfg := config.LoadConfig() //Load configuration from .env
	database.InitDB(cfg)       // init database

	defer database.CloseDB()

	// Start Scheduler
	// scheduler.StartCron()

	// Listen and Server in 0.0.0.0:8080
	router.Run(":8080")
}
