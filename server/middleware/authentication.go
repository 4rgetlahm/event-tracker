package middleware

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

var httpClient = &http.Client{}

var creds *google.Credentials = nil

func InitAuthenticationMiddleware() {
	base64Creds := os.Getenv("EVENT_TRACKER_GOOGLE_APPLICATION_CREDENTIALS")

	authJson, err := base64.RawStdEncoding.DecodeString(base64Creds)
	if err != nil {
		panic(err)
	}

	creds, err = google.CredentialsFromJSON(context.Background(), authJson, "https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile")
	if err != nil {
		panic(err)
	}
}

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.Request.Header.Get("Authorization")
		if !strings.HasPrefix(authorizationHeader, "Bearer ") {
			c.AbortWithStatus(401)
			return
		}
		tokenString := authorizationHeader[7:]
		tokenInfo, err := verifyIdToken(tokenString)
		if err != nil {
			c.AbortWithStatus(401)
			return
		}
		c.Set("tokenInfo", tokenInfo)
	}
}

func verifyIdToken(accessToken string) (*oauth2.Tokeninfo, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	oauth2Service, err := oauth2.NewService(ctx, option.WithCredentials(creds))
	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfoCall.AccessToken(accessToken)
	tokenInfo, err := tokenInfoCall.Do()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return tokenInfo, nil
}
