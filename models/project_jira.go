package models

import (
	"errors"
	"fmt"
	jira "github.com/andygrunwald/go-jira"
)

func (p *Project) createproject() error {

	log.Debugf("[DEBUG] Creating Project %s", p.Name)
	log.Debugf("Printing Project")
	p.PrintProject()

	pro := jira.Project{}
	pro.Name = p.Name
	pro.Description = p.Description

	log.Infof("[INFO] createproject Jira Base URL %v",
		p.jira.GetBaseURL().Host)

	req, err := p.jira.NewRequest("POST", "/rest/api/2/project", nil)
	if err != nil {

		panic(err)

	}

	log.Debugf("[DEBUG] Request %v", req)
	_, err = p.jira.Do(req, pro)
	if err != nil {
		log.Errorf("[ERROR] Creating Project: %v", err)
		return err
	}

	defer func() {
		// recover from panic if one occurred. Set err to nil otherwise.
		if recover() != nil {
			err = errors.New("error creating project")
			log.Errorf("[ERROR] Creating Project: %v", err)
		}
	}()
	return nil
}

func (p *Project) getAllProjects() ([]ProjectRow, error) {

	log.Debugf("[DEBUG] Creating Project %s", p.Name)
	log.Debugf("Printing Project")

	req, _ := p.jira.NewRequest("GET", "/rest/api/2/project", nil)

	projects := new([]jira.Project)
	_, err := p.jira.Do(req, projects)
	if err != nil {
		panic(err)
	}

	plist := []ProjectRow{}
	pRow := &ProjectRow{}

	for _, project := range *projects {
		log.Debugf("[DEBUG] %s: %s\n", project.Key, project.Name)

		log.Debugf("[DEBUG]\n %v/n", project)

		pRow.Name = project.Name
		pRow.Description = fmt.Sprintf(project.Description, project.Key, project.ID)

		plist = append(plist, *pRow)

		pRow = &ProjectRow{}
	}

	return plist, err
}
