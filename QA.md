# 記錄遇到的問題

1. 一開始 websocket package 取名跟 gorila 一樣叫做 websocket 會報錯,重複 package 名稱
2. 因為把 main 跟 websocket 拆開來成為兩個 package,並在 main 做 import websocketTool,那原本只有 main package 時成立的 manager.serveWS,會報錯,因為 package 要導出的資料必須為大寫開頭,所以要把 package websocketTool 那邊的 manager.serveWS 改成大寫 ServeWs.
3. 如果上下語法有錯誤會影響到 ide 的提示判斷,例如 defer 寫錯,後面 c.connection 的物件方法會提示不出來！
4. select 本身也是一個 loop,使用 break 打斷不會斷掉外面的 for loop

```
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

5. 如果type裡面值不建立就初始化物件,會發生啥？ 一樣可以建成物件,沒放值的屬性做成對應空值！
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