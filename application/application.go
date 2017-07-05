package application

import (
	"github.com/spf13/viper"
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"gopkg.in/op/go-logging.v1"

	"errors"
)

var (
	log = logging.MustGetLogger("gbs")


)


//Application - Struct the contains the Config, DSN for database, and the db structure
//
type Application struct {
	Config       *viper.Viper
	DSN          string
	DB 	*sqlx.DB
	SessionStore sessions.Store
}


//NewApplication - Creates a new applications and populates information based on the config.
//
func NewApplication(config *viper.Viper) (*Application, error) {
	dsn := config.Get("dsn").(string)
	if dsn == "" {
		log.Errorf("No DSN in config file")
		return nil, errors.New("No DSN in config file")
	}

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
	app.Config = config
	app.DSN = dsn
	app.DB = db
	app.SessionStore = sessions.NewCookieStore([]byte(cookieStoreSecret))

	return app, nil
}