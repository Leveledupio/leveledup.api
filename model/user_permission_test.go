package models

import (
	"testing"
)

func newPermissionForTest(t *testing.T) *Permission {
	return NewPermission(newDbForTest(t))
}

func newUserPermissionForTest(t *testing.T) *UserPermission {
	return NewUserPermission(newDbForTest(t))

}
func TestNewUserPermission(t *testing.T) {

	user := userForTesting(t)

	permission := newPermissionForTest(t)

	permission.Name = "Scrum Master"
	permission.PermissionRole = "project"

	err := permission.CreatePermission(nil)
	if err != nil {
		t.Fatalf("Creating a Permission should not fail %v", err)
	}

	user_permission := newUserPermissionForTest(t)

	user_permission.PermissionID = permission.PermissionID
	user_permission.UserID = user.UserID

	err = user_permission.CreateUserPermission(nil)
	if err != nil {
		t.Fatalf("Creating a User permission should not fail %v", err)
	}

	_, err = user.DeleteById(nil, user.UserID)
	if err != nil {
		t.Fatalf("Deleting a User should not fail %v", err)
	}

	_, err = permission.DeleteById(nil, permission.PermissionID)
	if err != nil {
		t.Fatalf("Deleting a Permission should not fail %v", err)
	}

	_, err = user_permission.DeleteById(nil, user_permission.PermissionID)
	if err != nil {
		t.Fatalf("Deleting a User permission should not fail %v", err)
	}

}
