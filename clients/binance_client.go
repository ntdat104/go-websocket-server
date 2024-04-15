package clients

import (
	"encoding/json"
	"go-websocket-server/models"
	"go-websocket-server/utils"
	"log"

	"github.com/gorilla/websocket"
)

// BinanceClient struct to hold WebSocket connection
type BinanceClient struct {
	Conn    *websocket.Conn
	params  utils.HashSet
	session map[*websocket.Conn][]string
}

// NewBinanceClient creates a new BinanceClient instance
func NewBinanceClient() (*BinanceClient, error) {
	connection_url := "wss://stream.binance.com/stream"
	conn, _, err := websocket.DefaultDialer.Dial(connection_url, nil)
	if err != nil {
		return nil, err
	}
	log.Println("Binance socket is connected!")
	return &BinanceClient{
		Conn:    conn,
		params:  make(utils.HashSet),
		session: make(map[*websocket.Conn][]string),
	}, nil
}

// Close closes the WebSocket connection
func (bc *BinanceClient) Close() error {
	return bc.Conn.Close()
}

// ReadMessage reads a message from the WebSocket connection
func (bc *BinanceClient) ReadMessage() ([]byte, error) {
	_, message, err := bc.Conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (bc *BinanceClient) SendMessage(conn *websocket.Conn, request models.BinanceRequest) {
	data, err := json.Marshal(request)
	if err != nil {
		log.Println(err)
	}
	if err := bc.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Println(err)
	}
}
