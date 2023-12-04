package websocketutil

import (
	"github.com/gorilla/websocket"
	"sync"
)

type Manager struct {
	clients ClientMap
	// Using a syncMutex here to be able to lcok state before editing clients
	// Could also use Channels to block
	//如果兩個 goroutine 同時嘗試將同一個客戶端從 ClientMap 中刪除，
	//那麼一個 goroutine 可能會成功，而另一個 goroutine 可能會失敗。
	// 這可能會導致客戶端的連接保持打開狀態，並可能導致資源泄漏。
	sync.RWMutex
	handlers map[string]EventHandler
}

type Client struct {
	manager *Manager

	egress chan Event

	connection *websocket.Conn

	chatroom string
}

type ClientMap map[*Client]bool
