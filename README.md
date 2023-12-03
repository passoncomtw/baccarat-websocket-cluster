# socket 練習

使用 gorilla/websocket 進行基本 websocket 練習,同時研究使用哪款框架作為後續使用

## 後續工作

loginCheck(jwt)->ws(origin check)->限制 payload 大小規則的 SetReadLimit->readMessages->routeEvent->Eventhandler->some goroutine function(整個過程都會用 ping pong check ws 連線對方活著)

## web socket 是啥

如果你要做一個網路聊天室,就會需要去 server 查詢對面發了啥訊息給你,之後你再發給他,比較笨的方法就是一直輪尋 server,看資訊是否有更新,於是後來開發的雙向通訊,例如 RTC,websockets 這些,其中一個目的就是要做到讓 server 可以主動告知 client 訊息到了,而不用讓 client 端一直來問！

WebSocket 標準在 RFC 645 中定義。

WebSockets 使用 HTTP 向伺服器發送初始請求。這是一個常規的 HTTP 請求，但它包含一個特殊的 HTTP 標頭 Connection: Upgrade。這告訴伺服器客戶端正在嘗試將 HTTP 請求 TCP 連線升級為長時間運行的 WebSocket。如果伺服器使用 HTTP 101 交換協定進行回應，則連線將保持活動狀態，從而使用戶端和伺服器可以雙向、全雙工發送訊息。一旦此連接達成一致，我們就可以從雙方發送和接收資料。

## websocket 常見應用

聊天程式,多人遊戲(所有客戶可以同步廣播推資料),基本需要即時數據,websocket 是很好解決方案.

##

Go 語言是由 Google 開發的，而 net/http 是標準庫的一部分，由 Go 語言的核心團隊負責維護。

## 前置準備建議

安裝 gvm,用以切換 go 版本到適合版本.

1. 安裝 gvc
   `bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)`
2. 可以依據提示重開 terminal 或者執行他給的 source 指令
3. 安裝指定版本
   ` gvm install go1.20.4 -B` `gvm install go1.21.4 -B`
4. 確定安裝`gvm list`

```zsh

gvm gos (installed)

   go1.20.4
=> go1.21.4
   system
```

5. 切換 go 版本
   `gvm use go1.20.4 ` 或者該變預設使用版本`gvm use go1.20.4 --default`

6. 用`go version`確認到指定版本.

   ```zsh
   go version go1.21.4 darwin/amd64
   ```

## 大致開發順序

endpoint->ws-> readMessages-> sendMessages 確認東西來回,之後變成 endpoint->ws->readMessages->routeEvent->Eventhandler->some goroutine function

// 負責 web socket 的升級,以及追蹤 client 狀態
// 故過 manager 管理所有 client
// 跟 REST 風格比 websocket 是一個生命週期的管理 為方便管理會區分 manager 跟 client
// 聊天室 會出現群組 個別聊天 確認是否上下線 這些東東 因此納管 client 狀態很重要
// 同時也要確認連線存活,判斷是否放出資源,雲端廠商限制連線 及 cpu ram
// 所以要有一個 addclient deleteclient 來管理 clientlist
// 會有 client 跟 manager 物件
// manager 有 serveWs addClient removeClient 等物件方法

## 循環影用問題

type 間的合理引用

```
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

// Manager is used to hold references to all Clients Registered, and Broadcasting etc
type Manager struct {
	clients ClientList
	// Using a syncMutex here to be able to lcok state before editing clients
	// Could also use Channels to block
	sync.RWMutex
}
```

這邊 Client 會用到 Manager 指標,而 Manager 用到 Client 指標,而不是直接引用對方 Client 本身,所以不算循環引用

## 參考資料

[Mastering WebSockets With Go](https://programmingpercy.tech/blog/mastering-websockets-with-go/)
