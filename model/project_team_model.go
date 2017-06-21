package models

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

//Database attributes
var (
	ProjectID = "project_id"
	TeamID = "team_id"
	ProjectURL = "project_url"
	ProjectTeamID = "project_team_id"
	ProjectTable = "project_team"
)

func NewProjectTeam(db *sqlx.DB) *ProjectTeam {
	ProjectTeam := &ProjectTeam{}
	ProjectTeam.ProjectTeamRow = &ProjectTeamRow{}
	ProjectTeam.db = db
	ProjectTeam.table = ProjectTable
	ProjectTeam.hasID = true
	ProjectTeam.tableID = ProjectTeamID
	return ProjectTeam
}

type ProjectTeamRow struct {
	ProjectID     int64  `db:"project_id"`
	TeamID        int64  `db:"team_id"`
	ProjectURL    string `db:"project_url"`
	ProjectTeamID int64  `db:"project_team_id"`
}

type ProjectTeam struct {
	Base
	*ProjectTeamRow
}

func (p *ProjectTeam) ProjectTeamRowFromSqlResult(tx *sqlx.Tx, sqlResult sql.Result) error {
	ProjectTeamId, err := sqlResult.LastInsertId()
	if err != nil {
		return err
	}

	return p.GetProjectTeamById(tx, ProjectTeamId)
}

func (p *ProjectTeam) CreateProjectTeam(tx *sqlx.Tx) error {

	data := make(map[string]interface{})

	if p.ProjectURL == "" {
		return errors.New("ProjectURL is invalid.")
	}

	if p.TeamID == 0 {
		return errors.New("TeamID is invalid.")
	}
	if p.ProjectID == 0 {
		return errors.New("ProjectID is invalid.")
	}
	data[ProjectURL] = p.ProjectURL
	data[TeamID] = p.TeamID
	data[ProjectID] = p.ProjectID

	sqlResult, err := p.InsertIntoTable(tx, data)
	if err != nil {
		return err
	}

	ProjectTeamId, err := sqlResult.LastInsertId()
	if err != nil {
		return err
	}
	log.Infof("Sql Result %v", ProjectTeamId)

	return p.ProjectTeamRowFromSqlResult(tx, sqlResult)
}

// GetById returns record by id.
func (p *ProjectTeam) GetProjectTeamById(tx *sqlx.Tx, project_team_id int64) error {
	ProjectTeamRow := &ProjectTeamRow{}

	query := fmt.Sprintf("SELECT * FROM %v WHERE %v=?", p.table, p.tableID)
	err := p.db.Get(ProjectTeamRow, query, project_team_id)

	p.ProjectTeamRow = ProjectTeamRow

	return err
}

func (p *ProjectTeam) GetProjectTeamByTeamID(tx *sqlx.Tx) ([]*ProjectTeamRow, error) {
	ProjectTeam := []*ProjectTeamRow{}

	query := fmt.Sprintf("SELECT * FROM %v WHERE %v=?", p.table, TeamID)

	log.Infof("Project Team query %v", query)

	err := p.db.Get(ProjectTeam, query, p.TeamID)

	return ProjectTeam, err
}

func (p *ProjectTeam) GetProjectTeamByProjectID(tx *sqlx.Tx) ([]*ProjectTeamRow, error) {
	ProjectTeam := []*ProjectTeamRow{}

	query := fmt.Sprintf("SELECT * FROM %v WHERE %v=?", p.table, ProjectID)

	log.Infof("Project Team query %v", query)

	err := p.db.Get(ProjectTeam, query, p.ProjectID)

	return ProjectTeam, err
}

//Update

//Delete
