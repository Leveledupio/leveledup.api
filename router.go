package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
)

type Application struct {
	config       *viper.Viper
	dsn          string
	db           *sqlx.DB
	sessionStore sessions.Store
}

func RequestIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Request-Id", uuid.NewV4().String())
		c.Next()
	}
}

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
	app.config = config
	app.dsn = dsn
	app.db = db
	app.sessionStore = sessions.NewCookieStore([]byte(cookieStoreSecret))

	return app, nil
}

func LoginEndpoint(c *gin.Context) {

	log.Debug("LoginEndpoint")

	c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})

}

func SignupEndpoint(c *gin.Context) {

	log.Debug("SignupEndpoint")

	c.JSON(http.StatusOK, gin.H{"status": "Sign up endpoint"})

}

func ProjectCreateEP(c *gin.Context) {
	log.Debug("ProjectCreateEP")

	c.JSON(http.StatusOK, gin.H{"status": "Create Project endpoint"})

}

func DeleteUserEP(c *gin.Context) {
	log.Debug("DeleteUserEP")

	c.JSON(http.StatusOK, gin.H{"status": "User Removed"})

}

func GetProjectEP(c *gin.Context) {
	log.Debug("GetProjectEP")

	c.JSON(http.StatusOK, gin.H{"status": "GetProjectEP"})

}

func DeleteProjectEP(c *gin.Context) {
	log.Debug("DeleteProjectEP")

	c.JSON(http.StatusOK, gin.H{"status": "DeleteProjectEP"})

}

func UpdateUserEP(c *gin.Context) {
	log.Debug("UpdateUserEP")

	c.JSON(http.StatusOK, gin.H{"status": "UpdateUserEP"})

}

func CreateTeamEP(c *gin.Context) {

	log.Debug("CreateTeamEP ID:%v", c.Keys["X-Request-Id"])

	c.JSON(http.StatusOK, gin.H{"status": "CreateTeamEP"})

}

func GetTeamEP(c *gin.Context) {
	log.Debug("GetTeamEP")

	c.JSON(http.StatusOK, gin.H{"status": "GetTeamEP"})

}

func RouteSetup() *gin.Engine {

	r := gin.Default()

	r.Use(RequestIdMiddleware())

	group := r.Group("/v1")

	{

		//User Actions
		group.POST("/login", LoginEndpoint)

		group.PUT("/user/:email", UpdateUserEP)
		group.DELETE("/user/:email", DeleteUserEP)

		group.POST("/signup", SignupEndpoint)

		//Project Actions
		group.POST("/project", ProjectCreateEP)
		group.GET("/project/:projectID", GetProjectEP)
		group.DELETE("/project/:projectID", DeleteProjectEP)

		//Team Actions
		group.POST("/team", CreateTeamEP)
		group.GET("/team", GetTeamEP)

		//Account Actions

		group.GET("/ping", func(c *gin.Context) {

			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
	}

	return r
}
