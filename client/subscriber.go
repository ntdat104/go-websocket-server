package client

import (
	"encoding/json"
	"go-websocket-server/model"

	"github.com/gorilla/websocket"
)

type Subscriber struct {
	bc     *BinanceClient
	Conn   *websocket.Conn
	Params map[string]bool
}

func NewSubscriber(bc *BinanceClient, conn *websocket.Conn) *Subscriber {
	return &Subscriber{
		bc:     bc,
		Conn:   conn,
		Params: make(map[string]bool),
	}
}

func (s *Subscriber) Unsubscribe() {
	delete(s.bc.Session, s.Conn)
}

func (s *Subscriber) WriteMessage(request model.BinanceRequest) error {
	if request.Method == "SUBSCRIBE" {
		newParams := make([]string, 0)
		for _, value := range request.Params {
			if _, ok := s.Params[value]; !ok {
				s.Params[value] = true
			}
			if s.bc.Params[value] == 0 {
				s.bc.Params[value] = 1
				newParams = append(newParams, value)
			} else {
				s.bc.Params[value]++
			}
		}
		if len(newParams) > 0 {
			data, err := json.Marshal(model.BinanceRequest{
				Method: request.Method,
				Params: newParams,
				Id:     1,
			})
			if err != nil {
				return err
			}
			s.bc.Conn.WriteMessage(websocket.TextMessage, data)
		}
	}

	if request.Method == "UNSUBSCRIBE" {
		deleteParams := make([]string, 0)
		for _, value := range request.Params {
			if s.Params[value] {
				delete(s.Params, value)
			}
			if s.bc.Params[value] > 1 {
				s.bc.Params[value]--
			} else {
				delete(s.bc.Params, value)
				deleteParams = append(deleteParams, value)
			}
		}
		if len(deleteParams) > 0 {
			data, err := json.Marshal(model.BinanceRequest{
				Method: request.Method,
				Params: deleteParams,
				Id:     1,
			})
			if err != nil {
				return err
			}
			s.bc.Conn.WriteMessage(websocket.TextMessage, data)
		}
	}
	return nil
}
