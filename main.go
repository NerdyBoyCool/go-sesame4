package main

import (
	"bytes"
	"crypto/aes"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/chmike/cmac-go"
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

	apiKey := os.Getenv("SESAME_DEVELOPER_API_KEY")
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

type Param struct {
	Cmd     int16  `json:"cmd"`
	History string `json:"history"`
	Sign    string `json:"sign"`
}

func toLittleEndian(t int64) []byte {
	i := int32(t)
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(i))
	return b
}

func toCMAC(secretKey string) (string, error) {
	message := toLittleEndian(time.Now().Unix())
	byteKey, err := hex.DecodeString(secretKey)
	if err != nil {
		return "", err
	}
	cm, err := cmac.New(aes.NewCipher, byteKey)
	if err != nil {
		return "", err
	}
	cm.Write(message[1:4])
	m := cm.Sum(nil)
	return hex.EncodeToString(m), nil
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

	key := os.Getenv("SESAME_SECRET_KEY")
	m, err := toCMAC(key)
	if err != nil {
		fmt.Errorf("cannot get cmac: %w", err)
		return
	}
	apiKey := os.Getenv("SESAME_DEVELOPER_API_KEY")
	from := base64.StdEncoding.EncodeToString([]byte("by API"))
	param := &Param{Cmd: 88, Sign: m, History: from}
	data, _ := json.Marshal(param)
	deviseUUID := os.Getenv("SESAME_UUID")
	req, err := http.NewRequest("POST", "https://app.candyhouse.co/api/sesame2/"+deviseUUID+"/cmd", bytes.NewBuffer(data))
	if err != nil {
		fmt.Errorf("cannot create HTTP request: %w", err)
		return
	}
	req.Header.Set("x-api-key", apiKey)
	var client *http.Client = &http.Client{}
	res, err := client.Do(req)
	defer res.Body.Close()
	if err != nil {
		fmt.Errorf("Request Error: %w", err)
		return
	}

	// body, _ := io.ReadAll(res.Body)
	// fmt.Println(string(body))

}
