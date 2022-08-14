package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"golang.org/x/oauth2/clientcredentials"
)

var (
	clientID = os.Getenv("CLIENT_ID")

	// Consider storing the secret in an environment variable or a dedicated storage system.
	clientSecret = os.Getenv("CLIENT_SECRET")
	oauth2Config *clientcredentials.Config

	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:7001",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
)

func main() {

	if clientID == "" {
		log.Fatalf("clientID not set")
		os.Exit(3)
	}

	if clientSecret == "" {
		log.Fatalf("clientSecret not set")
		os.Exit(3)
	}

	router := gin.Default()

	router.GET("/streamer", getStreamers)

	router.Run("localhost:8080")

}
