package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

type ConnectionManager struct {
	clients map[*websocket.Conn][]string
	mutex   sync.Mutex
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		clients: make(map[*websocket.Conn][]string),
	}
}

func (cm *ConnectionManager) AddConnection(conn *websocket.Conn) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.clients[conn] = nil
}

func (cm *ConnectionManager) RemoveConnection(conn *websocket.Conn) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	delete(cm.clients, conn)
}

func (cm *ConnectionManager) Broadcast(message []byte) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	for conn := range cm.clients {
		conn.WriteMessage(websocket.TextMessage, message)
	}
}
