package models

type StockMessage struct {
	Type      string   `json:"type"`
	Topic     string   `json:"topic"`
	Component string   `json:"component"`
	Variables []string `json:"variables"`
}