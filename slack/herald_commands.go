package slack

import (
	"github.com/valentim/ag-herald/slack/content"
)

// ListCommands will show the list of commands that Herald has
func ListCommands() IncomingWebhook {

	blocks := append([]content.Block{}, content.Block{
		Type: "section",
		Text: &content.Text{
			Type:  "plain_text",
			Emoji: true,
			Text:  "/account",
		},
	})

	blocks = append(blocks, content.Block{
		Type: "divider",
	})

	blocks = append(blocks, content.Block{
		Type: "section",
		Text: &content.Text{
			Type:  "plain_text",
			Emoji: true,
			Text:  "/cfd",
		},
	})

	blocks = append(blocks, content.Block{
		Type: "divider",
	})

	blocks = append(blocks, content.Block{
		Type: "section",
		Text: &content.Text{
			Type:  "plain_text",
			Emoji: true,
			Text:  "/projects",
		},
	})

	return IncomingWebhook{
		Blocks: blocks,
	}
}
