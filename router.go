package main

import (

	"github.com/satori/go.uuid"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"github.com/strongjz/leveledup-api/handlers"
)


//RequestIDMiddleware - Sets the Request ID header for request tracking
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Request-Id", uuid.NewV4().String())
		c.Next()
	}
}



//LoginEndpoint - Logs in a user
//
func LoginEndpoint(c *gin.Context) {

	log.Debug("LoginEndpoint")

	login, err := handlers.Login(app)

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
