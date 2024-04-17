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

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Error upgrading to WebSocket:", err)
			return
		}

		defer func() {
			bc.CloseConn(conn)
			conn.Close()
		}()

		for {
			_, p, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					log.Println("WebSocket closed by client:", err)
				} else {
					log.Println("Error reading message:", err)
				}
				break
			}

			if err := json.Unmarshal(p, &requestMessage); err != nil {
				log.Println("Error decoding JSON:", err)
				continue
			}

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
					if err := bc.ReadMessage(); err != nil {
						log.Println("Error reading message from Binance client:", err)
						conn.Close()
						return
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
