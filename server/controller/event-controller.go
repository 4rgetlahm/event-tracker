package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/4rgetlahm/event-tracker/server/middleware"
	"github.com/4rgetlahm/event-tracker/server/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/api/oauth2/v2"
)

func BindEventRoutes(router *gin.Engine) {
	router.POST("/v1/event", middleware.RequireAdministrator(), func(c *gin.Context) {
		var createRequest service.EventCreateRequest
		c.BindJSON(&createRequest)

		email, err := getEmail(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		event, err := service.CreateEvent(email, &createRequest)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, event)
	})

	router.GET("/v1/event/:id", middleware.RequireAdministrator(), func(c *gin.Context) {
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

		email, err := getEmail(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if !isAdministrator(c) {
			events, err := service.GetEventsForUser(email, from, to)

			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})

				return
			}

			c.JSON(http.StatusOK, gin.H{
				"events": events,
			})
			return
		}

		events, err := service.GetDetailedEvents(email, from, to)
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
		userEmail, err := getEmail(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

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

	router.POST("/v1/event/:id/cancel-registration", func(c *gin.Context) {
		uuidParam := c.Param("id")
		userEmail, err := getEmail(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

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

func isAdministrator(c *gin.Context) bool {
	isAdministrator, exists := c.Get("isAdministrator")
	if !exists || !isAdministrator.(bool) {
		return false
	}
	return true
}

func getEmail(c *gin.Context) (string, error) {
	tokenInfo, exists := c.Get("tokenInfo")
	if !exists {
		return "", errors.New("Token info not found")
	}

	tokenInfoCast := tokenInfo.(*oauth2.Tokeninfo)
	return tokenInfoCast.Email, nil
}
