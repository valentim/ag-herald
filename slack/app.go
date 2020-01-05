package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// IncomingWebhook is the message format to send messages to the slack channel
type IncomingWebhook struct {
	Blocks []Block `json:"blocks"`
}

// Block is the array with the block of message in Slack concept
type Block struct {
	Type     string `json:"type"`
	Title    Title  `json:"title"`
	ImageURL string `json:"image_url"`
	AltText  string `json:"alt_text"`
}

// Title is the object with the text message in Slack concept
type Title struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
	Emoji bool   `json:"emoji"`
}

// Send is the method to send data to the Slack webhook
func (i *IncomingWebhook) Send() error {
	url := os.Getenv("WEBHOOK_URL")

	data, err := json.Marshal(i)
	if err != nil {
		fmt.Println(err)
		return err
	}

	payload := bytes.NewReader(data)

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		fmt.Println(err)
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	return doRequest(req)
}

func doRequest(req *http.Request) error {
	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))

	return err
}
