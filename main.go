package main

import (
	"encoding/json"
	"go-websocket-server/models"
	"log"
	"math/rand"
	"net/http"

	"github.com/gorilla/websocket"
)

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	var (
		upgrader                  = websocket.Upgrader{}
		conn                      *websocket.Conn
		requestMessage            models.RequestMessage
		clientCrytoConn           *websocket.Conn
		clientStockConn           *websocket.Conn
		cryptoRequestMessageQueue = make(chan []byte, 10)
		stockRequestMessageQueue = make(chan []byte, 10)
		responseMessageQueue      = make(chan []byte, 10)
	)

	// Khởi tạo websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}

	// Đóng websocket khi xong
	defer func() {
		if clientCrytoConn != nil {
			if err := clientCrytoConn.Close(); err != nil {
				log.Println("Error closing Binance WebSocket connection:", err)
			}
		}
		conn.Close()
	}()

	for {
		// Đọc message từ client -> websocket server
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}

		// Convert dữ liệu đọc binary -> json, gán vào địa chỉ của requestMessage
		if err := json.Unmarshal(p, &requestMessage); err != nil {
			log.Println("Error decoding JSON:", err)
			continue
		}

		// Kiểm tra requestMessage.Type == CRYPTO | STOCK
		switch {
		case requestMessage.Type == "CRYPTO":
			cryptoMessage, err := json.Marshal(models.CryptoMessage{
				Method: requestMessage.Method, 
				Params: requestMessage.Params, 
				Id: rand.Int(),
			})
			if err != nil {
				log.Println("Error encoding JSON:", err)
				break
			}
			cryptoRequestMessageQueue <- cryptoMessage
		case requestMessage.Type == "STOCK":
			stockMessage, err := json.Marshal(models.StockMessage{
				Type: "sub",
				Topic: "stockRealtimeByListV2",
				Variables: requestMessage.Params,
				Component: "priceTableEquities",
			})
			if err != nil {
				log.Println("Error encoding JSON:", err)
				break
			}
			stockRequestMessageQueue <- stockMessage
		}

		// Nếu có msg trong responseMessageQueue thì trả về cho client
		go func() {
			for msg := range responseMessageQueue {
				if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
					log.Println("Error forwarding message to client:", err)
					break
				}
			}
		}()

		// Nếu requestMessage.Type == CRYPTO
		if requestMessage.Type == "CRYPTO" {

			// Nếu chưa khởi tạo connection tới binance thì khởi tạo
			if clientCrytoConn == nil {
				log.Println("Create a connection to wss://stream.binance.com/stream")
				clientCrytoConn, _, err = websocket.DefaultDialer.Dial("wss://stream.binance.com/stream", nil)
				if err != nil {
					log.Println("Error connecting to Binance WebSocket:", err)
					break
				}
			}

			// Nếu có msg trong cryptoRequestMessageQueue thì bắn msg cho binance
			go func() {
				for msg := range cryptoRequestMessageQueue {
					if err := clientCrytoConn.WriteMessage(websocket.TextMessage, msg); err != nil {
						log.Println("Error sending message to Binance WebSocket:", err)
						break
					}
				}
			}()

			// Nếu binance có trả về message thì đẩy message vào responseMessageQueue
			go func() {
				for {
					_, message, err := clientCrytoConn.ReadMessage()
					if err != nil {
						log.Println("Error reading from Binance WebSocket:", err)
						break
					}
					var messageJson models.BinanceResponse
					if err := json.Unmarshal(message, &messageJson); err != nil {
						log.Println("Error decoding JSON BinanceResponse:", err)
						break
					}
					if messageJson.Stream != "" {
						responseMessageQueue <- message
					}
				}
			}()
		}

		// Nếu requestMessage.Type == STOCK
		if requestMessage.Type == "STOCK" {

			// Nếu chưa khởi tạo connection tới SSI thì khởi tạo
			if clientStockConn == nil {
				log.Println("Create a connection to wss://iboard-pushstream.ssi.com.vn/realtime")
				clientStockConn, _, err = websocket.DefaultDialer.Dial("wss://iboard-pushstream.ssi.com.vn/realtime", nil)
				if err != nil {
					log.Println("Error connecting to SSI WebSocket:", err)
					break
				}
			}

			// Nếu có msg trong stockRequestMessageQueue thì bắn msg cho SSI
			go func() {
				for msg := range stockRequestMessageQueue {
					if err := clientStockConn.WriteMessage(websocket.TextMessage, msg); err != nil {
						log.Println("Error sending message to SSI WebSocket:", err)
						break
					}
				}
			}()

			// Nếu SSI có trả về message thì đẩy message vào responseMessageQueue
			go func() {
				for {
					_, message, err := clientStockConn.ReadMessage()
					if err != nil {
						log.Println("Error reading from SSI WebSocket:", err)
						break
					}
					log.Printf("value: %v", message)
					// responseMessageQueue <- message
				}
			}()
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	log.Println("Websocket server started on: http://localhost:8888/ws")
	log.Fatal(http.ListenAndServe(":8888", nil))
}
