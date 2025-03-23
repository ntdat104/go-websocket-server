package main

import (
	"encoding/json"
	"go-websocket-server/client"
	"go-websocket-server/model"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

func handleWebSocket(bc *client.BinanceClient) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			upgrader = websocket.Upgrader{}
			conn     *websocket.Conn
			request  model.BinanceRequest
		)

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Error upgrading to WebSocket:", err)
			return
		}
		defer conn.Close()

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

			if err := json.Unmarshal(p, &request); err != nil {
				log.Println("Error decoding JSON:", err)
				continue
			}

			bcRequest := model.BinanceRequest{
				Method: request.Method,
				Params: request.Params,
			}
			subscriber := bc.AddSubscriber(conn)
			subscriber.WriteMessage(bcRequest)

			defer subscriber.Unsubscribe()
		}
	}
}

func main() {
	bc := client.NewBinanceClient()
	defer bc.Close()

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}

	http.HandleFunc("/ws", handleWebSocket(bc))
	log.Printf("Websocket server started on: http://localhost:%v/ws", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}
