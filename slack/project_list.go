package slack

import (
	"fmt"
	"log"

	"github.com/valentim/ag-herald/provider/agstats"
	"github.com/valentim/ag-herald/slack/content"
)

// ProjectList will return the slack structure to show project as a list
func ProjectList(accountGUID string) IncomingWebhook {
	projects, err := agstats.Projects(accountGUID)

	if err != nil {
		log.Println(err)
	}

	var blocks = []content.Block{}

	var title string = "Your project list is bellow"

	if !projects.WasSync {
		title = "Your account is in sync process yet"
	}

	if projects.WasSync && len(projects.Data) > 0 {
		title = "Your account does not have projects yet"
	}

	blocks = append(blocks, content.Block{
		Type: "section",
		Text: &content.Text{
			Type:  "plain_text",
			Emoji: true,
			Text:  title,
		},
	})

	blocks = append(blocks, content.Block{
		Type: "divider",
	})

	for _, project := range projects.Data {
		options := []content.Option{}

		options = append(options, content.Option{
			Text: &content.Text{
				Type: "plain_text",
				Text: "Configure",
			},
			Value: fmt.Sprintf("configure_%s", project.GUID),
		})

		options = append(options, content.Option{
			Text: &content.Text{
				Type: "plain_text",
				Text: "Cumulative flow diagram",
			},
			Value: fmt.Sprintf("cfd_%s", project.GUID),
		})

		block := &content.Block{
			Type: "section",
			Text: &content.Text{
				Type: "plain_text",
				Text: project.Name,
			},
			Accessory: &content.Accessory{
				Type: "static_select",
				Placeholder: &content.Placeholder{
					Type: "plain_text",
					Text: "Project options",
				},
				Options: &options,
			},
		}

		blocks = append(blocks, *block)
	}

	return IncomingWebhook{
		Blocks: blocks,
	}
}
