package models

type BinanceRequest struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
	Id     int      `json:"id"`
}
