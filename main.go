package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	//"github.com/strongjz/leveledup-api/model"
	"github.com/spf13/viper"
	"os"
	"github.com/fsnotify/fsnotify"
	"fmt"
	"gopkg.in/op/go-logging.v1"
	//"github.com/jmoiron/sqlx"
)


var (
	log           = logging.MustGetLogger("gbs")
	defaultFormat = "%{color}%{time:2006-01-02T15:04:05.000Z07:00} %{level:-5s} [%{shortfile}]%{color:reset} %{message}"
	 ENV = "dev"

)

func init(){
	log.Debugf("Init: Default ENV: %v", ENV)

	if os.Getenv("ENV") != ""{
		ENV = os.Getenv("ENV")
	}

	log.Debugf("Init: ENV Set: %v", ENV)

}

func newConfig() (*viper.Viper, error) {


	configName := fmt.Sprintf("config/%v-config.yaml", ENV)

	log.Debugf("Loading Config file %v", configName)

	c := viper.New()
	c.AddConfigPath(".")
	c.SetConfigFile(configName)
	c.SetConfigType("yaml")
	//c.Debug()
	c.WatchConfig()

	c.OnConfigChange(func(e fsnotify.Event) {
		log.Infof("Config file changed:", e.Name)
	})

	c.AutomaticEnv()

	err := c.ReadInConfig() // Find and read the config file
	if err != nil {
		// Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	c.Debug()

	return c, nil
}


func main() {

	logBackend := logging.NewLogBackend(os.Stderr, "", 0)
	logging.SetBackend(
		logging.NewBackendFormatter(
			logBackend,
			logging.MustStringFormatter(defaultFormat),
		),
	)

	config, err := newConfig()
	if err != nil {
		log.Fatal(err)
	}

	log.Debugf("Config: %v", config)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {

		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}