package agwayland

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/valentim/ag-herald/provider/agstats"

	"github.com/valentim/ag-herald/infrastructure"
)

// SaveCFDResponse will contain the status and the response message.
type SaveCFDResponse struct {
	Data data `json:"data"`
}

type data struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// GenerateCFD will generate the cumulative flow diagram
func GenerateCFD(requestData agstats.ProjectCumulativeFlowDataResponse) (*SaveCFDResponse, error) {
	url := os.Getenv("AGWAYLAND_URL")
	body, bodyError := json.Marshal(requestData)

	if bodyError != nil {
		fmt.Println(bodyError)
		return nil, bodyError
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/diagrams", url), bytes.NewReader(body))

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	response, requestError := infrastructure.DoRequest(req)

	if requestError != nil {
		fmt.Println(requestError)
		return nil, requestError
	}

	responseData := &SaveCFDResponse{}
	jsonParserError := json.Unmarshal(response, responseData)

	if jsonParserError != nil {
		return nil, jsonParserError
	}

	responseData.Data.Message = fmt.Sprintf("https://d2420fd2.ngrok.io/diagrams/%s", responseData.Data.Message)
	return responseData, nil
}
