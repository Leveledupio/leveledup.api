package models

import "testing"

func newTeamMembersTest(t *testing.T) *TeamMembers {
	return NewTeamMembers(newDbForTest(t))
}

func TestNewTeamMembers(t *testing.T) {

	team, user, err := NewTeamForTest(t)
	if err != nil {
		t.Fatalf("Project Team: Creating a team should work. Error: %v", err)
	}

	t.Logf("[DEBUG][TestNewTeamMembers] Team Name %v", team.Name)

	users := []*User{}

	userA := newUserForTest(t)
	users = append(users, user)
	users = append(users, userA)

	for i := 0; i < 10; i++ {
		userA = newUserForTest(t)
		users = append(users, userA)
	}

	for _, u := range users {
		t.Logf("[DEBUG][TestNewTeamMembers] User Email: %v", u.Email)

	}

}
