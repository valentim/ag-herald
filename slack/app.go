package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/valentim/ag-herald/infrastructure"
	"github.com/valentim/ag-herald/infrastructure/database"
	"github.com/valentim/ag-herald/slack/content"
)

// IncomingWebhook is the message format to send messages to the slack channel
type IncomingWebhook struct {
	Blocks []content.Block `json:"blocks"`
}

// IncomingAPI is the top structure to send data to Slack API
type IncomingAPI struct {
	Token     string        `json:"token,omitempty"`
	TriggerID string        `json:"trigger_id,omitempty"`
	UserID    string        `json:"user_id,omitempty"`
	View      *content.View `json:"view"`
}

// AccessTokenResponse contains the access token response
type AccessTokenResponse struct {
	AccessToken string             `json:"access_token"`
	Team        content.Team       `json:"team"`
	AuthedUser  content.AuthedUser `json:"authed_user"`
}

// SendToWebhook is the method to send data to the Slack webhook
func (i *IncomingWebhook) SendToWebhook() error {
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

	_, requestError := infrastructure.DoRequest(req)

	return requestError
}

// SendToAPI is the method to send data to the Slack API
func (i *IncomingAPI) SendToAPI(accessToken string) error {
	url := os.Getenv("SLACK_API_URL")

	data, err := json.Marshal(i)
	if err != nil {
		fmt.Println(err)
		return err
	}

	payload := bytes.NewReader(data)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/views.open", url), payload)
	if err != nil {
		fmt.Println(err)
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	_, requestError := infrastructure.DoRequest(req)

	return requestError
}

// RequestAccessToken will retrieve the access token data
func RequestAccessToken(data url.Values) error {
	url := os.Getenv("SLACK_API_URL")

	data.Set("client_id", os.Getenv("SLACK_CLIENT_ID"))
	data.Set("client_secret", os.Getenv("SLACK_CLIENT_SECRET"))

	payload := strings.NewReader(data.Encode())

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/oauth.v2.access", url), payload)
	if err != nil {
		fmt.Println(err)
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, requestError := infrastructure.DoRequest(req)

	if requestError != nil {
		fmt.Println(requestError)
		return requestError
	}

	accessTokenResponse := &AccessTokenResponse{}

	jsonErr := json.Unmarshal(response, accessTokenResponse)

	if jsonErr != nil {
		fmt.Println(err)
		return err
	}

	d := database.Database{
		Name: "herald.db",
	}

	localAccount := database.LocalAccount{
		TeamID:      accessTokenResponse.Team.ID,
		AccessToken: accessTokenResponse.AccessToken,
	}

	_, insertErr := localAccount.Insert(d)

	var accessTokenInsertErr error

	if insertErr.Error() == "UNIQUE constraint failed: account.teamID" {
		_, accessTokenInsertErr = localAccount.UpdateAccessToken(d)
	}

	if accessTokenInsertErr != nil {
		fmt.Println("accessTokenInsertErr", accessTokenInsertErr)
	}

	var blocks []content.Block

	blocks = append([]content.Block{}, content.Block{
		Type: "section",
		Text: &content.Text{
			Type:  "plain_text",
			Emoji: true,
			Text:  "This project needs to be configured first",
		},
	})

	incomingAPI := &IncomingAPI{
		Token:  accessTokenResponse.AccessToken,
		UserID: accessTokenResponse.AuthedUser.ID,
		View: &content.View{
			Type: "section",
			Title: &content.Title{
				Type: "plain_text",
				Text: "Foda-se",
			},
			Blocks: blocks,
		},
	}

	sendToUserErr := incomingAPI.SendToUser()

	return sendToUserErr
}

// SendToUser will send message to a user
func (i *IncomingAPI) SendToUser() error {
	url := os.Getenv("SLACK_API_URL")

	data, err := json.Marshal(i)
	if err != nil {
		fmt.Println(err)
		return err
	}

	payload := bytes.NewReader(data)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/views.publish", url), payload)
	if err != nil {
		fmt.Println(err)
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	_, requestError := infrastructure.DoRequest(req)

	if requestError != nil {
		fmt.Println(requestError)
		return requestError
	}

	return nil
}
