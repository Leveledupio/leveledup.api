package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/jmoiron/sqlx"
	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"github.com/strongjz/leveledup-api/application"
	app "github.com/strongjz/leveledup-api/application"
	"github.com/strongjz/leveledup-api/handlers"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/op/go-logging.v1"

	"os"
)

var (
	log = logging.MustGetLogger("gbs")

	defaultFormat = "%{color}%{time:2006-01-02T15:04:05.000Z07:00} %{level:-5s} [%{shortfile}]%{color:reset} %{message}"
	//ENV sets the environment so the right config is loaded
	ENV = "dev"
)

//Init function sets the default env to dev
func init() {
	log.Debugf("Init: Default ENV: %v", ENV)

	if os.Getenv("ENV") != "" {
		ENV = os.Getenv("ENV")
	}

	log.Debugf("Init: ENV Set: %v", ENV)

}

//RequestIDMiddleware - Sets the Request ID header for request tracking
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Request-Id", uuid.NewV4().String())
		c.Next()
	}
}

func ApiMiddleware(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("databaseConn", db)
		c.Next()
	}
}

// RouteSetup - Set the http routes for the api
//
func RouteSetup(a *application.Application) *gin.Engine {

	r := gin.Default()
	r.Use(RequestIDMiddleware())
	r.Use(ApiMiddleware(a.DB))

	api := &handlers.ApiResource{DB: a.DB}

	//User Actions
	r.POST("/user/login", api.UserLogin)
	r.PUT("/user/:email", api.UserUpdate)
	r.GET("/user/:email", api.UserRetrieve)
	r.DELETE("/user/:email", api.UserDelete)
	r.POST("/user/signup", api.UserSignup)

	/*
		//Project Actions
		r.POST("/project", ProjectCreateEP)
		r.GET("/project/:projectID", GetProjectEP)
		r.DELETE("/project/:projectID", DeleteProjectEP)

		//Team Actions
		r.POST("/team", CreateTeamEP)
		r.GET("/team/:teamID", GetTeamEP)

		//Account Actions

	*/
	r.GET("/ping", func(c *gin.Context) {

		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	return r
}

//newConfig function parsers the config file for database connection etc
//
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
		panic(fmt.Errorf("fatal error config file: %s", err))
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

	a, err := app.NewApplication(config)
	if err != nil {
		log.Fatal(err)
	}

	log.Debugf("Config: %v", config)
	log.Debug("App DSN: %s", a.DSN)

	r := RouteSetup(a)

	port := config.Get("port").(string)
	if port == "" {
		log.Fatal("No DSN in config file")
	}

	address := fmt.Sprintf("0.0.0.0:%s", port)
	r.Run(address) // listen and serve on 0.0.0.0:8080
}
