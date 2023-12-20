package middleware

import (
	"encoding/base64"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/oauth2/v2"
)

var administratorEmails = []string{}

func InitAdministratorMiddleware() {
	hashedAdministrators := os.Getenv("EVENT_TRACKER_ADMINISTRATORS")
	administratorEmailsString, err := base64.StdEncoding.DecodeString(hashedAdministrators)
	if err != nil {
		panic(err)
	}

	administratorEmails = strings.Split(string(administratorEmailsString), ";")
}

func Administrator() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenInfo, exists := ctx.Get("tokenInfo")
		if !exists {
			ctx.AbortWithStatus(401)
			return
		}

		tokenInfoCast := tokenInfo.(*oauth2.Tokeninfo)
		ctx.Set("isAdministrator", isAdministrator(tokenInfoCast.Email))
	}
}

func RequireAdministrator() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		isAdministrator, exists := ctx.Get("isAdministrator")
		if !exists || !isAdministrator.(bool) {
			ctx.AbortWithStatus(403)
			return
		}
	}
}

func isAdministrator(email string) bool {
	for _, administratorEmail := range administratorEmails {
		if email == administratorEmail {
			return true
		}
	}
	return false
}