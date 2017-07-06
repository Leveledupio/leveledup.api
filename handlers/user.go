package handlers

import (
	//log "gopkg.in/op/go-logging.v1"
	_ "github.com/go-sql-driver/mysql"

	"errors"
	model "github.com/strongjz/leveledup-api/model"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
)

var (
	UserID        = "user_id"
	Email         = "email"
	Password      = "password"
	PasswordAgain = "password"
	FirstName     = "first_name"
	LastName      = "last_name"
	GithubName    = "github_name"
	SlackName     = "slack_name"
	DateCustomer  = "date_became_customer"
	UserTable     = "user"
)

func (h *ApiResource) UserUpdate(c *gin.Context) {

	log.Debugf("Handler UserUpdate")

	user := model.NewUser(h.DB)

	email := c.Param("email")

	err := c.Bind(&user.UserRow)
	if err != nil {
		//	log.Errorf("Problem decoding JSON body %s", err)
		c.JSON(400, errors.New("problem decoding body"))
		return
	}

	u, err := user.GetByEmail(nil, email)
	if err != nil {

		c.JSON(400, errors.New("Email not found"))
		return
	}

	u, err = user.UpdateUser(nil)

	log.Debugf("User Row: %v", u)

	//zero out the password
	u.Password = ""
	u.PasswordAgain = ""

	c.JSON(http.StatusOK, u)

}

func (h *ApiResource) UserRetrieve(c *gin.Context) {

	log.Debugf("Handler UserRetrieve")

	user := model.NewUser(h.DB)

	email := c.Param("email")

	err := c.Bind(&user.UserRow)
	if err != nil {
		//	log.Errorf("Problem decoding JSON body %s", err)
		c.JSON(400, errors.New("problem decoding body"))
		return
	}

	u, err := user.GetByEmail(nil, email)
	if err != nil {

		c.JSON(400, errors.New("Email not found"))
		return
	}

	//zero out the password
	u.Password = ""
	u.PasswordAgain = ""

	c.JSON(http.StatusOK, u)
}

func (h *ApiResource) GetUserEP(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *ApiResource) UserDelete(c *gin.Context) {
	log.Debugf("Handler UserSignup")

	user := model.NewUser(h.DB)

	err := c.Bind(&user.UserRow)
	if err != nil {
		//	log.Errorf("Problem decoding JSON body %s", err)
		c.JSON(400, errors.New("problem decoding body"))
		return
	}

	log.Debugf("Deleting User Email %s, Password: %s", user.Email, user.Password)

	err = user.DeleteUser(nil, user.Email, user.Password)
	if err != nil {
		log.Errorf("Error deleting user %s err %s", user.Email, err)
		c.JSON(500, errors.New("Deleting user email"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *ApiResource) UserSignup(c *gin.Context) {

	log.Debugf("Handler UserSignup")

	user := model.NewUser(h.DB)

	err := c.Bind(&user.UserRow)
	if err != nil {
		//	log.Errorf("Problem decoding JSON body %s", err)
		c.JSON(400, errors.New("problem decoding body"))
		return
	}

	log.Debugf(user.PrintUser())

	u, err := user.Signup(nil)
	if err != nil {

		c.JSON(500, errors.New("Error Signing up user"))
		return
	}

	//zero out the password
	u.Password = ""
	u.PasswordAgain = ""

	c.JSON(http.StatusOK, u)
}
