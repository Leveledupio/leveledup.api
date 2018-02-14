package handlers

import (
	//Mysql requires a blank import
	_ "github.com/go-sql-driver/mysql"

	"errors"
	"net/http"

	"github.com/Leveledupio/leveledup.api/models"
	"github.com/gin-gonic/gin"
)

//TeamCreate - API endpoint to create a team
func (h *APIResource) TeamCreate(c *gin.Context) {
	log.Debug("TeamCreate Endpoint")
	team := models.NewTeam(h.DB)

	err := c.Bind(&team.TeamRow)
	if err != nil {
		log.Errorf("Problem decoding JSON body %s", err)
		c.JSON(400, errors.New("problem decoding body"))
		return
	}

	t, err := team.CreateTeam(nil)
	if err != nil {
		log.Errorf("Error creating team %s", err)
		error := "Error creating team"
		c.JSON(400, gin.H{"error": error})
		return
	}

	log.Debugf("Name %s ID %v Created by %v , Description %v", t.Name, t.ID, t.CreatedBy, t.Description)

	c.JSON(http.StatusOK, t)
}

//TeamGet - API endpoint to create a new team
func (h *APIResource) TeamGet(c *gin.Context) {
	log.Debug("TeamGet Endpoint")

	team := models.NewTeam(h.DB)

	teamName := c.Param("team")

	t, err := team.GetTeamByName(nil, teamName)
	if err != nil {
		log.Errorf("Error Retrieving team %s", err)
		error := "Error Retrieving Team"
		c.JSON(400, gin.H{"error": error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "team": t})
}

//TeamUpdate - API endpoint to update a team
func (h *APIResource) TeamUpdate(c *gin.Context) {
	log.Debug("TeamUpdate Endpoint")
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

//TeamDelete - API endpoint to delete a team
func (h *APIResource) TeamDelete(c *gin.Context) {
	log.Debug("TeamDelete Endpoint")
	team := models.NewTeam(h.DB)

	teamName := c.Param("team")

	t, err := team.GetTeamByName(nil, teamName)
	if err != nil {
		log.Errorf("Error Retrieving team %s", err)
		error := "Team does not exist"
		c.JSON(404, gin.H{"error": error})
		return
	}

	_, err = team.DeleteById(nil, t.ID)
	if err != nil {
		log.Errorf("Error Retrieving team %s", err)
		error := "Error removing Team"
		c.JSON(404, gin.H{"error": error})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

//ProjectTeam - API endpoint to tie a project to a team
func (h *APIResource) ProjectTeam(c *gin.Context) {

}
