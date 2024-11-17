package types

type JiraLinkDataSource struct {
	ID         string             `json:"id"`
	Parameters JiraLinkParameters `json:"parameters"`
	Views      []JiraLinkView     `json:"views"`
}

type JiraLinkParameters struct {
	JQL     string `json:"jql"`
	CloudID string `json:"cloudId"`
}

type JiraLinkView struct {
	Type       string           `json:"type"`
	Properties JiraViewProperty `json:"properties"`
}

type JiraViewProperty struct {
	Columns []JiraPropertyKey `json:"columns"`
}
type JiraPropertyKey struct {
	Key       string `json:"key"`
	Width     uint   `json:"width,omitempty"`
	IsWrapped bool   `json:"isWrapped,omitempty"`
}
