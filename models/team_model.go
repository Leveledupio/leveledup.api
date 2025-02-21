package models

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

func NewTeam(db *sqlx.DB) *Team {
	Team := &Team{}
	Team.db = db
	Team.table = "team"
	Team.hasID = true
	Team.CreatedBy = 0
	Team.tableID = "team_id"
	return Team
}

type TeamRow struct {
	ID          int64  `db:"team_id" json:"team_id,omitempty"`
	Name        string `db:"team_name" json:"team_name"`
	Description string `db:"team_desc" json:"team_desc"`
	CreatedBy   int64  `db:"created_by" json:"created_by"`
}

type Team struct {
	Base
	TeamRow
}

func (t *Team) PrintTeam() string {
	return fmt.Sprintf("Team Name: '%s'  Team Description: '%s' Created by ID: '%v' ", t.Name, t.Description, t.ID)
}

func (t *Team) TeamRowFromSqlResult(tx *sqlx.Tx, sqlResult sql.Result) (*TeamRow, error) {
	TeamId, err := sqlResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	return t.GetTeamById(tx, TeamId)
}

func (t *Team) CreateTeam(tx *sqlx.Tx) (*TeamRow, error) {

	data := make(map[string]interface{})

	if t.Name == "" {
		return nil, errors.New("Team Name is empty or invalid.")
	}

	//Checking for unique team names
	//GetTeamByName returns nil if the team exists
	_, err := t.GetTeamByName(tx, t.Name)
	if err == nil {
		return nil, errors.New("Name must unique")
	}

	if t.Description == "" {
		return nil, errors.New("Description is invalid.")
	}

	if t.CreatedBy == 0 {
		return nil, errors.New("CreatedBy is invalid.")
	}

	data["team_name"] = t.Name
	data["team_desc"] = t.Description
	data["created_by"] = t.CreatedBy

	sqlResult, err := t.InsertIntoTable(tx, data)

	if err != nil {
		return nil, err
	}

	return t.TeamRowFromSqlResult(tx, sqlResult)
}

// GetById returns record by id.
func (t *Team) GetTeamById(tx *sqlx.Tx, id int64) (*TeamRow, error) {
	Team := &TeamRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE %v=?", t.table, t.tableID)
	err := t.db.Get(Team, query, id)

	return Team, err
}

// GetById returns record by id.
func (t *Team) GetTeamByName(tx *sqlx.Tx, name string) (*TeamRow, error) {
	Team := &TeamRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE %v=?", t.table, "team_name")
	err := t.db.Get(Team, query, name)
	if err != nil {
		return nil, err
	}

	return Team, nil
}

//Update
