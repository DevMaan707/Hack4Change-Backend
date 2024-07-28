package main

/*
add /user/:id/update-profile
*/
import (
	db "Hack4Change/database"
	routes "Hack4Change/routes"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func main() {

	dbConn, err := db.ConnectPostgreSQL()
	if err != nil {
		slog.Error("Error with postgresql", "Error", err.Error())
	}
	router := gin.Default()

	routes.InitializeRoutes(router, dbConn)

	router.Run(":8080")
}
