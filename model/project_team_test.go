package models

import (
	"fmt"
	"testing"
)

func NewTeamForTest(t *testing.T) (*Team, *User, error) {
	team_row := &TeamRow{}

	team := NewTeam(newDbForTest(t))
	user := userForTesting(t)

	team.Name = fmt.Sprintf("Testing-%v", randomIntforTest())
	team.Description = "testing"
	team.CreatedBy = user.UserID
	team_row, err := team.CreateTeam(nil)
	if err != nil {
		t.Errorf("Project Team: Creating a team should work. Error: %v", err)
	}

	team.TeamRow = *team_row

	return team, user, nil

}

func NewProjectTeamForTest(t *testing.T) *ProjectTeam {
	return NewProjectTeam(newDbForTest(t))
}

func NewProjectForTest(t *testing.T) *Project {
	return NewProject(newDbForTest(t))
}

func TestProjectTeam_CreateProjectTeam(t *testing.T) {

	t.Log("Project Team: Testing Creating Project Team")

	project := NewProjectForTest(t)

	project.Name = fmt.Sprintf("Testing-%v", randomIntforTest())
	project.Difficulty = 1
	project.Description = "Testing Project"

	err := project.CreateProject(nil)
	if err != nil {
		t.Fatalf("Project Team: Creating a project should work. Error: %v", err)
	}

	team, user, err := NewTeamForTest(t)
	if err != nil {
		t.Fatalf("Project Team: Creating a team should work. Error: %v", err)
	}

	pro_team := NewProjectTeamForTest(t)

	pro_team.ProjectID = project.ProjectID
	pro_team.TeamID = team.ID
	pro_team.ProjectURL = "google.com"

	err = pro_team.CreateProjectTeam(nil)

	if err != nil {
		t.Fatalf("Project Team: Creating a project team should work. Error: %v", err)
	}

	//cleaning up
	_, err = pro_team.DeleteById(nil, pro_team.ProjectTeamID)
	if err != nil {
		t.Fatalf("Project Team: Deleting project team by id should not fail. Error: %v", err)
	}
	_, err = project.DeleteById(nil, project.ProjectID)
	if err != nil {
		t.Fatalf("Project Team: Deleting project by id should not fail. Error: %v", err)
	}

	_, err = team.DeleteById(nil, team.ID)
	if err != nil {
		t.Fatalf("Project Team: Deleting team by id should not fail. Error: %v", err)
	}
	_, err = user.DeleteById(nil, user.UserID)
	if err != nil {
		t.Fatalf("Project Team: Deleting user by id should not fail. Error: %v", err)
	}
	t.Log("Project Team: TestProjectTeam_CreateProjectTeam completed")

}
