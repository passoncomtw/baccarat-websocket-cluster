package websocketutil

import (
	"encoding/json"
	"time"
)

/*event目的是將前端來的json,這邊叫他event,解析為資料結構
會在manager放一個map,key部分是Event.type名稱,value會是function對應該type的處理
於此會定義所有event的type名稱. handler會另外放在handler資料夾
後續如果type變多在拆分.
做以下:
 1. 基本event 2. 分組寫出type名稱 3. 分類後event4.定義EventHandler格式
 5. 如果key不在的報錯
*/
// 客戶端原始數據,也就是基本Event
type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// 目前處理的Event的type名稱,屬於流量往server端
const (
	EventSendMessage = "send_message"
	EventChangeRoom  = "change_room"
)

// 目前處理的Event的type名稱,屬於流量往client端
const (
	EventNewMessage = "new_message"
)

// 定義handler function格式
type EventHandler func(event Event, c *Client) error

// EventSendMessage會用的資料結構,轉成json
type SendMessageEvent struct {
	Message string `json:"message"`
	From    string `json:"from"`
}

type NewMessageEvent struct {
	SendMessageEvent
	Sent time.Time `json:"sent"`
}

type ChangeRoomEvent struct {
	Name string `json:"name"`
}
