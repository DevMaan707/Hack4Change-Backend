package main

import (
	"Hack4Change/database"
	db "Hack4Change/database"
	routes "Hack4Change/routes"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func main() {
	dynamoClient := db.ConnectDynamoDB()
	dbConn, err := database.ConnectPostgreSQL()
	if err != nil {
		slog.Error("Error with postgresql", "Error", err.Error())
	}
	router := gin.Default()

	routes.InitializeRoutes(router, dynamoClient, dbConn)

	router.Run(":8080")
}
