package model

type TScraping struct {
	Keyword string `json:"keyword,omitempty"`
	TotalAdvertised   int         `json:"totalAdvertised,omitempty"`
	TotalLink int `json:"totalLink,omitempty"`
	TotalSearch string `json:"totalSearch,omitempty"`
	Html string `json:"html,omitempty"`
}
