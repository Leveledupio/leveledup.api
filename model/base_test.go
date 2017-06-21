package models

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/strongjz/gbs/libstring"
	//"os"
	//"github.com/spf13/viper"

	//"strings"

	"testing"
)

func newEmailForTest() string {
	return fmt.Sprintf("user-%v@example.com", libstring.RandString(32))
}

func randomIntforTest() string {
	return fmt.Sprintf(libstring.RandString(32))
}
func newDbForTest(t *testing.T) *sqlx.DB {
	/*
		configName := os.Getenv("ENV") + "-config.yaml"

		t.Logf("Testing Config %s", configName)

		config := viper.New()

		config.AddConfigPath("./config/")

		config.SetConfigFile(configName)
		config.SetConfigType("yaml")

		config.Debug()

		dsn := config.GetString("dsn")

		t.Logf("DSN: %s\n", dsn)
	*/
	db, err := sqlx.Connect("mysql", "gbs:gbs@tcp(localhost:3306)/gbs")
	if err != nil {
		t.Fatalf("Connecting to local MySQL should never fail. Error: %v", err)
	}

	db.DB.Ping()
	return db
}

func newBaseForTest(t *testing.T) *Base {
	base := &Base{}
	base.db = newDbForTest(t)

	return base
}

func TestNewTransactionIfNeeded(t *testing.T) {
	base := newBaseForTest(t)

	// New Transaction block
	tx, wrapInSingleTransaction, err := base.newTransactionIfNeeded(nil)
	if err != nil {
		t.Fatalf("Creating new transaction block should not fail. Error: %v", err)
	}
	if wrapInSingleTransaction != true {
		t.Fatalf("Creating new transaction block should set wrapInSingleTransaction == true.")
	}
	if tx == nil {
		t.Fatalf("Creating new transaction block should not fail. Error: %v", err)
	}

	// Existing Transaction block
	tx2, wrapInSingleTransaction, err := base.newTransactionIfNeeded(tx)
	if err != nil {
		t.Fatalf("Receiving existing transaction block should not fail. Error: %v", err)
	}
	if wrapInSingleTransaction != false {
		t.Fatalf("Receiving existing transaction block should set wrapInSingleTransaction == false.")
	}
	if tx2 == nil {
		t.Fatalf("Receiving existing transaction block should not fail. Error: %v", err)
	}
	if tx2 != tx {
		t.Fatalf("Receiving existing transaction block should not fail. Error: %v", err)
	}
}

func TestCreateDeleteGeneric(t *testing.T) {
	base := newBaseForTest(t)
	base.table = "user"
	base.tableID = "user_id"

	// INSERT INTO user (email) VALUES (...)
	data := make(map[string]interface{})
	data["email"] = newEmailForTest()
	data["password"] = "abc123"
	data["first_name"] = "joe"
	data["last_name"] = "smith"
	data["github_name"] = "jsmith"

	result, err := base.InsertIntoTable(nil, data)
	if err != nil {
		t.Fatalf("Inserting new row should not fail. Error: %v", err)
	}

	lastInsertedId, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("Inserting new row should not fail. Error: %v", err)
	}

	// DELETE WHERE user_id=...
	where := fmt.Sprintf("user_id=%v", lastInsertedId)

	_, err = base.DeleteFromTable(nil, where)
	if err != nil {
		t.Fatalf("Deleting row by id should not fail. Error: %v", err)
	}

}

func TestCreateDeleteById(t *testing.T) {
	base := newBaseForTest(t)
	base.table = "user"

	data := make(map[string]interface{})
	data["email"] = newEmailForTest()
	data["password"] = "abc123"
	data["first_name"] = "joe"
	data["last_name"] = "smith"
	data["github_name"] = "jsmith"

	result, err := base.InsertIntoTable(nil, data)
	if err != nil {
		t.Fatalf("Inserting new row should not fail. Error: %v", err)
	}

	lastInsertedId, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("Inserting new row should not fail. Error: %v", err)
	}

	// DELETE WHERE id=...
	_, err = base.DeleteById(nil, lastInsertedId)
	if err != nil {
		t.Fatalf("Deleting user by id should not fail. Error: %v", err)
	}

}

func TestCreateUpdateGenericDelete(t *testing.T) {
	base := newBaseForTest(t)
	base.table = "user"
	base.tableID = "user_id"

	// INSERT INTO user (...) VALUES (...)
	data := make(map[string]interface{})
	data["email"] = newEmailForTest()
	data["password"] = "abc123"
	data["first_name"] = "joe"
	data["last_name"] = "smith"
	data["github_name"] = "jsmith"
	data["date_became_customer"] = base.todayDate()

	result, err := base.InsertIntoTable(nil, data)
	if err != nil {
		t.Fatalf("Inserting new row should not fail. Error: %v", err)
	}

	lastInsertedId, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("Inserting new row should not fail. Error: %v", err)
	}

	// UPDATE email=$1 WHERE id=$2
	data["email"] = "yo dawg"
	where := fmt.Sprintf("user_id=%v", lastInsertedId)

	_, err = base.UpdateFromTable(nil, data, where)
	if err != nil {
		t.Errorf("Updating existing row should not fail. Error: %v", err)
	}

	// DELETE WHERE id=...
	_, err = base.DeleteById(nil, lastInsertedId)
	if err != nil {
		t.Fatalf("Deleting row by id should not fail. Error: %v", err)
	}

}

func TestCreateUpdateByIDDelete(t *testing.T) {
	base := newBaseForTest(t)
	base.table = "user"
	base.tableID = "user_id"

	// INSERT INTO user (...) VALUES (...)
	data := make(map[string]interface{})
	data["email"] = newEmailForTest()
	data["password"] = "abc123"
	data["first_name"] = "joe"
	data["last_name"] = "smith"
	data["github_name"] = "jsmith"
	data["date_became_customer"] = base.todayDate()

	result, err := base.InsertIntoTable(nil, data)
	if err != nil {
		t.Fatalf("Inserting new row should not fail. Error: %v", err)
	}

	lastInsertedId, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("Inserting new row should not fail. Error: %v", err)
	}

	// UPDATE SET email=$1 WHERE id=$2
	data["email"] = "yo dawg"

	_, err = base.UpdateByID(nil, data, lastInsertedId)
	if err != nil {
		t.Errorf("Updating existing row should not fail. Error: %v", err)
	}

	// DELETE WHERE id=...
	_, err = base.DeleteById(nil, lastInsertedId)
	if err != nil {
		t.Fatalf("Deleting row by id should not fail. Error: %v", err)
	}

}

func TestCreateUpdateByKeyValueStringDelete(t *testing.T) {
	base := newBaseForTest(t)
	base.table = "user"
	base.tableID = "user_id"

	originalEmail := newEmailForTest()

	// INSERT INTO user (...) VALUES (...)
	data := make(map[string]interface{})
	data["email"] = originalEmail
	data["password"] = "abc123"
	data["first_name"] = "joe"
	data["last_name"] = "smith"
	data["github_name"] = "jsmith"
	data["date_became_customer"] = base.todayDate()

	result, err := base.InsertIntoTable(nil, data)
	if err != nil {
		t.Fatalf("Inserting new row should not fail. Error: %v", err)
	}

	lastInsertedId, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("Inserting new row should not fail. Error: %v", err)
	}

	// UPDATE SET email=$1 WHERE id=$2
	data["email"] = newEmailForTest()

	_, err = base.UpdateByKeyValueString(nil, data, "email", originalEmail)
	if err != nil {
		t.Errorf("Updating existing user should not fail. Error: %v", err)
	}

	// DELETE FROM user WHERE id=...
	_, err = base.DeleteById(nil, lastInsertedId)
	if err != nil {
		t.Fatalf("Deleting user by id should not fail. Error: %v", err)
	}

}
