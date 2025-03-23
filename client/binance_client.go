package client

import (
	"encoding/json"
	"go-websocket-server/model"
	"log"

	"github.com/gorilla/websocket"
)

type BinanceClient struct {
	Conn    *websocket.Conn
	Params  map[string]int
	Session map[*websocket.Conn]*Subscriber
}

func NewBinanceClient() *BinanceClient {
	connection_url := "wss://stream.binance.com/stream"
	conn, _, err := websocket.DefaultDialer.Dial(connection_url, nil)
	if err != nil {
		log.Fatalln("Failed to connect to Binance WebSocket:", err)
	}
	log.Println("Binance socket is connected!")
	bc := BinanceClient{
		Conn:    conn,
		Params:  make(map[string]int),
		Session: make(map[*websocket.Conn]*Subscriber),
	}
	go bc.ReadMessage()
	return &bc
}

func (bc *BinanceClient) Close() error {
	log.Println("Closing Binance WebSocket connection...")
	return bc.Conn.Close()
}

func (bc *BinanceClient) AddSubscriber(conn *websocket.Conn) *Subscriber {
	if _, ok := bc.Session[conn]; ok {
		return bc.Session[conn]
	}
	subscriber := NewSubscriber(bc, conn)
	bc.Session[conn] = subscriber
	return subscriber
}

func (bc *BinanceClient) ReadMessage() {
	for {
		_, message, err := bc.Conn.ReadMessage()
		if err != nil {
			log.Println("Error reading from Binance WebSocket:", err)
			break
		}

		var binanceData model.BinanceKlineData
		if err := json.Unmarshal(message, &binanceData); err != nil {
			log.Println("Error unmarshaling Binance data:", err)
			continue
		}

		for conn, subscriber := range bc.Session {
			if subscriber.Params[binanceData.Stream] {
				if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
					conn.Close()
					delete(bc.Session, conn)
				}
			}
		}
	}
}
