package models

import 	(
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

var (
	UserID = "user_id"
	Email = "email"
	Password = "password"
	PasswordAgain = "password"
	FirstName = "first_name"
	LastName = "last_name"
	GithubName = "github_name"
	SlackName = "slack_name"
	DateCustomer = "date_became_customer"
	UserTable = "user"
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
	UserID        int64  `db:"user_id"`
	Email         string `db:"email"`
	Password      string `db:"password"`
	PasswordAgain string
	FirstName     string `db:"first_name"`
	LastName      string `db:"last_name"`
	GithubName    string `db:"github_name"`
	SlackName     string `db:"slack_name"`
	DateCustomer  string `db:"date_became_customer"`
}

func (u *User) PrintUser() string {
	return fmt.Sprintf("Email: %s  First Name: %s Last Name %s Github: %s", u.Email, u.FirstName, u.LastName, u.GithubName)
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

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	return user, err
}

// Signup create a new record of user.
func (u *User) Signup(tx *sqlx.Tx) (*UserRow, error) {

	log.Debugf("Handler User Signup %v", u.PrintUser())

	if u.Email == "" {
		return nil, errors.New("Email cannot be blank.")
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
	data[Email] = u.Email
	data[Password] = hashedPassword
	data[FirstName] = u.FirstName
	data[LastName] = u.LastName
	data[GithubName] = u.GithubName
	data[SlackName] = u.SlackName
	data[DateCustomer] = u.todayDate()

	sqlResult, err := u.InsertIntoTable(tx, data)
	if err != nil {
		return nil, err
	}

	return u.userRowFromSqlResult(tx, sqlResult)
}

// UpdateEmailAndPasswordById updates user email and password.
func (u *User) UpdateEmailAndPasswordById(tx *sqlx.Tx, userId int64, email, password, passwordAgain string) (*UserRow, error) {
	data := make(map[string]interface{})

	if email != "" {
		data[Email] = email
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
