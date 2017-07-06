package models

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

var (
	PermissionID    = "permission_id"
	PermissionName  = "permission_name"
	PermissionRole  = "permissions"
	PermissionTable = "permission"
)

func NewPermission(db *sqlx.DB) *Permission {
	Permission := &Permission{}
	Permission.PermissionRow = &PermissionRow{}

	Permission.db = db
	Permission.table = PermissionTable
	Permission.hasID = true
	Permission.tableID = PermissionID
	Permission.PermissionID = 0
	return Permission
}

type PermissionRow struct {
	PermissionID   int64  `db:"permission_id"`
	Name           string `db:"permission_name"`
	PermissionRole string `db:"permissions"`
}

type Permission struct {
	Base
	*PermissionRow
}

func (p *Permission) PermissionRowFromSqlResult(tx *sqlx.Tx, sqlResult sql.Result) error {
	permissionId, err := sqlResult.LastInsertId()
	if err != nil {
		return err
	}

	return p.GetPermissionById(tx, permissionId)
}

func (p *Permission) CreatePermission(tx *sqlx.Tx) error {

	data := make(map[string]interface{})

	if p.Name == "" {
		return errors.New("Permission Name is invalid.")
	}

	if p.PermissionRole == "" {
		return errors.New("Permission is invalid.")
	}

	data[PermissionName] = p.Name
	data[PermissionRole] = p.PermissionRole

	sqlResult, err := p.InsertIntoTable(tx, data)

	if err != nil {
		return err
	}

	return p.PermissionRowFromSqlResult(tx, sqlResult)
}

// GetById returns record by id.
func (p *Permission) GetPermissionById(tx *sqlx.Tx, id int64) error {
	permission := &PermissionRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE %v=?", p.table, p.tableID)
	err := p.db.Get(permission, query, id)

	p.PermissionRow = permission

	return err
}

//Update

//Delete
