package models

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

func NewTeamMembers(db *sqlx.DB) *TeamMembers {
	TeamMembers := &TeamMembers{}
	TeamMembers.db = db
	TeamMembers.table = "team_members"
	TeamMembers.hasID = true
	TeamMembers.tableID = "team_members_id"
	return TeamMembers
}

type TeamMembersRow struct {
	ID     int64  `db:"team_members_id"`
	UserID string `db:"user_id"`
	TeamID string `db:"team_id"`
}

type TeamMembers struct {
	Base
	TeamMembersRow
}

func (p *TeamMembers) TeamMembersRowFromSqlResult(tx *sqlx.Tx, sqlResult sql.Result) (*TeamMembersRow, error) {
	TeamMembersId, err := sqlResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	return p.GetTeamMembersById(tx, TeamMembersId)
}

func (p *TeamMembers) CreateTeamMembers(tx *sqlx.Tx) (*TeamMembersRow, error) {

	data := make(map[string]interface{})

	if p.UserID == "" {
		return nil, errors.New("UserID  is invalid.")
	}

	if p.TeamID == "" {
		return nil, errors.New("TeamID is invalid.")
	}

	data["user_id"] = p.UserID
	data["team_id"] = p.TeamID

	sqlResult, err := p.InsertIntoTable(tx, data)

	if err != nil {
		return nil, err
	}

	return p.TeamMembersRowFromSqlResult(tx, sqlResult)
}

// GetById returns record by id.
func (p *TeamMembers) GetTeamMembersById(tx *sqlx.Tx, id int64) (*TeamMembersRow, error) {
	TeamMembers := &TeamMembersRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE %v=?", p.table, p.tableID)
	err := p.db.Get(TeamMembers, query, id)

	return TeamMembers, err
}

func (p *TeamMembers) GetTeamMembersByUserId(tx *sqlx.Tx, userID int64) ([]*TeamMembersRow, error) {
	TeamMembers := []*TeamMembersRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE %v=?", p.table, p.tableID)
	err := p.db.Get(TeamMembers, query, userID)

	return TeamMembers, err
}

func (p *TeamMembers) GetTeamMembersByTeamID(tx *sqlx.Tx, teamID int64) ([]*TeamMembersRow, error) {
	TeamMembers := []*TeamMembersRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE %v=?", p.table, "team_id")
	err := p.db.Get(TeamMembers, query, teamID)

	return TeamMembers, err
}

//Update

//Delete
