package models

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

var (
	CustAddID    = "customer_address_id"
	CustAddTable = "customer_addresses"
)

func NewCustomerAdd(db *sqlx.DB) *CustomerAdd {
	CustomerAdd := &CustomerAdd{}
	CustomerAdd.CustomerAddRow = &CustomerAddRow{}
	CustomerAdd.db = db
	CustomerAdd.table = CustAddTable
	CustomerAdd.hasID = true
	CustomerAdd.tableID = CustAddID
	CustomerAdd.AddressId = 0
	return CustomerAdd
}

type CustomerAddRow struct {
	CustAddID   int64 `db:"customer_address_id"`
	AddressId   int64 `db:"address_id"`
	AddressType int64 `db:"address_type_code"`
}

type CustomerAdd struct {
	Base
	*CustomerAddRow
}

func (p *CustomerAdd) CustomerAddRowFromSqlResult(tx *sqlx.Tx, sqlResult sql.Result) error {
	CustomerAddId, err := sqlResult.LastInsertId()
	if err != nil {
		return err
	}

	return p.GetCustomerAddById(tx, CustomerAddId)
}

func (p *CustomerAdd) CreateCustomerAdd(tx *sqlx.Tx) error {

	data := make(map[string]interface{})

	if p.AddressId == 0 {
		return errors.New("AddressId  is invalid.")
	}

	if p.AddressType == 0 {
		return errors.New("AddressType is invalid.")
	}

	data[AddressID] = p.AddressId
	data[AddressTypeID] = p.AddressType

	sqlResult, err := p.InsertIntoTable(tx, data)

	if err != nil {
		return err
	}

	return p.CustomerAddRowFromSqlResult(tx, sqlResult)
}

// GetById returns record by id.
func (p *CustomerAdd) GetCustomerAddById(tx *sqlx.Tx, id int64) error {
	CustomerAdd := &CustomerAddRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE %v=?", p.table, p.tableID)
	err := p.db.Get(CustomerAdd, query, id)

	p.CustomerAddRow = CustomerAdd

	return err
}

//Update

//Delete
