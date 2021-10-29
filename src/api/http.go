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
	"github.com/google/uuid"
)

// Not handling auth or resource checking as per requirements.

type HttpRouter struct {
	Engine *gin.Engine
	mysql storage.Database
}

const PathPrefix = "/user"

func NewHttpRouter(mysql storage.Database) *HttpRouter {
	engine := gin.Default()
	hr := &HttpRouter{
		Engine: engine,
		mysql:  mysql,
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


	user, err = CreateUser(user, hr.mysql)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, user)
	return
}

func (hr *HttpRouter) updateUserHandler(c *gin.Context) {
	var user *models.User
	userID, err := uuid.Parse(c.Param("userID")) // Extra check for valid UUID.

	err = c.BindJSON(&user)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user.ID = userID.String()

	user, err = UpdateUser(user, hr.mysql)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, user)
	return
}

func (hr *HttpRouter) removeUserHandler(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userID"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = RemoveUser(userID.String(), hr.mysql)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent) // Semantic choice for this code because no entity is returned.
	return
}

func (hr *HttpRouter) getUserHandler(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userID"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user, err := GetUser(userID.String(), hr.mysql)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, user)
	return
}

func (hr *HttpRouter) getUsersHandler(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	users, err := GetUsers(limit, offset, hr.mysql)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, users)
	return
}
