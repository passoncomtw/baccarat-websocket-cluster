// 負責web socket的升級,以及追蹤client狀態
// 故過manager 管理所有client
// 跟REST風格比 websocket是一個生命週期的管理 為方便管理會區分manager 跟client
// 聊天室 會出現群組 個別聊天 確認是否上下線 這些東東 因此納管client狀態很重要
// 同時也要確認連線存活,判斷是否放出資源,雲端廠商限制連線 及cpu ram
// 所以要有一個addclient deleteclient 來管理clientlist
// 會有client跟manager物件
// manager有serveWs addClient removeClient等物件方法

package websocketutil

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"errors"
)

var (
	/**
	websocketUpgrader is used to upgrade incomming HTTP requests into a persitent websocket connection
	*/
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	ErrEventNotSupported = errors.New("this event type is not supported")
)

func NewManager() *Manager {
	m:= &Manager{
		clients: make(ClientMap),
		handlers: make(map[string]EventHandler),
	}
	m.setupEventHandlers()
	return m
}

func(m *Manager)setupEventHandlers(){
	// 現在還模做eventhandler先做ｌｏｇ而已
	m.handlers[EventSendMessage]=SendMessageHandler
}

func(m *Manager)routeEvent(e Event, c *Client) error{
	if handler,ok:=m.handlers[e.Type]; ok{
		if err:=handler(e,c);err!= nil{
			log.Println("some error with handler",e.Type)
			return err //記得return
		}
		return nil
	} else {
		return ErrEventNotSupported
	}
	
}


// serveWS is a HTTP Handler that the has the Manager that allows connections
// 使用指針操作實際物件 用來管理ws連線升級 客戶端邏輯也可以在這
func (m *Manager) ServeWS(w http.ResponseWriter, r *http.Request) {

	log.Println("New connection")
	// Begin by upgrading the HTTP request
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Create New Client
	client := NewClient(conn, m)
	// 將client加入manager管理內,後續ws要刪除或新增會透過manager操作
	m.addClient(client)
	// 確定新增client以後就開始讀客戶給的東西囉
	// Start the read / write processes
	go client.GetClientEvent()
	go client.SendMessages()
	log.Println("ServeWs結束囉,不過因為主進程還在跑所以goroutine還活著")
}

func(m *Manager) addClient(client *Client) {
	// 讓多個goroutine依順序寫入
	m.Lock()
	defer m.Unlock()

	// 把client加進來到Map 對應bool拿來做刪除確認
	m.clients[client]=true

}

func (m *Manager) removeClient(client *Client) {
	// 讓多個goroutine依順序寫入
	m.Lock()
	defer m.Unlock()
	// map內建判斷 如果沒有該ｋｅｙ 會在第二個return值給false
	if _,ok:= m.clients[client]; ok{
		client.connection.Close()
		delete(m.clients,client)
	}
}



