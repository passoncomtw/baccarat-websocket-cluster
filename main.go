package main

import (
	"log"
	"net/http"
	"flag"
	"time"
)
// 定義全域變數
// flag是golang命令行參數寫法,可以用-addr :8080這樣用法來取代默認值:3000
//所以default值是:3000
var addr =  flag.String("addr",":3000","for ws connection port")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "./frontend/websocket.html")
}
func main() {
	// 解析命令行參數
	flag.Parse()
	// 建立Hub
	hub:=newHub()
	// 建立預定意的handler
	hub.setEventHandlers()
	// 要把run 跑下去go routone ,
	// 他會把通道規則包含註冊那些建立出來(register, unregister...)
	go hub.run()
	// 根目錄
	http.HandleFunc("/",serveHome)
	// ws 訪問
	http.HandleFunc("/ws",func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub,w,r)
	})
	server := &http.Server{
		Addr:              *addr,
		ReadHeaderTimeout: 3 * time.Second,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
