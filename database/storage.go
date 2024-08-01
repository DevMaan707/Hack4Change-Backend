package database

import (
	"Hack4Change/models"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type PostQreSQLCon struct {
	dbCon *sqlx.DB
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
    social_accounts JSONB,
    badges JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);`, `

CREATE TABLE IF NOT EXISTS socials (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    github VARCHAR(255),
    linkedin VARCHAR(255),
    instagram VARCHAR(255),
    noobs_social VARCHAR(255)
);`, `

CREATE TABLE IF NOT EXISTS projects (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    project_name VARCHAR(50) NOT NULL,
    project_description VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);`, `

CREATE TABLE IF NOT EXISTS files (
    id UUID PRIMARY KEY,
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    file_name VARCHAR(255) NOT NULL,
    file_content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);`, `

CREATE TABLE IF NOT EXISTS folders (
    id UUID PRIMARY KEY,
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    folder_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);`, `

CREATE TABLE IF NOT EXISTS skills (
    id UUID PRIMARY KEY,
    topic VARCHAR(255) NOT NULL,
    intro TEXT NOT NULL,
    data JSONB NOT NULL,
    user_ids TEXT[] NOT NULL
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
	query := `INSERT INTO projects (id, user_id, project_name, project_description, created_at, updated_at)
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

func (con *PostQreSQLCon) FetchHashedPassword(email string) (string, error) {
	var hashedPassword string
	query := `SELECT password FROM users WHERE email = $1`
	err := con.dbCon.QueryRow(query, email).Scan(&hashedPassword)
	if err != nil {
		return "", err
	}
	return hashedPassword, nil
}

func (con *PostQreSQLCon) FetchUserIdByEmail(email string) (string, error) {
	var userID string
	query := `SELECT id FROM users WHERE email = $1;`
	err := con.dbCon.QueryRow(query, email).Scan(&userID)
	if err != nil {
		return "", err
	}
	return userID, nil
}

func (con *PostQreSQLCon) FetchProjectsByUserId(userId string) ([]models.ProjectDetails, error) {
	query := `SELECT project_id, user_id,project_name, project_description,created_at,updated_at FROM projects WHERE user_id =$1;`
	rows, err := con.dbCon.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var projects []models.ProjectDetails
	for rows.Next() {
		var project models.ProjectDetails
		if err := rows.Scan(&project.ProjectID, &project.OwnerID, &project.ProjectName, &project.ProjectDescription, &project.CreatedAt, &project.UpdatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, nil
}
func (con *PostQreSQLCon) FetchFoldersByProjectId(projectId string) ([]models.FolderDetails, error) {
	query := `SELECT id,project_id,folder_name,created_at,updated_at FROM folders WHERE project_id =$1;`
	rows, err := con.dbCon.Query(query, projectId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var folders []models.FolderDetails
	for rows.Next() {
		var folder models.FolderDetails
		if err := rows.Scan(&folder.ID, &folder.ProjectID, &folder.FolderName, &folder.CreatedAt, &folder.UpdatedAt); err != nil {
			return nil, err
		}
		folders = append(folders, folder)
	}
	return folders, nil
}
func (con *PostQreSQLCon) FetchFilesByProjectId(projectId string) ([]models.File, error) {
	query := `SELECT id, project_id, file_name, file_content, created_at, updated_at FROM files WHERE project_id = $1`
	rows, err := con.dbCon.Query(query, projectId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var files []models.File
	for rows.Next() {
		var file models.File
		if err := rows.Scan(&file.ID, &file.ProjectID, &file.FileName, &file.FileContent, &file.CreatedAt, &file.UpdatedAt); err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return files, nil
}

func (con *PostQreSQLCon) SaveContent(content string) error {
	var req models.SaveFileRequest
	query := `UPDATE files SET file_content = $1, updated_at = NOW() WHERE id = $2 AND project_id = $3;`
	if _, err := con.dbCon.Exec(query, req.Content, req.FileID, req.ProjectID); err != nil {

		return err
	}
	return nil
}
func (con *PostQreSQLCon) FetchUserDetails(userId string) (*models.UserDetails, error) {
	query := `SELECT id, username, email, phone, first_name, last_name, social_accounts, badges, created_at, updated_at FROM users WHERE id=$1;`
	var user models.UserDetails
	var socialAccountsJSON, badgesJSON []byte

	err := con.dbCon.QueryRow(query, userId).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Phone,
		&user.FirstName,
		&user.LastName,
		&socialAccountsJSON,
		&badgesJSON,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(socialAccountsJSON, &user.SocialAccounts); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(badgesJSON, &user.Badges); err != nil {
		return nil, err
	}

	return &user, nil
}

func (con *PostQreSQLCon) AddSkill(skill *models.Skill) error {
	query := `INSERT INTO skills (id, topic, intro, data, user_ids) VALUES ($1, $2, $3, $4, $5);`

	id := uuid.New().String()
	dataJSON, err := json.Marshal(skill.Data)
	if err != nil {
		return err
	}

	_, err = con.dbCon.Exec(query, id, skill.Topic, skill.Intro, dataJSON, pq.Array(skill.UserIds))
	if err != nil {
		return err
	}
	return nil
}
func (con *PostQreSQLCon) FetchSkillIdAndNameByUserID(userID string) ([]models.SkillDetails, error) {
	query := `SELECT id, topic FROM skills WHERE $1 = ANY(user_ids);`
	rows, err := con.dbCon.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []models.SkillDetails
	for rows.Next() {
		var skill models.SkillDetails
		if err := rows.Scan(&skill.SkillId, &skill.Topic); err != nil {
			return nil, err
		}
		skills = append(skills, skill)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return skills, nil
}

func (con *PostQreSQLCon) FetchSkillsBySkillID(skillID string) (*models.SkillDetails, error) {
	query := `SELECT id, topic, intro, data, user_ids FROM skills WHERE id=$1;`
	var skill models.SkillDetails
	var dataJSON []byte

	err := con.dbCon.QueryRow(query, skillID).Scan(&skill.SkillId, &skill.Topic, &skill.Intro, &dataJSON, pq.Array(&skill.UserIds))
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(dataJSON, &skill.Data); err != nil {
		return nil, err
	}

	return &skill, nil
}
func (con *PostQreSQLCon) SubmitSolutionByQIDandSkillID(qid string, skillid string) error {
	return nil
}
