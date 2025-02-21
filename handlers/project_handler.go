package handlers

import (
	//Mysql reqiures a blank import
	_ "github.com/go-sql-driver/mysql"

	"errors"
	"net/http"

	"github.com/Leveledupio/leveledup.api/models"
	"github.com/gin-gonic/gin"
)

//ProjectCreate - Creates projects
func (h *APIResource) ProjectCreate(c *gin.Context) {
	log.Debug("ProjectCreate Endpoint")

	log.Debug("BASE URL %v", h.Jira.GetBaseURL())

	Project := models.NewProject(h.DB, h.Jira)

	err := c.Bind(&Project.ProjectRow)
	if err != nil {
		log.Errorf("Problem decoding JSON body %s", err)
		c.JSON(400, errors.New("problem decoding body"))
		return
	}

	err = Project.CreateProject()
	if err != nil {
		log.Errorf("Error creating Project %s", err)
		error := "Error creating Project"
		c.JSON(400, gin.H{"error": error})
		return
	}

	c.JSON(http.StatusOK, Project.ProjectRow)
}

//ProjectGet API endpoint to get a project
func (h *APIResource) ProjectGet(c *gin.Context) {
	log.Debug("ProjectGet Endpoint")
	Project := models.NewProject(h.DB, h.Jira)

	project := c.Param("project")

	Project.Name = project

	err := c.Bind(&Project.ProjectRow)
	if err != nil {
		log.Errorf("Problem decoding JSON body %s", err)
		c.JSON(400, errors.New("problem decoding body"))
		return
	}

	err = Project.GetProjectByName(nil)
	if err != nil {
		log.Errorf("Error retrieving Project Name: %s %s", Project.Name, err)
		error := "Error Project Does not exist"

		c.JSON(400, gin.H{"error": error})

		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "project": Project.ProjectRow})
}

//ProjectGetAll - API endpoint to get all projects
func (h *APIResource) ProjectGetAll(c *gin.Context) {

	log.Debug("ProjectGetall Endpoint")
	Project := models.NewProject(h.DB, h.Jira)

	projects, err := Project.GetAllProject()
	if err != nil {
		log.Errorf("Error retrieving Project Name: %s %s", Project.Name, err)
		error := "Error Project Does not exist"

		c.JSON(400, gin.H{"error": error})

		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "projects": projects})

}

//ProjectUpdate - API endpoint to update a project
func (h *APIResource) ProjectUpdate(c *gin.Context) {
	log.Debug("ProjectUpdate Endpoint")
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

//ProjectDelete - API endpoint to delete a project
func (h *APIResource) ProjectDelete(c *gin.Context) {
	log.Debug("ProjectDelete Endpoint")
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
