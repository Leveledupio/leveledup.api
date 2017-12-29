package models

import (
	"database/sql"
	"errors"
	"fmt"
	"net/mail"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

var (
	UserID        = "user_id"
	UserEmail     = "email"
	Password      = "password"
	PasswordAgain = "password"
	FirstName     = "first_name"
	LastName      = "last_name"
	GithubName    = "github_name"
	SlackName     = "slack_name"
	DateCustomer  = "date_became_customer"
	UserTable     = "user"
)

func NewUser(db *sqlx.DB) *User {
	user := &User{}
	user.db = db
	user.table = UserTable
	user.hasID = true
	user.tableID = UserID
	return user
}

type UserRow struct {
	UserID        int64  `db:"user_id" json:"user_id,omitempty"`
	Email         string `db:"email" json:"email,omitempty"`
	Password      string `db:"password" json:"password"`
	PasswordAgain string `json:"password_again"`
	FirstName     string `db:"first_name" json:"first_name,omitempty"`
	LastName      string `db:"last_name" json:"last_name,omitempty"`
	GithubName    string `db:"github_name" json:"github_name,omitempty"`
	SlackName     string `db:"slack_name" json:"slack_name,omitempty"`
	DateCustomer  string `db:"date_became_customer" json:"data_became_customer,omitempty"`
}

func (u *User) PrintUser() string {
	return fmt.Sprintf("Email: '%s'  First Name: '%s' Last Name: '%s' Github: '%s'", u.Email, u.FirstName, u.LastName, u.GithubName)
}

type User struct {
	Base
	UserRow
}

func (u *User) userRowFromSqlResult(tx *sqlx.Tx, sqlResult sql.Result) (*UserRow, error) {

	userId, err := sqlResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	log.Debugf("userRowFromSqlResult SQL Result User ID %v", userId)

	return u.GetUserById(tx, userId)
}

// AllUsers returns all user rows.
func (u *User) AllUsers(tx *sqlx.Tx) ([]*UserRow, error) {
	users := []*UserRow{}
	query := fmt.Sprintf("SELECT * FROM %v", u.table)
	err := u.db.Select(&users, query)
	if err != nil {
		return nil, err
	}

	return users, err
}

// GetById returns record by id.
func (u *User) GetUserById(tx *sqlx.Tx, id int64) (*UserRow, error) {

	user := &UserRow{}

	query := fmt.Sprintf("SELECT * FROM %v WHERE %v=?", u.table, u.tableID)

	err := u.db.Get(user, query, id)

	if err != nil {
		return nil, err
	}

	u.PasswordAgain = ""
	u.Password = ""
	return user, err
}

// GetByEmail returns record by email.
func (u *User) GetByEmail(tx *sqlx.Tx, email string) (*UserRow, error) {
	user := &UserRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE email=?", u.table)
	err := u.db.Get(user, query, email)
	if err != nil {
		return nil, err
	}

	return user, err
}

// GetByEmail returns record by email but checks password first.
func (u *User) GetUserByEmailAndPassword(tx *sqlx.Tx, email, password string) (*UserRow, error) {

	log.Debugf("Model GetUserByEmailAndPassword Email: %s", email)
	user, err := u.GetByEmail(tx, email)
	if err != nil {
		return nil, err
	}

	log.Debugf("Model GetUserByEmailAndPassword Email: %s USERID: %d", user.Email, user.UserID)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Debugf("model GetUserByEmailAndPassword email: %s bcrypt error: %s", email, err)
		return nil, err
	}

	//only the model needs access to the password, the controller does not
	user.Password = ""
	user.PasswordAgain = ""

	return user, err
}

func (u *User) DeleteUser(tx *sqlx.Tx, email, password string) error {

	log.Debugf("Model DeleteUser Email: %s", email)

	user, err := u.GetUserByEmailAndPassword(tx, email, password)
	if err != nil {
		log.Debugf("error retrieving GetUserByEmailAndPassword %s", err)
		return err
	}
	log.Debugf("Model Deleting Retrieving user Email: %s ID: %v", user.Email, user.UserID)

	_, err = u.DeleteById(nil, user.UserID)
	if err != nil {
		log.Debugf("Error DeleteById %s", err)
		return err
	}

	return nil

}

// Signup create a new record of user.
func (u *User) Signup(tx *sqlx.Tx) (*UserRow, error) {

	log.Debugf("Handler User Signup %v", u.PrintUser())

	if u.Email == "" {
		return nil, errors.New("Email cannot be blank.")
	}

	e, err := mail.ParseAddress(u.Email)
	if err != nil {
		return nil, errors.New("Email is invalid format")
	}

	if u.Password == "" {
		return nil, errors.New("Password cannot be blank.")
	}
	if u.Password != u.PasswordAgain {
		return nil, errors.New("Password is invalid.")
	}

	if u.FirstName == "" {
		return nil, errors.New("First Name is invalid.")
	}
	if u.LastName == "" {
		return nil, errors.New("Last name is invalid.")
	}
	if u.GithubName == "" {
		return nil, errors.New("Github Name is invalid.")
	}
	if u.SlackName == "" {
		u.SlackName = u.GithubName
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 5)
	if err != nil {

		return nil, err
	}

	data := make(map[string]interface{})
	data[UserEmail] = e.Address //ParseAddress returns e
	data[Password] = hashedPassword
	data[FirstName] = u.FirstName
	data[LastName] = u.LastName
	data[GithubName] = u.GithubName
	data[SlackName] = u.SlackName
	data[DateCustomer] = u.todayDate()

	sqlResult, err := u.InsertIntoTable(tx, data)
	if err != nil {
		log.Errorf("Handler User Signup Error after insert into table %v", err)
		return nil, err
	}

	return u.userRowFromSqlResult(tx, sqlResult)
}

// UpdateEmailAndPasswordById updates user email and password.
func (u *User) UpdateEmailAndPasswordById(tx *sqlx.Tx, userId int64, email, password, passwordAgain string) (*UserRow, error) {
	data := make(map[string]interface{})

	if email != "" {
		data[UserEmail] = email
	}

	if password != "" && passwordAgain != "" && password == passwordAgain {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 5)
		if err != nil {
			return nil, err
		}

		data[Password] = hashedPassword
	}

	if len(data) > 0 {
		_, err := u.UpdateByID(tx, data, userId)
		if err != nil {
			return nil, err
		}
	}

	return u.GetUserById(tx, userId)
}

func (u *User) UpdateUser(tx *sqlx.Tx) (*UserRow, error) {

	data := make(map[string]interface{})
	data[UserEmail] = u.Email
	data[FirstName] = u.FirstName
	data[LastName] = u.LastName
	data[GithubName] = u.GithubName
	data[SlackName] = u.SlackName

	sqlResult, err := u.InsertIntoTable(tx, data)
	if err != nil {
		return nil, err
	}

	return u.userRowFromSqlResult(tx, sqlResult)
}
