package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/oauth2/twitch"
)

func getToken() (string, error) {
	oauth2Config = &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     twitch.Endpoint.TokenURL,
	}

	token, err := oauth2Config.Token(context.Background())
	if err != nil {
		return "", err
	}

	return token.AccessToken, nil
}

func getStreamers(c *gin.Context) {

	token, err := getToken()
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": err})
		return
	}

	var data TwitchResponse
	mylist := StreamerList{}
	lang := c.Query("lang")
	minViewers := c.Query("minviewers")
	maxViewers := c.Query("maxviewers")

	if minViewers == "" {
		minViewers = "1"
	}

	if maxViewers == "" {
		maxViewers = "10"
	}

	intMinViewers, err := strconv.Atoi(minViewers)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"could not convert minviewers param to int": err})
	}

	intMaxViewers, err := strconv.Atoi(maxViewers)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"could not convert maxviewers param to int": err})
	}

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err})
	}

	previousPagination := ""
	pagination := true

	// val, err := rdb.Get(lang).Result() TO DO
	// if err != nil {
	// 	c.JSON(http.StatusBadGateway, gin.H{"error": err})
	// }

	for pagination {
		requestStreamer(&data, data.Pagination.Cursor, lang, token)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err})
			return
		}
		nextPagination := data.Pagination.Cursor

		if previousPagination == nextPagination {
			pagination = false
		} else {
			previousPagination = nextPagination
			nextPagination = ""
		}
		for _, v := range data.Data {
			mylist.AddStreamer(Streamer{UserName: v.UserName, ViewerCount: v.ViewerCount, GameName: v.GameName})
			log.Printf(v.UserName)
		}

	}

	e, err := json.Marshal(mylist)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": err})
		return
	}

	rdb.Set("fr", e, 10*time.Minute).Err()
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": err})
		return
	}

	streamer, err := mylist.GetStreamer(intMinViewers, intMaxViewers)

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"streamer": streamer.UserName,
		"viewers":  streamer.ViewerCount})

}
