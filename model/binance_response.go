package model

type BinanceKlineData struct {
	ID     int         `json:"id"`
	Result interface{} `json:"result"`
	Stream string      `json:"stream"`
	Data   KlineData   `json:"data"`
}

type KlineData struct {
	Event  string `json:"e"`
	ETime  int64  `json:"E"`
	Symbol string `json:"s"`
	Kline  Kline  `json:"k"`
}

type Kline struct {
	StartTime           int64  `json:"t"`
	EndTime             int64  `json:"T"`
	Symbol              string `json:"s"`
	Interval            string `json:"i"`
	FirstTradeID        int64  `json:"f"`
	LastTradeID         int64  `json:"L"`
	OpenPrice           string `json:"o"`
	ClosePrice          string `json:"c"`
	HighPrice           string `json:"h"`
	LowPrice            string `json:"l"`
	Volume              string `json:"v"`
	NumberOfTrades      int    `json:"n"`
	IsFinal             bool   `json:"x"`
	QuoteAssetVolume    string `json:"q"`
	TakerBuyBaseVolume  string `json:"V"`
	TakerBuyQuoteVolume string `json:"Q"`
	Ignore              string `json:"B"`
}
