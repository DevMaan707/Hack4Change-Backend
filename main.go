package main

import (
	db "Hack4Change/database"

	"github.com/gin-gonic/gin"
)

//code spaces:
//create-space(Project name , description, public/private, technology used)
//create-dir
//create-file
//compile
//runcode
//save-space
//list-all-spaces
//get-space

func main() {

	dynamoClient := db.ConnectDynamoDB()
	router := gin.Default()

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
		fmt.Println("Success")
	})
	router.GET(
		"/get-files", func(c *gin.Context){

		}
	)
	router.POST(
		"/create-folder", func(c *gin.Context){

		}
	)
	router.POST(
		"/create-file", func(c *gin.Context){

		}
	)

}
