package models

type RequestMessage struct {
	Type   string   `json:"type"`
	Method string   `json:"method"`
	Params []string `json:"params"`
}