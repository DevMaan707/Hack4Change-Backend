package routes

import (
	"Hack4Change/database"
	"Hack4Change/handlers"
	"Hack4Change/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine, dbConn *database.PostQreSQLCon) {
	// Test route
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Auth APIs
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", func(c *gin.Context) {
			handlers.Register(c, dbConn)
		})
		authGroup.GET("/login", func(c *gin.Context) {
			handlers.Login(c, dbConn)
		})
	}

	// Space APIs
	spaceGroup := router.Group("/space")
	spaceGroup.Use(middleware.AuthMiddleware())
	{
		spaceGroup.GET("/:id/get-files", func(c *gin.Context) {
			handlers.FetchFilesByProjectId(c, dbConn)
		})

		spaceGroup.POST("/:id/create-folder", func(c *gin.Context) {
			handlers.CreateFolder(c, dbConn)
		})

		spaceGroup.POST("/:id/create-file", func(c *gin.Context) {
			handlers.CreateFile(c, dbConn)
		})

		spaceGroup.POST("/:id/save-file", func(c *gin.Context) {
			handlers.SaveFileContent(c, dbConn)
		})

		spaceGroup.POST("/create-space", func(c *gin.Context) {
			handlers.CreateProject(c, dbConn)
		})
	}
	userGroup := router.Group("/user")
	{
		userGroup.GET("/profile", func(c *gin.Context) {
			handlers.FetchUserData(c, dbConn)
		})

		academyGroup := userGroup.Group("/academy")
		{
			academyGroup.GET("/dashboard", func(c *gin.Context) {
				handlers.Dashboard(c, dbConn)
			})
			academyGroup.GET("/:id/status", func(c *gin.Context) {
				handlers.Status(c, dbConn)
			})

			academyGroup.POST("/:id/status/:qid/submit", func(c *gin.Context) {
				handlers.SubmitSol(c, dbConn)

			})
			academyGroup.POST("/generate", func(ctx *gin.Context) {
				handlers.GenerateSkill(ctx, dbConn)
			})
		}

	}
}
