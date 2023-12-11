package main

import (
	"github.com/4rgetlahm/event-tracker/server/controller"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	controller.BindEventRoutes(server)
	server.Run(":5000")
}
