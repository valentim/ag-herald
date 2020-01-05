package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/valentim/ag-heralt/slack"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	http.HandleFunc("/tell", sendDiagramHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func sendDiagramHandler(w http.ResponseWriter, r *http.Request) {
	title := &slack.Title{
		Type:  "plain_text",
		Text:  "Example Image",
		Emoji: true,
	}

	block := &slack.Block{
		Type:     "image",
		Title:    *title,
		ImageURL: "https://api.slack.com/img/blocks/bkb_template_images/goldengate.png",
		AltText:  "Example Image",
	}

	var blocks []slack.Block
	blocks = append(blocks, *block)

	incomingMessage := &slack.IncomingWebhook{
		Blocks: blocks,
	}

	incomingMessage.Send()

	jsonResp(w, incomingMessage)
}

func jsonResp(w http.ResponseWriter, message *slack.IncomingWebhook) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	response, err := json.Marshal(message)
	if err != nil {
		log.Println("Couldn't marshal hook response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(response)
}
