package websocketTool

import (
	"github.com/gorilla/websocket"
	"sync"
)

type Manager struct {
	clients ClientList
	// Using a syncMutex here to be able to lcok state before editing clients
	// Could also use Channels to block
	sync.RWMutex
}

type Client struct {
	manager *Manager

	egress chan []byte

	connection *websocket.Conn
}

type ClientList map[*Client]bool
