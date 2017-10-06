package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	jira "github.com/andygrunwald/go-jira"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func NewProject(db *sqlx.DB, jira *jira.Client) *Project {
	Project := &Project{}
	Project.ProjectRow = &ProjectRow{}
	Project.db = db
	Project.jira = jira
	Project.table = "project"
	Project.hasID = true
	Project.tableID = "project_id"
	Project.Difficulty = 0

	log.Debugf("[DEBUG] New Project creation Porject Jira URL %v", Project.jira.GetBaseURL())

	return Project
}

type ProjectRow struct {
	ProjectID   int64  `db:"project_id" json:"project_id"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
	Difficulty  int    `db:"difficulty" json:"difficulty"`
	Data        []Data `json:"data"`
}

type Data struct {
	Epics   []Epic   `json:"epics"`
	Sprints []Sprint `json:"sprints"`
	Stories []Story  `json:"stories"`
	Tasks   []Task   `json:"tasks"`
}

type Epic struct {
	ProjectID   int64    `json:"projectID"`
	Sprints     []Sprint `json:"sprints"`
	Description string   `json:"description"`
	Name        string   `json:"name"`
}

type Sprint struct {
	Epic        Epic
	Stories     []Story
	Description string
	Name        string
}

type Story struct {
	Sprint      Sprint
	Points      int
	Tasks       []Task
	Description string
	Name        string
}

type Task struct {
	Story       Story
	Description string
	Name        string
}

type Project struct {
	Base
	*ProjectRow
}

func (p *Project) PrintProject() {

	b, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	log.Debugf("Project: %s", string(b))
}

func (p *Project) ProjectRowFromSqlResult(tx *sqlx.Tx, sqlResult sql.Result) error {
	ProjectId, err := sqlResult.LastInsertId()
	if err != nil {
		return err
	}

	return p.GetProjectById(tx, ProjectId)
}

func (p *Project) CreateProject() error {

	if p.Name == "" {
		return errors.New("Project Name is invalid.")
	}

	if p.Description == "" {
		return errors.New("Description is invalid.")
	}

	if p.Difficulty == 0 {
		return errors.New("Difficulty is invalid.")
	}

	err := p.createproject()
	if err != nil {
		return err
	}

	return nil
}

func (p *Project) GetAllProject() ([]ProjectRow, error) {

	pList, err := p.getAllProjects()
	if err != nil {
		return nil, err
	}

	return pList, nil
}

// GetById returns record by id.
func (p *Project) GetProjectById(tx *sqlx.Tx, id int64) error {
	ProjectRow := &ProjectRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE %v=?", p.table, p.tableID)
	err := p.db.Get(ProjectRow, query, id)

	p.ProjectRow = ProjectRow

	return err
}

// GetById returns record by id.
func (p *Project) GetProjectByName(tx *sqlx.Tx) error {
	ProjectRow := &ProjectRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE %v=?", p.table, "name")
	err := p.db.Get(ProjectRow, query, p.Name)

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
