package models

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

func NewProject(db *sqlx.DB) *Project {
	Project := &Project{}
	Project.ProjectRow = &ProjectRow{}
	Project.db = db
	Project.table = "project"
	Project.hasID = true
	Project.tableID = "project_id"
	Project.Difficulty = 0
	return Project
}

type ProjectRow struct {
	ProjectID   int64  `db:"project_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Difficulty  int    `db:"difficulty"`
}

type Project struct {
	Base
	*ProjectRow
}

func (p *Project) ProjectRowFromSqlResult(tx *sqlx.Tx, sqlResult sql.Result) error {
	ProjectId, err := sqlResult.LastInsertId()
	if err != nil {
		return err
	}

	return p.GetProjectById(tx, ProjectId)
}

func (p *Project) CreateProject(tx *sqlx.Tx) error {

	data := make(map[string]interface{})

	if p.Name == "" {
		return errors.New("Project Name is invalid.")
	}

	if p.Description == "" {
		return errors.New("Description is invalid.")
	}

	if p.Difficulty == 0 {
		return errors.New("Difficulty is invalid.")
	}

	data["name"] = p.Name
	data["description"] = p.Description
	data["difficulty"] = p.Difficulty

	sqlResult, err := p.InsertIntoTable(tx, data)

	if err != nil {
		return err
	}

	return p.ProjectRowFromSqlResult(tx, sqlResult)
}

// GetById returns record by id.
func (p *Project) GetProjectById(tx *sqlx.Tx, id int64) error {
	ProjectRow := &ProjectRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE %v=?", p.table, p.tableID)
	err := p.db.Get(ProjectRow, query, id)

	p.ProjectRow = ProjectRow

	return err
}

func (p *Project) GetProjectByDifficulty(tx *sqlx.Tx, diff int) ([]*ProjectRow, error) {
	Project := []*ProjectRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE %v=?", p.table, "difficulty")
	err := p.db.Get(Project, query, diff)

	return Project, err
}

//Update

//Delete
