package main

import (
	"fmt"
	"math/rand"
	"time"
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
		Cursor string `json:"cursor,omitempty"`
	} `json:"pagination,omitempty"`
}

type StreamerList struct {
	Streamers []Streamer `json:"streamers"`
}

func (r *StreamerList) AddStreamer(item Streamer) []Streamer {
	r.Streamers = append(r.Streamers, item)
	return r.Streamers
}

func (r *StreamerList) GetStreamer(minViewers int, maxViewers int) (Streamer, error) {

	var newList []Streamer

	if len(r.Streamers) == 0 {
		return Streamer{}, fmt.Errorf("no streamer found")
	}

	for _, streamer := range r.Streamers {
		if streamer.ViewerCount >= minViewers && streamer.ViewerCount <= maxViewers {
			newList = append(newList, streamer)
		}

	}

	chosenOne := rand.Intn(len(newList))

	streamer := newList[chosenOne]
	return streamer, nil

}

type Streamer struct {
	UserName    string
	ViewerCount int
	GameName    string
}
