package routes

import (
	"Hack4Change/database"
	"Hack4Change/handlers"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine, dbConn *database.PostQreSQLCon) {

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
		fmt.Println("Success")
	})

	// Space APIs
	router.GET("/space/:id/get-files", func(c *gin.Context) {
		// Add your handler logic here
	})

	router.POST("/space/:id/create-folder", func(c *gin.Context) {
		handlers.CreateFolder(c, dbConn)
	})

	router.POST("/space/:id/create-file", func(c *gin.Context) {
		handlers.CreateFile(c, dbConn)
	})

	router.POST("/space/create-space", func(c *gin.Context) {
		handlers.CreateProject(c, dbConn)
	})

	// Auth APIs
	router.POST("/auth/register", func(c *gin.Context) {
		handlers.Register(c, dbConn)
	})

	router.GET("/auth/login", func(c *gin.Context) {
		handlers.Login(c, dbConn)
	})
}
