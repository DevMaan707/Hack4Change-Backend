package handlers

import (
	"Hack4Change/database"
	"Hack4Change/helpers"

	"Hack4Change/models"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Login(c *gin.Context, db *database.PostQreSQLCon) {
	var login models.Login
	if err := c.BindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	hashedPassword, err := db.FetchHashedPassword(login.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	check := helpers.CheckPasswordHash(hashedPassword, login.Password)
	if !check {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password"})
		return
	}
	userID, err := db.FetchUserIdByEmail(login.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	token, err := helpers.GenerateJWT(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
func Register(c *gin.Context, dbConn *database.PostQreSQLCon) {

	var payload models.CreateAccountReq
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if payload.Password != payload.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
		return
	}

	userID := uuid.New().String()
	passwordHash, err := helpers.HashPassword(payload.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.UserDetails{
		ID:        userID,
		Username:  payload.Username,
		Email:     payload.Email,
		Phone:     payload.Phone,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = dbConn.InsertUser(user, passwordHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert user into database"})
		return
	}

	token, err := helpers.GenerateJWT(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func CreateProject(c *gin.Context, dbCon *database.PostQreSQLCon) {
	var payload models.CreateProjectReq
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userId, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	projectID := uuid.New().String()

	project := models.ProjectDetails{
		ProjectID:          projectID,
		ProjectName:        payload.ProjectName,
		ProjectDescription: payload.ProjectDescription,
		OwnerID:            userId.(string),
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
	err := dbCon.InsertProject(project)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func CreateFile(c *gin.Context, dbCon *database.PostQreSQLCon) {
	var payload models.CreateFileReq
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}
	fileID := uuid.New().String()
	file := models.File{
		ProjectID:      payload.ProjectID,
		FileName:       payload.FileName,
		ID:             fileID,
		ParentFolderId: payload.ParentFolderId,
	}
	err := dbCon.InsertFile(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
func CreateFolder(c *gin.Context, dbCon *database.PostQreSQLCon) {
	var payload models.CreateFolderReq
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	folderID := uuid.New().String()

	folder := models.Folder{
		ProjectID:      payload.ProjectID,
		ParentFolderId: payload.ParentFolderId,
		ID:             folderID,
		FolderName:     payload.FolderName,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err := dbCon.InsertFolder(folder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func SaveFileContent(c *gin.Context, db *database.PostQreSQLCon) {
	var req models.SaveFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.SaveContent(req.Content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "File content saved successfully"})
}

func FetchFilesByProjectId(c *gin.Context, db *database.PostQreSQLCon) {
	projectID := c.Param("id")

	fileDetails, err := db.FetchFilesByProjectId(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": fileDetails})

}
func FetchFoldersByProjectId(c *gin.Context, db *database.PostQreSQLCon) {
	projectID := c.Param("id")

	folderDetails, err := db.FetchFoldersByProjectId(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": folderDetails})

}
func FetchProjectsByUserId(c *gin.Context, db *database.PostQreSQLCon) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}
	projectDetails, err := db.FetchProjectsByUserId(userId.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": projectDetails})
}
func FetchUserData(c *gin.Context, db *database.PostQreSQLCon) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}
	userDetails, err := db.FetchProjectsByUserId(userId.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": userDetails})
}
func UserProfile(c * gin.Context, db *database.PostQreSQLCon){
	
}
