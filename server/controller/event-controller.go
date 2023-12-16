package controller

import (
	"net/http"
	"strconv"

	"github.com/4rgetlahm/event-tracker/server/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	router.GET("/v1/event/:uuid", func(c *gin.Context) {
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

	router.POST("/v1/event/:uuid/register", func(c *gin.Context) {
		uuidParam := c.Param("uuid")
		userEmail := "testemail@gmail.com"

		parsedUUID, err := uuid.Parse(uuidParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid UUID",
			})
			return
		}

		event, err := service.AddUserToEvent(parsedUUID, userEmail)
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

	router.POST("/v1/event/:uuid/cancel", func(c *gin.Context) {
		uuidParam := c.Param("uuid")
		userEmail := "testemail@gmail.com"

		parsedUUID, err := uuid.Parse(uuidParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid UUID",
			})
			return
		}

		event, err := service.RemoveUserFromEvent(parsedUUID, userEmail)
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
}
