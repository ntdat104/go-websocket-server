package main

import (
	"encoding/json"
	"go-websocket-server/clients"
	"go-websocket-server/models"
	"log"
	"math/rand"
	"net/http"

	"github.com/gorilla/websocket"
)

func handleWebSocket(bc *clients.BinanceClient) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			upgrader       = websocket.Upgrader{}
			conn           *websocket.Conn
			requestMessage models.RequestMessage
		)

		// Khởi tạo websocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Error upgrading to WebSocket:", err)
			return
		}

		// Đóng websocket khi xong
		defer conn.Close()

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
				bcRequest := models.BinanceRequest{
					Method: requestMessage.Method,
					Params: requestMessage.Params,
					Id:     rand.Int(),
				}
				bc.SendMessage(conn, bcRequest)
			}

			go func() {
				for {
					message, err := bc.ReadMessage()
					if err != nil {
						log.Println("Error reading from Binance WebSocket:", err)
						break
					}
					if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
						log.Println("Error write message from Binance to Client", err)
						break
					}
				}
			}()
		}
	}
}

func main() {
	bc, err := clients.NewBinanceClient()
	if err != nil {
		log.Println("Fail to connect Binance")
	}

	http.HandleFunc("/ws", handleWebSocket(bc))
	log.Println("Websocket server started on: http://localhost:8888/ws")
	log.Fatal(http.ListenAndServe(":8888", nil))
}
