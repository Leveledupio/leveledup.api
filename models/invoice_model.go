package models

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

func NewInvoice(db *sqlx.DB) *Invoice {
	Invoice := &Invoice{}
	Invoice.db = db
	Invoice.table = "invoice"
	Invoice.hasID = true
	Invoice.tableID = "invoice_id"
	Invoice.UnitPrice = 0
	Invoice.Units = 0
	return Invoice
}

type InvoiceRow struct {
	ID            int64  `db:"invoice_id"`
	Date          string `db:"invoice_date"`
	DueDate       string `db:"due_date"`
	PayDate       string `db:"pay_date"`
	Units         int    `db:"units"`
	UnitPrice     int64  `db:"unit_price"`
	Description   string `db:"description"`
	AmountDue     string `db:"amount_due"`
	PaymentAmount string `db:"payment_amount"`
	Notes         string `db:"notes"`
	NextBillDate  string `db:"next_bill_date"`
}

type Invoice struct {
	Base
	InvoiceRow
}

func (p *Invoice) InvoiceRowFromSqlResult(tx *sqlx.Tx, sqlResult sql.Result) (*InvoiceRow, error) {
	InvoiceId, err := sqlResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	return p.GetInvoiceById(tx, InvoiceId)
}

func (p *Invoice) CreateInvoice(tx *sqlx.Tx) (*InvoiceRow, error) {

	data := make(map[string]interface{})

	if p.Date == "" {
		return nil, errors.New("Date is invalid.")
	}

	if p.DueDate == "" {
		return nil, errors.New("Invoice Name is invalid.")
	}

	if p.PayDate == "" {
		return nil, errors.New("Invoice Name is invalid.")
	}

	if p.Units == 0 {
		return nil, errors.New("Invoice Name is invalid.")
	}

	if p.UnitPrice == 0 {
		return nil, errors.New("Invoice Name is invalid.")
	}

	if p.Description == "" {
		return nil, errors.New("Invoice Name is invalid.")
	}

	if p.AmountDue == "" {
		return nil, errors.New("Invoice Name is invalid.")
	}

	if p.PaymentAmount == "" {
		return nil, errors.New("Invoice Name is invalid.")
	}

	if p.NextBillDate == "" {
		return nil, errors.New("Invoice Name is invalid.")
	}

	data["invoice_date"] = p.Date
	data["due_date"] = p.DueDate
	data["pay_date"] = p.PayDate
	data["units"] = p.Units
	data["unit_price"] = p.UnitPrice
	data["description"] = p.Description
	data["amount_due"] = p.AmountDue
	data["payment_amount"] = p.PaymentAmount
	data["notes"] = p.Notes
	data["next_bill_date"] = p.NextBillDate

	sqlResult, err := p.InsertIntoTable(tx, data)

	if err != nil {
		return nil, err
	}

	return p.InvoiceRowFromSqlResult(tx, sqlResult)
}

// GetById returns record by id.
func (p *Invoice) GetInvoiceById(tx *sqlx.Tx, id int64) (*InvoiceRow, error) {
	Invoice := &InvoiceRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE %v=?", p.table, p.tableID)
	err := p.db.Get(Invoice, query, id)

	return Invoice, err
}

func (p *Invoice) GetInvoiceByDate(tx *sqlx.Tx, date string) (*InvoiceRow, error) {
	Invoice := &InvoiceRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE %v=?", p.table, "date")
	err := p.db.Get(Invoice, query, date)

	return Invoice, err
}

//Update

//Delete
