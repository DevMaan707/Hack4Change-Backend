package main

import (
	"log"
	"log/slog"

	db "Hack4Change/database"
	routes "Hack4Change/routes"

	"github.com/gin-gonic/gin"
	
)

func main() {
	dbConn, err := db.ConnectPostgreSQL()
	if err != nil {
		slog.Error("Error with postgresql", slog.String("error", err.Error()))
		log.Fatalf("Error with postgresql: %v", err)
	}

	router := gin.Default()
	routes.InitializeRoutes(router, dbConn)
	router.Run(":7563")
}
