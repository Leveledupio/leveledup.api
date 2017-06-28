package application

import (
	"github.com/spf13/viper"
	"github.com/jmoiron/sqlx"
	"github.com/gorilla/sessions"
	"gopkg.in/op/go-logging.v1"
)

var (
	log = logging.MustGetLogger("gbs")


)


//Application - Struct the contains the Config, DSN for database, and the db structure
//
type Application struct {
	Config       *viper.Viper
	Dsn          string
	DB 	*sqlx.DB
	SessionStore sessions.Store
}


//NewApplication - Creates a new applications and populates information based on the config.
//
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
	app.Config = config
	app.Dsn = dsn
	app.DB = db
	app.SessionStore = sessions.NewCookieStore([]byte(cookieStoreSecret))

	return app, nil
}