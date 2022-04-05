package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const notifyMessage string = "\n施錠に失敗しました。\n"

func notifyLine(m string) error {
	const endpointURL string = "https://notify-api.line.me/api/notify"

	c := &http.Client{}

	v := url.Values{}
	v.Add("message", m)

	body := strings.NewReader(v.Encode())

	req, err := http.NewRequest("POST", endpointURL, body)
	if err != nil {
		log.Fatal(err)
		return err
	}

	accessToken := os.Getenv("LINE_NOTIFY_TOKEN")
	if accessToken == "" {
		panic("LINE Notify Access Token is not set.")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	_, err = c.Do(req)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

type LockStatus int

const (
	Locked LockStatus = iota
	Unlocked
	Moved
)

type Sesame struct {
	BatteryPercentage int
	BatteryVoltage    float64
	Position          int
	Status            LockStatus
	TimeStamp         int
}

func Device() (*Sesame, error) {
	deviseUUID := os.Getenv("SESAME_UUID")
	if deviseUUID == "" {
		panic("SESAME_UUID is not set.")
	}

	c := &http.Client{}

	req, err := http.NewRequest("GET", "https://app.candyhouse.co/api/sesame2/"+deviseUUID, nil)
	if err != nil {
		return nil, err
	}

	apiKey := os.Getenv("SESAME_DEVELOPE_API_KEY")
	if apiKey == "" {
		panic("SESAME_DEVELOPE_API_KEY is not set.")
	}
	req.Header.Set("x-api-key", apiKey)

	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var s Sesame
	err = json.Unmarshal(body, &s)
	return &s, err
}

func main() {
	var s *Sesame
	s, err := Device()
	if err != nil {
		log.Fatal(err)
		notifyLine(notifyMessage)
	}
	if s.BatteryPercentage < 20 {
		notifyLine("バッテリー残量が少なくなっています。")
	}

}
