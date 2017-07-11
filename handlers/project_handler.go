package handlers

import (
	_ "github.com/go-sql-driver/mysql"

	"errors"
	"github.com/strongjz/leveledup-api/models"
	"github.com/gin-gonic/gin"
	"net/http"

)

func (h *ApiResource) ProjectCreate(c *gin.Context) {
	log.Debug("ProjectCreate Endpoint")
	Project := models.NewProject(h.DB)


	err := c.Bind(&Project.ProjectRow)
	if err != nil {
		log.Errorf("Problem decoding JSON body %s", err)
		c.JSON(400, errors.New("problem decoding body"))
		return
	}

	err = Project.CreateProject(nil)
	if err != nil {
		log.Errorf("Error creating Project %s", err)
		error := "Error creating Project"
		c.JSON(400, gin.H{"error": error})
		return
	}

	c.JSON(http.StatusOK, Project.ProjectRow)
}

func (h *ApiResource) ProjectGet(c *gin.Context) {
	log.Debug("ProjectGet Endpoint")
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *ApiResource) ProjectUpdate(c *gin.Context) {
	log.Debug("ProjectUpdate Endpoint")
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *ApiResource) ProjectDelete(c *gin.Context) {
	log.Debug("ProjectDelete Endpoint")
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}