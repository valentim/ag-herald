package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/joho/godotenv"
	"github.com/valentim/ag-herald/infrastructure"
	"github.com/valentim/ag-herald/infrastructure/database"
	"github.com/valentim/ag-herald/slack"
	"github.com/valentim/ag-herald/slack/content"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	infrastructure.Setup()
}

func main() {
	http.HandleFunc("/v1/command", listCommandHandler)
	http.HandleFunc("/v1/accounts", setupAccountHandler)
	http.HandleFunc("/v1/cfd", listProjectsHandler)
	http.HandleFunc("/v1/projects", listProjectsHandler)
	http.HandleFunc("/v1/events", eventHandler)
	http.HandleFunc("/v1/actions", actionHandler)
	http.HandleFunc("/v1/oauth2/callback", oauth2CallbackHandler)
	log.Fatal(http.ListenAndServe(":4000", nil))
}

func oauth2CallbackHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fmt.Println(r.Form)

	fmt.Println("code", r.Form.Get("code"))
	fmt.Println("state", r.Form.Get("state"))

	values := url.Values{}
	values.Set("code", r.Form.Get("code"))

	err := slack.RequestAccessToken(values)

	if err != nil {
		fmt.Println(err)
		return
	}
}

func setupAccountHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	triggerID := r.Form.Get("trigger_id")
	teamID := r.Form.Get("team_id")

	fmt.Println(triggerID)

	statement := fmt.Sprintf("SELECT accessToken from account WHERE teamID = '%s'", teamID)
	query := database.Query{
		Statement: statement,
	}

	accessToken := query.Get(database.Database{
		Name: "herald.db",
	})

	request := slack.ShowAccountConfiguration(triggerID)
	response, err := json.Marshal(request)
	if err != nil {
		log.Println("Couldn't marshal hook response:", err)
		// w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("request", string(response))
	request.SendToAPI(accessToken)
}

func listProjectsHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	triggerID := r.Form.Get("trigger_id")
	teamID := r.Form.Get("team_id")

	statement := fmt.Sprintf("SELECT accountGUID from account WHERE teamID = '%s'", teamID)
	query := database.Query{
		Statement: statement,
	}

	accountGUID := query.Get(database.Database{
		Name: "herald.db",
	})

	if len(accountGUID) > 0 {
		jsonResp(w, slack.ProjectList(accountGUID))
		return
	}

	request := slack.ShowAccountConfiguration(triggerID)
	response, err := json.Marshal(request)
	if err != nil {
		log.Println("Couldn't marshal hook response:", err)
		// w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("request", string(response))
	request.SendToAPI("asd")
}

func listCommandHandler(w http.ResponseWriter, r *http.Request) {
	responseData := slack.ListCommands()

	jsonResp(w, responseData)
}

type validation struct {
	Challenge string `json:"challenge"`
}

func eventHandler(w http.ResponseWriter, r *http.Request) {
	response, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println(err)
	}

	responseJSON := []byte(response)
	fmt.Println(string(responseJSON))
	validation := &validation{}
	jsonError := json.Unmarshal(responseJSON, validation)

	if jsonError != nil {
		log.Println(err)
		return
	}

	a := []byte(validation.Challenge)

	w.Write(a)
}

func actionHandler(w http.ResponseWriter, r *http.Request) {

	response, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println(err)
	}

	responseJSONStr, err := url.QueryUnescape(string(response)[8:])

	if err != nil {
		log.Println(err)
		return
	}

	log.Println(responseJSONStr)

	responseJSON := []byte(responseJSONStr)
	payload := &content.Payload{}
	jsonError := json.Unmarshal(responseJSON, payload)

	if jsonError != nil {
		log.Println(err)
		return
	}

	if payload.Type == "view_submission" {

		if payload.View.CallbackID == "account" {
			slack.SaveAccountConfiguration(responseJSON)
			return
		}

		slack.SaveProjectConfigurationContent(responseJSON)
		return
	}

	projectActionRawValue := strings.SplitN(payload.Actions[0].SelectedOption.Value, "_", 2)
	projectAction := projectActionRawValue[0]
	projectGUID := projectActionRawValue[1]
	teamID := payload.Team.ID

	statement := fmt.Sprintf("SELECT accessToken from account WHERE teamID = '%s'", teamID)
	query := database.Query{
		Statement: statement,
	}

	accessToken := query.Get(database.Database{
		Name: "herald.db",
	})

	if projectAction == "cfd" {
		cfd := slack.ShowCumulativeFlowDiagram(payload.TriggerID, projectGUID)
		cfd.SendToAPI(accessToken)

		return
	}

	request := slack.ListProjectConfigurationContent(*payload, projectGUID)

	if request != nil {
		request.SendToAPI(accessToken)
		jsonResp(w, request)
	}
}

func jsonResp(w http.ResponseWriter, message interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	response, err := json.Marshal(message)
	if err != nil {
		log.Println("Couldn't marshal hook response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("request", string(response))
	w.Write(response)
}
