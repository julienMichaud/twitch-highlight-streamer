package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/oauth2/twitch"
)

type TwitchResponse struct {
	Data []struct {
		ID          string    `json:"id"`
		UserID      string    `json:"user_id"`
		UserLogin   string    `json:"user_login"`
		UserName    string    `json:"user_name"`
		GameID      string    `json:"game_id"`
		GameName    string    `json:"game_name"`
		Type        string    `json:"type"`
		Title       string    `json:"title"`
		ViewerCount int       `json:"viewer_count"`
		StartedAt   time.Time `json:"started_at"`
		Language    string    `json:"language"`
	} `json:"data,omitempty"`
	Pagination struct {
		Cursor string `json:"cursor"`
	} `json:"pagination,omitempty"`
}

type StreamerList struct {
	Streamers []Streamer
}

func (r *StreamerList) AddStreamer(item Streamer) []Streamer {
	r.Streamers = append(r.Streamers, item)
	return r.Streamers
}

type Streamer struct {
	UserName    string
	ViewerCount int
	GameName    string
}

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

func requestStreamer(data *TwitchResponse, pagination string, lang, token string) error {

	if lang == "" {
		lang = "fr"
	}

	url := fmt.Sprintf("https://api.twitch.tv/helix/streams?language=%s&first=100", lang)
	var bearer = "Bearer " + token

	if pagination != "" {
		url = fmt.Sprintf(url+"&after=%s", pagination)
	}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return err
	}
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Client-Id", clientID)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	log.Printf(resp.Status)

	body, err := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &data)

	if err != nil {
		return err
	}
	return nil

}

func getStreamers(c *gin.Context) {
	var data TwitchResponse
	mylist := StreamerList{}
	lang := c.Query("lang")
	minViewers := c.Query("minviewers")
	maxViewers := c.Query("maxviewers")

	log.Printf("%s", minViewers)
	log.Printf("%s", maxViewers)

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

	token, err := getToken()
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err})
	}

	previousPagination := ""
	pagination := true

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

			if v.ViewerCount >= intMinViewers && v.ViewerCount <= intMaxViewers {
				mylist.AddStreamer(Streamer{UserName: v.UserName, ViewerCount: v.ViewerCount, GameName: v.GameName})
			}

		}

	}

	if len(mylist.Streamers) == 0 {
		c.JSON(http.StatusBadGateway, gin.H{
			"streamer": "no streamer found"})
		return
	}

	chosenOne := rand.Intn(len(mylist.Streamers))

	streamer := mylist.Streamers[chosenOne]

	c.JSON(http.StatusOK, gin.H{
		"streamer": streamer.UserName,
		"viewers":  streamer.ViewerCount})

}
