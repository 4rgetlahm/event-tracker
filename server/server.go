package main

import (
	"github.com/4rgetlahm/event-tracker/server/controller"
	"github.com/4rgetlahm/event-tracker/server/database"
	"github.com/gin-gonic/gin"
)

func Init() {
	database.GetDatabase()
}

func main() {
	Init()
	server := gin.Default()
	controller.BindEventRoutes(server)

	server.Run(":5000")
}
