package slack

import (
	"fmt"
	"log"

	"github.com/stretchr/objx"
	"github.com/valentim/ag-herald/provider/agstats"
	"github.com/valentim/ag-herald/provider/agstats/request"
	"github.com/valentim/ag-herald/provider/agwayland"
	"github.com/valentim/ag-herald/slack/content"
)

// ListProjectConfigurationContent returns the slack modal content
func ListProjectConfigurationContent(payload content.Payload, projectGUID string) *IncomingAPI {
	steps, err := agstats.ProjectSteps(projectGUID)

	if err != nil {
		log.Println(err)
	}

	var blocks = []content.Block{}

	blocks = append(blocks, content.Block{
		Type:    "section",
		BlockID: projectGUID,
		Text: &content.Text{
			Type:  "plain_text",
			Emoji: true,
			Text:  "Configure the project step order",
		},
	})

	blocks = append(blocks, content.Block{
		Type: "divider",
	})

	options := []content.Option{}
	for i := range steps.Data {
		order := i + 1
		options = append(options, content.Option{
			Text: &content.Text{
				Type: "plain_text",
				Text: fmt.Sprintf("%d", order),
			},
			Value: fmt.Sprintf("%d", order),
		})
	}

	stepOptions := []content.Option{}

	for _, step := range steps.Data {

		stepOptions = append(stepOptions, content.Option{
			Text: &content.Text{
				Type: "plain_text",
				Text: step.Name,
			},
			Value: step.GUID,
		})

		block := &content.Block{
			Type: "input",
			Element: &content.Element{
				Type: "static_select",
				Placeholder: &content.Placeholder{
					Type: "plain_text",
					Text: "Step order",
				},
				Options:  &options,
				ActionID: step.GUID,
			},
			Label: &content.Label{
				Type: "plain_text",
				Text: step.Name,
			},
		}

		blocks = append(blocks, *block)
	}

	blocks = append(blocks, content.Block{
		Type: "divider",
	})

	blocks = append(blocks, content.Block{
		Type: "section",
		Text: &content.Text{
			Type:  "plain_text",
			Emoji: true,
			Text:  "Configure the cycle time first step",
		},
	})

	blocks = append(blocks, content.Block{
		Type: "input",
		Element: &content.Element{
			Type: "static_select",
			Placeholder: &content.Placeholder{
				Type: "plain_text",
				Text: "Step name",
			},
			Options:  &stepOptions,
			ActionID: "startingCycleTime",
		},
		Label: &content.Label{
			Type: "plain_text",
			Text: "Cycle time starting at",
		},
	})

	triggerID := payload.TriggerID

	return &IncomingAPI{
		TriggerID: triggerID,
		View: &content.View{
			Type: "modal",
			Title: &content.Title{
				Type: "plain_text",
				Text: "Test",
			},
			Submit: &content.Submit{
				Type: "plain_text",
				Text: "Save",
			},
			Blocks: blocks,
		},
	}
}

// SaveProjectConfigurationContent will save the options chosen by the user
func SaveProjectConfigurationContent(jsonContent []byte) {

	m, objxJSONError := objx.FromJSON(string(jsonContent))

	if objxJSONError != nil {
		log.Println(objxJSONError)
	}

	var startingCycleTime string
	projectSteps := []request.ProjectStep{}

	for _, block := range m.Get("view.blocks").ObjxMapSlice() {
		if len(block.Get("element.action_id").String()) != 0 {
			identifier := block.Get("element.action_id").String()
			value := m.Get(fmt.Sprintf("view.state.values.%s.%s.selected_option.value", block.Get("block_id"), identifier)).String()

			if identifier != "startingCycleTime" {
				projectStep := &request.ProjectStep{}
				projectStep.GUID = identifier
				projectStep.Order = value

				projectSteps = append(projectSteps, *projectStep)
			} else {
				startingCycleTime = value
			}
		}
	}

	project := &agstats.ProjectStepRequest{
		StartingCycleTime: startingCycleTime,
		ProjectSteps:      projectSteps,
	}

	projectGUID := m.Get("view.blocks[0].block_id").String()
	project.ConfigureProject(projectGUID)
}

// ShowCumulativeFlowDiagram will print the CFD image
func ShowCumulativeFlowDiagram(triggerID string, projectGUID string) *IncomingAPI {
	projectCumulativeFlowDataResponse, agstatsErr := agstats.ProjectCumulativeFlowData(projectGUID)

	if agstatsErr != nil {
		fmt.Println("AgstatsErr error", agstatsErr)

		return nil
	}

	var blocks []content.Block

	if projectCumulativeFlowDataResponse.WasSetupDone {
		response, waylandErr := agwayland.GenerateCFD(projectCumulativeFlowDataResponse)

		fmt.Println("Wayland response", response)
		if waylandErr != nil {
			fmt.Println("Wayland error", waylandErr)
			return nil
		}

		blocks = append([]content.Block{}, content.Block{
			Type: "image",
			Title: &content.Title{
				Type:  "plain_text",
				Text:  "Test",
				Emoji: true,
			},
			ImageURL: response.Data.Message,
			AltText:  "Example Image",
		})
	} else {
		blocks = append([]content.Block{}, content.Block{
			Type: "section",
			Text: &content.Text{
				Type:  "plain_text",
				Emoji: true,
				Text:  "This project needs to be configured first",
			},
		})
	}

	return &IncomingAPI{
		TriggerID: triggerID,
		View: &content.View{
			Type: "modal",
			Title: &content.Title{
				Type: "plain_text",
				Text: "Test",
			},
			Blocks: blocks,
		},
	}
}
