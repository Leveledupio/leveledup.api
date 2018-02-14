package handlers

import (
	"os"

	"github.com/andygrunwald/go-jira"
	"github.com/aws/aws-sdk-go/aws/session"
	//Mysql setup requres a blank import
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"gopkg.in/op/go-logging.v1"
)

var (
	log           = logging.MustGetLogger("handlers")
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
	log.Infof("Logging Backend initailzated ")
}

//APIResource  Struct for the API
type APIResource struct {
	DB         *sqlx.DB
	AWSSession *session.Session
	Jira       *jira.Client
}
