package handlers

import (
	"gopkg.in/op/go-logging.v1"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	log = logging.MustGetLogger("gbs")
)

func Login(db *sqlx.DB) error {

	log.Debugf("Handler Login %v", db)

	return nil
}