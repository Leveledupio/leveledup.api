package models

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

var (
	AddressID    = "address_id"
	Line1        = "line_1"
	Line2        = "line_2"
	Line3        = "line_3"
	Line4        = "line_4"
	City         = "city"
	Zip          = "zip_or_post"
	Country      = "country"
	State        = "state"
	AddressTable = "address"
)

func NewAddress(db *sqlx.DB) *Address {
	Address := &Address{}
	Address.AddressRow = &AddressRow{}

	Address.db = db
	Address.table = AddressTable
	Address.hasID = true
	Address.tableID = AddressID
	Address.UserID = 0
	Address.Zip = 0
	return Address
}

type AddressRow struct {
	AddressID int64  `db:"address_id"`
	UserID    int64  `db:"user_id"`
	Line1     string `db:"line_1"`
	Line2     string `db:"line_2"`
	Line3     string `db:"line_3"`
	Line4     string `db:"line_4"`
	City      string `db:"city"`
	Zip       int    `db:"zip_or_post"`
	Country   string `db:"country"`
	State     string `db:"state"`
}

type Address struct {
	Base
	*AddressRow
}

func (a *Address) AddressRowFromSqlResult(tx *sqlx.Tx, sqlResult sql.Result) error {
	AddressId, err := sqlResult.LastInsertId()
	if err != nil {
		return err
	}

	return a.GetAddressById(tx, AddressId)
}

func (a *Address) CreateAddress(tx *sqlx.Tx) error {

	data := make(map[string]interface{})

	if a.UserID == 0 {
		return errors.New("User ID is invalid.")
	}

	if a.Line1 == "" {
		return errors.New("Line 1 is invalid.")
	}

	if a.City == "" {
		return errors.New("City is invalid.")
	}

	if a.Zip == 0 {
		return errors.New("Zip is invalid.")
	}

	if a.Country == "" {
		return errors.New("Country is invalid.")
	}

	if a.State == "" {
		return errors.New("State is invalid.")
	}

	data[UserID] = a.UserID
	data[City] = a.City
	data[Zip] = a.Zip
	data[Country] = a.Country
	data[State] = a.State
	data[Line1] = a.Line1
	data[Line2] = a.Line2
	data[Line3] = a.Line3
	data[Line4] = a.Line4

	sqlResult, err := a.InsertIntoTable(tx, data)

	if err != nil {
		return err
	}

	return a.AddressRowFromSqlResult(tx, sqlResult)
}

// GetById returns record by id.
func (a *Address) GetAddressById(tx *sqlx.Tx, id int64) error {
	AddressRow := &AddressRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE %v=?", a.table, a.tableID)
	err := a.db.Get(AddressRow, query, id)
	if err != nil {
		return err
	}

	a.AddressRow = AddressRow

	return nil
}

// GetById returns record by id.
func (a *Address) GetAddressByZip(tx *sqlx.Tx, zip int) ([]*AddressRow, error) {
	Address := []*AddressRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE %v=?", a.table, Zip)
	err := a.db.Get(Address, query, zip)

	return Address, err
}

func (a *Address) GetAddressByState(tx *sqlx.Tx, state string) ([]*AddressRow, error) {
	Address := []*AddressRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE %v=?", a.table, State)
	err := a.db.Get(Address, query, state)

	return Address, err
}

func (a *Address) GetAddressByCountry(tx *sqlx.Tx, country string) ([]*AddressRow, error) {
	Address := []*AddressRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE %v=?", a.table, Country)
	err := a.db.Get(Address, query, country)

	return Address, err
}

func (a *Address) GetAddressByUserID(tx *sqlx.Tx, userID int64) ([]*AddressRow, error) {
	Address := []*AddressRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE %v=?", a.table, UserID)
	err := a.db.Get(Address, query, userID)

	return Address, err
}

//Update

//Delete
