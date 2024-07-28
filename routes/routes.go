package routes

import (
	"Hack4Change/database"
	"Hack4Change/handlers"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine, dbConn *database.PostQreSQLCon) {
	// Test route
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
		fmt.Println("Success")
	})

	// Space APIs
	router.GET("/space/:id/get-files", func(c *gin.Context) {
		// Add your handler logic here
	})

	router.POST("/space/:id/create-folder", func(c *gin.Context) {
		// Add your handler logic here
	})

	router.POST("/space/:id/create-file", func(c *gin.Context) {
		// Add your handler logic here
	})

	router.GET("/space/activate-space", func(c *gin.Context) {
		// Add your handler logic here
	})

	router.POST("/space/create-space", func(c *gin.Context) {
		// Add your handler logic here
	})

	// Auth APIs
	router.POST("/auth/register", func(c *gin.Context) {
		handlers.Register(c, dbConn)
	})

	router.GET("/auth/login", func(c *gin.Context) {
		// Add your handler logic here
	})
}
