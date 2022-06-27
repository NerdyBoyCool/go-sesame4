package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/NerdyBoyCool/sesame"
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

func main() {
	key1 := os.Getenv("SESAME_SECRET_KEY1")
	key2 := os.Getenv("SESAME_SECRET_KEY2")
	apiKey := os.Getenv("SESAME_DEVELOPER_API_KEY")
	sesame1 := os.Getenv("SESAME_UUID1")
	sesame2 := os.Getenv("SESAME_UUID2")
	cli1 := sesame.NewClient(apiKey, key1, sesame1)
	cli2 := sesame.NewClient(apiKey, key2, sesame2)
	ctx := context.Background()
	s1, err := cli1.Device(ctx)
	s2, err := cli2.Device(ctx)
	if err != nil {
		log.Fatal(err)
		notifyLine(notifyMessage)
	}
	if s1.BatteryPercentage < 20 || s2.BatteryPercentage < 20 {
		log.Fatal(err)
		notifyLine("バッテリー残量が少なくなっています。")
	}
	if s1.CHSesame2Status == "unlocked" {
		err = cli1.Lock(ctx, "From Gihub Actions")
		if err != nil {
			log.Fatal(err)
			notifyLine("施錠に失敗しました。")
		}
	}
	if s2.CHSesame2Status == "unlocked" {
		err = cli2.Lock(ctx, "From Gihub Actions")
		if err != nil {
			log.Fatal(err)
			notifyLine("施錠に失敗しました。")
		}
	}
}
