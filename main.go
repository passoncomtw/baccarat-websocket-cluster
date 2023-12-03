package main

import (
	"gowebsocket/websocketutil"
	"log"
	"net/http"
)

// setupAPI will start all Routes and their Handlers
func setupAPI() {
	// Create a Manager instance used to handle WebSocket Connections
	manager := websocketutil.NewManager()

	// Serve the ./frontend directory at Route /
	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/ws", manager.ServeWS) //
}
func main() {
	setupAPI()

	// Serve on port :3000, fudge yeah hardcoded port
	log.Fatal(http.ListenAndServe(":3000", nil))
}
