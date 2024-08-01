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

func (pg *PostQreSQLCon) CreateTablesX() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			user_uid UUID PRIMARY KEY,
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
		);`,
		`CREATE TABLE IF NOT EXISTS socials (
			socials_uid UUID PRIMARY KEY,
			user_id UUID REFERENCES users(user_uid) ON DELETE CASCADE,
			github VARCHAR(255),
			linkedin VARCHAR(255),
			instagram VARCHAR(255),
			noobs_social VARCHAR(255)
		);`,
		`CREATE TABLE IF NOT EXISTS projects (
			project_uid UUID PRIMARY KEY,
			user_id UUID REFERENCES users(user_uid) ON DELETE CASCADE,
			project_name VARCHAR(50) NOT NULL,
			project_description VARCHAR(255),
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);`,
		`CREATE TABLE IF NOT EXISTS files (
			file_uid UUID PRIMARY KEY,
			project_id UUID REFERENCES projects(project_uid) ON DELETE CASCADE,
			file_name VARCHAR(255) NOT NULL,
			file_content TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);`,
		`CREATE TABLE IF NOT EXISTS folders (
			folder_uid UUID PRIMARY KEY,
			project_id UUID REFERENCES projects(project_uid) ON DELETE CASCADE,
			folder_name VARCHAR(255) NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);`,
		`CREATE TABLE IF NOT EXISTS skills (
			skill_uid UUID PRIMARY KEY,
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

func (pg *PostQreSQLCon) CreateTables() error {
	if err := pg.CreateUsersTable(); err != nil {
		return err
	}
	if err := pg.CreateSocialsTable(); err != nil {
		return err
	}
	if err := pg.CreateProjectsTable(); err != nil {
		return err
	}
	if err := pg.CreateFilesTable(); err != nil {
		return err
	}
	if err := pg.CreateFoldersTable(); err != nil {
		return err
	}
	if err := pg.CreateSkillsTable(); err != nil {
		return err
	}
	return nil
}

func (pg *PostQreSQLCon) CreateUsersTable() error {
	query := `CREATE TABLE IF NOT EXISTS users (
		user_uid UUID PRIMARY KEY,
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
	);`
	_, err := pg.dbCon.Exec(query)
	return err
}

func (pg *PostQreSQLCon) CreateSocialsTable() error {
	query := `CREATE TABLE IF NOT EXISTS socials (
		socials_uid UUID PRIMARY KEY,
		user_id UUID REFERENCES users(user_uid) ON DELETE CASCADE,
		github VARCHAR(255),
		linkedin VARCHAR(255),
		instagram VARCHAR(255),
		noobs_social VARCHAR(255)
	);`
	_, err := pg.dbCon.Exec(query)
	return err
}

func (pg *PostQreSQLCon) CreateProjectsTable() error {
	query := `CREATE TABLE IF NOT EXISTS projects (
		project_uid UUID PRIMARY KEY,
		user_id UUID REFERENCES users(user_uid) ON DELETE CASCADE,
		project_name VARCHAR(50) NOT NULL,
		project_description VARCHAR(255),
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`
	_, err := pg.dbCon.Exec(query)
	return err
}

func (pg *PostQreSQLCon) CreateFilesTable() error {
	query := `CREATE TABLE IF NOT EXISTS files (
		file_uid UUID PRIMARY KEY,
		project_id UUID REFERENCES projects(project_uid) ON DELETE CASCADE,
		file_name VARCHAR(255) NOT NULL,
		file_content TEXT NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`
	_, err := pg.dbCon.Exec(query)
	return err
}

func (pg *PostQreSQLCon) CreateFoldersTable() error {
	query := `CREATE TABLE IF NOT EXISTS folders (
		folder_uid UUID PRIMARY KEY,
		project_id UUID REFERENCES projects(project_uid) ON DELETE CASCADE,
		folder_name VARCHAR(255) NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`
	_, err := pg.dbCon.Exec(query)
	return err
}

func (pg *PostQreSQLCon) CreateSkillsTable() error {
	query := `CREATE TABLE IF NOT EXISTS skills (
		skill_uid UUID PRIMARY KEY,
		topic VARCHAR(255) NOT NULL,
		intro TEXT NOT NULL,
		data JSONB NOT NULL,
		user_ids TEXT[] NOT NULL
	);`
	_, err := pg.dbCon.Exec(query)
	return err
}

func (pg *PostQreSQLCon) InsertUser(user models.UserDetails, passwordHash string) error {
	// Initialize dummy values for social_accounts and badges
	emptySocials := models.Socials{}
	emptyBadges := []models.Badge{}

	// Convert these dummy values to JSON
	socialAccountsJSON, err := json.Marshal(emptySocials)
	if err != nil {
		return err
	}

	badgesJSON, err := json.Marshal(emptyBadges)
	if err != nil {
		return err
	}

	query := `INSERT INTO users (user_uid, username, email, phone, first_name, last_name, password_hash, social_accounts, badges, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())`
	_, err = pg.dbCon.Exec(query, user.ID, user.Username, user.Email, user.Phone, user.FirstName, user.LastName, passwordHash, socialAccountsJSON, badgesJSON)
	return err
}

func (pg *PostQreSQLCon) InsertSocialAccounts(userID string, socials models.Socials) error {
	query := `INSERT INTO socials (socials_uid, user_id, github, linkedin, instagram, noobs_social)
              VALUES ($1, $2, $3, $4, $5)`
	_, err := pg.dbCon.Exec(query, uuid.New().String(), userID, socials.GitHub, socials.LinkedIn, socials.Instagram, socials.NoobsSocial)
	return err
}

func (pg *PostQreSQLCon) InsertProject(project models.ProjectDetails) error {
	query := `INSERT INTO projects (project_uid, user_id, project_name, project_description, created_at, updated_at)
              VALUES ($1, $2, $3, $4, NOW(), NOW())`
	_, err := pg.dbCon.Exec(query, project.ProjectID, project.OwnerID, project.ProjectName, project.ProjectDescription)
	return err
}

func (pg *PostQreSQLCon) InsertFile(file models.File) error {
	query := `INSERT INTO files (file_uid, project_id, file_name, file_content, created_at, updated_at)
              VALUES ($1, $2, $3, $4, NOW(), NOW())`
	_, err := pg.dbCon.Exec(query, file.ID, file.ProjectID, file.FileName, file.FileContent)
	return err
}

func (pg *PostQreSQLCon) InsertFolder(folder models.Folder) error {
	query := `INSERT INTO folders (folder_uid, project_id, folder_name, created_at, updated_at)
              VALUES ($1, $2, $3, NOW(), NOW())`
	_, err := pg.dbCon.Exec(query, folder.ID, folder.ProjectID, folder.FolderName)
	return err
}

func (con *PostQreSQLCon) FetchHashedPassword(email string) (string, error) {
	var hashedPassword string
	query := `SELECT password_hash FROM users WHERE email = $1`
	err := con.dbCon.QueryRow(query, email).Scan(&hashedPassword)
	if err != nil {
		return "", err
	}
	return hashedPassword, nil
}

func (con *PostQreSQLCon) FetchUserIdByEmail(email string) (string, error) {
	var userID string
	query := `SELECT user_uid FROM users WHERE email = $1;`
	err := con.dbCon.QueryRow(query, email).Scan(&userID)
	if err != nil {
		return "", err
	}
	return userID, nil
}

func (con *PostQreSQLCon) FetchProjectsByUserId(userId string) ([]models.ProjectDetails, error) {
	query := `SELECT project_uid, user_id, project_name, project_description FROM projects WHERE user_id = $1`
	rows, err := con.dbCon.Queryx(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := []models.ProjectDetails{}
	for rows.Next() {
		var project models.ProjectDetails
		err := rows.Scan(&project.ProjectID, &project.OwnerID, &project.ProjectName, &project.ProjectDescription)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, nil
}

func (con *PostQreSQLCon) FetchFilesByProjectId(projectId string) ([]models.File, error) {
	query := `SELECT file_uid, file_name, file_content FROM files WHERE project_id = $1`
	rows, err := con.dbCon.Queryx(query, projectId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	files := []models.File{}
	for rows.Next() {
		var file models.File
		err := rows.Scan(&file.ID, &file.FileName, &file.FileContent)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return files, nil
}
func (con *PostQreSQLCon) FetchFoldersByProjectId(projectId string) ([]models.FolderDetails, error) {
	query := `SELECT folder_uid, project_id, folder_name, created_at, updated_at FROM folders WHERE project_id =$1;`
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

func (con *PostQreSQLCon) SaveContent(content string) error {
	var req models.SaveFileRequest
	query := `UPDATE files SET file_content = $1, updated_at = NOW() WHERE id = $2 AND project_id = $3;`
	if _, err := con.dbCon.Exec(query, req.Content, req.FileID, req.ProjectID); err != nil {

		return err
	}
	return nil
}
func (con *PostQreSQLCon) FetchUserDetails(userId string) (*models.UserDetails, error) {
	query := `SELECT user_uid, username, email, phone, first_name, last_name, social_accounts, badges, created_at, updated_at FROM users WHERE user_uid=$1`
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
	query := `INSERT INTO skills (skill_uid, topic, intro, data, user_ids) VALUES ($1, $2, $3, $4, $5);`

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
	query := `SELECT skill_uid, topic FROM skills WHERE $1 = ANY(user_ids);`
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
	query := `SELECT skill_uid, topic, intro, data, user_ids FROM skills WHERE id=$1;`
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

// delete tables
func (pg *PostQreSQLCon) DropUsersTable() error {
	query := `DROP TABLE IF EXISTS users;`
	_, err := pg.dbCon.Exec(query)
	return err
}

// DropSocialsTable drops the socials table
func (pg *PostQreSQLCon) DropSocialsTable() error {
	query := `DROP TABLE IF EXISTS socials;`
	_, err := pg.dbCon.Exec(query)
	return err
}

// DropProjectsTable drops the projects table
func (pg *PostQreSQLCon) DropProjectsTable() error {
	query := `DROP TABLE IF EXISTS projects;`
	_, err := pg.dbCon.Exec(query)
	return err
}

// DropFilesTable drops the files table
func (pg *PostQreSQLCon) DropFilesTable() error {
	query := `DROP TABLE IF EXISTS files;`
	_, err := pg.dbCon.Exec(query)
	return err
}

// DropFoldersTable drops the folders table
func (pg *PostQreSQLCon) DropFoldersTable() error {
	query := `DROP TABLE IF EXISTS folders;`
	_, err := pg.dbCon.Exec(query)
	return err
}

// DropSkillsTable drops the skills table
func (pg *PostQreSQLCon) DropSkillsTable() error {
	query := `DROP TABLE IF EXISTS skills;`
	_, err := pg.dbCon.Exec(query)
	return err
}

func (pg *PostQreSQLCon) DropAllTables() error {
	// Drop tables in order to satisfy foreign key constraints
	tables := []string{"files", "folders", "projects", "socials", "skills", "users"}
	for _, table := range tables {
		query := `DROP TABLE IF EXISTS ` + table + ` CASCADE;`
		_, err := pg.dbCon.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (con *PostQreSQLCon) UpdateSocialAccounts(userId string, socials models.Socials) error {
	socialAccountsJSON, err := json.Marshal(socials)
	if err != nil {
		return err
	}

	query := `UPDATE users SET social_accounts = $1, updated_at = NOW() WHERE user_uid = $2`
	_, err = con.dbCon.Exec(query, socialAccountsJSON, userId)
	return err
}
