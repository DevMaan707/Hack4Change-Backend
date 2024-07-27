package models

import (
	"time"

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
}

type CreateProjectReq struct {
	ProjectID          string `json:"project_id" validate:"required"`
	ProjectName        string `json:"project_name" validate:"required,min=3,max=50"`
	ProjectDescription string `json:"project_description" validate:"omitempty,max=255"`
}

type ProjectDetails struct {
	ID                 string    `json:"id"`
	ProjectID          string    `json:"project_id"`
	OwnerID            string    `json:"owner_id"`
	ProjectName        string    `json:"project_name"`
	ProjectDescription string    `json:"project_description"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type File struct {
	ID          string    `json:"id"`
	ProjectID   string    `json:"project_id"`
	FileName    string    `json:"file_name" validate:"required,min=1,max=255"`
	FileContent string    `json:"file_content" validate:"required"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Folder struct {
	ID         string    `json:"id"`
	ProjectID  string    `json:"project_id"`
	FolderName string    `json:"folder_name" validate:"required,min=1,max=255"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
