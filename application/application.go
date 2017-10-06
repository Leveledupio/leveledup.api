package application

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"gopkg.in/op/go-logging.v1"

	"errors"
	jira "github.com/andygrunwald/go-jira"
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
	JiraClient   *jira.Client
}

//NewApplication - Creates a new applications and populates information based on the config.
//
func NewApplication(config *viper.Viper) (*Application, error) {
	dsn := config.Get("dsn").(string)

	if dsn == "" {
		log.Errorf("No DSN in config file")
		return nil, errors.New("no DSN in config file")
	}
	log.Info("[INFO] Connecting to Database")

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Errorf("[ERROR] connecting to Database: %v", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Errorf("[ERROR] Database ping failed: %v", err)
		return nil, err
	}

	log.Info("[INFO] connected to Database")

	cookieStoreSecret := config.Get("cookie_secret").(string)
	region := config.Get("aws_region").(string)

	jiraURL := config.Get("jira_url").(string)
	jiraUSER := config.Get("jira_user").(string)
	jiraPassword := config.Get("jira_password").(string)

	log.Info("[INFO] Connecting to JIRA ")

	jiraClient, err := jira.NewClient(nil, jiraURL)
	if err != nil {
		log.Errorf("[ERROR] Creating Jira client %s ", err)
		return nil, err

	}

	log.Infof("Jira Base URL %v", jiraClient.GetBaseURL())

	jiraClient.Authentication.SetBasicAuth(jiraUSER, jiraPassword)

	if jiraClient.Authentication.Authenticated() {
		log.Debug("[DEBUG] Jira Authenticated")
	}

	app := &Application{}
	app.Config = config
	app.DSN = dsn
	app.DB = db
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))

	app.AWSSession = sess

	app.SessionStore = sessions.NewCookieStore([]byte(cookieStoreSecret))
	app.JiraClient = jiraClient

	return app, nil
}
