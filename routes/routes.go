package routes

import (
	"Hack4Change/database"
	"Hack4Change/handlers"
	"Hack4Change/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine, dbConn *database.PostQreSQLCon) {
	// Tested
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})
	router.GET("/create-tables", func(c *gin.Context) {
		handlers.CreateTables(c, dbConn)
	})
	// Auth APIs
	//Tested
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", func(c *gin.Context) {
			handlers.Register(c, dbConn)
		})
		authGroup.POST("/login", func(c *gin.Context) {
			handlers.Login(c, dbConn)
		})
	}

	// Space APIs
	spaceGroup := router.Group("/space")
	spaceGroup.Use(middleware.AuthMiddleware())
	{
		//Tested
		spaceGroup.GET("/:id/get-files", func(c *gin.Context) {
			handlers.FetchFilesByProjectId(c, dbConn)
		})
		spaceGroup.GET("/:id/details", func(c *gin.Context) {
			handlers.FetchFilesAndFoldersByProjectId(c, dbConn)
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
		//Tested
		spaceGroup.POST("/create-space", func(c *gin.Context) {
			handlers.CreateProject(c, dbConn)
		})
		//Tested
		spaceGroup.GET("/details", func(c *gin.Context) {
			handlers.FetchProjectsByUserId(c, dbConn)

		})
		spaceGroup.GET("/:id/structure", func(c *gin.Context) {
			handlers.GetProjectStructureHandler(c, dbConn)
		})
	}
	userGroup := router.Group("/user")
	userGroup.Use(middleware.AuthMiddleware())
	{
		//Tested
		userGroup.GET("/profile", func(c *gin.Context) {
			handlers.FetchUserData(c, dbConn)
		})
		//Tested
		userGroup.POST("/update-socials", func(c *gin.Context) {
			handlers.UpdateUserProfile(c, dbConn)
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
	//Tested
	deleteGroup := router.Group("/delete")
	deleteGroup.Use(middleware.AuthMiddleware())
	{
		deleteGroup.GET("/:name", func(c *gin.Context) {
			handlers.DeleteTables(c, dbConn)
		})
	}
}
