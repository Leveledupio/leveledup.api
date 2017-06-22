package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
)

type Application struct {
	config       *viper.Viper
	dsn          string
	db           *sqlx.DB
	sessionStore sessions.Store
}

func NewApplication(config *viper.Viper) (*Application, error) {
	dsn := config.Get("dsn").(string)

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Errorf("Connecting to Database: %v", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Errorf("Database ping failed: %v", err)
		return nil, err
	}

	log.Info("Connected to Database")

	cookieStoreSecret := config.Get("cookie_secret").(string)

	app := &Application{}
	app.config = config
	app.dsn = dsn
	app.db = db
	app.sessionStore = sessions.NewCookieStore([]byte(cookieStoreSecret))

	return app, nil
}

func LoginEndpoint(c *gin.Context) {

	log.Debug("LoginEndpoint")

	c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})

}
