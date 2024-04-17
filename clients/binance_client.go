package clients

import (
	"encoding/json"
	"go-websocket-server/data"
	"go-websocket-server/models"
	"go-websocket-server/utils"
	"log"

	"github.com/gorilla/websocket"
)

type BinanceClient struct {
	Conn    *websocket.Conn
	Params  []string
	Session map[*websocket.Conn][]string
}

func NewBinanceClient() (*BinanceClient, error) {
	connection_url := "wss://stream.binance.com/stream"
	conn, _, err := websocket.DefaultDialer.Dial(connection_url, nil)
	if err != nil {
		return nil, err
	}
	log.Println("Binance socket is connected!")
	return &BinanceClient{
		Conn:    conn,
		Params:  make([]string, 0),
		Session: make(map[*websocket.Conn][]string),
	}, nil
}

func (bc *BinanceClient) Close() error {
	return bc.Conn.Close()
}

func (bc *BinanceClient) CloseConn(conn *websocket.Conn) {
	delete(bc.Session, conn)
}

func (bc *BinanceClient) ReadMessage() error {
	_, message, err := bc.Conn.ReadMessage()
	if err != nil {
		return err
	}
	var binnace_data data.BinanceKlineData
	if err := json.Unmarshal(message, &binnace_data); err != nil {
		return err
	}

	for key, value := range bc.Session {
		if utils.IndexOf(value, binnace_data.Stream) >= 0 {
			if err := key.WriteMessage(websocket.TextMessage, message); err != nil {
				delete(bc.Session, key)
				return err
			}
		}
	}

	return nil
}

func (bc *BinanceClient) SendMessage(conn *websocket.Conn, request models.BinanceRequest) {
	log.Println(request)

	if bc.Session[conn] != nil {
		params := bc.Session[conn]

		if request.Method == "SUBSCRIBE" {
			params = mergeArrays(params, request.Params)
		} else {
			for _, value := range request.Params {
				params = utils.Filter(params, func(c string, _ int) bool {
					return c != value
				})
			}
		}

		bc.Session[conn] = params
	} else {
		bc.Session[conn] = request.Params
	}

	if request.Method == "SUBSCRIBE" {
		bc.Params = mergeArrays(bc.Params, request.Params)
		data, err := json.Marshal(request)
		if err != nil {
			log.Println(err)
		}
		if err := bc.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Println(err)
		}
	}

}

func mergeArrays[T comparable](arr1, arr2 []T) []T {
	uniqueElements := make(map[T]bool)

	for _, num := range arr1 {
		uniqueElements[num] = true
	}

	for _, num := range arr2 {
		uniqueElements[num] = true
	}

	mergedArray := make([]T, 0, len(uniqueElements))

	for num := range uniqueElements {
		mergedArray = append(mergedArray, num)
	}

	return mergedArray
}
