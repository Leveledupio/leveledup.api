package models

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"

)

var (
	UserPermissionTable = "user_permission"
	UserPermID = "user_permission_id"
)

func NewUserPermission(db *sqlx.DB) *UserPermission {
	UserPermission := &UserPermission{}
	UserPermission.UserPermissionRow = &UserPermissionRow{}
	UserPermission.db = db
	UserPermission.table = UserPermissionTable
	UserPermission.hasID = true
	UserPermission.tableID = UserPermID
	return UserPermission
}

type UserPermissionRow struct {
	UserPermID   int64 `db:"user_permission_id"`
	UserID       int64 `db:"user_id"`
	PermissionID int64 `db:"permission_id"`
}

type UserPermission struct {
	Base
	*UserPermissionRow
}

// GetByUserId returns user permission record that contains user_id.
func (up *UserPermission) GetPermissionByUserId(tx *sqlx.Tx, id int64, tableID string) ([]*UserPermissionRow, error) {

	permissions := []*UserPermissionRow{}

	query := fmt.Sprintf("SELECT * FROM %v WHERE %v=?", up.table, tableID)
	err := up.db.Get(permissions, query, id)

	return permissions, err
}

// GetPermissionByPermissionId returns user permission record that contains permission_id.
func (up *UserPermission) GetUserPermissionByPermissionId(tx *sqlx.Tx, id int64, tableID string) ([]*UserPermissionRow, error) {

	permissions := []*UserPermissionRow{}

	query := fmt.Sprintf("SELECT * FROM %v WHERE %v=?", up.table, tableID)
	err := up.db.Get(&permissions, query, id)

	return permissions, err
}

func (up *UserPermission) CreateUserPermission(tx *sqlx.Tx) error {

	data := make(map[string]interface{})

	if up.UserID == 0 {
		return errors.New("User ID is invalid.")
	}

	if up.PermissionID == 0 {
		return errors.New("Permission is invalid.")
	}

	data[UserID] = up.UserID
	data[PermissionID] = up.PermissionID

	sqlResult, err := up.InsertIntoTable(tx, data)
	if err != nil {
		log.Errorf("[ERROR][USER_PERMISSION] - %v", err)
		return err
	}

	return up.userPermissionRowFromSqlResult(tx, sqlResult)
}

func (up *UserPermission) userPermissionRowFromSqlResult(tx *sqlx.Tx, sqlResult sql.Result) error {
	userPermissionId, err := sqlResult.LastInsertId()
	if err != nil {
		return err
	}

	//log.Infof("[DEBUG][USER_PERMISSION] SQL Result - Permission ID %v", userPermissionId)
	return up.GetUserPermissionById(tx, userPermissionId)
}

// GetById returns record by id.
func (up *UserPermission) GetUserPermissionById(tx *sqlx.Tx, permissionId int64) error {

	//log.Print("[DEBUG][USER_PERMISSION]) GetUserPermissionById")

	userPermissionRow := &UserPermissionRow{}

	query := fmt.Sprintf("SELECT * FROM %v WHERE %v=?", up.table, UserPermID)

	log.Debugf("[DEBUG][USER_PERMISSION]) QUERY %v", query)
	err := up.db.Get(userPermissionRow, query, permissionId)
	if err != nil {
		log.Errorf("[ERROR][USER_PERMISSION] Error retrieving %v from Database %v", permissionId, err)
		return err
	}

	//log.Infof("[ERROR][USER_PERMISSION] User Permission Row: %v", userPermissionRow)

	up.UserPermissionRow = userPermissionRow

	return nil
}

//Update

//Delete
