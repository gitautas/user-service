package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
	"user-service/src/api"
	"user-service/src/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	// Using environment variables for simplicity, would use a config file instead.
	httpAddr := os.Getenv("HTTP_ADDRESS") // Would also add checks for values but omitted for brevity.

	rpcAddr := os.Getenv("RPC_ADDRESS")

	redisEndpoint := os.Getenv("REDIS_ENDPOINT")

	mongoEndpoint := os.Getenv("MONGO_ENDPOINT")
	mongoDBName := os.Getenv("MONGO_DATABASE")
	mongoCollectionName := os.Getenv("MONGO_COLLECTION")
	queryTimeoutSecondsStr := os.Getenv("MONGO_TIMEOUT")

	queryTimeoutSeconds, _ := strconv.Atoi(queryTimeoutSecondsStr)
	queryTimeout := time.Second * time.Duration(queryTimeoutSeconds)

	db := storage.NewMongo(mongoEndpoint, mongoDBName, mongoCollectionName, queryTimeout)
	err := db.Connect()
	if err != nil {
		panic(fmt.Errorf("rror initializing mongodb: %v", err))
	}

	pubSub := storage.NewRedis(redisEndpoint, queryTimeout) // Using same timeout for brevity
	pubSub.Connect()

	userService := api.NewUserService(db, pubSub)

	gin.SetMode(gin.DebugMode) // Would have a check if this was ever deployed

	httpRouter       := api.NewHttpRouter(userService)

	rpcListener, err := net.Listen("tcp", rpcAddr)
	rpcServer        := api.NewRpcServer(userService)

	go httpRouter.Engine.Run(httpAddr)
	fmt.Printf("[gRPC] Listening and serving on %v", httpAddr)
	go rpcServer.Server.Serve(rpcListener)
}
