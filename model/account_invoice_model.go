package models

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/go-bootstrap/go-bootstrap/project-templates/postgresql/models"
)

var (
	InvoiceID = "invoice_id"
	AccountINVTable = "account_invoice"
)

func NewAccountINV(db *sqlx.DB) *AccountINV {
	AccountINV := &AccountINV{}
	AccountINV.db = db
	AccountINV.table = AccountINVTable
	AccountINV.hasID = true
	AccountINV.tableID = AccountID
	AccountINV.AccountID = 0
	AccountINV.InvoiceID = 0
	return AccountINV
}

type AccountINVRow struct {
	AccountID int64 `db:"account_id"`
	InvoiceID int64 `db:"invoice_id"`
}

type AccountINV struct {
	Base
	AccountINVRow
}

func (a *AccountINV) AccountINVRowFromSqlResult(tx *sqlx.Tx, sqlResult sql.Result) (*AccountINVRow, error) {
	AccountINVId, err := sqlResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	return a.GetAccountINVById(tx, AccountINVId)
}

func (a *AccountINV) CreateAccountINV(tx *sqlx.Tx) (*AccountINVRow, error) {

	data := make(map[string]interface{})

	if a.AccountID == 0 {
		return nil, errors.New("Account ID is invalid.")
	}

	if a.InvoiceID == 0 {
		return nil, errors.New("Invoice ID is invalid.")
	}

	data[AccountID] = a.AccountID
	data[InvoiceID] = a.InvoiceID

	sqlResult, err := a.InsertIntoTable(tx, data)

	if err != nil {
		return nil, err
	}

	return a.AccountINVRowFromSqlResult(tx, sqlResult)
}

// GetById returns record by id.
func (a *AccountINV) GetAccountINVById(tx *sqlx.Tx, id int64) (*AccountINVRow, error) {
	AccountINV := &AccountINVRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE %v=?", a.table, a.tableID)
	err := a.db.Get(AccountINV, query, id)

	return AccountINV, err
}

func (a *AccountINV) GetAccountINVByInvoiceId(tx *sqlx.Tx, id int64) (*AccountINVRow, error) {
	AccountINV := &AccountINVRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE %v=?", a.table, InvoiceID)
	err := a.db.Get(AccountINV, query, id)

	return AccountINV, err
}

func (a *AccountINV) GetAccountINVByAccountId(tx *sqlx.Tx, id int64) ([]AccountINVRow, error) {
	AccountINV := []AccountINVRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE %v=?", a.table, AccountID)

	err := a.db.Get(AccountINV, query, id)

	return AccountINV, err
}

//Update

//Delete
