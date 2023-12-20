package controller

import (
	"net/http"
	"time"

	"github.com/4rgetlahm/event-tracker/server/middleware"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/oauth2/v2"
)

type ReducedTokenInfo struct {
	Email      string    `json:"email"`
	Expiration time.Time `json:"expiration"`
}

func BindTokenController(router *gin.Engine) {
	router.GET("/v1/token", middleware.Authentication(), func(c *gin.Context) {
		c.AbortWithStatus(http.StatusOK)
	})

	router.GET("/v1/identity", middleware.Authentication(), func(c *gin.Context) {
		tokenInfo, exists := c.Get("tokenInfo")
		if !exists {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenInfoCast := tokenInfo.(*oauth2.Tokeninfo)

		reducedTokenInfo := ReducedTokenInfo{
			Email:      tokenInfoCast.Email,
			Expiration: time.Now().Add(time.Duration(tokenInfoCast.ExpiresIn) * time.Second),
		}

		c.JSON(http.StatusOK, reducedTokenInfo)
	})
}
