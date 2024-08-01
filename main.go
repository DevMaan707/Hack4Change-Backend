package main

import (
	"log"
	"log/slog"

	db "Hack4Change/database"
	routes "Hack4Change/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file", slog.String("error", err.Error()))
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbConn, err := db.ConnectPostgreSQL()
	if err != nil {
		slog.Error("Error with postgresql", slog.String("error", err.Error()))
		log.Fatalf("Error with postgresql: %v", err)
	}

	router := gin.Default()
	routes.InitializeRoutes(router, dbConn)
	router.Run(":8080")
}
