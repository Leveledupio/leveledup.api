package handlers

import (
	"gopkg.in/op/go-logging.v1"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	model "github.com/strongjz/leveledup-api/model"
	"errors"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
)

var (
	log = logging.MustGetLogger("gbs")
)

type ApiResource struct {
	DB *sqlx.DB
}


func (h *ApiResource) LoginEndpoint (c *gin.Context)  {


	log.Debugf("Handler Login")

	user := model.User{}

	if c.Bind(&user) != nil {
		c.JSON(400, errors.New("problem decoding body"))
		return
	}


	u, err := user.GetUserByEmailAndPassword(h.DB, user.Email, user.Password)

	if err != nil{

		c.JSON(401, errors.New("Username and Password did not match"))
		return
	}



	c.JSON(http.StatusOK, user)
}