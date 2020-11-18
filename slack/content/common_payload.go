package content

// EventPayload represents the wrapper with the event dat
type EventPayload struct {
	TeamID string `json:"team_id"`
	Event  Event  `json:"event"`
}

// Event represents the content of the event
type Event struct {
	Type string `json:"type"`
}

// View is the object that contains the blocks
type View struct {
	CallbackID string  `json:"callback_id,omitempty"`
	Type       string  `json:"type,omitempty"`
	Submit     *Submit `json:"submit,omitempty"`
	Title      *Title  `json:"title,omitempty"`
	Blocks     []Block `json:"blocks,omitempty"`
}

// Submit is the object responsible for the submit text and action
type Submit struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// Block is the array with the block of message in Slack concept
type Block struct {
	Type      string     `json:"type"`
	Text      *Text      `json:"text,omitempty"`
	Title     *Title     `json:"title,omitempty"`
	Accessory *Accessory `json:"accessory,omitempty"`
	Element   *Element   `json:"element,omitempty"`
	ImageURL  string     `json:"image_url,omitempty"`
	AltText   string     `json:"alt_text,omitempty"`
	Label     *Label     `json:"label,omitempty"`
	BlockID   string     `json:"block_id,omitempty"`
}

// Label is the way to represent the field
type Label struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
	Emoji bool   `json:"emoji"`
}

// Element is part of interactive components
type Element struct {
	Type        string       `json:"type"`
	Placeholder *Placeholder `json:"placeholder"`
	Options     *[]Option    `json:"options,omitempty"`
	ActionID    string       `json:"action_id,omitempty"`
}

// Accessory is the object with some specific king of resources in Slack concept, for example, buttons...
type Accessory struct {
	Type        string       `json:"type"`
	ActionID    string       `json:"action_id,omitempty"`
	Text        *Text        `json:"text,omitempty"`
	Value       string       `json:"value,omitempty"`
	Placeholder *Placeholder `json:"placeholder,omitempty"`
	Options     *[]Option    `json:"options,omitempty"`
}

// Text is the object with text message in Slack concept
type Text struct {
	Type  string `json:"type"`
	Emoji bool   `json:"emoji"`
	Text  string `json:"text"`
}

// Placeholder is the object with the placeholder data for selectbox
type Placeholder struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
	Emoji bool   `json:"emoji"`
}

// Option is the object with the option value and text for selectbox
type Option struct {
	Text  *Text  `json:"text"`
	Value string `json:"value,omitempty"`
}

// Title is the object with the title message in Slack concept
type Title struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
	Emoji bool   `json:"emoji"`
}

// ResponseAction is the root structure of the action response
type ResponseAction struct {
	Payload *Payload `json:"payload"`
}

// Payload is the data from action response
type Payload struct {
	Type      string    `json:"type"`
	TriggerID string    `json:"trigger_id"`
	Actions   []Actions `json:"actions"`
	Team      Team      `json:"team,omitempty"`
	View      View      `json:"view"`
}

// Team has the unique identifier of the company/team
type Team struct {
	ID string `json:"id,omitempty"`
}

// AuthedUser has the unique identifier of the user that did authorization
type AuthedUser struct {
	ID string `json:"id,omitempty"`
}

// Actions is the data from the interactive component of the action response
type Actions struct {
	Type           string         `json:"type"`
	ActionID       string         `json:"action_id"`
	SelectedOption SelectedOption `json:"selected_option"`
}

// SelectedOption is the selected value of the option field
type SelectedOption struct {
	Text  Text   `json:"text"`
	Value string `json:"value"`
}
