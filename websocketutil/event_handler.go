package websocketutil

import (
	"encoding/json"
	"errors"
	"time"
	"fmt"
)

var (
	ErrPayloadNotSupported=errors.New("payload not supported")
)

// 定義handler會被放到manager裡面的handlers
func SendMessageHandler(e Event,c *Client) error{
	// 把payload解出來(要發訊息的會格式不同)
	var sendMessage SendMessageEvent
	if err:= json.Unmarshal(e.Payload,&sendMessage); err!=nil{
		return ErrPayloadNotSupported
	}
	// 定義新payload格式
	var newMessageBroad NewMessageEvent
	// 把資料放進去
	newMessageBroad.Sent=time.Now()
	newMessageBroad.From=sendMessage.From
	newMessageBroad.Message=sendMessage.Message
	// 換成json
	data,err:=json.Marshal(newMessageBroad)
	if err != nil{
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}
	/* sendMessage會發json給客戶,所以他的通道也是json資料
	,重新包裝成type,payload的Event回前端 */
	var outputEvent Event
	outputEvent.Payload= data
	outputEvent.Type = EventNewMessage

	// 把data往通道送 注意這邊是用廣播的 對所有clients
	for client,_:=range c.manager.clients {
		client.egress <- outputEvent
	}
	return nil
}