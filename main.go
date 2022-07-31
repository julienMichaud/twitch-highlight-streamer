package main

import (
	"os"

	"golang.org/x/oauth2/clientcredentials"
)

var (
	clientID = os.Getenv("CLIENT_ID")
	// Consider storing the secret in an environment variable or a dedicated storage system.
	clientSecret = os.Getenv("CLIENT_SECRET")
	oauth2Config *clientcredentials.Config
)

func main() {

	getStreamers()

}
