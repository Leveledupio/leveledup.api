package handlers

import (
	_ "github.com/go-sql-driver/mysql"

	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/strongjz/leveledup.api/models"
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

		c.JSON(400, errors.New("email not found"))
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
		log.Errorf("Problem decoding JSON body %s", err)
		c.JSON(400, errors.New("problem reading request body, Please try request again later"))
		return
	}

	u, err := user.GetByEmail(nil, email)
	if err != nil {
		log.Errorf("Get by email error %s", err)

		if err.Error() == "sql: no rows in result set" {
			c.JSON(http.StatusOK, gin.H{"status": "Email not found"})
			return

		}

		c.JSON(http.StatusBadRequest, gin.H{"status": "Please try request again later"})
		return
	}

	//zero out the password
	u.Password = ""
	u.PasswordAgain = ""

	c.JSON(http.StatusOK, u)
}

func (h *ApiResource) UserDelete(c *gin.Context) {
	log.Debugf("Handler UserDelete")

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
			c.JSON(404, errors.New("user does not exist"))
			return
		}

		if err.Error() == "crypto/bcrypt: hashedPassword is not the hash of the given password" {

			c.JSON(http.StatusBadRequest, gin.H{"error": "Username and Password did not match"})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": "Error deleting user please try again later"})
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

		log.Debugf("error signing up ERROR %s", err)
		c.JSON(400, errors.New("error Signing up user"))
		return
	}

	//zero out the password
	u.Password = ""
	u.PasswordAgain = ""

	c.JSON(http.StatusOK, u)
}
