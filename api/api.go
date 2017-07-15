package api

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"github.com/strongjz/leveledup-api/application"
	"github.com/strongjz/leveledup-api/handlers"
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
func RouteSetup(a *application.Application) *gin.Engine {

	r := gin.Default()
	r.Use(RequestIDMiddleware())
	r.Use(ApiMiddleware(a.DB))

	cors_config := cors.DefaultConfig()
	cors_config.AllowAllOrigins = true
	r.Use(cors.New(cors_config))

	api := &handlers.ApiResource{DB: a.DB}

	//User Actions
	r.POST("/user/login", api.UserLogin)
	r.PUT("/user/:email", api.UserUpdate)
	r.GET("/user/:email", api.UserRetrieve)
	r.DELETE("/user", api.UserDelete)
	r.POST("/user/signup", api.UserSignup)

	//Team Actions
	r.POST("/team", api.TeamCreate)
	r.GET("/team/:teamID", api.TeamGet)
	r.PUT("/team/:teamID", api.TeamUpdate)
	r.DELETE("/team/:teamID", api.TeamDelete)


	//Project Actions
	r.POST("/project", api.ProjectCreate)
	r.GET("/project/:projectID", api.ProjectGet)
	r.PUT("/project/:projectID", api.ProjectUpdate)
	r.DELETE("/project/:projectID", api.ProjectDelete)


	r.GET("/health", func(c *gin.Context) {

		err := a.DB.Ping()
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

	pwd, _ := os.Getwd()
	log.Debugf("Current working dir %s", pwd)
	log.Debugf("Loading Config file %v", configName)

	c := viper.New()
	c.SetConfigFile(configName)
	c.SetConfigType("yaml")
	c.WatchConfig()
	c.Debug()

	c.OnConfigChange(func(e fsnotify.Event) {
		log.Infof("Config file changed:", e.Name)
	})

	err := c.ReadInConfig() // Find and read the config file
	if err != nil {
		// Handle errors reading the config file
		log.Panicf("fatal error config file: %s", err)

	}

	c.Debug()

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

	log.Debugf("Config: %v", config)
	log.Debugf("App DSN: %s", app.DSN)

	r := RouteSetup(app)

	port := config.Get("port").(string)
	if port == "" {
		log.Fatal("No Port in config file")
	}

	address := fmt.Sprintf("0.0.0.0:%s", port)

	r.Run(address) // listen and serve on 0.0.0.0:8080
}
