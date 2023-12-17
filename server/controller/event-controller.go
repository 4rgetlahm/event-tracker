package controller

import (
	"net/http"
	"strconv"

	"github.com/4rgetlahm/event-tracker/server/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func BindEventRoutes(router *gin.Engine) {
	router.POST("/v1/event", func(c *gin.Context) {
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

	router.GET("/v1/event/:id", func(c *gin.Context) {
		uuidParam := c.Param("id")
		objectId, err := primitive.ObjectIDFromHex(uuidParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid event ID",
			})
			return
		}

		event, err := service.GetEvent(objectId)
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

	router.GET("/v1/events/:from/:to", func(c *gin.Context) {
		fromParam := c.Param("from")
		toParam := c.Param("to")

		from, err := strconv.Atoi(fromParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid from",
			})
			return
		}

		to, err := strconv.Atoi(toParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid to",
			})
			return
		}

		events, err := service.GetEvents(from, to)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"events": events,
		})
	})

	router.POST("/v1/event/:id/register", func(c *gin.Context) {
		uuidParam := c.Param("id")
		userEmail := "testemail@gmail.com"

		objectId, err := primitive.ObjectIDFromHex(uuidParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid event ID",
			})
			return
		}

		_, err = service.AddUserToEvent(objectId, userEmail)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.Status(http.StatusOK)

	})

	router.POST("/v1/event/:id/cancel", func(c *gin.Context) {
		uuidParam := c.Param("id")
		userEmail := "testemail@gmail.com"

		objectId, err := primitive.ObjectIDFromHex(uuidParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid event ID",
			})
			return
		}

		_, err = service.RemoveUserFromEvent(objectId, userEmail)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.Status(http.StatusOK)
	})
}
