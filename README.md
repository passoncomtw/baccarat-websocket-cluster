# socket練習
使用gorilla/websocket進行基本websocket練習,同時研究使用哪款框架作為後續使用

## web socket是啥
如果你要做一個網路聊天室,就會需要去server查詢對面發了啥訊息給你,之後你再發給他,比較笨的方法就是一直輪尋server,看資訊是否有更新,於是後來開發的雙向通訊,例如RTC,websockets這些,其中一個目的就是要做到讓server可以主動告知client訊息到了,而不用讓client端一直來問！

WebSocket 標準在RFC 645中定義。

WebSockets 使用 HTTP 向伺服器發送初始請求。這是一個常規的 HTTP 請求，但它包含一個特殊的 HTTP 標頭Connection: Upgrade。這告訴伺服器客戶端正在嘗試將 HTTP 請求 TCP 連線升級為長時間運行的 WebSocket。如果伺服器使用HTTP 101 交換協定進行回應，則連線將保持活動狀態，從而使用戶端和伺服器可以雙向、全雙工發送訊息。一旦此連接達成一致，我們就可以從雙方發送和接收資料。

## websocket常見應用
聊天程式,多人遊戲(所有客戶可以同步廣播推資料),基本需要即時數據,websocket是很好解決方案.

## 
Go 語言是由 Google 開發的，而 net/http 是標準庫的一部分，由 Go 語言的核心團隊負責維護。

## 前置準備建議

安裝gvm,用以切換go版本到適合版本.
1. 安裝gvc
   `bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)`
2. 可以依據提示重開terminal或者執行他給的source指令
3. 安裝指定版本
   ` gvm install go1.20.4 -B` `gvm install go1.21.4 -B`
4. 確定安裝`gvm list`


```zsh

gvm gos (installed)

   go1.20.4
=> go1.21.4
   system
```
5. 切換go版本
   `gvm use go1.20.4 ` 或者該變預設使用版本`gvm use go1.20.4 --default`

6. 用`go version`確認到指定版本.
   
   ```zsh
   go version go1.21.4 darwin/amd64
   ```

## 參考資料
[Mastering WebSockets With Go](https://programmingpercy.tech/blog/mastering-websockets-with-go/)