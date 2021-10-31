//RIP in pieces you glorious bastard
package api

import (
	"log"
	"net/http"
	"strconv"
	"user-service/src/api/health"
	"user-service/src/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Not handling auth or resource checking as per requirements.

type HttpRouter struct {
	Engine *gin.Engine
	Address string
	healthChannel chan health.HealthCheckResponse_ServingStatus
	us *UserService
}

const PathPrefix = "/user"

func NewHttpRouter(userService *UserService, healthChan chan health.HealthCheckResponse_ServingStatus) *HttpRouter {
	engine := gin.Default()
	hr := &HttpRouter{
		Engine: engine,
		us: userService,
		healthChannel: healthChan,
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

	group.GET("/ping", func(c *gin.Context) {
		c.Status(200)
	}) // Health check handler

	return hr
}

func (hr *HttpRouter) Connect(addr string) {
	hr.Address = addr
	err := hr.Engine.Run(addr)
	if err != nil {
		hr.healthChannel <- health.HealthCheckResponse_NOT_SERVING
		panic("Could not start HTTP listener.")
	}
}

func (hr *HttpRouter) createUserHandler(c *gin.Context) {
	var user *models.User

	err := c.BindJSON(&user)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}


	user, err = hr.us.CreateUser(user)
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

	user, err = hr.us.UpdateUser(user)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, user)
	return
}

func (hr *HttpRouter) removeUserHandler(c *gin.Context) {
	userID := c.Param("userID")

	err := hr.us.RemoveUser(userID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent) // Semantic choice for this code because no entity is returned.
	return
}

func (hr *HttpRouter) getUserHandler(c *gin.Context) {
	userID := c.Param("userID")

	user, err := hr.us.GetUser(userID)
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

	users, err := hr.us.GetUserList(limit, skip, params)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, users)
	return
}
