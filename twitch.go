package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

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

func requestStreamer(data *TwitchResponse, pagination string, token string) error {

	url := "https://api.twitch.tv/helix/streams?language=fr&first=100"
	var bearer = "Bearer " + token

	if pagination != "" {
		url = fmt.Sprintf("https://api.twitch.tv/helix/streams?language=fr&first=100&after=%s", pagination)
	}

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Client-Id", clientID)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &data)

	if err != nil {
		panic(err.Error())
	}

	return err
}

func getStreamers() error {
	var data TwitchResponse
	mylist := StreamerList{}

	token, err := getToken()
	if err != nil {
		return err
	}

	previousPagination := ""
	pagination := true

	requestStreamer(&data, "", token)

	for pagination {
		requestStreamer(&data, data.Pagination.Cursor, token)
		nextPagination := data.Pagination.Cursor

		if previousPagination == nextPagination {
			pagination = false
		} else {
			previousPagination = nextPagination
			nextPagination = ""
		}
		for _, v := range data.Data {

			if v.ViewerCount <= 10 {
				mylist.AddStreamer(Streamer{UserName: v.UserName, ViewerCount: v.ViewerCount, GameName: v.GameName})
			}

		}

	}

	chosenOne := rand.Intn(len(mylist.Streamers))

	streamer := mylist.Streamers[chosenOne]

	fmt.Printf("https://twitch.tv/%s\n", streamer.UserName)

	return nil
}
