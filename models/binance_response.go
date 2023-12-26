package models

type BinanceResponse struct {
	Id     int    `json:"id"`
	Result string `json:"result"`
	Stream string `json:"stream"`
	Data   struct {
		E      string `json:"e"`
		S      int64  `json:"E"`
		Symbol string `json:"s"`
		K      struct {
			T                 int64  `json:"t"`
			Timestamp         int64  `json:"T"`
			Symbol            string `json:"s"`
			Interval          string `json:"i"`
			FirstTradeID      int64  `json:"f"`
			LastTradeID       int64  `json:"L"`
			Open              string `json:"o"`
			Close             string `json:"c"`
			High              string `json:"h"`
			Low               string `json:"l"`
			Volume            string `json:"v"`
			NumberOfTrades    int    `json:"n"`
			IsFinal           bool   `json:"x"`
			QuoteVolume       string `json:"q"`
			VolumeActive      string `json:"V"`
			QuoteVolumeActive string `json:"Q"`
			Ignore            string `json:"B"`
		} `json:"k"`
	} `json:"data"`
}
