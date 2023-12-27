// client.go規範client物件內容,有ws連線 manager物件 通道
// 通道是拿來buffer要寫入的資訊 在既有的同步下的解決方案
//

package websocketutil

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

// NewClient is used to initialize a new Client with all required values initialized
func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan Event),
	}
}

/*
	這邊是ws接受部分邏輯

readMessages會操作client(ws)讀取js過來的訊息,因此用for卡住,所以一定要用goroutine
需要具備ws讀訊息,並判斷是否err非nil,並記錄

	收到訊息後往後邊處理邏輯丟
*/
func (c *Client) GetClientEvent() {
	// 當func return時,代表client端斷線執行操作remove client
	defer func() {
		c.manager.removeClient(c)
	}()
	for {
		// 讀取ws訊息,為了code可讀性就不跟if結合！
		messageType, payload, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break
		}
		log.Println("MessageType: ", messageType)
		log.Println("Payload: ", string(payload))
		var data Event
		// 解payload到data
		if err := json.Unmarshal(payload, &data); err != nil {
			log.Fatal("error marshalling message: %v", err)
			break // 這段要怎麼改再看看 有點硬 直接刪client
		}
		if err := c.manager.routeEvent(data, c); err != nil {
			log.Println("No handle to deal with ", err)
		}

	}
}

/*
	當需要向client端發送訊息會用此發送

為一個goroutine,用for卡住循環讀取通道以拿到訊息,因為main在跑就不會退出
*/
func (c *Client) SendMessages() {
	// 如果ws connection,代表客戶沒了,把通道關掉放出資源
	defer func() {
		// 只要把通道關掉就好,不能把client砍了
		c.connection.Close()
		log.Println("close the client ")
	}()

	/* 做循環通道讀取 , 從通道拿資料並向每個客戶端發資料*/
	for {
		select {
		case wsMessage, ok := <-c.egress:
			// 如果來源通道被關需要return關閉go routine
			if !ok {
				// 定義型態及內容
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("ws connection has close,err is", err)
					// 這邊一定要寫return,寫break只會跳出select,造成無窮迴圈
					return
				}
			}
			// 把tunnel資料解出來到json
			data, err := json.Marshal(wsMessage)
			if err != nil {
				log.Println(err)
				return
			}

			// 如果ws連線寫入錯誤要紀錄
			if err := c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println("ws connection has some err:", err)
			} else {
				// 都沒問題,資訊就會發給指定客戶
				log.Println("資訊已寄出給", c)
				log.Println("資訊是", wsMessage)
			}

		}

	}
}
