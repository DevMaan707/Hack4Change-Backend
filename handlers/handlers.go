package handlers

import (
	"Hack4Change/helpers"
	"Hack4Change/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Login(c *gin.Context) {

}
func Register(c *gin.Context) {

	var payload models.CreateAccountReq
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if payload.Password != payload.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
	}

	userId := uuid.New().String()
	passwordHash, err := helpers.HashPassword(payload.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

}
