package main

import (
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

	access_token := os.Getenv("LINE_NOTIFY_TOKEN")
	if access_token == "" {
		panic("LINE Notify Access Token is not set.")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+access_token)

	_, err = c.Do(req)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func main() {
	notifyLine(notifyMessage)
}
