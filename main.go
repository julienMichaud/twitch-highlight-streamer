// @title          twitch-highlight-streamer
// @version        1.0
// @description    This is an api to retrieve a random streamer on Twtich.
// @termsOfService http://swagger.io/terms/

// @contact.name  mymail
// @contact.url   http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url  http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080

package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	docs "github.com/julienMichaud/twitch-highlight-streamer/docs"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"golang.org/x/oauth2/clientcredentials"
	// swagger embed files
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
	docs.SwaggerInfo.BasePath = "/"

	router.GET("/streamer", getStreamers)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run("localhost:8080")

}
