# 記錄遇到的問題

1. 一開始執行`git init xxx`指令時, package 取名跟 gorila 一樣叫做 websocket ,這樣會報錯重複 package 名稱. 取 module 名稱不能跟 import 進來的撞名！
2. 因為把 main 跟 websocket 拆開來成為兩個 package,並在 main 做 import websocketTool,那原本只有 main package 時成立的 manager.serveWS,會報錯,因為 package 要導出的資料必須為大寫開頭,所以要把 package websocketTool 那邊的 manager.serveWS 改成大寫 ServeWs.
3. 如果上下語法有錯誤會影響到 ide 的提示判斷,例如 defer 寫錯,後面 c.connection 的物件方法會提示不出來！
4. select 本身也是一個 loop,使用 break 打斷不會斷掉外面的 for loop, 要關掉循環 goroutine 要用 return 關.

```go
	go func(){
		for {
			select{
			case t1,ok:=<-c:
				if !ok {
					fmt.Println("tunnel closed")
                    // 這邊用break只會斷掉select,導致for繼續執行,打印一大堆tunnel closed
					<!-- break  -->
                    // 要關掉go routine用return
                    return
				}

			t2 := time.Now()
			fmt.Println("總執行時間為:",t2.Sub(t1))
			}
		}
	}()
```

5. 如果 type 裡面值不建立就初始化物件,會發生啥？ 一樣可以建成物件,沒放值的屬性做成對應空值！

```go
type Event struct {
	Type    string
	Message string
}

func NewEvent() *Event {
	return &Event{}
}

func (e *Event) ReadAttribute() {
	fmt.Println(e.Message)
	fmt.Println(e.Type)
	fmt.Println("以上為所有屬性")
}

func main() {
	fmt.Println("hello world")
	e := NewEvent()
	e.ReadAttribute()
}
//程式執行結果:
// hello world


// 以上為所有屬性
```

7. 規劃用 manager 管理所有 client,而 client 本身也會包含 manager 物件去做斷線移除 client,這邊會擔心有循環引用問題. 後續使用指標代替 type 就沒這個問題.

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

8. 紀錄在重構 ws 方程式遇到的報錯,
   參考資料https://myapollo.com.tw/blog/golang-go-module-tutorial/

重新初始化

```go title="重新初始化"
// 1. main.go:4:2: package gowebsocket/websocketutil is not in std (/usr/local/go/src/gowebsocket/websocketutil)
// 這是因為之前初始化一個package現在要換名
// 先刪掉之前的go.mod檔案
rm -f go.mod
// 建立新module 這邊跟github串接
go mod init github.com/passoncomtw/baccarat-websocket-cluster
```

安裝相依套件,因為相依套件被刪了,這邊要重裝

```go
go get github.com/gorilla/websocket
// 也可以加指定版好 go get github.com/gorilla/websocket@v1.x.x

```
