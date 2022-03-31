package main

import (
	"log"
	"net/http"
	"net/url"
	"strings"
)

const LINE_NOTIFY_URL string = "https://notify-api.line.me/api/notify"
const ACCESS_TOKEN string = ""

func notifyLine(m string) error {
	c := &http.Client{}

	v := url.Values{}
	v.Add("message", m)

	body := strings.NewReader(v.Encode())

	req, err := http.NewRequest("POST", LINE_NOTIFY_URL, body)
	if err != nil {
		log.Fatal(err)
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+ACCESS_TOKEN)

	_, err = c.Do(req)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func main() {
	notifyLine("HI")
}
