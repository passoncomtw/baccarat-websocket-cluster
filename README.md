# socket 使用指南

使用 gorilla/websocket 進行基本 websocket 練習,同時研究使用哪款框架作為後續使用.

後續記得說明 handler 為啥不選擇用通道,以及可能會用到通道的情境.
單純設計 room,功用在 type 是帶入 room id 的話,把同 room id 放給有同 roomid 的人,這部分是不用的,也就是說 client 加 roomid,發訊息對同 roomid 發送

1. 完成 roomid 設定
2.

使用步驟

- 安裝套件

```
go mod download
```

- 執行程式,目前起在 3000 port,去 localhost:3000 看

```
go run *.go
```

目前可以接收格式以及功能

- 填入訊息到 message 後廣播訊息給所有人,會用 type 判斷要解析的 payload 格式

```json
{ "type": "send_message", "payload": { "message": "sdfsdf", "from": "percy" } }
```

回應

```json
{
  "type": "new_message",
  "payload": {
    "message": "sdfsdf",
    "from": "percy",
    "sent": "2023-12-04T01:35:59.83419+08:00"
  }
}
```

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

## web socket 是啥

如果你要做一個網路聊天室,就會需要去 server 查詢對面發了啥訊息給你,之後你再發給他,比較笨的方法就是一直輪尋 server,看資訊是否有更新,於是後來開發的雙向通訊,例如 RTC,websockets 這些,其中一個目的就是要做到讓 server 可以主動告知 client 訊息到了,而不用讓 client 端一直來問！

WebSocket 標準在 RFC 645 中定義。

WebSockets 使用 HTTP 向伺服器發送初始請求。這是一個常規的 HTTP 請求，但它包含一個特殊的 HTTP 標頭 Connection: Upgrade。這告訴伺服器客戶端正在嘗試將 HTTP 請求 TCP 連線升級為長時間運行的 WebSocket。如果伺服器使用 HTTP 101 交換協定進行回應，則連線將保持活動狀態，從而使用戶端和伺服器可以雙向、全雙工發送訊息。一旦此連接達成一致，我們就可以從雙方發送和接收資料。

## 大致開發順序

endpoint->ws-> readMessages-> sendMessages 確認東西來回,之後變成 endpoint->ws->readMessages->routeEvent->Eventhandler->some goroutine function -> rabbitmq 加從 API 端拿訊息示範 以及給訊息示範.

跟 REST 風格比 websocket 是一個生命週期的管理 為方便管理會區分 manager 跟 client,聊天室 會出現群組 個別聊天 確認是否上下線 這些東東 因此納管 client 狀態很重要,同時也要確認連線存活,判斷是否放出資源,雲端廠商限制連線 及 cpu ram

## 後續工作

- loginCheck
- ws(origin check)
- payload set read limit
- heartbeat
- room
- 客戶資料擴充

```go
type wsClients struct {
	Conn *websocket.Conn `json:"conn"`

	RemoteAddr string `json:"remote_addr"`

	Uid string `json:"uid"`

	Username string `json:"username"`

	RoomId string `json:"room_id"`

	AvatarId string `json:"avatar_id"`
}
type msgData struct {
	Uid      string        `json:"uid"`
	Username string        `json:"username"`
	AvatarId string        `json:"avatar_id"`
	ToUid    string        `json:"to_uid"`
	Content  string        `json:"content"`
	ImageUrl string        `json:"image_url"`
	RoomId   string        `json:"room_id"`
	Count    int           `json:"count"`
	List     []interface{} `json:"list"`
	Time     int64         `json:"time"`
}
```

## 參考資料

1. [Mastering WebSockets With Go](https://programmingpercy.tech/blog/mastering-websockets-with-go/)

2. https://github.com/hezhizheng/go-gin-chat/blob/master/ws/go_ws/serve.go
3. https://github.com/gorilla/websocket/blob/main/examples/chat/client.go
