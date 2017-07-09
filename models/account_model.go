package models

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

var (
	AccountID       = "account_id"
	BillingAddress  = "billing_address"
	ShippingAddress = "shipping_address"
	AccountTable    = "account"
)

//new object
func NewAccount(db *sqlx.DB) *Account {
	Account := &Account{}
	Account.db = db
	Account.table = AccountTable
	Account.hasID = true
	Account.tableID = AccountID
	Account.UserID = 0
	return Account
}

//row
type AccountRow struct {
	AccountID       int64  `db:"account_id"`
	UserID          int64  `db:"user_id"`
	BillingAddress  string `db:"billing_address"`
	ShippingAddress string `db:"shipping_address"`
}

//struct
type Account struct {
	Base
	*AccountRow
}

func (a *Account) AccountRowFromSqlResult(tx *sqlx.Tx, sqlResult sql.Result) error {
	AccountId, err := sqlResult.LastInsertId()
	if err != nil {
		return err
	}

	return a.GetAccountById(tx, AccountId)
}

func (a *Account) CreateAccount(tx *sqlx.Tx) error {

	data := make(map[string]interface{})

	if a.UserID == 0 {
		return errors.New("Account User ID is empty.")
	}

	if a.BillingAddress == "" {
		return errors.New("Billing Address is empty.")
	}

	if a.ShippingAddress == "" {
		return errors.New("Shipping Address is empty.")
	}

	data[BillingAddress] = a.BillingAddress
	data[ShippingAddress] = a.ShippingAddress
	data[UserID] = a.UserID

	sqlResult, err := a.InsertIntoTable(tx, data)

	if err != nil {
		return err
	}

	return a.AccountRowFromSqlResult(tx, sqlResult)
}

// GetById returns record by id.
func (a *Account) GetAccountById(tx *sqlx.Tx, id int64) error {
	AccountRow := &AccountRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE %v=?", a.table, a.tableID)
	err := a.db.Get(AccountRow, query, id)
	if err != nil {
		return err
	}

	a.AccountRow = AccountRow

	return nil
}

//Update

//Delete
