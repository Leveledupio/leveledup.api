package api

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"github.com/strongjz/leveledup.api/application"
	"github.com/strongjz/leveledup.api/handlers"
	"gopkg.in/op/go-logging.v1"
	"os"
)

type Api struct {
	Router *gin.Engine
}

var (
	log = logging.MustGetLogger("api")

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
	logging.SetLevel(logging.INFO, "main")

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
func RouteSetup(app *application.Application) *gin.Engine {

	r := gin.Default()
	r.Use(RequestIDMiddleware())
	r.Use(ApiMiddleware(app.DB))

	cors_config := cors.DefaultConfig()
	cors_config.AllowAllOrigins = true
	r.Use(cors.New(cors_config))

	api := &handlers.ApiResource{DB: app.DB, AWSSession: app.AWSSession, Jira: app.JiraClient}

	//User Actions
	r.POST("/user/login", api.UserLogin)
	r.PUT("/user/:email", api.UserUpdate)
	r.GET("/user/:email", api.UserRetrieve)
	r.DELETE("/user", api.UserDelete)
	r.POST("/user/signup", api.UserSignup)

	//Team Actions
	r.POST("/team", api.TeamCreate)
	r.GET("/team/:team", api.TeamGet)
	r.PUT("/team/:team", api.TeamUpdate)
	r.DELETE("/team/:team", api.TeamDelete)

	//Project Actions
	r.POST("/project", api.ProjectCreate)
	r.GET("/project/:projectID", api.ProjectGet)
	r.GET("/projects", api.ProjectGetAll)
	r.PUT("/project/:projectID", api.ProjectUpdate)
	r.DELETE("/project/:projectID", api.ProjectDelete)

	//Contact Us
	r.POST("contact", api.Contact)

	//Loadbalancer Health check
	r.GET("/health", func(c *gin.Context) {

		err := app.DB.Ping()
		if err != nil {
			e := fmt.Sprintf("Database health ping failed: %v", err)
			log.Error(e)
			c.JSON(500, gin.H{
				"error": e,
			})
			c.Abort()
			return
		}

		log.Debug("Connected to Database")

		c.JSON(200, gin.H{
			"message": "Connected!",
		})
	})

	return r
}

//newConfig function parsers the config file for database connection etc
//
func newConfig() (*viper.Viper, error) {

	configName := fmt.Sprint("config/config.yaml")

	c := viper.New()
	c.SetConfigFile(configName)
	c.SetConfigType("yaml")
	c.WatchConfig()

	c.OnConfigChange(func(e fsnotify.Event) {
		log.Infof("Config file changed:", e.Name)
	})

	err := c.ReadInConfig() // Find and read the config file
	if err != nil {
		// Handle errors reading the config file
		log.Panicf("fatal error config file: %s", err)

	}

	return c, nil
}

func (a *Api) LevelUp() {

	logBackend := logging.NewLogBackend(os.Stderr, "", 0)
	logging.SetBackend(
		logging.NewBackendFormatter(
			logBackend,
			logging.MustStringFormatter(defaultFormat),
		),
	)

	logging.SetLevel(logging.DEBUG, "")

	config, err := newConfig()
	if err != nil {
		log.Fatal(err)
	}

	app, err := application.NewApplication(config)
	if err != nil {
		log.Fatal(err)
	}

	r := RouteSetup(app)

	port := config.Get("port").(string)
	if port == "" {
		log.Warning("No Port in config file")
		port = "8080"
	}

	address := fmt.Sprintf("0.0.0.0:%s", port)

	r.Run(address) // listen and serve on 0.0.0.0:8080
}
