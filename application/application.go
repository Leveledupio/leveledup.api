package application

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"gopkg.in/op/go-logging.v1"

	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

var (
	log = logging.MustGetLogger("gbs")
)

//Application - Struct the contains the Config, DSN for database, and the db structure
//
type Application struct {
	Config       *viper.Viper
	DSN          string
	DB           *sqlx.DB
	AWSSession   *session.Session
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
	region := config.Get("aws_region").(string)

	app := &Application{}
	app.Config = config
	app.DSN = dsn
	app.DB = db
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))

	//log.Debugf("Application Session %v", sess)
	app.AWSSession = sess

	app.SessionStore = sessions.NewCookieStore([]byte(cookieStoreSecret))

	return app, nil
}
