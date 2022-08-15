package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

// @BasePath /

// getStreamers godoc
// @Summary     /streamer?lang=fr&minviewers=1&maxviewers=10
// @Schemes     http https
// @Description Request a random streamer from a country
// @Produce     json
// @Param       lang       query    string false "Language of the streamer, should be an ISO code like fr,de,it. Default to fr"
// @Param       minviewers query    string false "Minimum number of viewers you want the streamer to have. Default to 1"
// @Param       maxviewers query    string false "Maximum number of viewers you want the streamer to have. Default to 10"
// @Success     200        {object} Streamer
// @Router      /streamer [get]
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

	if lang == "" {
		lang = "fr"
	}

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

	val, err := rdb.Get(lang).Result()

	if err == redis.Nil {
		previousPagination := ""
		pagination := true
		log.Printf("streamer list %s not on redis, querying twitch api and add them on redis after", lang)
		for pagination {
			err := requestStreamer(&data, data.Pagination.Cursor, lang, token)
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
			}

		}
		putStreamerListOnRedis(mylist, lang)

	} else if err != nil {
		log.Print(err)
	}

	valToByte := []byte(val)

	err = json.Unmarshal(valToByte, &mylist)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err})
	}

	streamer, err := mylist.GetStreamer(intMinViewers, intMaxViewers)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": err})
		return
	}

	c.JSON(http.StatusOK, streamer)

}
