package clients

import (
	"encoding/json"
	"go-websocket-server/data"
	"go-websocket-server/models"
	"log"

	"github.com/gorilla/websocket"
)

type BinanceClient struct {
	Conn *websocket.Conn
}

func NewBinanceClient() (*BinanceClient, error) {
	connection_url := "wss://stream.binance.com/stream"
	conn, _, err := websocket.DefaultDialer.Dial(connection_url, nil)
	if err != nil {
		return nil, err
	}
	log.Println("Binance socket is connected!")
	return &BinanceClient{
		Conn: conn,
	}, nil
}

func (bc *BinanceClient) Close() error {
	return bc.Conn.Close()
}

func (bc *BinanceClient) ReadMessage() (*data.BinanceKlineData, error) {
	_, message, err := bc.Conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	var binnace_data data.BinanceKlineData
	if err := json.Unmarshal(message, &binnace_data); err != nil {
		return nil, err
	}
	return &binnace_data, nil
}

func (bc *BinanceClient) SendMessage(request models.BinanceRequest) {
	data, err := json.Marshal(request)
	if err != nil {
		log.Println(err)
	}
	if err := bc.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Println(err)
	}
}
