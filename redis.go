package main

import (
	"encoding/json"
	"time"
)

func putStreamerListOnRedis(mylist StreamerList, lang string) error {
	e, err := json.Marshal(mylist)
	if err != nil {
		return err
	}

	rdb.Set(lang, e, 10*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil

}
