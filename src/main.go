package main

import (
	"fmt"
	"net"
	"os"
	"user-service/src/api"
	"user-service/src/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	// Using environment variables for simplicity, would use a config file instead.
	httpAddr := os.Getenv("HTTP_ADDRESS") // Would also add checks for values but omitted for brevity.

	rpcAddr := os.Getenv("RPC_ADDRESS")

	sqlUsername := os.Getenv("SQL_USERNAME")
	sqlPassword := os.Getenv("SQL_PASSWORD")
	sqlEndpoint := os.Getenv("SQL_ENDPOINT")
	sqlDBName := os.Getenv("SQL_DATABASE")

	mysql := storage.Mysql{
		UserName:     sqlUsername,
		Password:     sqlPassword,
		Endpoint:     sqlEndpoint,
		DatabaseName: sqlDBName,
	}

	err := mysql.Connect()
	if err != nil {
		panic(fmt.Errorf("error initializing mysql: %v", err))
	}

	gin.SetMode(gin.DebugMode) // Would have a check if this was ever deployed

	httpRouter       := api.NewHttpRouter(&mysql)

	rpcListener, err := net.Listen("tcp", rpcAddr)
	rpcServer        := api.NewRpcServer(&mysql)

	go httpRouter.Engine.Run(httpAddr)
	fmt.Printf("[gRPC] Listening and serving on %v", httpAddr) //FIXME
	rpcServer.Server.Serve(rpcListener)
}
