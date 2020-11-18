package agstats

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/valentim/ag-herald/infrastructure"
	"github.com/valentim/ag-herald/provider/agstats/request"
	"github.com/valentim/ag-herald/provider/agstats/response"
)

// ProjectListResponse is the data structure that represents the project list
type ProjectListResponse struct {
	Data    []response.Project `json:"data"`
	WasSync bool               `json:"was_sync"`
}

// ProjectStepsResponse is the data structure that represents the project steps of response
type ProjectStepsResponse struct {
	Data []response.Step `json:"data"`
}

// ProjectStepRequest is the data structure that represents the project steps of request
type ProjectStepRequest struct {
	StartingCycleTime string                `json:"startingCycleTime"`
	ProjectSteps      []request.ProjectStep `json:"projectSteps"`
}

// AccountRequest the account data structure
type AccountRequest struct {
	ClientKey    string `json:"clientKey"`
	SharedSecret string `json:"sharedSecret"`
	PublicKey    string `json:"publicKey"`
	BaseURL      string `json:"baseUrl"`
	CompanyName  string `json:"companyName"`
}

// AccountResponse is the data structure returned from API
type AccountResponse struct {
	Status  string           `json:"status"`
	Account response.Account `json:"account"`
}

// ProjectCumulativeFlowDataResponse is the data structure that represents the cumulative flow data
type ProjectCumulativeFlowDataResponse struct {
	Data         []response.CumulativeFlowPerStep `json:"data"`
	WasSetupDone bool                             `json:"was_setup_done"`
}

// Projects will list the projects of an account
func Projects(accountGUID string) (ProjectListResponse, error) {
	url := os.Getenv("AGSTATS_URL")
	data := ProjectListResponse{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/projects?accountGuid=%s", url, accountGUID), nil)

	if err != nil {
		fmt.Println(err)
		return data, err
	}

	req.Header.Add("Content-Type", "application/json")

	response, requestError := infrastructure.DoRequest(req)

	if requestError != nil {
		return data, requestError
	}

	jsonParserError := json.Unmarshal(response, &data)

	if jsonParserError != nil {
		return data, jsonParserError
	}

	return data, nil
}

// ProjectSteps will return the project steps
func ProjectSteps(projectGUID string) (ProjectStepsResponse, error) {
	url := os.Getenv("AGSTATS_URL")
	data := ProjectStepsResponse{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/projects/%s", url, projectGUID), nil)

	if err != nil {
		fmt.Println(err)
		return data, err
	}

	req.Header.Add("Content-Type", "application/json")

	response, requestError := infrastructure.DoRequest(req)

	if requestError != nil {
		return data, requestError
	}

	jsonParserError := json.Unmarshal(response, &data)

	if jsonParserError != nil {
		return data, jsonParserError
	}

	return data, nil
}

// ProjectCumulativeFlowData will retrieve the cumulative flow data
func ProjectCumulativeFlowData(projectGUID string) (ProjectCumulativeFlowDataResponse, error) {
	url := os.Getenv("AGSTATS_URL")
	data := ProjectCumulativeFlowDataResponse{}
	lastWeek := time.Now().AddDate(0, -1, 0).Format("2006-01-02")

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/projects/%s/stats/cumulativeflow?from=%s", url, projectGUID, lastWeek), nil)

	if err != nil {
		fmt.Println(err)
		return data, err
	}

	req.Header.Add("Content-Type", "application/json")

	response, requestError := infrastructure.DoRequest(req)

	if requestError != nil {
		return data, requestError
	}

	jsonParserError := json.Unmarshal(response, &data)

	if jsonParserError != nil {
		return data, jsonParserError
	}

	return data, nil
}

// CreateAccount will register a new account
func (a AccountRequest) CreateAccount() (*AccountResponse, error) {
	url := os.Getenv("AGSTATS_URL")
	body, bodyError := json.Marshal(a)

	response := &AccountResponse{}

	if bodyError != nil {
		fmt.Println(bodyError)
		return response, bodyError
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/accounts", url), bytes.NewReader(body))

	if err != nil {
		fmt.Println(err)
		return response, err
	}

	req.Header.Add("Content-Type", "application/json")

	byteResponse, requestError := infrastructure.DoRequest(req)

	if requestError != nil {
		return response, requestError
	}
	jsonErr := json.Unmarshal(byteResponse, response)

	if jsonErr != nil {
		return response, jsonErr
	}

	return response, nil
}

// ConfigureProject will update the project
func (p ProjectStepRequest) ConfigureProject(projectGUID string) error {
	url := os.Getenv("AGSTATS_URL")
	body, bodyError := json.Marshal(p)

	if bodyError != nil {
		fmt.Println(bodyError)
		return bodyError
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/projects/%s", url, projectGUID), bytes.NewReader(body))

	if err != nil {
		fmt.Println(err)
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	_, requestError := infrastructure.DoRequest(req)

	if requestError != nil {
		return requestError
	}

	return nil
}
