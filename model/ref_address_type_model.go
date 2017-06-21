package models

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

var (
	AddressTypeID = "address_type_code"
	AddressDesc = "address_description"
	AddressName = "name"
	AddressTypeTable = "ref_address_type"
)

func NewAddressType(db *sqlx.DB) *AddressType {
	AddressType := &AddressType{}
	AddressType.AddressTypeRow = &AddressTypeRow{}

	AddressType.db = db
	AddressType.tableID = AddressTypeID
	AddressType.hasID = true
	AddressType.table = AddressTypeTable
	return AddressType
}

type AddressTypeRow struct {
	AddressTypeID int64  `db:"address_type_code"`
	AddressDesc   string `db:"address_description"`
	AddressName   string `db:"name"`
}

type AddressType struct {
	Base
	*AddressTypeRow
}

func (at *AddressType) AddressTypeRowFromSqlResult(tx *sqlx.Tx, sqlResult sql.Result) error {
	AddressTypeId, err := sqlResult.LastInsertId()
	if err != nil {
		return err
	}

	return at.GetAddressTypeById(tx, AddressTypeId)
}

func (at *AddressType) CreateAddressType(tx *sqlx.Tx) error {

	data := make(map[string]interface{})

	if at.AddressName == "" {
		return errors.New("Address Name is invalid.")
	}

	if at.AddressDesc == "" {
		return errors.New("Description is invalid.")
	}

	data[AddressName] = at.AddressName
	data[AddressDesc] = at.AddressDesc

	sqlResult, err := at.InsertIntoTable(tx, data)

	if err != nil {
		return err
	}

	return at.AddressTypeRowFromSqlResult(tx, sqlResult)
}

// GetById returns record by id.
func (at *AddressType) GetAddressTypeById(tx *sqlx.Tx, id int64) error {
	AddressTypeRow := &AddressTypeRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE %v=?", at.table, at.tableID)
	err := at.db.Get(AddressTypeRow, query, id)
	if err != nil {
		return err
	}

	at.AddressTypeRow = AddressTypeRow

	return nil
}

//Update

//Delete
