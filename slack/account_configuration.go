package slack

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"

	"github.com/stretchr/objx"
	"github.com/valentim/ag-herald/infrastructure/database"
	"github.com/valentim/ag-herald/provider/agstats"
	"github.com/valentim/ag-herald/slack/content"
)

// ShowAccountConfiguration the fields to configure the account
func ShowAccountConfiguration(triggerID string) *IncomingAPI {

	var blocks = []content.Block{}

	blocks = append(blocks, content.Block{
		Type: "section",
		Text: &content.Text{
			Type:  "plain_text",
			Emoji: true,
			Text:  "Integration - Jira account",
		},
	})

	blocks = append(blocks, content.Block{
		Type: "divider",
	})

	blocks = append(blocks, content.Block{
		Type:    "input",
		BlockID: "email",
		Element: &content.Element{
			Type: "plain_text_input",
			Placeholder: &content.Placeholder{
				Type: "plain_text",
				Text: "Jira access e-mail",
			},
			ActionID: "email",
		},
		Label: &content.Label{
			Type: "plain_text",
			Text: "E-mail",
		},
	})

	blocks = append(blocks, content.Block{
		Type:    "input",
		BlockID: "token",
		Element: &content.Element{
			Type: "plain_text_input",
			Placeholder: &content.Placeholder{
				Type: "plain_text",
				Text: "Jira token",
			},
			ActionID: "token",
		},
		Label: &content.Label{
			Type: "plain_text",
			Text: "Token",
		},
	})

	blocks = append(blocks, content.Block{
		Type:    "input",
		BlockID: "company",
		Element: &content.Element{
			Type: "plain_text_input",
			Placeholder: &content.Placeholder{
				Type: "plain_text",
				Text: "Company name",
			},
			ActionID: "company",
		},
		Label: &content.Label{
			Type: "plain_text",
			Text: "Company",
		},
	})

	blocks = append(blocks, content.Block{
		Type:    "input",
		BlockID: "jira",
		Element: &content.Element{
			Type: "plain_text_input",
			Placeholder: &content.Placeholder{
				Type: "plain_text",
				Text: "https://company_path.atlassian.net",
			},
			ActionID: "jira",
		},
		Label: &content.Label{
			Type: "plain_text",
			Text: "Jira url",
		},
	})

	return &IncomingAPI{
		TriggerID: triggerID,
		View: &content.View{
			CallbackID: "account",
			Type:       "modal",
			Title: &content.Title{
				Type: "plain_text",
				Text: "Account configuration",
			},
			Submit: &content.Submit{
				Type: "plain_text",
				Text: "Save",
			},
			Blocks: blocks,
		},
	}
}

// SaveAccountConfiguration will save the new account
func SaveAccountConfiguration(payload []byte) {

	m, objxJSONError := objx.FromJSON(string(payload))

	if objxJSONError != nil {
		log.Println(objxJSONError)
	}

	email := m.Get("view.state.values.email.email.value").String()
	token := m.Get("view.state.values.token.token.value").String()
	jira := m.Get("view.state.values.jira.jira.value").String()
	company := m.Get("view.state.values.company.company.value").String()
	secret := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", email, token)))

	dataJSON := fmt.Sprintf(
		"{\"companyName\": \"%s\", \"baseUrl\": \"%s\", \"publicKey\": \"%s\", \"sharedSecret\": \"%s\"}",
		company,
		jira,
		token,
		fmt.Sprintf("Basic %s", secret),
	)

	accountRequest := &agstats.AccountRequest{}

	jsonErr := json.Unmarshal([]byte(dataJSON), accountRequest)
	if jsonErr != nil {
		fmt.Println("jsonErr", jsonErr)
	}

	accountResponse, err := accountRequest.CreateAccount()

	if err != nil {
		fmt.Println("errSaveAccount", err)
	}

	d := database.Database{
		Name: "herald.db",
	}

	localAccount := database.LocalAccount{
		TeamID:      m.Get("team.id").String(),
		AccountGUID: accountResponse.Account.GUID,
	}

	localAccount.Insert(d)
}
