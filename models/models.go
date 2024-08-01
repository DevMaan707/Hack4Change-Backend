package models

import (
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

type CreateAccountReq struct {
	Username        string `json:"username" validate:"required,min=3,max=32"`
	Email           string `json:"email" validate:"required,email"`
	Phone           string `json:"phone" validate:"omitempty,e164"`
	FirstName       string `json:"first_name" validate:"omitempty,min=1,max=32"`
	LastName        string `json:"last_name" validate:"omitempty,min=1,max=32"`
	Password        string `json:"password" validate:"required,min=8,max=128"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8,max=128,eqfield=Password"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Socials struct {
	GitHub      string `json:"github" validate:"omitempty,url"`
	LinkedIn    string `json:"linkedin" validate:"omitempty,url"`
	Instagram   string `json:"instagram" validate:"omitempty,url"`
	NoobsSocial string `json:"noobs_social" validate:"omitempty,url"`
}

type UserDetails struct {
	ID             string    `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	SocialAccounts Socials   `json:"social_accounts"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Badges         []Badge   `json:"badges"`
}

type CreateProjectReq struct {
	ProjectName        string `json:"project_name" validate:"required,min=3,max=50"`
	ProjectDescription string `json:"project_description" validate:"omitempty,max=255"`
}

type ProjectDetails struct {
	ProjectID          string    `json:"project_id"`
	OwnerID            string    `json:"owner_id"`
	ProjectName        string    `json:"project_name"`
	ProjectDescription string    `json:"project_description"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type File struct {
	ID             string    `json:"id"`
	ProjectID      string    `json:"project_id"`
	ParentFolderId string    `json:"parent_folder_id"`
	FileName       string    `json:"file_name" validate:"required,min=1,max=255"`
	FileContent    string    `json:"file_content" validate:"required"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Folder struct {
	ID             string `json:"id"`
	ProjectID      string `json:"project_id"`
	FolderName     string `json:"folder_name" validate:"required,min=1,max=255"`
	ParentFolderId string `json:"parent_folder_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateAccountDetails struct {
	ID        string    `json:"id"`
	Username  string    `json:"username" validate:"required,min=3,max=32"`
	Email     string    `json:"email" validate:"required,email"`
	Phone     string    `json:"phone" validate:"omitempty,e164"`
	FirstName string    `json:"first_name" validate:"omitempty,min=1,max=32"`
	LastName  string    `json:"last_name" validate:"omitempty,min=1,max=32"`
	Password  string    `json:"password" validate:"required,min=8,max=128"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateFileReq struct {
	ProjectID      string `json:"project_id" validate:"required"`
	FileName       string `json:"file_name" validate:"required,min=1,max=255"`
	FileContent    string `json:"file_content" validate:"required"`
	ParentFolderId string `json:"parent_folder_id"`
}

type CreateFolderReq struct {
	ProjectID      string `json:"project_id" validate:"required"`
	FolderName     string `json:"folder_name" validate:"required,min=1,max=255"`
	ParentFolderId string `json:"parent_folder_id"`
}

type FolderDetails struct {
	ID             string    `json:"id"`
	ProjectID      string    `json:"project_id"`
	FolderName     string    `json:"folder_name"`
	ParentFolderId string    `json:"parent_folder_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
type SaveFileRequest struct {
	ProjectID uuid.UUID `json:"projectId"`
	FileID    uuid.UUID `json:"fileId"`
	Content   string    `json:"content"`
}

type Badge struct {
	Name        string   `json:"name"`
	Topics      []string `json:"topics"`
	Topics_done []string `json:"topics_done"`
	Completed   bool     `json:"completed"`
}

type Skill struct {
	Topic   string      `json:"topic"`
	Intro   string      `json:"intro"`
	UserIds []string    `json:"userids"`
	Data    []SkillData `json:"data"`
}

type SkillData struct {
	Question   string `json:"question"`
	QuestionId string `json:"question_id"`
	Tutorial   string `json:"tutorial"`
	ExpectedOP string `json:"expected_op"`
	UserCode   string `json:"user_code"`
	Completed  bool   `json:"completed"`
}

type SkillDetails struct {
	Topic   string      `json:"topic"`
	SkillId string      `json:"skill_id"`
	Intro   string      `json:"intro"`
	UserIds []string    `json:"userids"`
	Data    []SkillData `json:"data"`
}

type SkillRes struct {
	Topic string      `json:"topic"`
	Intro string      `json:"intro"`
	Data  []SkillData `json:"data"`
}
type StatusReq struct {
	SkillId string `json:"skill_id"`
}

type GenerateSkillsReq struct {
	Difficulty string `json:"difficulty"`
	Topic      string `json:"topic"`
}
type SubmitSolReq struct {
	Code string `json:"code"`
}
