package handlers

import (
	_ "github.com/go-sql-driver/mysql"

	"errors"
	"github.com/strongjz/leveledup.api/models"
	"github.com/gin-gonic/gin"
	"net/http"

)


func (h *ApiResource) TeamCreate(c *gin.Context) {
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

	log.Debugf("Name %s ID %v Created by %v , Description %v",t.Name, t.ID, t.CreatedBy, t.Description)

	c.JSON(http.StatusOK, t)
}

func (h *ApiResource) TeamGet(c *gin.Context) {
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

func (h *ApiResource) TeamUpdate(c *gin.Context) {
	log.Debug("TeamUpdate Endpoint")
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *ApiResource) TeamDelete(c *gin.Context) {
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

func (h *ApiResource) ProjectTeam(c *gin.Context) {


}

