package models

import (
	jira "github.com/andygrunwald/go-jira"
)

func (p *Project) createproject() error {

	log.Debugf("[DEBUG] Creating Project %s", p.Name)
	log.Debugf("Printing Project")
	p.PrintProject()

	pro := jira.Project{}
	pro.Name = p.Name
	pro.Description = p.Description

	req, err := p.jira.NewRequest("POST", "/rest/api/2/project", nil)
	if err != nil {
		log.Errorf("[ERROR] Creating Project: %v", err)
		return err
	}

	_, err = p.jira.Do(req, pro)
	if err != nil {
		log.Errorf("[ERROR] Creating Project: %v", err)
		return err
	}

	return nil
}
