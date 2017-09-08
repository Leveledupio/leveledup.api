package handlers

import (
	_ "github.com/go-sql-driver/mysql"

	"errors"
	"github.com/strongjz/leveledup.api/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *ApiResource) UserLogin(c *gin.Context) {

	log.Debugf("Handler Login")

	user := models.NewUser(h.DB)

	err := c.Bind(&user.UserRow)
	if err != nil {
		log.Errorf("Problem decoding JSON body %s", err)
		c.JSON(400, errors.New("problem decoding body"))
		return
	}

	log.Debugf("Handler Login User Email: %s", user.Email)
	log.Debugf("Handler Login User Password: %s", user.Password)

	u, err := user.GetUserByEmailAndPassword(nil, user.Email, user.Password)
	if err != nil {
		log.Errorf("Username and Password did not match %s", err)
		error := "Username and Password did not match"
		c.JSON(401, gin.H{"error": error})
		return
	}

	//zero out the password
	u.Password = ""
	u.PasswordAgain = ""

	c.JSON(http.StatusOK, u)
}

func (h *ApiResource) UserUpdate(c *gin.Context) {

	log.Debugf("Handler UserUpdate")

	user := models.NewUser(h.DB)

	email := c.Param("email")

	err := c.Bind(&user.UserRow)
	if err != nil {
		log.Errorf("Problem decoding JSON body %s", err)
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

	user := models.NewUser(h.DB)

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

func (h *ApiResource) UserDelete(c *gin.Context) {
	log.Debugf("Handler UserSignup")

	user := models.NewUser(h.DB)


	err := c.Bind(&user.UserRow)

	log.Debugf("User Data print after c.Bind %s", user.PrintUser())

	if err != nil {
		log.Errorf("Problem decoding JSON body %s", err)
		c.JSON(400, errors.New("problem decoding body"))
		return
	}

	log.Debugf("Deleting User Email %s, Password: %s", user.Email, user.Password)

	err = user.DeleteUser(nil, user.Email, user.Password)
	if err != nil {
		log.Errorf("Error deleting user %s err %s", user.Email, err)

		if err.Error() == "sql: no rows in result set" {
			c.JSON(404, errors.New("User does not exist"))
			return
		}

		c.JSON(400, errors.New("Deleting user email"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *ApiResource) UserSignup(c *gin.Context) {

	log.Debugf("Handler UserSignup")

	user := models.NewUser(h.DB)

	err := c.Bind(&user.UserRow)
	if err != nil {
		log.Errorf("Problem decoding JSON body %s", err)
		c.JSON(400, errors.New("problem decoding body"))
		return
	}

	log.Debugf("User Data print after c.Bind %s", user.PrintUser())
	u := &models.UserRow{}

	u, err = user.Signup(nil)

	if err != nil {

		log.Debugf("Erroring signing up ERROR %s", err)
		c.JSON(400, errors.New("Error Signing up user"))
		return
	}

	//zero out the password
	u.Password = ""
	u.PasswordAgain = ""

	c.JSON(http.StatusOK, u)
}
