package handlers

import (
	//Mysql import
	_ "github.com/go-sql-driver/mysql"

	"errors"
	"net/http"

	"github.com/Leveledupio/leveledup.api/models"
	"github.com/gin-gonic/gin"
)

//Contact - Create a new contact
func (h *APIResource) Contact(c *gin.Context) {

	log.Debugf("Contact Handler")

	email := models.NewEmail(h.DB, h.AWSSession)

	err := c.Bind(&email.EmailRow)
	if err != nil {
		log.Errorf("Problem decoding JSON body %s", err)
		c.JSON(400, errors.New("problem decoding body"))
		return
	}

	email.EmailTo = "clientservices@leveledup.io"
	email.Subject = "Contact Us Form"

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
