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

//Application - Struct the contains the Config, DSN for database, and the db structure
//
type Application struct {
	config       *viper.Viper
	dsn          string
	db           *sqlx.DB
	sessionStore sessions.Store
}

//RequestIDMiddleware - Sets the Request ID header for request tracking
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Request-Id", uuid.NewV4().String())
		c.Next()
	}
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
	app.config = config
	app.dsn = dsn
	app.db = db
	app.sessionStore = sessions.NewCookieStore([]byte(cookieStoreSecret))

	return app, nil
}

//LoginEndpoint - Logs in a user
//
func LoginEndpoint(c *gin.Context) {

	log.Debug("LoginEndpoint")

	c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})

}

//SignupEndpoint - Creates a new user for leveledup
//
func SignupEndpoint(c *gin.Context) {

	log.Debug("SignupEndpoint")

	c.JSON(http.StatusOK, gin.H{"status": "Sign up endpoint"})

}

//GetUserEP - Returns user specified by email
//
func GetUserEP(c *gin.Context) {

	log.Debug("SignupEndpoint")

	c.JSON(http.StatusOK, gin.H{"status": "Sign up endpoint"})
}

//ProjectCreateEP - Creates a project with specificied attributes
//
func ProjectCreateEP(c *gin.Context) {
	log.Debug("ProjectCreateEP")

	c.JSON(http.StatusOK, gin.H{"status": "Create Project endpoint"})

}

//DeleteUserEP - Deletes the specified User by email
//
func DeleteUserEP(c *gin.Context) {
	log.Debug("DeleteUserEP")

	c.JSON(http.StatusOK, gin.H{"status": "User Removed"})

}

// GetProjectEP - Retrieves the project by the specified ID
//
func GetProjectEP(c *gin.Context) {
	log.Debug("GetProjectEP")

	c.JSON(http.StatusOK, gin.H{"status": "GetProjectEP"})

}

//DeleteProjectEP - Deletes the specified project
//
func DeleteProjectEP(c *gin.Context) {
	log.Debug("DeleteProjectEP")

	c.JSON(http.StatusOK, gin.H{"status": "DeleteProjectEP"})

}

// UpdateUserEP - Updates User's data
//
func UpdateUserEP(c *gin.Context) {
	log.Debug("UpdateUserEP")

	c.JSON(http.StatusOK, gin.H{"status": "UpdateUserEP"})

}

// CreateTeamEP - Creates new team with data provider
//
func CreateTeamEP(c *gin.Context) {

	log.Debug("CreateTeamEP ID:%v", c.Keys["X-Request-Id"])

	c.JSON(http.StatusOK, gin.H{"status": "CreateTeamEP"})

}

// GetTeamEP - Gets the specified team from the ID data
//
func GetTeamEP(c *gin.Context) {
	log.Debug("GetTeamEP")

	c.JSON(http.StatusOK, gin.H{"status": "GetTeamEP"})

}

// RouteSetup - Set the http routes for the api
//
func RouteSetup() *gin.Engine {

	r := gin.Default()

	r.Use(RequestIDMiddleware())

	group := r.Group("/v1")

	{

		//User Actions
		group.POST("/login", LoginEndpoint)

		group.PUT("/user/:email", UpdateUserEP)
		group.GET("/user/:email", GetUserEP)
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
