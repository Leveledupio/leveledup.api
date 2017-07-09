package models

import (
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

func newUserForTest(t *testing.T) *User {
	return NewUser(newDbForTest(t))
}

func TestUserCRUD(t *testing.T) {

	t.Log("TestUserCRUD")
	u := newUserForTest(t)

	// Signup
	u.GithubName = "TestUserCRUD"
	u.Password = "abc123"
	u.PasswordAgain = u.Password
	u.FirstName = "Jeff"
	u.LastName = "Smith"

	u.Email = newEmailForTest()

	t.Logf("TestUserCRUD User: User for Test: %s", u.UserID)

	userRow, err := u.Signup(nil)

	if err != nil {
		t.Errorf("User: Signing up user should work. Error: %v", err)
	}
	if userRow == nil {
		t.Fatal("User: Signing up user should work.")
	}
	if userRow.UserID <= 0 {
		t.Fatal("User: signing up user should work.")
	}

	t.Log("TestUserCRUD User: DELETE FROM users WHERE id=%s", userRow.UserID)
	_, err = u.DeleteById(nil, userRow.UserID)
	if err != nil {
		t.Fatalf("Deleting user by id should not fail. Error: %v", err)
	}

}

func userForTesting(t *testing.T) *User {
	t.Log("userForTesting User: Creating user for testing")
	u := newUserForTest(t)

	// Signup
	u.GithubName = "userForTesting"
	u.Password = "abc123"
	u.PasswordAgain = u.Password
	u.FirstName = "Jeff"
	u.LastName = "Smith"
	u.DateCustomer = u.todayDate()
	u.Email = newEmailForTest()

	userRow, err := u.Signup(nil)
	t.Log("userForTesting userRow %v", userRow)

	t.Logf("userForTesting User for Test: Email %s ID %v", u.Email, userRow.UserID)

	if err != nil {
		t.Errorf("User: Signing up user should work. Error: %v", err)
	}
	if userRow == nil {
		t.Fatal("User: Signing up user should work.")
	}
	if userRow.UserID <= 0 {
		t.Fatal("User: Signing up user should work.")
	}

	u.UserRow = *userRow

	t.Logf("userForTesting Copy userRow to user %v", u.PrintUser())

	return u
}
