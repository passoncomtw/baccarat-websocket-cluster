
## 術語

Field（欄位）： The variables inside a struct are called fields.

Receiver（接收者）： In a method, the receiver is the struct instance on which the method operates.

Instantiation（實例化）： Creating an instance of a struct, i.e., creating a concrete object of the struct.

Struct Literal（結構體字面值）： Initializing a struct instance using struct literals.

Anonymous Field（匿名字段）： A field in a struct that doesn't have a specified name.

## 語法
物件
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


## interfacee功用
物件方法的抽象化集合,定義有哪些方法,有實現者就符合該interface,主要功用是使用 interface 實現抽象和多態。這允許你在代碼中創建抽象的組件，並使用它們的共同介面進行操作.

interface應用範例,主要就是當底層操作對象可能很多種,例如，我資料的儲存可能會放在db,文件檔案,記憶體這樣,他們每個物件都要實現save,retrive_data,如果沒用interface變成db要呼叫save,文件檔案也要叫save,然後如果有錯誤要把log記錄在各自的save方法裡面,如果用interface實現就可以抽象的使用interface,把return的err拿來另外處理.

interface部分
```
type DataStore interface {
    Save(data interface{}) error
    Retrieve(id int) (interface{}, error)
} 
```
實現interface的各個struct
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
主要用在計數器,保留上下文,函數工廠,主要就是一個function的return值是另一個function,如果該被return function有return值也可以飆出來
,所以參數也可以是function.(但是會遇到某些情況是在工廠內會return err,這時候會標示return err)
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

特例當可能return function或者err

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