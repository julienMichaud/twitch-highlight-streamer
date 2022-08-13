package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2/clientcredentials"
)

var (
	clientID = os.Getenv("CLIENT_ID")
	// Consider storing the secret in an environment variable or a dedicated storage system.
	clientSecret = os.Getenv("CLIENT_SECRET")
	oauth2Config *clientcredentials.Config
)

func main() {
	router := gin.Default()
	router.GET("/streamer", getStreamers)

	router.Run("localhost:8080")

}
