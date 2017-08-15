package models

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"gopkg.in/op/go-logging.v1"

	"strings"
	"time"
	"github.com/aws/aws-sdk-go/aws/session"
)

var (
	log = logging.MustGetLogger("gbs")
)

type Base struct {
	db      *sqlx.DB
	awsSession     *session.Session
	table   string
	tableID string
	hasID   bool
}

func (b *Base) todayDate() string {
	year, month, day := time.Now().Date()
	var month_int int = int(month)
	return fmt.Sprintf("%v-%v-%v", year, month_int, day)
}

func (b *Base) newTransactionIfNeeded(tx *sqlx.Tx) (*sqlx.Tx, bool, error) {
	var err error
	wrapInSingleTransaction := false

	if tx != nil {
		return tx, wrapInSingleTransaction, nil
	}

	tx, err = b.db.Beginx()
	if err == nil {
		wrapInSingleTransaction = true
	}

	if err != nil {
		return nil, wrapInSingleTransaction, err
	}

	//log.Infof("[DEBUG][BASE.NewTransaction] %v", tx)
	return tx, wrapInSingleTransaction, nil

}

func (b *Base) InsertIntoTable(tx *sqlx.Tx, data map[string]interface{}) (sql.Result, error) {
	if b.table == "" {
		return nil, errors.New("Table must not be empty.")
	}

	tx, wrapInSingleTransaction, err := b.newTransactionIfNeeded(tx)
	if tx == nil {
		return nil, errors.New("Transaction struct must not be empty.")
	}
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0)
	qMarks := make([]string, 0)
	values := make([]interface{}, 0)

	for key, value := range data {
		keys = append(keys, key)
		qMarks = append(qMarks, "?")
		values = append(values, value)
	}

	query := fmt.Sprintf(
		"INSERT INTO %v (%v) VALUES (%v)",
		b.table,
		strings.Join(keys, ","),
		strings.Join(qMarks, ","))

	result, err := tx.Exec(query, values...)
	if err != nil {
		log.Errorf("[ERROR][BASE.INSERTINTOTABLE] %v", err)
		return nil, err
	}

	if wrapInSingleTransaction == true {
		err = tx.Commit()
	}

	return result, err
}

func (b *Base) UpdateFromTable(tx *sqlx.Tx, data map[string]interface{}, where string) (sql.Result, error) {
	var result sql.Result

	if b.table == "" {
		return nil, errors.New("Table must not be empty.")
	}

	tx, wrapInSingleTransaction, err := b.newTransactionIfNeeded(tx)
	if tx == nil {
		return nil, errors.New("Transaction struct must not be empty.")
	}
	if err != nil {
		return nil, err
	}

	keysWithQuestionMarks := make([]string, 0)
	values := make([]interface{}, 0)

	for key, value := range data {
		keysWithQuestionMark := fmt.Sprintf("%v=?", key)
		keysWithQuestionMarks = append(keysWithQuestionMarks, keysWithQuestionMark)
		values = append(values, value)
	}

	query := fmt.Sprintf(
		"UPDATE %v SET %v WHERE %v",
		b.table,
		strings.Join(keysWithQuestionMarks, ","),
		where)

	result, err = tx.Exec(query, values...)

	if err != nil {
		return nil, err
	}

	if wrapInSingleTransaction == true {
		err = tx.Commit()
	}

	return result, err
}

func (b *Base) UpdateByID(tx *sqlx.Tx, data map[string]interface{}, id int64) (sql.Result, error) {
	var result sql.Result

	if b.table == "" {
		return nil, errors.New("Table must not be empty.")
	}

	tx, wrapInSingleTransaction, err := b.newTransactionIfNeeded(tx)
	if tx == nil {
		return nil, errors.New("Transaction struct must not be empty.")
	}
	if err != nil {
		return nil, err
	}

	keysWithQuestionMarks := make([]string, 0)
	values := make([]interface{}, 0)

	for key, value := range data {

		log.Infof("Key: %v Value: %v", key, value)

		keysWithQuestionMark := fmt.Sprintf("%v=?", key)
		keysWithQuestionMarks = append(keysWithQuestionMarks, keysWithQuestionMark)
		values = append(values, value)
	}

	log.Infof("ID Name: %v ID Value: %v", b.tableID, id)

	// Add id as part of values
	values = append(values, id)

	query := fmt.Sprintf(
		"UPDATE %v SET %v WHERE %v=?",
		b.table,
		strings.Join(keysWithQuestionMarks, ","),
		b.tableID)

	log.Infof("Query: %v", query)

	result, err = tx.Exec(query, values...)

	if err != nil {
		return nil, err
	}

	if wrapInSingleTransaction == true {
		err = tx.Commit()
	}

	return result, err
}

func (b *Base) UpdateByKeyValueString(tx *sqlx.Tx, data map[string]interface{}, key, value string) (sql.Result, error) {
	var result sql.Result

	if b.table == "" {
		return nil, errors.New("Table must not be empty.")
	}

	tx, wrapInSingleTransaction, err := b.newTransactionIfNeeded(tx)
	if tx == nil {
		return nil, errors.New("Transaction struct must not be empty.")
	}
	if err != nil {
		return nil, err
	}

	keysWithQuestionMarks := make([]string, 0)
	values := make([]interface{}, 0)

	for key, value := range data {
		keysWithQuestionMark := fmt.Sprintf("%v=?", key)
		keysWithQuestionMarks = append(keysWithQuestionMarks, keysWithQuestionMark)
		values = append(values, value)
	}

	// Add value as part of values
	values = append(values, value)

	query := fmt.Sprintf(
		"UPDATE %v SET %v WHERE %v=?",
		b.table,
		strings.Join(keysWithQuestionMarks, ","),
		key)

	result, err = tx.Exec(query, values...)

	if err != nil {
		return nil, err
	}

	if wrapInSingleTransaction == true {
		err = tx.Commit()
	}

	return result, err
}

func (b *Base) DeleteFromTable(tx *sqlx.Tx, where string) (sql.Result, error) {
	var result sql.Result

	if b.table == "" {
		return nil, errors.New("Table must not be empty.")
	}

	tx, wrapInSingleTransaction, err := b.newTransactionIfNeeded(tx)
	if tx == nil {
		return nil, errors.New("Transaction struct must not be empty.")
	}
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("DELETE FROM %v", b.table)

	if where != "" {
		query = query + " WHERE " + where
	}

	result, err = tx.Exec(query)

	if wrapInSingleTransaction == true {
		err = tx.Commit()
	}

	if err != nil {
		return nil, err
	}

	return result, err
}

func (b *Base) DeleteById(tx *sqlx.Tx, id int64) (sql.Result, error) {

	log.Debugf("DeleteById: Deleting ID from database %v", id)
	var result sql.Result

	if b.table == "" {
		return nil, errors.New("Table must not be empty.")
	}

	tx, wrapInSingleTransaction, err := b.newTransactionIfNeeded(tx)
	if tx == nil {
		return nil, errors.New("Transaction struct must not be empty.")
	}
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("DELETE FROM %v WHERE %v=?", b.table, b.tableID)

	result, err = tx.Exec(query, id)

	if wrapInSingleTransaction == true {
		err = tx.Commit()
	}

	if err != nil {
		return nil, err
	}

	return result, err
}
