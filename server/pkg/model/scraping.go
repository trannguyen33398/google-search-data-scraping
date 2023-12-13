package model

type TScraping struct {
	TotalAdvertised   int         `json:"totalAdvertised,omitempty"`
	TotalLink int `json:"totalLink,omitempty"`
	TotalSearch string `json:"totalSearch,omitempty"`
	Html string `json:"html,omitempty"`
}
