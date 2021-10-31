//RIP in pieces you glorious bastard
package api

import (
	"fmt"
	"net/http"
	"strconv"
	"user-service/src/models"
	"user-service/src/storage"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Not handling auth or resource checking as per requirements.

type HttpRouter struct {
	Engine *gin.Engine
	db storage.Database
}

const PathPrefix = "/user"

func NewHttpRouter(db storage.Database) *HttpRouter {
	engine := gin.Default()
	hr := &HttpRouter{
		Engine: engine,
		db:  db,
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"*"}
	engine.Use(cors.New(corsConfig))

	group := engine.Group(PathPrefix)
	group.Use(cors.New(corsConfig))

	group.POST("/create", hr.createUserHandler)
	group.PUT("/:userID", hr.updateUserHandler)
	group.DELETE("/:userID", hr.removeUserHandler)
	group.GET("/:userID", hr.getUserHandler)
	group.GET("/", hr.getUsersHandler)

	return hr
}

func (hr *HttpRouter) createUserHandler(c *gin.Context) {
	var user *models.User

	err := c.BindJSON(&user)
	if err != nil {
		fmt.Println(err) // FIXME
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}


	user, err = CreateUser(user, hr.db)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, user)
	return
}

func (hr *HttpRouter) updateUserHandler(c *gin.Context) {
	var user *models.User
	userID := c.Param("userID")

	err := c.BindJSON(&user)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user.Id = userID

	user, err = UpdateUser(user, hr.db)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, user)
	return
}

func (hr *HttpRouter) removeUserHandler(c *gin.Context) {
	userID := c.Param("userID")

	err := RemoveUser(userID, hr.db)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent) // Semantic choice for this code because no entity is returned.
	return
}

func (hr *HttpRouter) getUserHandler(c *gin.Context) {
	userID := c.Param("userID")

	user, err := GetUser(userID, hr.db)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, user)
	return
}

func (hr *HttpRouter) getUsersHandler(c *gin.Context) {
	params := make(map[string]string)
	for key, value := range c.Request.URL.Query() {
		// This is an ugly solution because it means that only the first value of a filter
		// Might come up with a better solution later down the line.
		params[key] = value[0]
	}
	// This is inefficient since I already have these values inside
	// the params variabe, but this is a much more graceful way to handle defaults.
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	skip, err := strconv.Atoi(c.DefaultQuery("skip", "0"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	delete(params, "limit")
	delete(params, "skip")

	users, err := GetUsers(limit, skip, params, hr.db)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, users)
	return
}
