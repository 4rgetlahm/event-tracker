package main

import (
	"github.com/4rgetlahm/event-tracker/server/controller"
	"github.com/4rgetlahm/event-tracker/server/database"
	"github.com/4rgetlahm/event-tracker/server/entity"
	"github.com/gin-gonic/gin"
)

func Init() {
	database.GetDatabase()
	entity.Init()
}

func main() {
	Init()
	server := gin.Default()
	controller.BindEventRoutes(server)

	server.Run(":5000")
}
