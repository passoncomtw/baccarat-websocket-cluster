package main


// function eventHandler for different type
type eventHandler func(string) error


type Hub struct {

	// Registered rooms
	rooms map[string][]*Client
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// handler different type
	handlers map[string]eventHandler
}


// 要用雙引號表示字串
func(hub *Hub) setEventHandlers(){
	hub.handlers["send_message"]= func(s string) error {}
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}



// go function
func (h *Hub) run() {
	// 每次for 都會讀一次通道,這是為了避免client concurrency寫入manager造成的同步問題
	for {
		select {
		// 新增客戶
		case client := <-h.register:
			h.clients[client] = true
		// 刪除顧客
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
			// 如果該客戶是阻塞狀態,也就是說前一訊息還卡住,就會走到default
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
