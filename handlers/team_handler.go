package handlers

import (
	_ "github.com/go-sql-driver/mysql"

	"errors"
	"github.com/strongjz/leveledup-api/models"
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

	t,err := team.CreateTeam(nil)
	if err != nil {
		log.Errorf("Error creating team %s", err)
		error := "Error creating team"
		c.JSON(400, gin.H{"error": error})
		return
	}

	c.JSON(http.StatusOK, t)
}

func (h *ApiResource) TeamGet(c *gin.Context) {
	log.Debug("TeamGet Endpoint")
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *ApiResource) TeamUpdate(c *gin.Context) {
	log.Debug("TeamUpdate Endpoint")
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *ApiResource) TeamDelete(c *gin.Context) {
	log.Debug("TeamDelete Endpoint")
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

