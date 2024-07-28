package database

import (
	"Hack4Change/models"
	"database/sql"
)

type PostQreSQLCon struct {
	dbCon *sql.DB
}

func (pg *PostQreSQLCon) CreateTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			username VARCHAR(32) NOT NULL UNIQUE,
			email VARCHAR(255) NOT NULL UNIQUE,
			phone VARCHAR(20),
			first_name VARCHAR(32),
			last_name VARCHAR(32),
			password_hash TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);`,
		`CREATE TABLE IF NOT EXISTS socials (
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			github VARCHAR(255),
			linkedin VARCHAR(255),
			instagram VARCHAR(255),
			noobs_social VARCHAR(255),
			PRIMARY KEY (user_id)
		);`,
		`CREATE TABLE IF NOT EXISTS projects (
			
			userid UUID REFERENCES users(id),
			project_id VARCHAR(50) NOT NULL,
			owner_id UUID REFERENCES users(id) ON DELETE CASCADE,
			project_name VARCHAR(50) NOT NULL,
			project_description VARCHAR(255),
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);`,
		`CREATE TABLE IF NOT EXISTS files (
			id UUID PRIMARY KEY,
			project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
			file_name VARCHAR(255) NOT NULL,
			file_content TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);`,
		`CREATE TABLE IF NOT EXISTS folders (
			id UUID PRIMARY KEY,
			project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
			folder_name VARCHAR(255) NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);`,
	}

	for _, query := range queries {
		if _, err := pg.dbCon.Exec(query); err != nil {
			return err
		}
	}
	return nil
}

func (pg *PostQreSQLCon) InsertUser(user models.UserDetails, passwordHash string) error {
	query := `INSERT INTO users (id, username, email, phone, first_name, last_name, password_hash, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())`
	_, err := pg.dbCon.Exec(query, user.ID, user.Username, user.Email, user.Phone, user.FirstName, user.LastName, passwordHash)
	return err
}

func (pg *PostQreSQLCon) InsertSocialAccounts(userID string, socials models.Socials) error {
	query := `INSERT INTO socials (user_id, github, linkedin, instagram, noobs_social)
              VALUES ($1, $2, $3, $4, $5)`
	_, err := pg.dbCon.Exec(query, userID, socials.GitHub, socials.LinkedIn, socials.Instagram, socials.NoobsSocial)
	return err
}

func (pg *PostQreSQLCon) InsertProject(project models.ProjectDetails) error {
	query := `INSERT INTO projects (id, project_id, owner_id, project_name, project_description, created_at, updated_at)
              VALUES ($1, $2, $3, $4, NOW(), NOW())`
	_, err := pg.dbCon.Exec(query, project.ProjectID, project.OwnerID, project.ProjectName, project.ProjectDescription)
	return err
}

func (pg *PostQreSQLCon) InsertFile(file models.File) error {
	query := `INSERT INTO files (id, project_id, file_name, file_content, created_at, updated_at)
              VALUES ($1, $2, $3, $4, NOW(), NOW())`
	_, err := pg.dbCon.Exec(query, file.ID, file.ProjectID, file.FileName, file.FileContent)
	return err
}

func (pg *PostQreSQLCon) InsertFolder(folder models.Folder) error {
	query := `INSERT INTO folders (id, project_id, folder_name, created_at, updated_at)
              VALUES ($1, $2, $3, NOW(), NOW())`
	_, err := pg.dbCon.Exec(query, folder.ID, folder.ProjectID, folder.FolderName)
	return err
}
