package handlers

import (
	"Hack4Change/database"
	"Hack4Change/helpers"
	"bytes"
	"encoding/json"
	"log/slog"

	"Hack4Change/models"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateTables(c *gin.Context, db *database.PostQreSQLCon) {
	if err := db.CreateTables(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "succcessful"})
}

func Login(c *gin.Context, db *database.PostQreSQLCon) {
	var login models.Login
	if err := c.BindJSON(&login); err != nil {
		slog.Error("Login failed: Invalid request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	hashedPassword, err := db.FetchHashedPassword(login.Email)
	if err != nil {
		slog.Error("Login failed: Error fetching hashed password", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	check := helpers.CheckPasswordHash(hashedPassword, login.Password)
	if !check {
		slog.Warn("Login failed: Invalid password", "email", login.Email)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password"})
		return
	}

	userID, err := db.FetchUserIdByEmail(login.Email)
	if err != nil {
		slog.Error("Login failed: Error fetching user ID", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := helpers.GenerateJWT(userID)
	if err != nil {
		slog.Error("Login failed: Error generating JWT", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Login successful", "userID", userID)
	c.JSON(http.StatusOK, gin.H{"token": token, "userID": userID})
}

func Register(c *gin.Context, dbConn *database.PostQreSQLCon) {
	var payload models.CreateAccountReq
	if err := c.ShouldBindJSON(&payload); err != nil {
		slog.Error("Registration failed: Invalid request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if payload.Password != payload.ConfirmPassword {
		slog.Warn("Registration failed: Passwords do not match", "email", payload.Email)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
		return
	}

	userID := uuid.New().String()
	passwordHash, err := helpers.HashPassword(payload.Password)
	if err != nil {
		slog.Error("Registration failed: Error hashing password", "error", err)
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
		slog.Error("Registration failed: Error inserting user into database", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert user into database"})
		return
	}

	token, err := helpers.GenerateJWT(userID)
	if err != nil {
		slog.Error("Registration failed: Error generating JWT", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	slog.Info("Registration successful", "userID", userID)
	c.JSON(http.StatusOK, gin.H{"token": token, "userId": userID})
}

func CreateProject(c *gin.Context, dbCon *database.PostQreSQLCon) {
	var payload models.CreateProjectReq
	if err := c.ShouldBindJSON(&payload); err != nil {
		slog.Error("CreateProject failed: Invalid request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, exist := c.Get("userID")
	if !exist {
		slog.Warn("CreateProject failed: Unauthorized access")
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
		slog.Error("CreateProject failed: Error inserting project", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Project created successfully", "projectID", projectID)
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func CreateFile(c *gin.Context, dbCon *database.PostQreSQLCon) {
	var payload models.CreateFileReq
	if err := c.ShouldBindJSON(&payload); err != nil {
		slog.Error("CreateFile failed: Invalid request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, exist := c.Get("userID")
	if !exist {
		slog.Warn("CreateFile failed: Unauthorized access")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
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
		slog.Error("CreateFile failed: Error inserting file", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Info("File created successfully", "fileID", fileID)
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func CreateFolder(c *gin.Context, dbCon *database.PostQreSQLCon) {
	var payload models.CreateFolderReq
	if err := c.ShouldBindJSON(&payload); err != nil {
		slog.Error("CreateFolder failed: Invalid request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, exist := c.Get("userID")
	if !exist {
		slog.Warn("CreateFolder failed: Unauthorized access")
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
		slog.Error("CreateFolder failed: Error inserting folder", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Folder created successfully", "folderID", folderID)
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func SaveFileContent(c *gin.Context, db *database.PostQreSQLCon) {
	var req models.SaveFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("SaveFileContent failed: Invalid request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.SaveContent(req.Content); err != nil {
		slog.Error("SaveFileContent failed: Error saving content", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Info("File content saved successfully")
	c.JSON(http.StatusOK, gin.H{"message": "File content saved successfully"})
}

func FetchFilesByProjectId(c *gin.Context, db *database.PostQreSQLCon) {
	projectID := c.Param("id")

	fileDetails, err := db.FetchFilesByProjectId(projectID)
	if err != nil {
		slog.Error("FetchFilesByProjectId failed: Error fetching files", "projectID", projectID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Files fetched successfully", "projectID", projectID)
	c.JSON(http.StatusOK, gin.H{"data": fileDetails})
}

func FetchFoldersByProjectId(c *gin.Context, db *database.PostQreSQLCon) {
	projectID := c.Param("id")

	folderDetails, err := db.FetchFoldersByProjectId(projectID)
	if err != nil {
		slog.Error("FetchFoldersByProjectId failed: Error fetching folders", "projectID", projectID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Folders fetched successfully", "projectID", projectID)
	c.JSON(http.StatusOK, gin.H{"data": folderDetails})
}

func FetchProjectsByUserId(c *gin.Context, db *database.PostQreSQLCon) {
	userId, exists := c.Get("userID")
	if !exists {
		slog.Warn("FetchProjectsByUserId failed: Unauthorized access")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	projectDetails, err := db.FetchProjectsByUserId(userId.(string))
	if err != nil {
		slog.Error("FetchProjectsByUserId failed: Error fetching projects", "userId", userId, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Projects fetched successfully", "userId", userId)
	c.JSON(http.StatusOK, gin.H{"data": projectDetails})
}

func FetchUserData(c *gin.Context, db *database.PostQreSQLCon) {
	userId, exists := c.Get("userID")
	if !exists {
		slog.Warn("FetchUserData failed: Unauthorized access")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	userDetails, err := db.FetchProjectsByUserId(userId.(string))
	if err != nil {
		slog.Error("FetchUserData failed: Error fetching user data", "userId", userId, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Info("User data fetched successfully", "userId", userId)
	c.JSON(http.StatusOK, gin.H{"data": userDetails})
}

func Dashboard(c *gin.Context, db *database.PostQreSQLCon) {
	userId, exists := c.Get("userID")
	if !exists {
		slog.Warn("Dashboard failed: Unauthorized access")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}

	payload, err := db.FetchSkillIdAndNameByUserID(userId.(string))
	if err != nil {
		slog.Error("Dashboard failed: Error fetching skill ID and name", "userId", userId, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Dashboard data fetched successfully", "userId", userId)
	c.JSON(http.StatusOK, gin.H{"data": payload})
}

func Status(c *gin.Context, db *database.PostQreSQLCon) {
	var payload models.StatusReq
	if err := c.ShouldBindJSON(&payload); err != nil {
		slog.Error("Status failed: Invalid request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong Request"})
		return
	}

	skillId := payload.SkillId

	payload_, err := db.FetchSkillIdAndNameByUserID(skillId)
	if err != nil {
		slog.Error("Status failed: Error fetching skill details", "skillId", skillId, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Skill details fetched successfully", "skillId", skillId)
	c.JSON(http.StatusOK, gin.H{"data": payload_})
}

func SubmitSol(c *gin.Context, db *database.PostQreSQLCon) {
	qid := c.Param("qid")
	skillid := c.Param("id")
	var payload models.SubmitSolReq
	if err := c.ShouldBindJSON(&payload); err != nil {
		slog.Error("SubmitSol failed: Invalid request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong Request"})
		return
	}

	if err := db.SubmitSolutionByQIDandSkillID(qid, skillid); err != nil {
		slog.Error("SubmitSol failed: Error submitting solution", "qid", qid, "skillid", skillid, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Solution submitted successfully", "qid", qid, "skillid", skillid)
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func GenerateSkill(c *gin.Context, db *database.PostQreSQLCon) {
	var payload models.GenerateSkillsReq
	if err := c.ShouldBindJSON(&payload); err != nil {
		slog.Error("GenerateSkill failed: Invalid request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		slog.Error("GenerateSkill failed: Error marshalling payload", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp, err := http.Post("http://localhost:5868/ai/generate-skill", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		slog.Error("GenerateSkill failed: Error sending POST request", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	var payloadRes models.SkillRes
	if err := json.NewDecoder(resp.Body).Decode(&payloadRes); err != nil {
		slog.Error("GenerateSkill failed: Error decoding response", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Skill generated successfully")
	c.JSON(http.StatusOK, gin.H{"data": payloadRes})
}

func BadgeHandler(c *gin.Context, db *database.PostQreSQLCon) {
	// Implement logging as needed when the function is implemented
}
