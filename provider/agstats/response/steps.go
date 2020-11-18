package response

// Step is the representation of each stage of Kanban
type Step struct {
	GUID string `json:"guid"`
	Name string `json:"name"`
}

// CumulativeFlowPerStep contains the CF Data
type CumulativeFlowPerStep struct {
	Column string  `json:"column"`
	Values []value `json:"values"`
}

type value struct {
	X int64 `json:"x"`
	Y int16 `json:"y"`
}
