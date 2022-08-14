package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

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

	body, err := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &data)

	if err != nil {
		return err
	}
	return nil

}
