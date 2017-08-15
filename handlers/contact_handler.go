package handlers

import (
	_ "github.com/go-sql-driver/mysql"

	"errors"
	"github.com/strongjz/leveledup-api/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *ApiResource) Contact(c *gin.Context) {

	log.Debugf("Contact Handler")

	email := models.NewEmail(h.DB, h.AWSSession)


	err := c.Bind(&email.EmailRow)
	if err != nil {
		log.Errorf("Problem decoding JSON body %s", err)
		c.JSON(400, errors.New("problem decoding body"))
		return
	}

	log.Debugf("Contact Handler Email Subject: %s", email.Subject)
	log.Debugf("Contact Handler Email to: %s", email.EmailTo)

	err = email.SendEmail()
	if err != nil {
		log.Errorf("Contact Handler Error %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, "")
}