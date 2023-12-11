package controller

import (
	"net/http"

	"github.com/4rgetlahm/event-tracker/server/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func BindEventRoutes(router *gin.Engine) {
	router.GET("/v1/events/:uuid", func(c *gin.Context) {
		uuidParam := c.Param("uuid")
		parsedUUID, err := uuid.Parse(uuidParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid UUID",
			})
			return
		}

		event, err := service.GetEvent(parsedUUID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"event": event,
		})

	})

	router.POST("/v1/events", func(c *gin.Context) {
		var createRequest service.EventCreateRequest
		c.BindJSON(&createRequest)
		event, err := service.CreateEvent(&createRequest)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, event)
	})
}
