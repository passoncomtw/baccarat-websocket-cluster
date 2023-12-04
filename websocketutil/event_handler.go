// Event就是庫戶端來的原始數據 有type跟payload

package websocketutil

import (
	"encoding/json"
	"errors"
	"fmt"
	// "log"
	"time"
)

var (
	ErrPayloadNotSupported = errors.New("payload not supported")
)

// 定義handler會被放到manager裡面的handlers
func SendMessageHandler(e Event, c *Client) error {
	// 把payload解出來(要發訊息的會格式不同)
	var sendMessage SendMessageEvent
	if err := json.Unmarshal(e.Payload, &sendMessage); err != nil {
		return ErrPayloadNotSupported
	}
	// 定義新payload格式
	var newMessageBroad NewMessageEvent
	// 把資料放進去
	newMessageBroad.Sent = time.Now()
	newMessageBroad.From = sendMessage.From
	newMessageBroad.Message = sendMessage.Message
	// 換成json
	data, err := json.Marshal(newMessageBroad)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}
	/* sendMessage會發json給客戶,所以他的通道也是json資料
	,重新包裝成type,payload的Event回前端 */
	var outputEvent Event
	outputEvent.Payload = data
	outputEvent.Type = EventNewMessage

	// 把data往通道送 注意這邊是用廣播的 對所有clients
	for client, _ := range c.manager.clients {
		if client.chatroom == c.chatroom {
			client.egress <- outputEvent
		}
	}
	return nil
}

/*
	把chatroom值依據客戶端需求修改

1. 先建立空變數型態是ChangeRoomEvent
1. 就把payload值放到對應變數內
2. 依據name值修改物件client的屬性
*/
func ChangeRoomHandler(e Event, c *Client) error {
	var roomInfo ChangeRoomEvent
	if err := json.Unmarshal(e.Payload, &roomInfo); err != nil {
		return fmt.Errorf("failed to unmarshal message: %v", err)
	}
	c.chatroom = roomInfo.Name
	return nil
}
