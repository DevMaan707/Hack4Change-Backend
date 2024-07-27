package main

import (
	db "Hack4Change/database"
	routes "Hack4Change/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	dynamoClient := db.ConnectDynamoDB()
	router := gin.Default()

	routes.InitializeRoutes(router, dynamoClient)

	router.Run(":8080")
}
