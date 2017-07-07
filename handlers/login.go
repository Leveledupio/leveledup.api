package handlers

import (
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	model "github.com/strongjz/leveledup-api/model"
	"github.com/gin-gonic/gin"
	"gopkg.in/op/go-logging.v1"
	"net/http"
	//"io/ioutil"
	"os"
)

var (
	log = logging.MustGetLogger("handlers")
	defaultFormat = "%{color}%{time:2006-01-02T15:04:05.000Z07:00} %{level:-5s} [%{shortfile}]%{color:reset} %{message}"

)

func init() {
	logBackend := logging.NewLogBackend(os.Stderr, "", 0)
	logging.SetBackend(
		logging.NewBackendFormatter(
			logBackend,
			logging.MustStringFormatter(defaultFormat),
		),
	)

	logging.SetLevel(logging.DEBUG, "handlers")
}

type ApiResource struct {
	DB *sqlx.DB

}

func (h *ApiResource) UserLogin(c *gin.Context) {

	//body, _ := ioutil.ReadAll(c.Request.Body)

	//log.Debugf("Request body: %v", string(body))

	log.Debugf("Handler Login")

	user := model.NewUser(h.DB)

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
