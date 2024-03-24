package main

import (
	"encoding/json"
	"time"
	"errors"
	"fmt"
)


var (
	egressMessageType="newMessage"
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
// client ws data which will send in
type Event struct {
	Type string `json:"type"`
	// let client send any data
	Payload json.RawMessage `json:"payload"`

}



// type to send message to other client
type ingressMessageEvent struct{
	Message string `json:"message"`
	From  string `json:"from"`
	RoomID string `json:"roomID"`
}

type egressPayloadEvent struct {
	ingressMessageEvent
	SentTime time.Time `json:"sentTime"`
}
// 定義handler function格式
type EventHandler func(event Event, c *Client) error

func sendMessageEvent(event Event, c *Client) error{
	// 做空物件給json賦予值
	var ingressMessage ingressMessageEvent
	// 把近來資訊做解析
	if err:=json.Unmarshal(event.Payload,&ingressMessage);err != nil{
		return errors.New("field not supported")
	}
	// 包要出去的payload
	egressPayload:=egressPayloadEvent{
		ingressMessageEvent: ingressMessage,
		SentTime: time.Now(),
	}

	payload,err:=json.Marshal(egressPayload)
	if err != nil {
		return fmt.Errorf("failed to marshal  message: %v", err)
	}

	// 包要出去完整的json
	// 先做 struct of Event
	egressEvent:=Event{Payload: payload,Type: egressMessageType}
	egressMessage,err:=json.Marshal(egressEvent)
	if err != nil{
		return fmt.Errorf("failed to marshal  message: %v", err)
	}
	// 往client通道送資訊
	for client,_ := range c.hub.clients{
		if client.roomName == c.roomName{
			// 送進來的資訊會被goroutine那邊通套拿去做write message
			c.send <- egressMessage
		}
	}
	return nil
}

// // EventSendMessage會用的資料結構,轉成json
// type SendMessageEvent struct {
// 	Message string `json:"message"`
// 	From    string `json:"from"`
// }

// type NewMessageEvent struct {
// 	SendMessageEvent
// 	Sent time.Time `json:"sent"`
// }

// type ChangeRoomEvent struct {
// 	Name string `json:"name"`
// }
