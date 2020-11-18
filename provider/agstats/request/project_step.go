package request

// ProjectStep contains the order and the guid of the project step
type ProjectStep struct {
	Order string `json:"order"`
	GUID  string `json:"guid"`
}
