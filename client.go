// client.go規範client物件內容,有ws連線 manager物件 通道
// 通道是拿來buffer要寫入的資訊 在既有的同步下的解決方案
//

package main

import (
	"log"

	"github.com/gorilla/websocket"
)

// ClientList is a map used to help manage a map of clients
type ClientList map[*Client]bool

// Client is a websocket client, basically a frontend visitor
type Client struct {
	// the websocket connection
	connection *websocket.Conn

	// manager is the manager used to manage the client
	manager *Manager
	// egress is used to avoid concurrent writes on the WebSocket
	egress chan []byte
}

// NewClient is used to initialize a new Client with all required values initialized
func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan []byte),
	}
}

// readMessages will start the client to read messages and handle them
// appropriatly.
// This is suppose to be ran as a goroutine
func (c *Client) readMessages() {
	defer func() {
		// Graceful Close the Connection once this
		// function is done
		c.manager.removeClient(c)
	}()
	// Loop Forever
	for {
		// ReadMessage is used to read the next message in queue
		// in the connection
		messageType, payload, err := c.connection.ReadMessage()

		if err != nil {
			// If Connection is closed, we will Recieve an error here
			// We only want to log Strange errors, but simple Disconnection
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break // Break the loop to close conn & Cleanup
		}
		log.Println("MessageType: ", messageType)
		log.Println("Payload: ", string(payload))

		// Hack to test that WriteMessages works as intended
		// Will be replaced soon
		for wsclient := range c.manager.clients {
			wsclient.egress <- payload
			log.Println("go to write messages")
		}
	}
}

func (c *Client) writeMessages() {
	// 在 goroutine 完成時，確保執行 removeClient 方法來從 manager 中移除此客戶端
	defer func() {
		// Graceful close if this triggers a closing
		c.manager.removeClient(c)
		log.Println("i close the client")
	}()

	// 無限循環，等待要發送的消息
	for {
		// 使用 select 監聽多個 channel，這裡只監聽 egress channel
		select {
		case message, ok := <-c.egress:
			// 檢查 egress channel 是否已經關閉
			if !ok {
				// 如果 egress channel 被關閉，向前端發送關閉消息
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					// 如果寫入關閉消息時發生錯誤，記錄並結束 goroutine
					log.Println("connection closed: ", err)
				}
				// 返回以結束 goroutine
				return
			}
			// 將文本消息寫入 WebSocket 連線
			if err := c.connection.WriteMessage(websocket.TextMessage, message); err != nil {
				// 如果寫入消息時發生錯誤，記錄錯誤
				log.Println(err)
			}
			log.Println("i got ", message)
			log.Println("sent message")
		}
	}
}
