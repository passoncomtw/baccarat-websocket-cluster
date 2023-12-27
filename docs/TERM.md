## 語法特性

1. 大寫開頭代表可匯出,小寫就是指通於該package,被import到別的package會看不到
2. 大括號用來表示控制流程的範圍，例如 if 語句、for 語句、switch 語句和函數。
3. 基本上所有東西都要定義type,除非不確定是哪種型別可以用 interface{} 來接任任何type.

```golang
type MyInterface interface{}

func myFunction(value MyInterface) {
// 做一些操作
}
```

4. for迴圈需要取key,value 例如list跟map,要用range,`for i, v := range slice{}`,一般取值可以用`for i := 0; i < len(arr); i++ {}`

5. 閉包用來return function,主要用於保留上下文,應用情境:計數器.
6. make 函數主要用於創建切片、映射和通道，並初始化它們的內部數據結構, `make(Type, size)`

```go
slice := make([]int, 5) // 創建一個包含5個元素的整數切片
mymap := make(map[string]int) // 創建一個字符串到整數的映射
ch := make(chan int) // 創建一個整數通道
slice := make([]int, 0, 10) // 創建一個初始大小為0，容量為10的整數切片

```

7. 	map內建判斷 如果沒有該key會在第二個return值給false
	`if _, ok := m.clients[client]; ok { do something}`, `csharpRating, ok := rating["C#"]`
8. 通道使用原理,定義丟進來的物件型態, 資料傳入 `chan <- 物件`, 資料傳出有多個方式 `t1 := <-c`, `for t1:= range c {}`,


```golang
for {
    select {
        # 通道一
        case value:= <- c:
        # 使用通道一給進來的值做switch case判斷對應處理
        switch value {
        case "add":
            doSomething
        case "remove":
            doOtherthing
        }
        # 通道二
        case value2:= <-c2:
        # 通道二值判斷
        switch value2 {
            case "hekko":
                doHekko
            case "hello":
                doHello
        }

    }
}
```

## for 寶典

對 slice,array,有基本的跟用 range

```
// 只取值
func main() {
    arr := []int{1, 2, 3, 4, 5}

    for i := 0; i < len(arr); i++ {
        fmt.Println(arr[i])
    }
}
// 取角位,value
slice := []int{1, 2, 3, 4, 5}

for i, v := range slice {
    fmt.Println(i, v)
}

```

對 str

```
func main() {
    str := "Hello, world!"

    for i := 0; i < len(str); i++ {
        fmt.Println(str[i])
    }
}
```

對 map,主要都用 range

```
map := map[string]int{
    "a": 1,
    "b": 2,
    "c": 3,
}

for k, v := range map {
    fmt.Println(k, v)
}
```

## 通道使用

1. 任何傳送將會被阻塞，直到資料被讀出
   白話解釋：
   假設有兩個人 A 和 B，A 要給 B 一張紙。A 拿著紙，準備給 B 時，發現 B 還沒有準備好接收。因此，A 只能站在那裡等待 B 準備好。在 A 等待的這段時間內，A 不會做任何事情，也不會消耗任何資源。
2. 發送端用 goroutine
   基本示範:

```go
package main

import (
	"fmt"
	"time"
)
func main() {
    c := make(chan time.Time)
    fmt.Println("start test")
	// 送資料過來的那一端要用goroutine
    go func() {
        t1 := time.Now()
        fmt.Println("通道載送時間過來囉",t1)
        c <- t1
    }()
	fmt.Println("我要開始睡")
	time.Sleep(5*time.Second)
	t2 := time.Now()
	// 接收端不要用goroutine 不然整個程式跑完就中止了 不然就要循環通道
	t1 := <-c
	fmt.Println("總執行時間為:",t2.Sub(t1))
}
```

錯誤示範:在 goroutine 外面做 for 迴圈,會導致並行同步出問題,也不到何時該 close 通道.

```go
package main

import (
	"fmt"
	"time"
)
func main() {
    c := make(chan time.Time)
    fmt.Println("start test")
	for i:=0;i<5;i++ {
	// 送資料過來的那一端要用goroutine
        go func() {
            t1 := time.Now()
            fmt.Println("通道載送時間過來囉",t1)
            c <- t1
        }()
    }
	fmt.Println("我要開始睡")
	time.Sleep(2*time.Second)
	t2 := time.Now()
	// 接收端不要用goroutine 不然整個程式跑完就中止了
    for t1:= range c{
		fmt.Println("總執行時間為:",t2.Sub(t1))
	}
}
```

正確寫法: 在 goroutine 裡面寫 for 迴圈依序阻塞的往通道給資料！

```go
package main

import (
	"fmt"
	"time"
)
func main() {
    c := make(chan time.Time)
    fmt.Println("start test")

	// 送資料過來的那一端要用goroutine
    go func() {
		for i:=0; i<4; i++ {
        t1 := time.Now()
        fmt.Println("通道載送時間過來囉",t1)
        c <- t1
		}
		close(c)
    }()
	fmt.Println("我要開始睡")

	// 接收端不要用goroutine 不然整個程式跑完就中止了
    for t1:= range c{
		time.Sleep(2*time.Second)

		t2 := time.Now()
		fmt.Println("總執行時間為:",t2.Sub(t1))
	}
}
//start test
我要開始睡
通道載送時間過來囉 2023-12-03 02:15:40.993265 +0800 CST m=+0.000210042
通道載送時間過來囉 2023-12-03 02:15:40.993567 +0800 CST m=+0.000512334
總執行時間為: 2.001394834s
通道載送時間過來囉 2023-12-03 02:15:42.994909 +0800 CST m=+2.001794917
總執行時間為: 4.001366s
通道載送時間過來囉 2023-12-03 02:15:44.99513 +0800 CST m=+4.001955834
總執行時間為: 4.001191792s
總執行時間為: 4.002174125s
```

A->Ｂ送,確認 B 已拿資料(<-),就馬上回去做 A,又馬上往Ｂ送,所以第二次Ｂ會被前一次影響多而多了兩秒的執行時間,意思就是通道兩端的執行都會是被阻塞的！

## map

map 內建有判斷是否存在 key 的方式
make for slice array map 指標創建,new for 其他的做指標（int string 自定義的 type)

```
// 初始化一個字典
rating := map[string]float32{"C":5, "Go":4.5, "Python":4.5, "C++":2 }
// map 有兩個回傳值，第二個回傳值，如果不存在 key，那麼 ok 為 false，如果存在 ok 為 true
csharpRating, ok := rating["C#"]
```

## Embedding

golang 結構體 enbed 其他結構體,會同時具有 embed 的結構體的屬性跟方法,embed 之後會繼承該物件的所有方法.例如:Person 有發法 SayHello,在 Student embed Person 之後可以直經使用 SayHello 方法.

```
type Person struct {
    name string
}

type Student struct {
    Person // 內嵌 Person 結構體 這樣會繼承方法
    school string
}
// 傳入receiver p, ()代表無參數 ,執行{}操作
func (p *Person) SayHello() {
    fmt.Println("Hello, my name is", p.name)
}
s := Student{
    Person: Person{
        name: "John Doe",
    },
    school: "MIT",
}
s.SayHello() // "Hello, my name is John Doe"
```

某些情況一下也會把物件指標放入 以方便操作該物件

```
type Client struct {
	manager *Manager

	egress chan []byte

	connection *websocket.Conn
}
如果物件client要斷線失聯後會操作manager刪除client,這時候會呼叫子物件manager去做removeClient.
....一個前面直行讀資料的function,當function做return時會先叫
	defer func() {
		// Graceful Close the Connection once this
		// function is done
		c.manager.removeClient(c)
	}()
```

## defer

當函數返回時，會從堆棧中彈出所有 defer 語句，並按 LIFO 順序執行它們。

```
package main

import "fmt"

func openFile() (*os.File, error) {
    f, err := os.Open("test.txt")
    if err != nil {
        return nil, err
    }

    defer f.Close()

    fmt.Println("File opened")
    return f, nil
}

func main() {
    f, err := openFile()
    if err != nil {
        fmt.Println(err)
    }
}
//File opened
```

## 術語

Field（欄位）： The variables inside a struct are called fields.

Receiver（接收者）： In a method, the receiver is the struct instance on which the method operates.

Instantiation（實例化）： Creating an instance of a struct, i.e., creating a concrete object of the struct.

Struct Literal（結構體字面值）： Initializing a struct instance using struct literals.

Anonymous Field（匿名字段）： A field in a struct that doesn't have a specified name.

## 語法

### 物件

```

// 宣告一個新的型別
type person struct {
    name string
    age int
}
// 賦值初始化
tom.name, tom.age = "Tom", 18

// 兩個欄位都寫清楚的初始化
bob := person{age:25, name:"Bob"}

// 按照 struct 定義順序初始化值
paul := person{"Paul", 43}
```

### 物件方法

```
// 透過物件指標操作本物件,操作物件為m
// 物件方法參數為client
func (m *Manager) removeClient(client *Client) {
	// 讓多個goroutine依順序寫入
	m.Lock()
	defer m.Unlock()
	if _,ok:= m.clients[client]; ok{
		client.connection.Close()
		delete(m.clients,client)
	}
}
// 使用
...省略
manager.removeClient(c)
```

函式

只有以下方法

```
// 回傳 a、b 中最大值.
func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

func main(){
    max_xy := max(x, y) //呼叫函式 max(x, y)
}
```

匿名函示

```
func withFile(filename string, callback func(file *File) error) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    return callback(file)
}

func main() {
    err := withFile("example.txt", func(file *File) error {
        // 使用 file 執行一些操作
        return nil
    })

    if err != nil {
        fmt.Println("Error:", err)
    }
}
```

## interfacee 功用

物件方法的抽象化集合,定義有哪些方法,有實現者就符合該 interface,主要功用是使用 interface 實現抽象和多態。這允許你在代碼中創建抽象的組件，並使用它們的共同介面進行操作.

interface 應用範例,主要就是當底層操作對象可能很多種,例如，我資料的儲存可能會放在 db,文件檔案,記憶體這樣,他們每個物件都要實現 save,retrive_data,如果沒用 interface 變成 db 要呼叫 save,文件檔案也要叫 save,然後如果有錯誤要把 log 記錄在各自的 save 方法裡面,如果用 interface 實現就可以抽象的使用 interface,把 return 的 err 拿來另外處理.

interface 部分

```
type DataStore interface {
    Save(data interface{}) error
    Retrieve(id int) (interface{}, error)
}
```

實現 interface 的各個 struct

```
// **memory_store.go**

package main

import (
	"fmt"
	"sync"
)

// MemoryStore 實現 DataStore 接口
type MemoryStore struct {
	data  map[int]interface{}
	mutex sync.RWMutex
}

// Save 將數據保存到內存中
func (m *MemoryStore) Save(data interface{}) error {
	id := len(m.data) + 1
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.data[id] = data
	return nil
}

// Retrieve 從內存中檢索數據
func (m *MemoryStore) Retrieve(id int) (interface{}, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	if data, ok := m.data[id]; ok {
		return data, nil
	}
	return nil, fmt.Errorf("Data not found for ID: %d", id)
}

// **file_store.go**

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

// FileStore 實現 DataStore 接口
type FileStore struct {
	filename string
	mutex    sync.RWMutex
}

// Save 將數據保存到文件中
func (f *FileStore) Save(data interface{}) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	fileData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(f.filename, fileData, 0644)
}

// Retrieve 從文件中檢索數據
func (f *FileStore) Retrieve(id int) (interface{}, error) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	fileData, err := ioutil.ReadFile(f.filename)
	if err != nil {
		return nil, err
	}

	var storedData interface{}
	if err := json.Unmarshal(fileData, &storedData); err != nil {
		return nil, err
	}

	return storedData, nil
}

// **main.go**

package main

import "fmt"

func main() {
	// 使用 MemoryStore
	memoryStore := &MemoryStore{data: make(map[int]interface{})}
	useDataStore(memoryStore)

	// 使用 FileStore
	fileStore := &FileStore{filename: "data.json"}
	useDataStore(fileStore)
}

func useDataStore(ds DataStore) {
	// 保存數據
	err := ds.Save(map[string]interface{}{"key": "value"})
	if err != nil {
		fmt.Println("Save error:", err)
		return
	}

	// 檢索數據
	data, err := ds.Retrieve(1)
	if err != nil {
		fmt.Println("Retrieve error:", err)
		return
	}

	fmt.Println("Retrieved data:", data)
}



```

## 閉包使用

主要用在計數器,保留上下文,函數工廠,主要就是一個 function 的 return 值是另一個 function,如果該被 return function 有 return 值也可以飆出來
,所以參數也可以是 function.(但是會遇到某些情況是在工廠內會 return err,這時候會標示 return err)

- 計數器

```格式
// counter 是一個閉包，保留了內部狀態 count
func counter() func() int { // 會return func該func的return值是int
    count := 0
    return func() int {
        count++
        return count
    }
}

func main() {
    // 創建一個計數器
    increment := counter()

    // 使用計數器多次
    fmt.Println(increment()) // 1
    fmt.Println(increment()) // 2
    fmt.Println(increment()) // 3
}
```

- 函數工廠

```
package main

import "fmt"

// multiplier 是一個函數工廠，返回一個新的乘法函數
func multiplier(factor int) func(int) int {
    return func(x int) int {
        return x * factor
    }
}

func main() {
    // 創建一個乘法函數，以 5 為因子
    timesFive := multiplier(5)

    // 使用乘法函數
    result := timesFive(3)
    fmt.Println(result) // 15
}

```

- 保留上下文

```
package main

import "fmt"

// greeter 是一個閉包，可以使用外部的 greeting 變數
func greeter(greeting string) func(string) {
    return func(name string) {
        fmt.Println(greeting, name)
    }
}

func main() {
    // 創建一個問候函數，使用 "Hello" 作為 greeting
    sayHello := greeter("Hello")

    // 使用問候函數
    sayHello("Alice") // Hello Alice
}


```

特例當可能 return function 或者 err

```
// 傳入function 不過因為可能return err或function所以直接標示最後結果
func withFile(filename string, callback func(file *File) error) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    return callback(file)
}

func main() {
    err := withFile("example.txt", func(file *File) error {
        // 使用 file 執行一些操作
        return nil
    })

    if err != nil {
        fmt.Println("Error:", err)
    }
}

```
