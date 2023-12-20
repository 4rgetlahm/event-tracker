package main

import (
	"github.com/4rgetlahm/event-tracker/server/controller"
	"github.com/4rgetlahm/event-tracker/server/database"
	"github.com/4rgetlahm/event-tracker/server/middleware"
	"github.com/gin-gonic/gin"
)

func Init() {
	database.GetDatabase()
}

func main() {
	Init()
	middleware.InitAuthenticationMiddleware()
	middleware.InitAdministratorMiddleware()

	server := gin.Default()

	server.Use(gin.Logger())
	server.Use(gin.Recovery())
	server.Use(middleware.CORS())
	controller.BindTokenController(server)

	server.Use(middleware.Authentication(), middleware.Administrator())
	controller.BindEventRoutes(server)

	server.Run(":5000")
}
